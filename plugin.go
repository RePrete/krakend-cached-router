package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
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
	ttlConfig := extra["cache_ttl"]
	if ttlConfig == nil {
		ttlConfig = "60s"
	}
	ttl, _ := time.ParseDuration(ttlConfig.(string))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     extra["host"].(string),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	marshaller := NewMarshaller(redisClient, ctx)

	cachedHandler := CachedHandler{
		handler: handler,
		ttl:     ttl,
		client:  marshaller,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cachedHandler.ServeHTTP(w, req)
	}), nil
}

func createCacheKey(req *http.Request) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(req.URL.String())))
}

func init() {
	fmt.Println("cached-router-plugin handler loaded")
}

func main() {}
