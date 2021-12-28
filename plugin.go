package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var HandlerRegisterer = registerer("cached-router-plugin")

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     extra["host"].(string),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := fmt.Sprintf("%x", sha1.Sum([]byte(req.URL.String())))
		fmt.Println("key: ", key)

		val, err := rdb.Get(ctx, key).Result()
		fmt.Println(err)
		if err == redis.Nil {
			// TODO add header caching
			c := httptest.NewRecorder()
			handler.ServeHTTP(c, req)

			for k, v := range c.HeaderMap {
				w.Header()[k] = v
			}

			w.WriteHeader(c.Code)
			content := c.Body.String()
			fmt.Println("body: ", content)
			rdb.Set(ctx, key, content, 0)

			w.Write([]byte(content))
		} else if err != nil {
			panic(err)
		} else {
			w.Write([]byte(val))
		}
	}), nil
}

func init() {
	fmt.Println("cached-router-plugin handler loaded!!!")
}

func main() {}
