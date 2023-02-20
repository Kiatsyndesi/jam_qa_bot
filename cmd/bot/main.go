package main

import (
	"github.com/joho/godotenv"
	"jam_qa_bot/internal/app/mm-client"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	//Create goroutine for tracking work of websocket
	mm_client.StartDebuggingMMBot()
}
