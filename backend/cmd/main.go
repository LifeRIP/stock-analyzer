package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	log.Println("Testing air")
	log.Println("prueba", os.Getenv("TEST"))
}
