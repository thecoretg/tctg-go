package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thecoretg/tctg-go/salesforce"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatal("usage: local <soql query> [--simplify]")
	}

	q := os.Args[1]
	simplify := len(os.Args) > 2 && os.Args[2] == "--simplify"

	ctx := context.Background()

	client, err := salesforce.NewClient(ctx, salesforce.Config{
		ClientID:       os.Getenv("SALESFORCE_CLIENT_ID"),
		ClientSecret:   os.Getenv("SALESFORCE_CLIENT_SECRET"),
		CompanyURLName: os.Getenv("SALESFORCE_COMPANY_URL_NAME"),
	})
	if err != nil {
		log.Fatalf("creating salesforce client: %v", err)
	}

	records, err := salesforce.Query[map[string]any](ctx, client, q, simplify)
	if err != nil {
		log.Fatalf("query: %v", err)
	}

	out, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		log.Fatalf("marshaling records: %v", err)
	}

	if err := os.WriteFile("query_result.json", out, 0o644); err != nil {
		log.Fatalf("writing query_result.json: %v", err)
	}

	log.Printf("wrote %d records to query_result.json", len(records))
}
