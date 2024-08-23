package service

import (
	"lastro.co/go-challenge/controllers"
	"net/http"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/chat", controllers.HandleChat)
	return router
}
