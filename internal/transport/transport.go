package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/1tzArad/sms_ir/statuscode"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.sms.ir"

type envelope[T any] struct {
	Status  statuscode.Code `json:"status"`
	Message string          `json:"message"`
	Data    T               `json:"data"`
}

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func New(apiKey, baseURL string, httpClient *http.Client) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func Do[T any](ctx context.Context, c *Client, method, path string, body any, out *T) error {
	var reqBody io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("SMS_IR_SDK: failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("SMS_IR_SDK: failed to build request: %w", err)
	}

	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("SMS_IR_SDK: http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("SMS_IR_SDK: failed to read response body: %w", err)
	}

	var env envelope[T]
	if len(respBytes) > 0 {
		if jsonErr := json.Unmarshal(respBytes, &env); jsonErr != nil {
			return fmt.Errorf("SMS_IR_SDK: failed to decode response (http status %d): %w", resp.StatusCode, jsonErr)
		}
	}

	if env.Status != 1 {
		return &APIError{
			Code:    env.Status,
			Message: env.Message,
		}
	}

	if out != nil {
		*out = env.Data
	}

	return nil
}
