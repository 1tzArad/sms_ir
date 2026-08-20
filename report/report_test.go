package report_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1tzArad/sms_ir/internal/transport"
	"github.com/1tzArad/sms_ir/report"
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

func newService(t *testing.T, handler http.HandlerFunc) *report.Service {
	t.Helper()
	srv := newTestServer(t, handler)
	tc := transport.New("test-api-key", srv.URL, srv.Client())
	return report.NewService(tc)
}

func int64Ptr(i int64) *int64        { return &i }
func deliveryStatePtr(d statuscode.DeliveryState) *statuscode.DeliveryState { return &d }

func TestNewService_ReturnsNonNil(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	tc := transport.New("key", srv.URL, srv.Client())
	if s := report.NewService(tc); s == nil {
		t.Fatal("expected non-nil service")
	}
}

// ---------------------------------------------------------------------------
// GetSentMessageStatus
// ---------------------------------------------------------------------------

func TestGetSentMessageStatus_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/send/msg-1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send/msg-1")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(report.SentMessageReport{
			MessageID:        123,
			Mobile:           "09120000000",
			MessageText:      "hello",
			SendDateTime:     1700000000,
			LineNumber:       100010001,
			Cost:             0.5,
			DeliveryState:    deliveryStatePtr(statuscode.DeliveredToDevice),
			DeliveryDateTime: int64Ptr(1700000100),
		})))
	})

	resp, err := svc.GetSentMessageStatus(context.Background(), "msg-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MessageID != 123 {
		t.Errorf("MessageID = %d, want 123", resp.MessageID)
	}
	if resp.Mobile != "09120000000" {
		t.Errorf("Mobile = %q", resp.Mobile)
	}
	if resp.DeliveryState == nil || *resp.DeliveryState != statuscode.DeliveredToDevice {
		t.Errorf("unexpected DeliveryState: %v", resp.DeliveryState)
	}
}

func TestGetSentMessageStatus_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.PackNotFound), "با این شناسه ارسالی ثبت نشده است")))
	})

	_, err := svc.GetSentMessageStatus(context.Background(), "msg-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.PackNotFound {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.PackNotFound)
	}
}

// ---------------------------------------------------------------------------
// ListTodaySendPacks
// ---------------------------------------------------------------------------

func TestListTodaySendPacks_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pack" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/pack")
		}
		if r.URL.Query().Get("pageSize") != "10" || r.URL.Query().Get("pageNumber") != "1" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]report.SendPackReport{
			{PackID: "p1", RecipientCount: 5, CreationDateTime: 1700000000},
		})))
	})

	resp, err := svc.ListTodaySendPacks(context.Background(), report.ListTodaySendPacksRequestParams{
		PageSize: 10, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*resp) != 1 {
		t.Fatalf("len = %d, want 1", len(*resp))
	}
	if (*resp)[0].PackID != "p1" || (*resp)[0].RecipientCount != 5 {
		t.Errorf("unexpected item: %+v", (*resp)[0])
	}
}

func TestListTodaySendPacks_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InsufficientCredit), "اعتبار کافی نمی‌باشد")))
	})

	_, err := svc.ListTodaySendPacks(context.Background(), report.ListTodaySendPacksRequestParams{})
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
// GetSendPackReport
// ---------------------------------------------------------------------------

func TestGetSendPackReport_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/pack/pack-1" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/pack/pack-1")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]report.SentMessageReport{
			{MessageID: 1, Mobile: "09120000000", Cost: 1.0},
		})))
	})

	resp, err := svc.GetSendPackReport(context.Background(), report.GetSendPackReportRequestParams{
		PackID: "pack-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*resp) != 1 || (*resp)[0].MessageID != 1 {
		t.Errorf("unexpected response: %+v", *resp)
	}
}

func TestGetSendPackReport_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.PackNotFound), "با این شناسه ارسالی ثبت نشده است")))
	})

	_, err := svc.GetSendPackReport(context.Background(), report.GetSendPackReportRequestParams{PackID: "pack-1"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.PackNotFound {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.PackNotFound)
	}
}

// ---------------------------------------------------------------------------
// ListTodaySentMessages
// ---------------------------------------------------------------------------

func TestListTodaySentMessages_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/send/live" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send/live")
		}
		if r.URL.Query().Get("pageSize") != "20" || r.URL.Query().Get("pageNumber") != "2" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]report.SentMessageReport{
			{MessageID: 1, Mobile: "09120000000"},
		})))
	})

	resp, err := svc.ListTodaySentMessages(context.Background(), report.ListTodaySentMessagesRequestParams{
		PageSize: 20, PageNumber: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*resp) != 1 {
		t.Fatalf("len = %d, want 1", len(*resp))
	}
}

func TestListTodaySentMessages_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.SystemError), "مشکلی در سامانه رخ داده است، لطفا با پشتیبانی در تماس باشید")))
	})

	_, err := svc.ListTodaySentMessages(context.Background(), report.ListTodaySentMessagesRequestParams{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.SystemError {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.SystemError)
	}
}

