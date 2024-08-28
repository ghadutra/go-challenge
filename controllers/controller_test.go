package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"lastro.co/go-challenge/db"
)

var mock sqlmock.Sqlmock

func setup(t *testing.T) {
	var err error
	var dbMock *sql.DB
	dbMock, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Erro ao criar o mock do banco de dados: %v", err)
	}
	db.DB = dbMock
}

func teardown() {
	db.DB.Close()
}

func TestHandleChat(t *testing.T) {
	setup(t)
	defer teardown()

	var wg sync.WaitGroup

	// Mock para verificar a existência do chat
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM chats WHERE id=\\$1\\)").
		WithArgs("test_chat").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Mock para criar um novo chat
	mock.ExpectExec("INSERT INTO chats \\(id\\) VALUES \\(\\$1\\)").
		WithArgs("test_chat").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock para inserir a mensagem do usuário
	mock.ExpectExec("INSERT INTO messages \\(chat_id, author, content\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs("test_chat", "user", "Hello, World!").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock para inserir a resposta do assistente após o delay
	mock.ExpectExec("INSERT INTO messages \\(chat_id, author, content\\) VALUES \\(\\$1, \\$2, \\$3\\)").
		WithArgs("test_chat", "assistant", "RESPOSTA DO ASSISTENTE PARA MENSAGEM: Hello, World! ").
		WillReturnResult(sqlmock.NewResult(1, 1))

	chatMessage := ChatMessage{
		ChatID:  "test_chat",
		Content: "Hello, World!",
	}

	jsonData, _ := json.Marshal(chatMessage)

	req, err := http.NewRequest("POST", "/chat", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleChat)

	wg.Add(1)
	go func() {
		defer wg.Done()
		handler.ServeHTTP(rr, req)
	}()

	wg.Wait()

	time.Sleep(11 * time.Second)

	// Check if the status code is 202 (Accepted)
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

	expected := "Mensagem recebida com sucesso"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandleGetMessages(t *testing.T) {
	setup(t)
	defer teardown()

	// Mock da consulta de mensagens
	rows := sqlmock.NewRows([]string{"author", "content"}).
		AddRow("user", "Hello").
		AddRow("assistant", "RESPOSTA DO ASSISTENTE PARA MENSAGEM: Hello")

	// Ajuste para garantir que a consulta SQL corresponde corretamente
	mock.ExpectQuery(`SELECT author, content FROM messages WHERE chat_id = \$1`).
		WithArgs("test_chat").
		WillReturnRows(rows)

	// Criar uma nova requisição GET
	req, err := http.NewRequest("GET", "/chat?chat_id=test_chat", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleGetMessages)

	handler.ServeHTTP(rr, req)

	// Verifique o status da resposta
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", rr.Body.String())
	}

	// Decodifique a resposta JSON
	var messages []MessageResponse
	err = json.NewDecoder(rr.Body).Decode(&messages)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	expectedMessages := []MessageResponse{
		{Author: "user", Content: "Hello"},
		{Author: "assistant", Content: "RESPOSTA DO ASSISTENTE PARA MENSAGEM: Hello"},
	}

	if !reflect.DeepEqual(messages, expectedMessages) {
		t.Errorf("handler returned unexpected body: got %v want %v", messages, expectedMessages)
	}

	// Verifique se todas as expectativas foram atendidas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
