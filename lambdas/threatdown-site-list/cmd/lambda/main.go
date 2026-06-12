package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	sites "github.com/thecoretg/threatdown-site-list-lambda"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if req.Headers["x-api-secret"] != os.Getenv("API_SECRET") {
		return errResponse(401, errors.New("unauthorized")), nil
	}

	opts := sites.FetchOptions{
		SkipSubs:   req.QueryStringParameters["skip_subs"] == "true",
		SkipAddOns: req.QueryStringParameters["skip_addons"] == "true",
	}

	data, err := sites.FetchSites(ctx, opts)
	if err != nil {
		return errResponse(500, err), nil
	}

	body, err := json.Marshal(data)
	if err != nil {
		return errResponse(500, err), nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}

func errResponse(status int, err error) events.LambdaFunctionURLResponse {
	var msgs []string
	type multiErr interface{ Unwrap() []error }
	if me, ok := err.(multiErr); ok {
		for _, e := range me.Unwrap() {
			msgs = append(msgs, e.Error())
		}
	} else {
		msgs = []string{err.Error()}
	}

	body, _ := json.Marshal(map[string][]string{"errors": msgs})
	return events.LambdaFunctionURLResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}
}

func main() {
	lambda.Start(handler)
}
