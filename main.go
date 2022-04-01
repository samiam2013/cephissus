package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const defaultPort = ":3000"
const debugPort = ":80"

func main() {
	http.HandleFunc("/", narcissus)
	http.HandleFunc("/testForm", testForm)
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
		if isFile(v) {
			var (
				b64 string
				err error
			)
			if b64, err = parseFile(r, k); err == nil {
				req.Post[k] = append(req.Post[k], b64)
				continue
			}
			log.Println(err.Error())
		}
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

func testForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./htdocs/index.html")
}

// isFile checks that there is a single argument and it looks like a file
func isFile(postArg []string) bool {
	if len(postArg) != 1 {
		return false
	}
	split := strings.Split(postArg[0], ".")
	if len(split) < 2 {
		return false
	}
	// this is restrictive and doesn't include .js, .ts etc. on purpose
	exts := map[string]bool{
		"css":  true,
		"html": true,
		"xml":  true,
		"json": true,
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
	}
	_, ok := exts[split[len(split)-1]]
	return ok
}

// parseFile unpacks multi-part form data into a b64 string
func parseFile(r *http.Request, formFilename string) (string, error) {
	in, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer in.Close()

	out := bytes.NewBuffer([]byte{})
	io.Copy(out, in)
	return out.String(), nil
}
