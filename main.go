package main

import (
	// "fmt"
	"log"
)

func main() {
	// store is a pointer to a PostgresStore instance and an error
	// NewPostgresStore creates a new PostgresStore
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// server is a pointer to an APIServer instance
	// NewAPIServer creates the api server
	server := NewAPIServer(":8080", store)
	server.Run()
}
