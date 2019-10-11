package main

import (
	"log"

	"github.com/albertwidi/kothak/server"
)

// main function to keep all controls
func main() {
	if err := server.Main(); err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Program exited")
}
