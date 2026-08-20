package transport_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/1tzArad/sms_ir/internal/transport"
	"github.com/1tzArad/sms_ir/statuscode"
)

type testPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

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

func TestNew_ReturnsNonNilClient(t *testing.T) {
	if c := transport.New("key", "https://api.test", nil); c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNew_DefaultBaseURL(t *testing.T) {
	var captured *http.Request
	rt := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		captured = req
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(successEnvelope(testPayload{ID: 1}))),
			Header:     make(http.Header),
		}, nil
	})
	httpClient := &http.Client{Transport: rt}

	// Empty baseURL should default to the production default base URL.
	c := transport.New("my-api-key", "", httpClient)

	var out testPayload
	if err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, &out); err != nil {
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
	if got, want := captured.Header.Get("X-API-KEY"), "my-api-key"; got != want {
		t.Errorf("X-API-KEY header = %q, want %q", got, want)
	}
	if got, want := captured.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("Content-Type header = %q, want %q", got, want)
	}
}

func TestDo_Success(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want %q", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/v1/credit" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/credit")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(testPayload{ID: 7, Name: "ok"})))
	})

	c := transport.New("key", srv.URL, srv.Client())
	var out testPayload
	if err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != 7 || out.Name != "ok" {
		t.Errorf("unexpected decoded data: %+v", out)
	}
}

func TestDo_SuccessWithBody(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(testPayload{ID: 1})))
	})

	c := transport.New("key", srv.URL, srv.Client())
	body := map[string]string{"hello": "world"}
	var out testPayload
	if err := transport.Do[testPayload](context.Background(), c, http.MethodPost, "/v1/send", body, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want %q", gotMethod, http.MethodPost)
	}
	if gotBody["hello"] != "world" {
		t.Errorf("unexpected request body: %v", gotBody)
	}
}

func TestDo_APIError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errorEnvelope(int(statuscode.InvalidAPIKey), "کلید وب سرویس نامعتبر است")))
	})

	c := transport.New("key", srv.URL, srv.Client())
	var out testPayload
	err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, &out)
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
	if apiErr.Message != "کلید وب سرویس نامعتبر است" {
		t.Errorf("Message = %q", apiErr.Message)
	}
	if out != (testPayload{}) {
		t.Errorf("out should remain zero value: %+v", out)
	}
}

func TestDo_HTTPTransportError(t *testing.T) {
	// A closed server makes the underlying dial fail, surfacing a transport error.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	addr := srv.URL
	srv.Close()

	c := transport.New("key", addr, &http.Client{})
	var out testPayload
	err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "http request failed") {
		t.Errorf("error = %q, want substring %q", err.Error(), "http request failed")
	}
	var apiErr *transport.APIError
	if errors.As(err, &apiErr) {
		t.Errorf("expected non-API error, got *transport.APIError: %v", err)
	}
}

func TestDo_UnmarshalError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not-json"))
	})

	c := transport.New("key", srv.URL, srv.Client())
	var out testPayload
	err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode response") {
		t.Errorf("error = %q, want substring %q", err.Error(), "failed to decode response")
	}
}

func TestDo_RequestBuildError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})

	// Malformed scheme makes http.NewRequest fail.
	c := transport.New("key", "://invalid-host", srv.Client())
	var out testPayload
	err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to build request") {
		t.Errorf("error = %q, want substring %q", err.Error(), "failed to build request")
	}
}

func TestDo_MarshalBodyError(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})

	c := transport.New("key", srv.URL, srv.Client())
	// A struct containing a channel cannot be JSON-marshaled.
	body := struct {
		Invalid chan int `json:"invalid"`
	}{}
	var out testPayload
	err := transport.Do[testPayload](context.Background(), c, http.MethodPost, "/v1/send", body, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to marshal request body") {
		t.Errorf("error = %q, want substring %q", err.Error(), "failed to marshal request body")
	}
}

func TestDo_EmptyBodyIsAllowed(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength != 0 {
			t.Errorf("expected empty body, got content-length %d", r.ContentLength)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(testPayload{ID: 9})))
	})

	c := transport.New("key", srv.URL, srv.Client())
	var out testPayload
	if err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.ID != 9 {
		t.Errorf("ID = %d, want 9", out.ID)
	}
}

func TestAPIError_ErrorAndFields(t *testing.T) {
	e := &transport.APIError{Code: statuscode.InsufficientCredit, Message: "اعتبار کافی نمی‌باشد"}

	want := "sms_ir: request failed [code=102]: اعتبار کافی نمی‌باشد"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestDo_SuccessNilOut(t *testing.T) {
	// When out is nil the decoded data is discarded without panicking.
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successEnvelope(testPayload{ID: 1})))
	})

	c := transport.New("key", srv.URL, srv.Client())
	err := transport.Do[testPayload](context.Background(), c, http.MethodGet, "/v1/credit", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
