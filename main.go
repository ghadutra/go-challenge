package main

import (
	"fmt"
	"log"
	"net/http"

	"lastro.co/go-challenge/db"
	"lastro.co/go-challenge/service"
)

func main() {
	// Se tratando de um teste apenas em ambiente local, não tem necessidade de armazenar as váriaveis em local seguro
	db.InitDB("user=youruser dbname=yourdbname sslmode=disable password=yourpassword")
	defer db.CloseDB()

	router := service.NewRouter()

	fmt.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
