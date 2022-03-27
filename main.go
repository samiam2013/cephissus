package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = ":3000"
const debugPort = ":80"

func main() {
	http.HandleFunc("/", narcissus)
	//debug := false
	port := defaultPort
	if _, err := os.Stat(".debug"); err == nil {
		//debug = true
		port = debugPort
	}

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println(err.Error())
	}
}

// Req allows for marshalling the request data for JSON
type Req struct {
	Headers   map[string][]string `json:"headers"`
	Get       map[string][]string `json:"get"`
	Post      map[string][]string `json:"post"`
	Body      []byte              `json:"body"`
	SourceURL string              `json:"source_code_available_at"`
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
	var err error
	if err = r.ParseForm(); err != nil {
		log.Fatal(err.Error())
	}
	for k, v := range r.PostForm {
		req.Post[k] = append(req.Post[k], v...)
	}

	// body request data (like JSON/XML)
	defer r.Body.Close()
	if req.Body, err = ioutil.ReadAll(r.Body); err != nil {
		log.Fatal(err.Error())
	}

	req.SourceURL = "https://github.com/samiam2013/cephissus"

	var toWrite []byte
	if toWrite, err = json.Marshal(req); err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	// because it's like the still lake he looked into
	w.Header().Add("x-powered-by", "narcissus()")
	// post to here from anywhere. I don't particularly care.
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Length", strconv.Itoa(len(toWrite)))
	written, err := w.Write(toWrite)
	if err != nil {
		log.Fatal(err.Error())
	}
	if written != len(toWrite) {
		log.Fatal("written length not equal to length of write content (unthinkable)")
	}
}
