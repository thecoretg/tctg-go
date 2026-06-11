package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thecoretg/tctg-go/salesforce"
)

var sfClient *salesforce.Client

func init() {
	var err error
	sfClient, err = salesforce.NewClient(context.Background(), salesforce.Config{
		ClientID:       os.Getenv("SALESFORCE_CLIENT_ID"),
		ClientSecret:   os.Getenv("SALESFORCE_CLIENT_SECRET"),
		CompanyURLName: os.Getenv("SALESFORCE_COMPANY_URL_NAME"),
	})
	if err != nil {
		panic("salesforce client: " + err.Error())
	}
}

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if req.Headers["x-api-secret"] != os.Getenv("API_SECRET") {
		return errResponse(401, "unauthorized"), nil
	}

	q := req.QueryStringParameters["q"]
	if q == "" {
		return errResponse(400, "missing query parameter: q"), nil
	}

	simplify := req.QueryStringParameters["simplify"] == "true"

	records, err := salesforce.Query[map[string]any](ctx, sfClient, q, simplify)
	if err != nil {
		return errResponse(500, err.Error()), nil
	}

	body, err := json.Marshal(records)
	if err != nil {
		return errResponse(500, err.Error()), nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func errResponse(status int, msg string) events.LambdaFunctionURLResponse {
	body, _ := json.Marshal(map[string]string{"error": msg})
	return events.LambdaFunctionURLResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}
}

func main() {
	lambda.Start(handler)
}
