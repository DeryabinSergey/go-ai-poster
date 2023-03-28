package main

import (
	"context"
	"github.com/DeryabinSergey/go-ai-poster"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterEventFunctionContext(ctx, "/", aiposter.EventHandler); err != nil {
		log.Fatalf("funcframework.RegisterEventFunctionContext: %v\n", err)
	}
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("framework.Start: %v\n", err)
	}
}
