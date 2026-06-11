package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	sites "github.com/thecoretg/threatdown-site-list-lambda"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	data, err := sites.FetchSites(ctx)
	if err != nil {
		log.Fatalf("fetching sites: %v", err)
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("marshaling sites: %v", err)
	}

	if err := os.WriteFile("sites.json", out, 0o644); err != nil {
		log.Fatalf("writing sites.json: %v", err)
	}

	log.Printf("wrote %d sites to sites.json", len(data))
}
