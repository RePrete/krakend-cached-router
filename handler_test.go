package main

import (
	"github.com/RePrete/krakend-cached-router/mocks"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const responseContent = `{"key":"value"}"`

func TestCachedHandler_ServeHTTP_NoCacheHit(t *testing.T) {
	marshaller, handler, responseWriter, cachedHandler, request, ttl := setup(t)
	marshaller.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, redis.Nil)
	handler.
		EXPECT().
		ServeHTTP(
			// Not ssure is correct
			gomock.AssignableToTypeOf(&httptest.ResponseRecorder{}),
			&request,
		).
		Do(func(recorder *httptest.ResponseRecorder, _ *http.Request) {
			recorder.WriteHeader(200)
			recorder.Header()[`Content-Type`] = []string{`application/json`}
			recorder.Write([]byte(`{"key":"value"}`))
		})

	responseWriter.
		EXPECT().
		Header().
		Return(http.Header{})
	responseWriter.
		EXPECT().
		WriteHeader(200)
	responseWriter.
		EXPECT().
		Write([]byte(`{"key":"value"}`))

	set1 := marshaller.
		EXPECT().
		Set(
			gomock.AssignableToTypeOf(``),
			gomock.AssignableToTypeOf(``), // A bit flaky, should test the body is the same
			ttl,
		)

	set2 := marshaller.
		EXPECT().
		Set(
			gomock.AssignableToTypeOf(``),
			http.Header{`Content-Type`: []string{`application/json`}},
			ttl,
		)

	gomock.InOrder(set1, set2)

	cachedHandler.ServeHTTP(responseWriter, &request)
}

func setup(t *testing.T) (
	*mocks.MockClientMarshalerInterface,
	*mocks.MockHandler,
	*mocks.MockResponseWriter,
	CachedHandler,
	http.Request,
	time.Duration,
) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockClientMarshalerInterface(ctrl)
	h := mocks.NewMockHandler(ctrl)
	w := mocks.NewMockResponseWriter(ctrl)
	ttl, _ := time.ParseDuration(`60s`)

	handler := CachedHandler{
		handler: h,
		client:  m,
		ttl:     ttl,
	}
	request := http.Request{
		URL: &url.URL{
			Host: `https://petstore.swagger.io`,
			Path: `/v2/store/inventory`,
		},
		Body: io.NopCloser(strings.NewReader(responseContent)),
	}
	return m, h, w, handler, request, ttl
}
