package main

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func Test_narcissus(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		args   args
		expect http.ResponseWriter
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				w: testResp{
					header: map[string][]string{},
					body:   bytes.Buffer{},
					Status: 0,
				},
				r: func() *http.Request {
					r, err := http.NewRequest("GET", "http://localhost/", bytes.NewBuffer([]byte{}))
					if err != nil {
						t.Fatal("could not build request for test:", err.Error())
					}
					return r
				}(),
			},
			expect: testResp{
				header: map[string][]string{
					"Access-Control-Allow-Origin": {"*"},
					"Content-Length":              {"112"},
					"Content-Type":                {"application/json"},
					"X-Powered-By":                {"narcissus()"}},
				body:   bytes.Buffer{},
				Status: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			narcissus(tt.args.w, tt.args.r)
			if !reflect.DeepEqual(tt.args.w, tt.expect) {
				t.Fatalf("narcissus()='%#v'; expected '%#v'", tt.args.w, tt.expect)
			}
		})
	}
}
