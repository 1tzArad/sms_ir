package sms_ir_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/1tzArad/sms_ir"
	"github.com/1tzArad/sms_ir/internal/transport"
	"github.com/1tzArad/sms_ir/statuscode"
)

func successEnvelope(data any) string {
	b, _ := json.Marshal(data)
	return `{"status":1,"message":"عملیات موفقیت‌آمیز بود","data":` + string(b) + `}`
}

func errorEnvelope(code int, message string) string {
	b, _ := json.Marshal(map[string]any{"status": code, "message": message})
	return string(b)
}

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func TestNew_ReturnsClientWithServices(t *testing.T) {
	c := sms_ir.New("api-key")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.Send == nil {
		t.Error("expected Send service to be initialized")
	}
	if c.Report == nil {
		t.Error("expected Report service to be initialized")
	}
	if c.Settings == nil {
		t.Error("expected Settings service to be initialized")
	}
}

func TestWithBaseURL_AppliesBaseURL(t *testing.T) {
	var seenPath string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(50.0)))
	})

	c := sms_ir.New("api-key", sms_ir.WithBaseURL(srv.URL), sms_ir.WithHTTPClient(srv.Client()))

	credit, err := c.Settings.GetCredit(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *credit != 50.0 {
		t.Errorf("credit = %v, want 50.0", *credit)
	}
	if seenPath != "/v1/credit" {
		t.Errorf("request path = %q, want %q", seenPath, "/v1/credit")
	}
}

func TestWithBaseURL_PropagatesAPIKey(t *testing.T) {
	var seenKey string
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seenKey = r.Header.Get("X-API-KEY")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(1.0)))
	})

	c := sms_ir.New("super-secret-key", sms_ir.WithBaseURL(srv.URL), sms_ir.WithHTTPClient(srv.Client()))
	if _, err := c.Settings.GetCredit(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seenKey != "super-secret-key" {
		t.Errorf("X-API-KEY = %q, want %q", seenKey, "super-secret-key")
	}
}

func TestWithHTTPClient_UsesProvidedClient(t *testing.T) {
	var calls int32
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(0.0)))
	})

	// A custom client with a short timeout still works against the test server.
	httpClient := srv.Client()
	c := sms_ir.New("api-key", sms_ir.WithBaseURL(srv.URL), sms_ir.WithHTTPClient(httpClient))

	if _, err := c.Settings.GetCredit(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("expected server to be called once, got %d", atomic.LoadInt32(&calls))
	}
}

func TestWithHTTPClient_NilKeepsDefault(t *testing.T) {
	// Passing a nil client should leave the default http.Client in place.
	c := sms_ir.New("api-key", sms_ir.WithBaseURL("http://localhost:0"), sms_ir.WithHTTPClient(nil))
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestClient_PropagatesAPIError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InvalidAPIKey), "کلید وب سرویس نامعتبر است")))
	})

	c := sms_ir.New("api-key", sms_ir.WithBaseURL(srv.URL), sms_ir.WithHTTPClient(srv.Client()))
	_, err := c.Settings.GetCredit(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// sms_ir.APIError is a type alias for transport.APIError.
	var apiErr *sms_ir.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *sms_ir.APIError, got %T", err)
	}
	if apiErr.Code != statuscode.InvalidAPIKey {
		t.Errorf("Code = %d, want %d", apiErr.Code, statuscode.InvalidAPIKey)
	}

	// errors.As should also resolve to the underlying transport.APIError type.
	var transportErr *transport.APIError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected *transport.APIError via alias, got %T", err)
	}
}

func TestNew_DefaultBaseURLIsProduction(t *testing.T) {
	// Without WithBaseURL the client targets the production host. We verify this
	// without real network by using a RoundTripper that captures the request URL.
	var captured *http.Request
	rt := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		captured = req
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(successEnvelope(0.0))),
			Header:     make(http.Header),
		}, nil
	})
	httpClient := &http.Client{Transport: rt}

	c := sms_ir.New("key", sms_ir.WithHTTPClient(httpClient))
	if _, err := c.Settings.GetCredit(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("expected request to be captured")
	}
	if got, want := captured.URL.Host, "api.sms.ir"; got != want {
		t.Errorf("base URL host = %q, want %q", got, want)
	}
	if got, want := captured.URL.Path, "/v1/credit"; got != want {
		t.Errorf("request path = %q, want %q", got, want)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
