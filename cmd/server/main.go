package main

import (
	"context"

	"task-tracker/internal/config"
	"task-tracker/internal/server"
)

func main() {
	config := config.LoadConfig()

	server := server.NewHTTPServer(*config)

	ctx := context.Background()

	err := server.ConfigureServer(ctx)
	if err != nil {
		panic(err)
	}

	err = server.Start(ctx)
	if err != nil {
		panic(err)
	}
}
