package send_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1tzArad/sms_ir/internal/transport"
	"github.com/1tzArad/sms_ir/send"
	"github.com/1tzArad/sms_ir/statuscode"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func successEnvelope(data any) string {
	b, _ := json.Marshal(data)
	return `{"status":1,"message":"عملیات موفقیت‌آمیز بود","data":` + string(b) + `}`
}

func errorEnvelope(code int, message string) string {
	b, _ := json.Marshal(map[string]any{"status": code, "message": message})
	return string(b)
}

func newService(t *testing.T, handler http.HandlerFunc) *send.Service {
	t.Helper()
	srv := newTestServer(t, handler)
	tc := transport.New("test-api-key", srv.URL, srv.Client())
	return send.NewService(tc)
}

func int64Ptr(i int64) *int64 { return &i }

func TestNewService_ReturnsNonNil(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	tc := transport.New("key", srv.URL, srv.Client())
	if s := send.NewService(tc); s == nil {
		t.Fatal("expected non-nil service")
	}
}

// ---------------------------------------------------------------------------
// Bulk
// ---------------------------------------------------------------------------

func TestBulk_Success(t *testing.T) {
	var gotBody send.BulkRequest
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/v1/send/bulk" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send/bulk")
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(send.BulkSendResponseData{
			PackID:     "pack-1",
			MessageIDs: []*int64{int64Ptr(10), int64Ptr(11)},
			Cost:       2.5,
		})))
	})

	resp, err := svc.Bulk(context.Background(), send.BulkRequest{
		LineNumber:  100010001,
		MessageText: "hello",
		Mobiles:     []string{"09120000000", "09121111111"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PackID != "pack-1" {
		t.Errorf("PackID = %q, want %q", resp.PackID, "pack-1")
	}
	if resp.Cost != 2.5 {
		t.Errorf("Cost = %v, want 2.5", resp.Cost)
	}
	if len(resp.MessageIDs) != 2 {
		t.Fatalf("MessageIDs len = %d, want 2", len(resp.MessageIDs))
	}
	if *resp.MessageIDs[0] != 10 || *resp.MessageIDs[1] != 11 {
		t.Errorf("MessageIDs = %v", resp.MessageIDs)
	}
	if gotBody.LineNumber != 100010001 {
		t.Errorf("request LineNumber = %d", gotBody.LineNumber)
	}
	if gotBody.MessageText != "hello" {
		t.Errorf("request MessageText = %q", gotBody.MessageText)
	}
}

func TestBulk_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InvalidMobileNumbers), "درخواست شما دارای موبایل(های) نادرست است")))
	})

	_, err := svc.Bulk(context.Background(), send.BulkRequest{
		LineNumber:  1,
		MessageText: "hi",
		Mobiles:     []string{"09120000000"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.InvalidMobileNumbers {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.InvalidMobileNumbers)
	}
}

// ---------------------------------------------------------------------------
// LikeToLike
// ---------------------------------------------------------------------------

func TestLikeToLike_Success(t *testing.T) {
	var gotBody send.LikeToLikeRequest
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/v1/send/likeToLike" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send/likeToLike")
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(send.LikeToLikeResponseData{
			PackID:     "ltk-1",
			MessageIDs: []*int64{int64Ptr(20)},
			Cost:       3.0,
		})))
	})

	resp, err := svc.LikeToLike(context.Background(), send.LikeToLikeRequest{
		LineNumber:   100010001,
		MessageTexts: []string{"first", "second"},
		Mobiles:      []string{"09120000000", "09121111111"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PackID != "ltk-1" {
		t.Errorf("PackID = %q, want %q", resp.PackID, "ltk-1")
	}
	if len(resp.MessageIDs) != 1 || *resp.MessageIDs[0] != 20 {
		t.Errorf("unexpected MessageIDs: %v", resp.MessageIDs)
	}
	if len(gotBody.MessageTexts) != 2 || len(gotBody.Mobiles) != 2 {
		t.Errorf("unexpected body: %+v", gotBody)
	}
}

func TestLikeToLike_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InsufficientCredit), "اعتبار کافی نمی‌باشد")))
	})

	_, err := svc.LikeToLike(context.Background(), send.LikeToLikeRequest{
		LineNumber:   1,
		MessageTexts: []string{"hi"},
		Mobiles:      []string{"09120000000"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.InsufficientCredit {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.InsufficientCredit)
	}
}

// ---------------------------------------------------------------------------
// Verify
// ---------------------------------------------------------------------------

func TestVerify_Success(t *testing.T) {
	var gotBody send.VerifyRequest
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/v1/send/verify" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send/verify")
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(send.VerifyResponseData{
			MessageID: 99,
			Cost:      1.0,
		})))
	})

	resp, err := svc.Verify(context.Background(), send.VerifyRequest{
		Mobile:     "09120000000",
		TemplateID: 123456,
		Parameters: []send.Parameter{{Name: "code", Value: "1234"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MessageID != 99 {
		t.Errorf("MessageID = %d, want 99", resp.MessageID)
	}
	if resp.Cost != 1.0 {
		t.Errorf("Cost = %v, want 1.0", resp.Cost)
	}
	if gotBody.Mobile != "09120000000" || gotBody.TemplateID != 123456 {
		t.Errorf("unexpected body: %+v", gotBody)
	}
}

func TestVerify_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.TemplateNotFound), "قالب یافت نشد")))
	})

	_, err := svc.Verify(context.Background(), send.VerifyRequest{
		Mobile:     "09120000000",
		TemplateID: 123456,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.TemplateNotFound {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.TemplateNotFound)
	}
}

