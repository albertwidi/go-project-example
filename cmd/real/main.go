package main

import (
	"log"

	"gitlab.com/kosanapp/kothak/server"
)

// main function to keep all controls
func main() {
	if err := server.Main(); err != nil {
		log.Fatal(err)
		return
	}
	log.Println("Program exited")
}
