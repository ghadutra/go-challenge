package controllers

import (
	"net/http"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chatId"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func HandleChat(w http.ResponseWriter, r *http.Request) {
	// Implementação da lógica de processamento de mensagens
}
