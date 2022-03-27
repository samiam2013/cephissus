package main

import (
	"bytes"
	"log"
	"net/http"
	"testing"
)

var _ http.ResponseWriter = testResp{}

type testResp struct {
	header map[string][]string
	body   bytes.Buffer
	Status int
}

// Header implements http.ResponseWriter
func (r testResp) Header() http.Header {
	return r.header
}

// Write implements http.ResponseWriter
func (r testResp) Write(b []byte) (n int, err error) {
	return r.body.Write(b)
}

// WriteHeader implements http.ResponseWriter
func (r testResp) WriteHeader(statusCode int) {
	r.Status = statusCode
}

// BenchmarkNarcissus _
func BenchmarkNarcissus(b *testing.B) {
	body := bytes.NewBuffer(nil)
	r, err := http.NewRequest("GET", "http://localhost/", body)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		w := testResp{
			header: map[string][]string{},
			body:   bytes.Buffer{},
			Status: 0,
		}
		narcissus(&w, r)
	}
}
