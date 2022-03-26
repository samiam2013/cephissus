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
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println(err.Error())
	}
}

// Req allows for marshalling the request data for JSON
type Req struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func narcissus(w http.ResponseWriter, r *http.Request) {
	req := Req{
		Headers: map[string]string{},
		Body:    "",
	}

	// Loop over header names
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			req.Headers[name] = value
		}
	}

	bytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	}
	req.Body = string(bytes)

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
