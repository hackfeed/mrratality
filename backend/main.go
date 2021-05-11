package main

import (
	server "backend/server"
)

func main() {
	server.SetupServer().Run()
}
