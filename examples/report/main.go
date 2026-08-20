// Example: report
//
// Fetches today's send-pack reports and prints a readable, formatted summary
// of every pack created today, including the number of recipients per pack and
// when it was created.
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

	resp, err := client.Report.ListTodaySendPacks(ctx, report.ListTodaySendPacksRequestParams{
		PageSize: 20,
	})
	if err != nil {
		fail(err)
	}

	packs := *resp
	fmt.Printf("Today's send packs (%d):\n\n", len(packs))
	fmt.Printf("%-20s %-11s %s\n", "PACK ID", "RECIPIENTS", "CREATED AT (UTC)")
	for _, p := range packs {
		fmt.Printf("%-20s %-11d %s\n",
			p.PackID, p.RecipientCount,
			time.Unix(p.CreationDateTime, 0).UTC().Format("2006-01-02 15:04:05"),
		)
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
