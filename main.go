package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"faceit-cli/internal/app"
	"faceit-cli/internal/config"

	"github.com/joho/godotenv"
)

// Version is set during build time via ldflags
var version = "dev"

func main() {
	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("faceit-cli version %s\n", version)
		os.Exit(0)
	}

	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't treat this as an error
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	
	application := app.NewApp(cfg)
	
	if err := application.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
