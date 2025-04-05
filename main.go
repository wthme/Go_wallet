package main

import (
	"log"
	"net/http"
	"os"
)

type Handler struct {}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request){
	_,err := w.Write([]byte("ti pidor"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	log.Println(http.ListenAndServe("127.0.0.1:8888", &Handler{}))

	os.Getenv("PASS")
}