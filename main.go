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
	Headers map[string][]string `json:"headers"`
	Get     map[string][]string `json:"get"`
	Post    map[string][]string `json:"post"`
	Body    []byte              `json:"body"`
}

func narcissus(w http.ResponseWriter, r *http.Request) {
	req := Req{
		Headers: map[string][]string{},
		Get:     map[string][]string{},
		Post:    map[string][]string{},
		Body:    []byte{},
	}

	// must accept multiple values for each header
	for name, values := range r.Header {
		req.Headers[name] = append(req.Headers[name], values...)
	}

	// GET variables
	for k, v := range r.URL.Query() {
		req.Get[k] = append(req.Get[k], v...)
	}

	// POST variables
	if err := r.ParseForm(); err != nil {
		log.Fatal(err.Error())
	}
	for k, v := range r.PostForm {
		req.Post[k] = append(req.Post[k], v...)
	}

	// body request data (like JSON/XML)
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

	// because it's like the still lake he looked into
	w.Header().Add("x-powered-by", "narcissus()")
	// post to here from anywhere. I don't particularly care.
	w.Header().Add("Access-Control-Allow-Origin", "*")
	strRepLen := strconv.Itoa(len(toWrite))
	w.Header().Add("Content-Length", strRepLen)
	written, err := w.Write(toWrite)
	if err != nil {
		log.Fatal(err.Error())
	}
	if written != len(toWrite) {
		log.Fatal("written length not equal to length of write content (unthinkable)")
	}
}
