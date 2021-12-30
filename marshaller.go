package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClientMarshalerInterface interface {
	Get(key string, returnObj interface{}) (interface{}, error)
	Set(key string, object interface{}, expiration time.Duration) error
}

type ClientMarshaler struct {
	client  redis.Cmdable
	context context.Context
}

func NewMarshaller(cache redis.Cmdable, ctx context.Context) *ClientMarshaler {
	return &ClientMarshaler{
		client:  cache,
		context: ctx,
	}
}

func (c *ClientMarshaler) Get(key string, returnObj interface{}) (interface{}, error) {
	result, err := c.client.Get(c.context, key).Result()
	if err != nil {
		return nil, err
	}
	switch returnObj.(type) {
	case string:
		return result, nil
	}

	// Default case, meaning desired return is a not a string
	unmarshallErr := json.Unmarshal([]byte(result), returnObj)
	if unmarshallErr != nil {
		return nil, unmarshallErr
	}
	return returnObj, nil
}

func (c *ClientMarshaler) Set(key string, object interface{}, expiration time.Duration) error {
	var value interface{}
	var err error

	switch object.(type) {
	case string:
		value = object.(string)
	default:
		value, err = json.Marshal(object)
		if err != nil {
			return err
		}
	}

	return c.client.Set(c.context, key, value, expiration).Err()
}
