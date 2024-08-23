package main

import (
	"lastro.co/go-challenge/service"
	"log"
	"net/http"
)

func main() {
	router := service.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
