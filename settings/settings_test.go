package settings_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1tzArad/sms_ir/internal/transport"
	"github.com/1tzArad/sms_ir/settings"
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

func newService(t *testing.T, handler http.HandlerFunc) *settings.Service {
	t.Helper()
	srv := newTestServer(t, handler)
	tc := transport.New("test-api-key", srv.URL, srv.Client())
	return settings.NewService(tc)
}

func TestNewService_ReturnsNonNil(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	tc := transport.New("key", srv.URL, srv.Client())
	if s := settings.NewService(tc); s == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestGetCredit_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/credit" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/credit")
		}
		if got := r.Header.Get("X-API-KEY"); got != "test-api-key" {
			t.Errorf("X-API-KEY = %q, want %q", got, "test-api-key")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(123.45)))
	})

	credit, err := svc.GetCredit(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if credit == nil {
		t.Fatal("expected non-nil credit pointer")
	}
	if *credit != 123.45 {
		t.Errorf("credit = %v, want 123.45", *credit)
	}
}

func TestGetCredit_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InvalidAPIKey), "کلید وب سرویس نامعتبر است")))
	})

	credit, err := svc.GetCredit(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if credit != nil {
		t.Errorf("expected nil credit, got %v", *credit)
	}

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.InvalidAPIKey {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.InvalidAPIKey)
	}
}

func TestListLines_Success(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/line" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/line")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope([]int64{100010001, 100010002})))
	})

	lines, err := svc.ListLines(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines == nil {
		t.Fatal("expected non-nil lines pointer")
	}
	if len(*lines) != 2 {
		t.Fatalf("lines len = %d, want 2", len(*lines))
	}
	if (*lines)[0] != 100010001 || (*lines)[1] != 100010002 {
		t.Errorf("unexpected lines: %v", *lines)
	}
}

func TestListLines_APIError(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InactiveAccount), "حساب کاربری غیرفعال است")))
	})

	lines, err := svc.ListLines(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if lines != nil {
		t.Errorf("expected nil lines, got %v", *lines)
	}

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *transport.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.InactiveAccount {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.InactiveAccount)
	}
}
