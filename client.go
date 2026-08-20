package sms_ir

import (
	"github.com/1tzArad/sms_ir/internal/transport"
	"github.com/1tzArad/sms_ir/report"
	"github.com/1tzArad/sms_ir/send"
	"github.com/1tzArad/sms_ir/settings"
)

type Client struct {
	Send     *send.Service
	Report   *report.Service
	Settings *settings.Service
}

func New(apiKey string, opts ...Option) *Client {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	tc := transport.New(apiKey, cfg.baseURL, cfg.httpClient)

	return &Client{
		Send:     send.NewService(tc),
		Report:   report.NewService(tc),
		Settings: settings.NewService(tc),
	}
}
