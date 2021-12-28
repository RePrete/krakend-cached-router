package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
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
		ttlConfig = "60"
	}
	ttl, _ := time.ParseDuration(ttlConfig.(string))

	rdb := redis.NewClient(&redis.Options{
		Addr:     extra["host"].(string),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := fmt.Sprintf("%x", sha1.Sum([]byte(req.URL.String())))

		content, contentError := rdb.Get(ctx, key+".content").Result()
		header, headerError := rdb.Get(ctx, key+".header").Result()

		if (contentError == redis.Nil || headerError == redis.Nil) && (contentError != nil || headerError != nil) {
			c := httptest.NewRecorder()
			handler.ServeHTTP(c, req)

			for k, v := range c.HeaderMap {
				w.Header()[k] = v
			}

			w.WriteHeader(c.Code)
			content := c.Body.String()

			header, _ := json.Marshal(c.Header())
			rdb.Set(ctx, key+".content", content, ttl)
			rdb.Set(ctx, key+".header", header, ttl)

			w.Write([]byte(content))
		} else {
			var responseHeader http.Header
			err := json.Unmarshal([]byte(header), &responseHeader)
			if err != nil {
				fmt.Println(err)
				panic("error in header unmarshall")
			}
			fmt.Println(responseHeader)
			for k, v := range responseHeader {
				w.Header()[k] = v
			}

			w.Write([]byte(content))
		}
	}), nil
}

func init() {
	fmt.Println("cached-router-plugin handler loaded!!!")
}

func main() {}
