package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

type CachedHandler struct {
	handler http.Handler
	client  ClientMarshalerInterface
	ttl     time.Duration
}

func (h *CachedHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key, keyError := createCacheKey(request)
	if keyError != nil {
		// Add log error
		h.handler.ServeHTTP(writer, request)
	}
	err := h.tryServeFromCache(key, writer)
	if err != nil {
		recorder := httptest.NewRecorder()
		h.handler.ServeHTTP(recorder, request)
		h.respond(writer, recorder)
		defer h.storeInCache(key, recorder, h.ttl)
	}
}

func (h *CachedHandler) tryServeFromCache(key string, writer http.ResponseWriter) error {
	var content string
	var header http.Header

	tmpContent, contentError := h.client.Get(key+".content", content)
	if contentError != nil {
		return contentError
	}

	_, headerError := h.client.Get(key+".header", &header)
	if headerError != nil {
		return headerError
	}

	content = tmpContent.(string)
	for k, v := range header {
		writer.Header()[k] = v
	}
	writer.Write([]byte(content))
	return nil
}

func (h *CachedHandler) storeInCache(key string, c *httptest.ResponseRecorder, ttl time.Duration) {
	content := c.Body.String()
	header := c.Header()
	h.client.Set(key+".content", content, ttl)
	h.client.Set(key+".header", header, ttl)
}

func (h *CachedHandler) respond(w http.ResponseWriter, c *httptest.ResponseRecorder) {
	content := c.Body.String()
	for k, v := range c.Header() {
		w.Header()[k] = v
	}

	w.WriteHeader(c.Code)
	w.Write([]byte(content))
}

func createCacheKey(req *http.Request) (string, error) {
	if req.URL == nil {
		return ``, errors.New(`url is null`)
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(req.URL.String()))), nil
}
