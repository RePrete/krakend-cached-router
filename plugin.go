package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var HandlerRegisterer = registerer("cached-router")

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	ttlConfig := extra["ttl"]
	if ttlConfig == nil {
		ttlConfig = "60s"
	}
	ttl, _ := time.ParseDuration(ttlConfig.(string))

	rdb := redis.NewClient(&redis.Options{
		Addr:     extra["host"].(string),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	m := NewMarshaller(rdb, ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := fmt.Sprintf("%x", sha1.Sum([]byte(req.URL.String())))
		content, header, err := readFromCache(m, key)

		if err != nil {
			// TODO create a new recorder that will cache things
			c := httptest.NewRecorder()
			handler.ServeHTTP(c, req)

			fromRecorderToCache(c, m, key, ttl)
			fromRecorderToResponse(w, c)
		} else {
			fromCacheToResposne(w, header, content)
		}
	}), nil
}

func readFromCache(m *Marshaler, key string) (string, http.Header, error) {
	var content string
	var header http.Header

	tmpContent, contentError := m.Get(key+".content", content)
	if contentError != nil {
		return "", nil, contentError
	}

	_, headerError := m.Get(key+".header", &header)
	if headerError != nil {
		return "", nil, headerError
	}

	content = tmpContent.(string)
	return content, header, nil
}

func fromCacheToResposne(w http.ResponseWriter, header http.Header, content string) {
	for k, v := range header {
		w.Header()[k] = v
	}

	w.Write([]byte(content))
}

func fromRecorderToResponse(w http.ResponseWriter, c *httptest.ResponseRecorder) {
	content := c.Body.String()
	for k, v := range c.HeaderMap {
		w.Header()[k] = v
	}

	w.WriteHeader(c.Code)
	w.Write([]byte(content))
}

func fromRecorderToCache(c *httptest.ResponseRecorder, m *Marshaler, key string, ttl time.Duration) {
	content := c.Body.String()
	header := c.Header()
	m.Set(key+".content", content, ttl)
	m.Set(key+".header", header, ttl)
}

func init() {
	fmt.Println("cached-router-plugin handler loaded!!!")
}

func main() {}