// ---------------------------------------------------------------------------
// SendByURL
// ---------------------------------------------------------------------------

func TestSendByURL_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/send" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send")
		}
		if r.URL.Query().Get("username") != "user1" {
			t.Errorf("username query = %q", r.URL.Query().Get("username"))
		}
		if r.URL.Query().Get("password") != "pass1" {
			t.Errorf("password query = %q", r.URL.Query().Get("password"))
		}
		if r.URL.Query().Get("line") != "100010001" {
			t.Errorf("line query = %q", r.URL.Query().Get("line"))
		}
		if r.URL.Query().Get("mobile") != "09120000000" {
			t.Errorf("mobile query = %q", r.URL.Query().Get("mobile"))
		}
		if r.URL.Query().Get("text") != "hello" {
			t.Errorf("text query = %q", r.URL.Query().Get("text"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(send.SendByURLResponseData{
			MessageID: 42,
			Cost:      1.25,
		})))
	})

	resp, err := svc.SendByURL(context.Background(), send.SendByURLRequest{
		Username: "user1",
		Password: "pass1",
		Line:     100010001,
		Mobile:   "09120000000",
		Text:     "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MessageID != 42 || resp.Cost != 1.25 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestSendByURL_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.MobilesBlacklisted), "شماره موبایل(ها) در لیست سیاه سامانه می‌باشند")))
	})

	_, err := svc.SendByURL(context.Background(), send.SendByURLRequest{
		Username: "user1", Password: "pass1",
		Line: 1, Mobile: "09120000000", Text: "hello",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.MobilesBlacklisted {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.MobilesBlacklisted)
	}
}

// ---------------------------------------------------------------------------
// DeleteScheduled
// ---------------------------------------------------------------------------

func TestDeleteScheduled_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want %q", r.Method, http.MethodDelete)
		}
		if r.URL.Path != "/v1/send/scheduled/abc-123" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send/scheduled/abc-123")
		}
		if r.Header.Get("X-API-KEY") != "test-api-key" {
			t.Errorf("X-API-KEY = %q", r.Header.Get("X-API-KEY"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(send.DeleteScheduledResponseData{
			ReturnedCreditCount: 5.0,
			SMSCount:            3,
		})))
	})

	resp, err := svc.DeleteScheduled(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ReturnedCreditCount != 5.0 || resp.SMSCount != 3 {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestDeleteScheduled_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.RecordNotFoundForDeletion), "رکوردی برای حذف یافت نشد")))
	})

	_, err := svc.DeleteScheduled(context.Background(), "abc-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.RecordNotFoundForDeletion {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.RecordNotFoundForDeletion)
	}
}

func TestDeleteScheduled_PathEscapesPackID(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		// url.PathEscape must turn the "/" inside the pack id into "%2F" so it
		// is treated as part of the last path segment rather than a separator.
		if got := r.URL.EscapedPath(); got != "/v1/send/scheduled/abc%2F123" {
			t.Errorf("escaped path = %q, want %q", got, "/v1/send/scheduled/abc%2F123")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(send.DeleteScheduledResponseData{})))
	})

	_, err := svc.DeleteScheduled(context.Background(), "abc/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
