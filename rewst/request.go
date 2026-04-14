package rewst

import (
	"context"
	"fmt"
	"net/http"

	"resty.dev/v3"
)

var ErrNotFound = fmt.Errorf("404 status returned")

type WebhookClient struct {
	rc *resty.Client
}

func NewWebhookClient(secret string) *WebhookClient {
	rc := resty.New()
	rc.SetHeader("Accept", "application/json")
	rc.SetHeader("x-rewst-secret", secret)
	rc.SetRetryCount(3)
	rc.SetDisableWarn(true)
	return &WebhookClient{rc: rc}
}

func Get[T any](ctx context.Context, wc *WebhookClient, url string, params map[string]string) (*T, error) {
	var target T
	res, err := wc.rc.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&target).
		Get(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from Rewst: %s [%d]", res.String(), res.StatusCode())
	}
	return &target, nil
}

func Send(ctx context.Context, wc *WebhookClient, url string, body any) error {
	_, err := Post[struct{}](ctx, wc, url, body)
	return err
}

func Post[T any](ctx context.Context, wc *WebhookClient, url string, body any) (*T, error) {
	var target T
	res, err := wc.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&target).
		Post(url)
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, fmt.Errorf("error response from Rewst: %s [%d]", res.String(), res.StatusCode())
	}
	return &target, nil
}
