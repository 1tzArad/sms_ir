// Example: basic
//
// Creates an sms.ir client and sends a bulk SMS to a list of mobile numbers.
// It reads the API key from the SMSIR_API_KEY environment variable and prints
// the resulting pack id, the number of messages queued and the total cost.
//
// The line number below is a placeholder — replace it with one of your own
// dedicated numbers as returned by the sms.ir dashboard.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	smsir "github.com/1tzArad/sms_ir"
	"github.com/1tzArad/sms_ir/send"
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

	resp, err := client.Send.Bulk(ctx, send.BulkRequest{
		LineNumber:  3000000000000,
		MessageText: "Hello from the sms.ir Go SDK example!",
		Mobiles:     []string{"09120000000", "09121111111"},
	})
	if err != nil {
		fail(err)
	}

	fmt.Println("Bulk send succeeded:")
	fmt.Printf("  pack id        : %s\n", resp.PackID)
	fmt.Printf("  messages queued: %d\n", len(resp.MessageIDs))
	fmt.Printf("  total cost     : %f\n", resp.Cost)
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
