package main

import (
	"books/database"
	"books/router"

	"log"
)

func main() {
	db := database.CreateDb()

	r := router.SetupRouter(db)

	// Start server on port 8080 (default)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