// ---------------------------------------------------------------------------
// ListArchivedSentMessages
// ---------------------------------------------------------------------------

func TestListArchivedSentMessages_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/send/archive" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/send/archive")
		}
		q := r.URL.Query()
		if q.Get("fromDate") != "1700000000" || q.Get("toDate") != "1700100000" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		if q.Get("pageSize") != "10" || q.Get("pageNumber") != "1" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]report.SentMessageReport{
			{MessageID: 5, Mobile: "09120000000"},
		})))
	})

	resp, err := svc.ListArchivedSentMessages(context.Background(), report.ListArchivedSentMessagesRequestParams{
		FromDate: 1700000000, ToDate: 1700100000, PageSize: 10, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*resp) != 1 || (*resp)[0].MessageID != 5 {
		t.Errorf("unexpected response: %+v", *resp)
	}
}

func TestListArchivedSentMessages_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.RateLimitExceeded), "تعداد درخواست بیشتر از حد مجاز است")))
	})

	_, err := svc.ListArchivedSentMessages(context.Background(), report.ListArchivedSentMessagesRequestParams{
		FromDate: 1, ToDate: 2,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.RateLimitExceeded {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.RateLimitExceeded)
	}
}

// ---------------------------------------------------------------------------
// ListLatestReceivedMessages
// ---------------------------------------------------------------------------

func TestListLatestReceivedMessages_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/receive/latest" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/receive/latest")
		}
		if r.URL.Query().Get("count") != "10" {
			t.Errorf("count query = %q, want %q", r.URL.Query().Get("count"), "10")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]report.ReceivedMessageReport{
			{ReceiveReturnID: 1, MessageText: "hi", Number: 100010001, Mobile: "09120000000", ReceivedDateTime: 1700000000},
		})))
	})

	resp, err := svc.ListLatestReceivedMessages(context.Background(), report.GetLatestReceivedMessagesRequestParams{Count: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*resp) != 1 {
		t.Fatalf("len = %d, want 1", len(*resp))
	}
	if (*resp)[0].MessageText != "hi" {
		t.Errorf("unexpected MessageText: %q", (*resp)[0].MessageText)
	}
}

func TestListLatestReceivedMessages_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InvalidAPIKey), "کلید وب سرویس نامعتبر است")))
	})

	_, err := svc.ListLatestReceivedMessages(context.Background(), report.GetLatestReceivedMessagesRequestParams{Count: 10})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.InvalidAPIKey {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.InvalidAPIKey)
	}
}

// ---------------------------------------------------------------------------
// ListTodayReceivedMessages
// ---------------------------------------------------------------------------

func TestListTodayReceivedMessages_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/receive/live" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/receive/live")
		}
		q := r.URL.Query()
		if q.Get("sortByNewest") != "true" || q.Get("pageSize") != "5" || q.Get("pageNumber") != "1" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]report.ReceivedMessageReport{
			{ReceiveReturnID: 7, MessageText: "ok"},
		})))
	})

	resp, err := svc.ListTodayReceivedMessages(context.Background(), report.ListTodayReceivedMessagesRequestParams{
		PageSize: 5, PageNumber: 1, SortByNewest: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*resp) != 1 || (*resp)[0].ReceiveReturnID != 7 {
		t.Errorf("unexpected response: %+v", *resp)
	}
}

func TestListTodayReceivedMessages_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InactiveAccount), "حساب کاربری غیرفعال است")))
	})

	_, err := svc.ListTodayReceivedMessages(context.Background(), report.ListTodayReceivedMessagesRequestParams{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.InactiveAccount {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.InactiveAccount)
	}
}

// ---------------------------------------------------------------------------
// ListArchivedReceivedMessages
// ---------------------------------------------------------------------------

func TestListArchivedReceivedMessages_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/receive/archive" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/receive/archive")
		}
		q := r.URL.Query()
		if q.Get("fromDate") != "1700000000" || q.Get("toDate") != "1700100000" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]report.ReceivedMessageReport{
			{ReceiveReturnID: 3, MessageText: "archived"},
		})))
	})

	resp, err := svc.ListArchivedReceivedMessages(context.Background(), report.ListArchivedReceivedMessagesRequestParams{
		FromDate: 1700000000, ToDate: 1700100000, PageSize: 10, PageNumber: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*resp) != 1 || (*resp)[0].ReceiveReturnID != 3 {
		t.Errorf("unexpected response: %+v", *resp)
	}
}

func TestListArchivedReceivedMessages_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.RateLimitExceeded), "تعداد درخواست بیشتر از حد مجاز است")))
	})

	_, err := svc.ListArchivedReceivedMessages(context.Background(), report.ListArchivedReceivedMessagesRequestParams{
		FromDate: 1, ToDate: 2,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.RateLimitExceeded {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.RateLimitExceeded)
	}
}
