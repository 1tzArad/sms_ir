package send

// ============================================================================
// Shared Models
// ============================================================================

// Parameter represents a dynamic parameter used in a Verify SMS template.
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// SendResponseData contains the common response data returned by SMS send
// operations that send a single message.
type SendResponseData struct {
	MessageID int64   `json:"messageId"`
	Cost      float64 `json:"cost"`
}

// BulkSendResponseData contains the response data returned by bulk and
// like-to-like SMS sending operations.
type BulkSendResponseData struct {
	PackID     string   `json:"packId"`
	MessageIDs []*int64 `json:"messageIds"`
	Cost       float64  `json:"cost"`
}

// ============================================================================
// Verify SMS
// POST /send/verify
// ============================================================================

type VerifyRequest struct {
	Mobile     string      `json:"mobile"`
	TemplateID int         `json:"templateId"`
	Parameters []Parameter `json:"parameters"`
}

type VerifyResponseData = SendResponseData

// ============================================================================
// Bulk SMS
// POST /send/bulk
// ============================================================================

type BulkRequest struct {
	LineNumber   int64    `json:"lineNumber"`
	MessageText  string   `json:"messageText"`
	Mobiles      []string `json:"mobiles"`
	SendDateTime *int64   `json:"sendDateTime,omitempty"`
}

type BulkResponseData = BulkSendResponseData

// ============================================================================
// Like-to-Like SMS
// POST /send/likeToLike
// ============================================================================

type LikeToLikeRequest struct {
	LineNumber   int64    `json:"lineNumber"`
	MessageTexts []string `json:"messageTexts"`
	Mobiles      []string `json:"mobiles"`
	SendDateTime *int64   `json:"sendDateTime,omitempty"`
}

type LikeToLikeResponseData = BulkSendResponseData

// ============================================================================
// Send SMS via URL
// GET, POST /send
// ============================================================================

type SendByURLRequest struct {
	Username string `url:"username"`
	Password string `url:"password"`
	Line     int64  `url:"line"`
	Mobile   string `url:"mobile"`
	Text     string `url:"text"`
}

type SendByURLResponseData = SendResponseData

// ============================================================================
// Delete Scheduled SMS
// DELETE /send/scheduled/{packId}
// ============================================================================

type DeleteScheduledRequest struct {
	PackID string `path:"packId"`
}

type DeleteScheduledResponseData struct {
	ReturnedCreditCount float64 `json:"returnedCreditCount"`
	SMSCount            int     `json:"smsCount"`
}
