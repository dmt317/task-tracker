package main

import (
	"task-tracker/internal/config"
	"task-tracker/internal/repository"
	"task-tracker/internal/server"
)

func main() {
	config := config.LoadConfig()

	repo := repository.NewMemoryTaskRepository()

	server := server.NewHTTPServer(*config, repo)

	server.StartHTTPServer()
}
