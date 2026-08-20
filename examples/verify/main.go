// Example: verify
//
// Sends a verification (OTP) message through a pre-approved template. The API
// key, template id and recipient mobile are read from the environment so the
// program can be re-run with different values without editing the source.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	smsir "github.com/1tzArad/sms_ir"
	"github.com/1tzArad/sms_ir/send"
)

func main() {
	apiKey := requiredEnv("SMSIR_API_KEY")
	templateID := requiredEnv("SMSIR_TEMPLATE_ID")
	mobile := requiredEnv("SMSIR_RECIPIENT")

	client := smsir.New(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tid, err := strconv.Atoi(templateID)
	if err != nil {
		failf("invalid SMSIR_TEMPLATE_ID %q: %v", templateID, err)
	}

	resp, err := client.Send.Verify(ctx, send.VerifyRequest{
		Mobile:     mobile,
		TemplateID: tid,
		Parameters: []send.Parameter{
			{Name: "code", Value: "123456"},
		},
	})
	if err != nil {
		fail(err)
	}

	fmt.Println("Verification OTP sent:")
	fmt.Printf("  message id : %d\n", resp.MessageID)
	fmt.Printf("  cost       : %f\n", resp.Cost)
}

func requiredEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "%s environment variable is required\n", key)
		os.Exit(1)
	}
	return v
}

func failf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
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
