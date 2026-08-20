# sms.ir Go SDK

A typed Go client for the [sms.ir](https://sms.ir) REST API. It exposes the
send, report, receive and settings endpoints through small, focused service
types built on top of a shared HTTP transport.

## Installation

```bash
go get github.com/1tzArad/sms_ir
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	smsir "github.com/1tzArad/sms_ir"
	"github.com/1tzArad/sms_ir/send"
)

func main() {
	client := smsir.New(os.Getenv("SMSIR_API_KEY"))
	resp, err := client.Send.Bulk(context.Background(), send.BulkRequest{
		LineNumber:  3000000000000,
		MessageText: "Hello from the sms.ir Go SDK!",
		Mobiles:     []string{"09120000000"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sent pack %s (%d messages)\n", resp.PackID, len(resp.MessageIDs))
}
```

## Configuration

`New` returns a `*smsir.Client` wiring together the `Send`, `Report` and
`Settings` services. Optional configuration is applied via functional options:

```go
client := smsir.New("YOUR-API-KEY",
    smsir.WithBaseURL("https://api.sms.ir"), // defaults to the production host
    smsir.WithHTTPClient(myHTTPClient),      // defaults to http.DefaultClient
)
```

## Services

All services live on `*smsir.Client` as exported fields: `Send`, `Report` and
`Settings`.

### Settings

```go
credit, err := client.Settings.GetCredit(ctx)   // *float64
lines,  err := client.Settings.ListLines(ctx)    // *[]int64
```

### Send

```go
// Bulk send
resp, err := client.Send.Bulk(ctx, send.BulkRequest{
    LineNumber:  100010001,
    MessageText: "hello",
    Mobiles:     []string{"09120000000", "09121111111"},
})

// Verify (template) send
resp, err := client.Send.Verify(ctx, send.VerifyRequest{
    Mobile:     "09120000000",
    TemplateID: 123456,
    Parameters: []send.Parameter{{Name: "code", Value: "1234"}},
})

// Send by URL (legacy credentials)
resp, err := client.Send.SendByURL(ctx, send.SendByURLRequest{
    Username: "u", Password: "p",
    Line: 100010001, Mobile: "09120000000", Text: "hello",
})

// Delete a scheduled pack
resp, err := client.Send.DeleteScheduled(ctx, "pack-id")
```

`Bulk` and `LikeToLike` validate their inputs (e.g. mobiles count ≤ 100,
matching lengths) before sending any request.

### Report & Receive

> The receive endpoints are grouped under the `Report` service (there is no
> separate `receive` package); `client.Report` covers both sent and received
> message reports.

```go
// Sent message status
status, err := client.Report.GetSentMessageStatus(ctx, messageID)

// Pack list & pack report
packs,  err := client.Report.ListTodaySendPacks(ctx, report.ListTodaySendPacksRequestParams{...})
items,  err := client.Report.GetSendPackReport(ctx, report.GetSendPackReportRequestParams{PackID: "pack-id"})

// Today's / archived sent messages
today,   err := client.Report.ListTodaySentMessages(ctx, report.ListTodaySentMessagesRequestParams{...})
archived, err := client.Report.ListArchivedSentMessages(ctx, report.ListArchivedSentMessagesRequestParams{
    FromDate: from, ToDate: to, PageSize: 20, PageNumber: 1,
})

// Received messages
latest,   err := client.Report.ListLatestReceivedMessages(ctx, report.GetLatestReceivedMessagesRequestParams{Count: 50})
liveRecv, err := client.Report.ListTodayReceivedMessages(ctx, report.ListTodayReceivedMessagesRequestParams{...})
archRecv, err := client.Report.ListArchivedReceivedMessages(ctx, report.ListArchivedReceivedMessagesRequestParams{...})
```

## Status codes & delivery states

Status codes and delivery states are typed enums in the
[`statuscode`](https://github.com/1tzArad/sms_ir/tree/main/statuscode) package:

```go
code := statuscode.InvalidAPIKey
fmt.Println(code)               // human-readable message
fmt.Println(code.IsAuthError()) // true
fmt.Println(code.IsRateLimit()) // false

state := statuscode.DeliveredToDevice
fmt.Println(state) // human-readable delivery state
```

## Error handling

Errors returned by the services are either client-side validation errors or
`*smsir.APIError` (a type alias for `*transport.APIError`) returned by the API:

```go
_, err := client.Settings.GetCredit(ctx)
if err != nil {
    var apiErr *smsir.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("api error: code=%d, message=%s\n", apiErr.Code, apiErr.Message)
    } else {
        // transport / serialization error
    }
}
```

## Examples

Runnable programs live under [`examples/`](examples). Each reads the API key
from the `SMSIR_API_KEY` environment variable (set it before running).

| Example | Description |
| --- | --- |
| [examples/basic](examples/basic) | Create a client and send a bulk SMS, printing the pack id and cost. |
| [examples/verify](examples/verify) | Send a verification/OTP message via a template with a parameter. |
| [examples/report](examples/report) | Fetch today's send packs and print a formatted summary. |
| [examples/custom-http-client](examples/custom-http-client) | Configure a custom HTTP client and base URL. |
| [examples/receive](examples/receive) | Fetch and print the latest received (inbound) SMS messages. |

## License

Distributed under the terms of the project license.
