package main

import (
	server "backend/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Failed to load .env file")
	}
	server.SetupServer().Run()
}
