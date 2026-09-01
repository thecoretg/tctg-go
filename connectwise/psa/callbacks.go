package psa

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func callbackIDEndpoint(callbackID int) string {
	return fmt.Sprintf("system/callbacks/%d", callbackID)
}

func (c *Client) PostCallback(ctx context.Context, webhook *Callback) (*Callback, error) {
	return c.Post[Callback](ctx, "system/callbacks", webhook)
}

func (c *Client) ListCallbacks(ctx context.Context, params map[string]string, opts ...ListOption) ([]Callback, error) {
	return c.GetMany[Callback](ctx, "system/callbacks", params, opts...)
}

func (c *Client) GetCallback(ctx context.Context, callbackID int, params map[string]string) (*Callback, error) {
	return c.Get[Callback](ctx, callbackIDEndpoint(callbackID), params)
}

func (c *Client) PutCallback(ctx context.Context, callbackID int, webhook *Callback) (*Callback, error) {
	return c.Put[Callback](ctx, callbackIDEndpoint(callbackID), webhook)
}

func (c *Client) PatchCallback(ctx context.Context, callbackID int, patchOps []PatchOp) (*Callback, error) {
	return c.Patch[Callback](ctx, callbackIDEndpoint(callbackID), patchOps)
}

func (c *Client) DeleteCallback(ctx context.Context, callbackID int) error {
	return c.Delete(ctx, callbackIDEndpoint(callbackID))
}

func ValidateWebhook(r *http.Request) (bool, error) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("reading request body: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing request body: %v", err)
		}
	}(r.Body)

	var meta struct {
		Metadata struct {
			KeyURL string `json:"key_url"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(payload, &meta); err != nil {
		return false, fmt.Errorf("unmarshaling request body: %w", err)
	}

	resp, err := http.Get(meta.Metadata.KeyURL)
	if err != nil {
		return false, fmt.Errorf("getting shared secret key from %s: %w", meta.Metadata.KeyURL, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing shared secret key response body: %v", err)
		}
	}(resp.Body)

	var keyResp struct {
		SigningKey string `json:"signing_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&keyResp); err != nil {
		return false, fmt.Errorf("decoding shared secret key response: %w", err)
	}

	sharedSecret := []byte(keyResp.SigningKey)
	hash := sha256.Sum256(sharedSecret)
	h := hmac.New(sha256.New, hash[:])
	h.Write(payload)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	expectedSig := r.Header.Get("x-content-signature")

	return signature == expectedSig, nil
}
