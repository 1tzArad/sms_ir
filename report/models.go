package report

import "github.com/1tzArad/sms_ir/statuscode"

// ============================================================================
// Sent Message Reports
// ============================================================================

// SentMessageReport represents the report information of a sent SMS.
type SentMessageReport struct {
	MessageID        int64                     `json:"messageId"`
	Mobile           string                    `json:"mobile"`
	MessageText      string                    `json:"messageText"`
	SendDateTime     int64                     `json:"sendDateTime"`
	LineNumber       int64                     `json:"lineNumber"`
	Cost             float64                   `json:"cost"`
	DeliveryState    *statuscode.DeliveryState `json:"deliveryState"`
	DeliveryDateTime *int64                    `json:"deliveryDateTime"`
}

// ----------------------------------------------------------------------------
// Get Sent Message Status
// GET /send/{messageId}
// ----------------------------------------------------------------------------

type GetSentMessageStatusResponseData = SentMessageReport

// ----------------------------------------------------------------------------
// List Today's Send Packs
// GET /send/pack
// ----------------------------------------------------------------------------

type ListTodaySendPacksRequestParams struct {
	PageSize   int `url:"pageSize"`
	PageNumber int `url:"pageNumber"`
}

type SendPackReport struct {
	PackID           string `json:"packId"`
	RecipientCount   int    `json:"recipientCount"`
	CreationDateTime int64  `json:"creationDateTime"`
}

type ListTodaySendPacksResponseData []SendPackReport

// ----------------------------------------------------------------------------
// Get Send Pack Report
// GET /send/pack/{packId}
// ----------------------------------------------------------------------------

type GetSendPackReportRequestParams struct {
	PackID string `path:"packId"`
}

type GetSendPackReportResponseData []SentMessageReport

// ----------------------------------------------------------------------------
// List Today's Sent Messages
// GET /send/live
// ----------------------------------------------------------------------------

type ListTodaySentMessagesRequestParams struct {
	PageSize   int `url:"pageSize"`
	PageNumber int `url:"pageNumber"`
}

type ListTodaySentMessagesResponseData []SentMessageReport

// ----------------------------------------------------------------------------
// List Archived Sent Messages
// GET /send/archive
// ----------------------------------------------------------------------------

type ListArchivedSentMessagesRequestParams struct {
	FromDate   int64 `url:"fromDate"`
	ToDate     int64 `url:"toDate"`
	PageSize   int   `url:"pageSize"`
	PageNumber int   `url:"pageNumber"`
}

type ListArchivedSentMessagesResponseData []SentMessageReport

// ============================================================================
// Received Message Reports
// ============================================================================

// ReceivedMessageReport represents the report information of a received SMS.
type ReceivedMessageReport struct {
	ReceiveReturnID  int64  `json:"receiveReturnId"`
	MessageText      string `json:"messageText"`
	Number           int64  `json:"number"`
	Mobile           string `json:"mobile"`
	ReceivedDateTime int64  `json:"receivedDateTime"`
}

// ----------------------------------------------------------------------------
// Get Latest Received Messages
// GET /receive/latest
// ----------------------------------------------------------------------------

type GetLatestReceivedMessagesRequestParams struct {
	Count int `url:"count"`
}

type GetLatestReceivedMessagesResponseData []ReceivedMessageReport

// ----------------------------------------------------------------------------
// List Today's Received Messages
// GET /receive/live
// ----------------------------------------------------------------------------

type ListTodayReceivedMessagesRequestParams struct {
	PageSize     int  `url:"pageSize"`
	PageNumber   int  `url:"pageNumber"`
	SortByNewest bool `url:"sortByNewest"`
}

type ListTodayReceivedMessagesResponseData []ReceivedMessageReport

// ----------------------------------------------------------------------------
// List Archived Received Messages
// GET /receive/archive
// ----------------------------------------------------------------------------

type ListArchivedReceivedMessagesRequestParams struct {
	FromDate   int64 `url:"fromDate"`
	ToDate     int64 `url:"toDate"`
	PageSize   int   `url:"pageSize"`
	PageNumber int   `url:"pageNumber"`
}

type ListArchivedReceivedMessagesResponseData []ReceivedMessageReport
