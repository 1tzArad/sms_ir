package sms_ir

import "net/http"

type config struct {
	baseURL    string
	httpClient *http.Client
}

func defaultConfig() *config {
	return &config{
		baseURL:    "",
		httpClient: http.DefaultClient,
	}
}

type Option func(*config)

func WithBaseURL(baseURL string) Option {
	return func(c *config) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}
