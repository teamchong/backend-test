package main

import (
	"context"
	"log"

	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/joho/godotenv"

	// Import all standard Benthos components
	_ "github.com/benthosdev/benthos/v4/public/components/all"

	// Import all Benthos Plugin packages
	_ "github.com/teamchong/backend-test/ratelimit"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	service.RunCLI(context.Background())
}
