package send_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/1tzArad/sms_ir/internal/transport"
	"github.com/1tzArad/sms_ir/send"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, wantMethod, wantPath string, respBody any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != wantMethod {
			t.Errorf("expected method %s, got %s", wantMethod, r.Method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
		}
		if r.Header.Get("X-API-KEY") == "" {
			t.Error("expected X-API-KEY header to be set")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(respBody)
	}))
}

func successEnvelope(data any) map[string]any {
	return map[string]any{
		"status":  1,
		"message": "OK",
		"data":    data,
	}
}

func errorEnvelope(code int, message string) map[string]any {
	return map[string]any{
		"status":  code,
		"message": message,
		"data":    nil,
	}
}

// ---------- Bulk ----------

func TestBulk_Success(t *testing.T) {
	srv := newTestServer(t, http.MethodPost, "/v1/send/bulk", successEnvelope(map[string]any{
		"packId":     "2b99e63c-9bf8-4a21-9bfe-3f72dc1b46f1",
		"messageIds": []any{86522023, 86522024},
		"cost":       2.0,
	}))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	resp, err := send.Bulk(context.Background(), tc, send.BulkRequest{
		LineNumber:  300004505000017,
		MessageText: "hello",
		Mobiles:     []string{"09120000000", "09190000000"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.PackID != "2b99e63c-9bf8-4a21-9bfe-3f72dc1b46f1" {
		t.Errorf("unexpected packID: %s", resp.PackID)
	}
	if resp.Cost != 2.0 {
		t.Errorf("unexpected cost: %v", resp.Cost)
	}
	if len(resp.MessageIDs) != 2 {
		t.Fatalf("expected 2 message ids, got %d", len(resp.MessageIDs))
	}
}

func TestBulk_APIError(t *testing.T) {
	srv := newTestServer(t, http.MethodPost, "/v1/send/bulk", errorEnvelope(102, "insufficient credit"))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	_, err := send.Bulk(context.Background(), tc, send.BulkRequest{
		LineNumber:  300004505000017,
		MessageText: "hello",
		Mobiles:     []string{"09120000000"},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T: %v", err, err)
	}
	if apiErr.Code != 102 {
		t.Errorf("expected code 102, got %d", apiErr.Code)
	}
}

func TestBulk_Validation(t *testing.T) {
	tc := transport.New("fake-api-key", "http://unused.invalid", nil)

	tests := []struct {
		name string
		req  send.BulkRequest
	}{
		{
			name: "missing lineNumber",
			req: send.BulkRequest{
				MessageText: "test",
				Mobiles:     []string{"09120000000"},
			},
		},
		{
			name: "missing messageText",
			req: send.BulkRequest{
				LineNumber: 123,
				Mobiles:    []string{"09120000000"},
			},
		},
		{
			name: "empty mobiles",
			req: send.BulkRequest{
				LineNumber:  123,
				MessageText: "test",
				Mobiles:     []string{},
			},
		},
		{
			name: "too many mobiles",
			req: send.BulkRequest{
				LineNumber:  123,
				MessageText: "test",
				Mobiles:     make([]string, 101),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := send.Bulk(context.Background(), tc, tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

// ---------- LikeToLike ----------

func TestLikeToLike_Success(t *testing.T) {
	srv := newTestServer(t, http.MethodPost, "/v1/send/likeToLike", successEnvelope(map[string]any{
		"packId":     "2b99e63c-9bf8-4a21-9bfe-3f72dc1b46f1",
		"messageIds": []any{86522023, 86522024},
		"cost":       2.0,
	}))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	resp, err := send.LikeToLike(context.Background(), tc, send.LikeToLikeRequest{
		LineNumber:   300004505000017,
		MessageTexts: []string{"text one", "text two"},
		Mobiles:      []string{"09120000000", "09190000000"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Cost != 2.0 {
		t.Errorf("unexpected cost: %v", resp.Cost)
	}
}

func TestLikeToLike_Validation(t *testing.T) {
	tc := transport.New("fake-api-key", "http://unused.invalid", nil)

	tests := []struct {
		name string
		req  send.LikeToLikeRequest
	}{
		{
			name: "missing lineNumber",
			req: send.LikeToLikeRequest{
				MessageTexts: []string{"a"},
				Mobiles:      []string{"09120000000"},
			},
		},
		{
			name: "mismatched lengths",
			req: send.LikeToLikeRequest{
				LineNumber:   123,
				MessageTexts: []string{"a", "b"},
				Mobiles:      []string{"09120000000"},
			},
		},
		{
			name: "empty lists",
			req: send.LikeToLikeRequest{
				LineNumber:   123,
				MessageTexts: []string{},
				Mobiles:      []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := send.LikeToLike(context.Background(), tc, tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

// ---------- Verify ----------

func TestVerify_Success(t *testing.T) {
	srv := newTestServer(t, http.MethodPost, "/v1/send/verify", successEnvelope(map[string]any{
		"messageId": 89545112,
		"cost":      1.0,
	}))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	resp, err := send.Verify(context.Background(), tc, send.VerifyRequest{
		Mobile:     "09190000000",
		TemplateID: 123456,
		Parameters: []send.Parameter{
			{Name: "Code", Value: "12345"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MessageID != 89545112 {
		t.Errorf("unexpected messageID: %d", resp.MessageID)
	}
}

func TestVerify_APIError(t *testing.T) {
	srv := newTestServer(t, http.MethodPost, "/v1/send/verify", errorEnvelope(113, "template not found"))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	_, err := send.Verify(context.Background(), tc, send.VerifyRequest{
		Mobile:     "09190000000",
		TemplateID: 999999,
	})

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T: %v", err, err)
	}
	if apiErr.Code != 113 {
		t.Errorf("expected code 113, got %d", apiErr.Code)
	}
}

func TestVerify_Validation(t *testing.T) {
	tc := transport.New("fake-api-key", "http://unused.invalid", nil)

	tests := []struct {
		name string
		req  send.VerifyRequest
	}{
		{name: "missing mobile", req: send.VerifyRequest{TemplateID: 123}},
		{name: "missing templateId", req: send.VerifyRequest{Mobile: "09120000000"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := send.Verify(context.Background(), tc, tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

// ---------- SendByURL ----------

func TestSendByURL_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/send" {
			t.Errorf("expected path /v1/send, got %s", r.URL.Path)
		}

		q := r.URL.Query()
		if q.Get("username") != "myuser" {
			t.Errorf("unexpected username: %s", q.Get("username"))
		}
		if q.Get("mobile") != "09120000000" {
			t.Errorf("unexpected mobile: %s", q.Get("mobile"))
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(successEnvelope(map[string]any{
			"messageId": 89545112,
			"cost":      1.0,
		}))
	}))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	resp, err := send.SendByURL(context.Background(), tc, send.ByURLRequest{
		Username: "myuser",
		Password: "mypass",
		Line:     300004505000017,
		Mobile:   "09120000000",
		Text:     "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MessageID != 89545112 {
		t.Errorf("unexpected messageID: %d", resp.MessageID)
	}
}

func TestSendByURL_Validation(t *testing.T) {
	tc := transport.New("fake-api-key", "http://unused.invalid", nil)

	tests := []struct {
		name string
		req  send.ByURLRequest
	}{
		{name: "missing username", req: send.ByURLRequest{Password: "p", Mobile: "m", Text: "t"}},
		{name: "missing password", req: send.ByURLRequest{Username: "u", Mobile: "m", Text: "t"}},
		{name: "missing mobile", req: send.ByURLRequest{Username: "u", Password: "p", Text: "t"}},
		{name: "missing text", req: send.ByURLRequest{Username: "u", Password: "p", Mobile: "m"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := send.SendByURL(context.Background(), tc, tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

// ---------- DeleteScheduled ----------

func TestDeleteScheduled_Success(t *testing.T) {
	const packID = "2b99e63c-9bf8-4a21-9bfe-3f72dc1b46f1"

	srv := newTestServer(t, http.MethodDelete, "/v1/send/scheduled/"+packID, successEnvelope(map[string]any{
		"returnedCreditCount": 10.0,
		"smsCount":            5,
	}))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	resp, err := send.DeleteScheduled(context.Background(), tc, packID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SmsCount != 5 {
		t.Errorf("unexpected smsCount: %d", resp.SmsCount)
	}
	if resp.ReturnedCreditCount != 10.0 {
		t.Errorf("unexpected returnedCreditCount: %v", resp.ReturnedCreditCount)
	}
}

func TestDeleteScheduled_APIError(t *testing.T) {
	const packID = "nonexistent-pack-id"

	srv := newTestServer(t, http.MethodDelete, "/v1/send/scheduled/"+packID, errorEnvelope(112, "record not found for deletion"))
	defer srv.Close()

	tc := transport.New("fake-api-key", srv.URL, nil)

	_, err := send.DeleteScheduled(context.Background(), tc, packID)

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T: %v", err, err)
	}
	if apiErr.Code != 112 {
		t.Errorf("expected code 112, got %d", apiErr.Code)
	}
}

func TestDeleteScheduled_Validation(t *testing.T) {
	tc := transport.New("fake-api-key", "http://unused.invalid", nil)

	_, err := send.DeleteScheduled(context.Background(), tc, "")
	if err == nil {
		t.Fatal("expected validation error for empty packID")
	}
}
