package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

// Marshaler is the struct that marshal and unmarshal client values
type Marshaler struct {
	client  redis.Cmdable
	context context.Context
}

// New creates a new marshaler that marshals/unmarshals client values
func NewMarshaller(cache redis.Cmdable, ctx context.Context) *Marshaler {
	return &Marshaler{
		client:  cache,
		context: ctx,
	}
}

// Get obtains a value from client and unmarshal value with given object
func (c *Marshaler) Get(key string, returnObj interface{}) (interface{}, error) {
	result, err := c.client.Get(c.context, key).Result()
	if err != nil {
		return nil, err
	}
	switch returnObj.(type) {
	case string:
		return result, nil
	}

	// Default case, meaning desired return is a not a string
	unmarshallErr := json.Unmarshal([]byte(result), &returnObj)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return returnObj, nil
}

func (c *Marshaler) Set(key string, object interface{}, expiration time.Duration) error {
	var value interface{}
	var err error

	switch object.(type) {
	case string:
		value = object.(string)
	default:
		value, err = json.Marshal(object)
		if err != nil {
			return nil
		}
	}

	return c.client.Set(c.context, key, value, expiration).Err()
}
