package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	sites "github.com/thecoretg/threatdown-site-list-lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if req.Headers["x-api-secret"] != os.Getenv("API_SECRET") {
		return events.LambdaFunctionURLResponse{StatusCode: 401, Body: "unauthorized"}, nil
	}

	data, err := sites.FetchSites(ctx)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	body, err := json.Marshal(data)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
