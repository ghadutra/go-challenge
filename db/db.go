package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type Message struct {
	Author  string
	Content string
}

// InitDB inicia o banco de dados e testa a conexão
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro ao pingar o banco de dados: %v", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso")
}

// CloseDB gera um wrapper para utilizar no shutdown
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// ChatExists verifica se um chat já existe com o id enviado
func ChatExists(chatID string) (bool, error) {
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM chats WHERE id=$1)", chatID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

// CreateChat cria um chat com o id enviado
func CreateChat(chatID string) error {
	_, err := DB.Exec("INSERT INTO chats (id) VALUES ($1)", chatID)
	if err != nil {
		return fmt.Errorf("erro ao criar chat: %v", err)
	}
	return nil
}

// CreateMessage insere uma mensagem na tabela de mensagens relacionando com o chat id
func CreateMessage(chatID, author, content string) error {
	_, err := DB.Exec("INSERT INTO messages (chat_id, author, content) VALUES ($1, $2, $3)", chatID, author, content)
	return err
}

// GetMessages seleciona todas as mensagens do chat id
func GetMessages(chatID string) ([]Message, error) {
	rows, err := DB.Query("SELECT author, content FROM messages WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var author, content string
		if err := rows.Scan(&author, &content); err != nil {
			return nil, err
		}
		messages = append(messages, Message{Author: author, Content: content})
	}

	return messages, nil
}
