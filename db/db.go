package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	db.Close()
}

// Funções de CRUD para chats e mensagens
