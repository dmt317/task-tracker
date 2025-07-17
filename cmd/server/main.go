package main

import (
	"context"

	"task-tracker/internal/config"
	"task-tracker/internal/server"
)

func main() {
	config := config.LoadConfig()

	server := server.NewHTTPServer(*config)

	ctx, cancel := context.WithCancel(context.Background())

	err := server.ConfigureServer(ctx)
	if err != nil {
		panic(err)
	}

	err = server.Start(cancel)
	if err != nil {
		panic(err)
	}
}
