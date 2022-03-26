package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", narcissus)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println(err.Error())
	}
}

// Req allows for marshalling the request data for JSON
type Req struct {
	Headers map[string][]string `json:"headers"`
	Get     map[string][]string `json:"get"`
	Body    []byte              `json:"body"`
}

func narcissus(w http.ResponseWriter, r *http.Request) {
	req := Req{
		Headers: map[string][]string{},
		Get:     map[string][]string{},
		Body:    []byte{},
	}

	// must accept multiple values for each header
	for name, values := range r.Header {
		req.Headers[name] = append(req.Headers[name], values...)
	}

	for k, v := range r.URL.Query() {
		req.Get[k] = append(req.Get[k], v...)
	}

	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Body = bytes

	toWrite, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	strRepLen := strconv.Itoa(len(toWrite))
	w.Header().Add("x-powered-by", "narcissus()")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Length", strRepLen)
	written, err := w.Write(toWrite)
	if err != nil {
		log.Fatal(err.Error())
	}
	if written != len(toWrite) {
		log.Fatal("written length not equal to length of write content (unthinkable)")
	}
}
