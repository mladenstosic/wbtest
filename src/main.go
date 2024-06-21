package main

import (
	"log"
)

func main() {
	app := App{}

	// Init server
	log.Println("Server starting...")
	if err := app.Init("users.db"); err != nil {
		log.Panicln("cannot initialize the server")
	}

	// Run the server
	log.Println("Server running...")
	if err := app.Run(":8080"); err != nil {
		log.Panicln("cannot run the server")
	}
}
