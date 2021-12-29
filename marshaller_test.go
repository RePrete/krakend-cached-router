package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"reflect"
	"testing"
	"time"
)

func TestMarshaler_Get(t *testing.T) {
	type fields struct {
		client  *redis.Client
		context context.Context
	}
	type args struct {
		key       string
		returnObj interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Marshaler{
				client:  tt.fields.client,
				context: tt.fields.context,
			}
			got, err := c.Get(tt.args.key, tt.args.returnObj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshaler_Set(t *testing.T) {
	type fields struct {
		client  *redis.Client
		context context.Context
	}
	type args struct {
		key        string
		object     interface{}
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Marshaler{
				client:  tt.fields.client,
				context: tt.fields.context,
			}
			if err := c.Set(tt.args.key, tt.args.object, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMarshaller(t *testing.T) {
	type args struct {
		cache *redis.Client
		ctx   context.Context
	}
	tests := []struct {
		name string
		args args
		want *Marshaler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMarshaller(tt.args.cache, tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMarshaller() = %v, want %v", got, tt.want)
			}
		})
	}
}
