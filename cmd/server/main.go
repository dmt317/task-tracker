package main

import (
	"task-tracker/internal/config"
	"task-tracker/internal/server"
)

func main() {
	config := config.LoadConfig()

	server := server.NewHTTPServer(*config)

	err := server.Start()
	if err != nil {
		panic(err)
	}
}
