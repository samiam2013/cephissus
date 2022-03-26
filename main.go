package main

import(
	"net/http"
	"log"
)

func main(){
	http.HandleFunc("/", narcissus)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println(err.Error())
	}
}

func narcissus( w http.ResponseWriter, r *http.Request ){
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		log.Fatal(err.Error())
	}
	r.Body.Close()
}
