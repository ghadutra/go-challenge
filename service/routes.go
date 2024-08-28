package service

import (
	"net/http"

	"lastro.co/go-challenge/controllers"
)

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Foi gerada duas rotas alternando entre POST e GET para enviar mensagens e pegar as mensagens
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.HandleChat(w, r)
		} else if r.Method == http.MethodGet {
			controllers.HandleGetMessages(w, r)
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	return router
}
