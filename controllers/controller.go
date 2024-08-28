package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"lastro.co/go-challenge/db"
)

type ChatMessage struct {
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
}

type MessageResponse struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

var messageBuffers = make(map[string][]string)
var timers = make(map[string]*time.Timer)
var mutex = &sync.Mutex{}

// HandleChat função principal para lidar com o endpoint de POST da mensagem, também onde está a lógica de go routines
func HandleChat(w http.ResponseWriter, r *http.Request) {
	var msg ChatMessage

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a mensagem", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Erro ao decodificar a mensagem", http.StatusBadRequest)
		return
	}

	if msg.ChatID == "" || msg.Content == "" {
		http.Error(w, "chat_id e conteúdo são obrigatórios", http.StatusBadRequest)
		return
	}

	go func(msg ChatMessage) {
		err = createChatIfNotExists(msg.ChatID)
		if err != nil {
			fmt.Println("Erro ao verificar ou criar o chat:", err)
			return
		}

		err = db.CreateMessage(msg.ChatID, "user", msg.Content)
		if err != nil {
			fmt.Println("Erro ao salvar a mensagem:", err)
			return
		}

		mutex.Lock()

		messageBuffers[msg.ChatID] = append(messageBuffers[msg.ChatID], msg.Content)

		if timer, exists := timers[msg.ChatID]; exists {
			timer.Stop()
		}

		timers[msg.ChatID] = time.AfterFunc(10*time.Second, func() {
			generateAndStoreAssistantResponse(msg.ChatID)
		})

		mutex.Unlock()

	}(msg)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Mensagem recebida com sucesso"))
}

// createChatIfNotExists valida se existe um chat, se não cria um chat novo
func createChatIfNotExists(chatID string) error {
	exists, err := db.ChatExists(chatID)
	if err != nil {
		return err
	}

	if !exists {
		return db.CreateChat(chatID)
	}

	return nil
}

// generateAndStoreAssistantResponse função principal da goroutine para gerar e armazenar a resposta do assistente
// O agente não responde sincronamente o usuário, para isso foi criado um endpoint de GET, caso contrario necessitáriamos de um WebSocket ou relativo
func generateAndStoreAssistantResponse(chatID string) {
	mutex.Lock()

	var responseContent string
	for _, msg := range messageBuffers[chatID] {
		responseContent += msg + " "
	}

	delete(messageBuffers, chatID)
	timer := timers[chatID]
	delete(timers, chatID)

	mutex.Unlock()

	if timer != nil {
		timer.Stop()
	}

	response := fmt.Sprintf("RESPOSTA DO ASSISTENTE PARA MENSAGEM: %s", responseContent)

	err := db.CreateMessage(chatID, "assistant", response)
	if err != nil {
		fmt.Println("Erro ao salvar resposta do assistente:", err)
	}
}

// HandleGetMessages lida com o endpoint de GET para pegar as mensagens
func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")

	if chatID == "" {
		http.Error(w, "chat_id é obrigatório", http.StatusBadRequest)
		return
	}

	messages, err := db.GetMessages(chatID)
	if err != nil {
		fmt.Println("Erro ao recuperar mensagens:", err)
		http.Error(w, "Erro ao recuperar mensagens", http.StatusInternalServerError)
		return
	}

	var response []MessageResponse
	for _, msg := range messages {
		response = append(response, MessageResponse{Author: msg.Author, Content: msg.Content})
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Erro ao gerar resposta JSON:", err)
		http.Error(w, "Erro ao gerar resposta JSON", http.StatusInternalServerError)
		return
	}

	fmt.Println("Mensagens recuperadas com sucesso")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
