// Example: receive
//
// Polls the latest received (inbound) SMS messages and prints a readable
// summary. The receive endpoints are exposed through the Report service.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	smsir "github.com/1tzArad/sms_ir"
	"github.com/1tzArad/sms_ir/report"
)

func main() {
	apiKey := os.Getenv("SMSIR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "SMSIR_API_KEY environment variable is required")
		os.Exit(1)
	}

	client := smsir.New(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := client.Report.ListLatestReceivedMessages(ctx, report.GetLatestReceivedMessagesRequestParams{
		Count: 10,
	})
	if err != nil {
		fail(err)
	}

	messages := *resp
	if len(messages) == 0 {
		fmt.Println("No received messages.")
		return
	}

	fmt.Printf("Latest received messages (last %d):\n\n", len(messages))
	fmt.Printf("%-14s %-14s %s\n", "MOBILE", "NUMBER", "TEXT")
	for _, m := range messages {
		fmt.Printf("%-14s %-14d %s\n", m.Mobile, m.Number, m.MessageText)
	}
}

func fail(err error) {
	var apiErr *smsir.APIError
	if errors.As(err, &apiErr) {
		fmt.Fprintf(os.Stderr, "sms.ir API error: code=%d %s\n", apiErr.Code, apiErr.Message)
	} else {
		fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
	}
	os.Exit(1)
}
