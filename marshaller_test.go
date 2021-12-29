package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RePrete/cached-router-plugin/mocks"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestMarshaler_GetWithString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockCmdable(ctrl)
	ctx := context.Background()
	m := NewMarshaller(client, ctx)

	var input string
	key := `key`
	value := `This is a simple string`
	status := redis.NewStringResult(value, nil)

	client.
		EXPECT().
		Get(ctx, key).
		Return(status).
		AnyTimes()

	result, err := m.Get(key, input)
	if !reflect.DeepEqual(value, result) {
		t.Errorf("m.Get() = %v, want %v", result, value)
	}
	if err != nil {
		t.Errorf("Unexpected err from m.Get()")
	}
}

func TestMarshaler_GetWithObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockCmdable(ctrl)
	ctx := context.Background()
	m := NewMarshaller(client, ctx)

	header := http.Header{}
	header.Set(`Content-Type`, `application/json; charset=utf-8`)
	header.Set(`X-Krakend-Completed`, `true`)

	key := `key`
	marshalledHeader, _ := json.Marshal(header)
	status := redis.NewStringResult(string(marshalledHeader), nil)
	var expectedObj http.Header

	client.
		EXPECT().
		Get(ctx, key).
		Return(status).
		AnyTimes()

	tmp, err := m.Get(key, expectedObj)
	if err != nil {
		t.Errorf("m.Get() unexpected err = %v", err)
	}
	result := tmp.(map[string]interface{})
	for k, v := range result {
		// Shitty workaround, needs to be improved
		expected := fmt.Sprintf("[%s]", header.Get(k))
		if fmt.Sprint(v) != expected {
			t.Errorf("heder property for key %v not equal:  %v, want %v", k, v, header.Get(k))
		}
	}
}

func TestMarshaler_SetWithObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockCmdable(ctrl)
	ctx := context.Background()
	m := NewMarshaller(client, ctx)

	input := []string{}
	key := `key`
	duration, _ := time.ParseDuration(`60s`)
	status := &redis.StatusCmd{}
	expected := status.Err()

	marshalledObject, _ := json.Marshal(input)
	client.
		EXPECT().
		Set(ctx, key, marshalledObject, duration).
		Return(status).
		AnyTimes()

	result := m.Set(key, input, duration)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("m.Set() = %v, want %v", result, expected)
	}
}

func TestMarshaler_SetWithString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockCmdable(ctrl)
	ctx := context.Background()
	m := NewMarshaller(client, ctx)

	input := `This is a sample string`
	key := `key`
	duration, _ := time.ParseDuration(`60s`)
	status := &redis.StatusCmd{}
	expected := status.Err()

	client.
		EXPECT().
		Set(ctx, key, input, duration).
		Return(status).
		AnyTimes()

	result := m.Set(key, input, duration)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("m.Set() = %v, want %v", result, expected)
	}
}

func TestMarshaler_SetWithErrorFromSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockCmdable(ctrl)
	ctx := context.Background()
	m := NewMarshaller(client, ctx)

	input := `This is a sample string`
	key := `key`
	duration, _ := time.ParseDuration(`60s`)
	status := &redis.StatusCmd{}
	status.SetErr(errors.New(`sample error`))
	expected := status.Err()

	client.
		EXPECT().
		Set(ctx, key, input, duration).
		Return(status).
		AnyTimes()

	result := m.Set(key, input, duration)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("m.Set() = %v, want %v", result, expected)
	}
}

func TestNewMarshaller(t *testing.T) {
	ctx := context.Background()
	client := &redis.Client{}
	m := NewMarshaller(client, ctx)
	if !reflect.DeepEqual(ctx, m.context) {
		t.Errorf("NewMarshaller().context = %v, want %v", m.context, ctx)
	}
	if !reflect.DeepEqual(client, m.client) {
		t.Errorf("NewMarshaller().client = %v, want %v", m.client, client)
	}
}
