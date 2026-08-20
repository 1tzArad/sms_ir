// Example: custom-http-client
//
// Demonstrates configuring the client with a custom HTTP client (timeouts,
// transport tuning) and a custom base URL, then reads the account credit
// balance. The base URL can optionally be overridden with SMSIR_BASE_URL.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	smsir "github.com/1tzArad/sms_ir"
)

func main() {
	apiKey := os.Getenv("SMSIR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "SMSIR_API_KEY environment variable is required")
		os.Exit(1)
	}

	// A custom HTTP client gives full control over timeouts, redirects and the
	// underlying transport (for tracing, custom TLS config, etc.).
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	baseURL := "https://api.sms.ir"
	if custom := os.Getenv("SMSIR_BASE_URL"); custom != "" {
		baseURL = custom
	}

	client := smsir.New(apiKey,
		smsir.WithBaseURL(baseURL),
		smsir.WithHTTPClient(httpClient),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	credit, err := client.Settings.GetCredit(ctx)
	if err != nil {
		var apiErr *smsir.APIError
		if errors.As(err, &apiErr) {
			fmt.Fprintf(os.Stderr, "sms.ir API error: code=%d %s\n", apiErr.Code, apiErr.Message)
		} else {
			fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Account credit: %f\n", *credit)
}
