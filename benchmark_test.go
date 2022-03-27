package main

import (
	"bytes"
	"log"
	"net/http"
	"testing"
)

var _ http.ResponseWriter = response{}

type response struct {
	header map[string][]string
	body   bytes.Buffer
	Status int
}

// Header implements http.ResponseWriter
func (r response) Header() http.Header {
	return r.header
}

// Write implements http.ResponseWriter
func (r response) Write(b []byte) (n int, err error) {
	return r.body.Write(b)
}

// WriteHeader implements http.ResponseWriter
func (r response) WriteHeader(statusCode int) {
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
		w := response{
			header: map[string][]string{},
			body:   bytes.Buffer{},
			Status: 0,
		}
		narcissus(&w, r)
	}
}
