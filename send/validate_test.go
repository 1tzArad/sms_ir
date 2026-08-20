package send_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/1tzArad/sms_ir/send"
)

func TestBulk_Validation(t *testing.T) {
	baseMobiles := []string{"09120000000"}

	tooManyMobiles := make([]string, 101)
	for i := range tooManyMobiles {
		tooManyMobiles[i] = "09120000000"
	}

	tests := []struct {
		name string
		req  send.BulkRequest
		want string
	}{
		{
			name: "missing line number",
			req: send.BulkRequest{
				MessageText: "hi",
				Mobiles:     baseMobiles,
			},
			want: "lineNumber is required",
		},
		{
			name: "missing message text",
			req: send.BulkRequest{
				LineNumber: 1,
				Mobiles:    baseMobiles,
			},
			want: "messageText is required",
		},
		{
			name: "empty mobiles",
			req: send.BulkRequest{
				LineNumber:  1,
				MessageText: "hi",
				Mobiles:     []string{},
			},
			want: "mobiles list cannot be empty",
		},
		{
			name: "too many mobiles",
			req: send.BulkRequest{
				LineNumber:  1,
				MessageText: "hi",
				Mobiles:     tooManyMobiles,
			},
			want: "mobiles list cannot exceed 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("server should not be called for validation failure")
			})
			_, err := svc.Bulk(context.Background(), tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}

func TestLikeToLike_Validation(t *testing.T) {
	tests := []struct {
		name string
		req  send.LikeToLikeRequest
		want string
	}{
		{
			name: "missing line number",
			req: send.LikeToLikeRequest{
				MessageTexts: []string{"hi"},
				Mobiles:      []string{"09120000000"},
			},
			want: "lineNumber is required",
		},
		{
			name: "empty mobiles and texts",
			req: send.LikeToLikeRequest{
				LineNumber: 1,
			},
			want: "mobiles and messageTexts cannot be empty",
		},
		{
			name: "length mismatch",
			req: send.LikeToLikeRequest{
				LineNumber:   1,
				MessageTexts: []string{"a", "b"},
				Mobiles:      []string{"09120000000"},
			},
			want: "must have equal length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("server should not be called for validation failure")
			})
			_, err := svc.LikeToLike(context.Background(), tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}

func TestVerify_Validation(t *testing.T) {
	tests := []struct {
		name string
		req  send.VerifyRequest
		want string
	}{
		{
			name: "missing mobile",
			req: send.VerifyRequest{
				TemplateID: 123,
			},
			want: "mobile is required",
		},
		{
			name: "missing template id",
			req: send.VerifyRequest{
				Mobile: "09120000000",
			},
			want: "templateId is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("server should not be called for validation failure")
			})
			_, err := svc.Verify(context.Background(), tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}

func TestSendByURL_Validation(t *testing.T) {
	tests := []struct {
		name string
		req  send.SendByURLRequest
		want string
	}{
		{
			name: "missing username and password",
			req: send.SendByURLRequest{
				Mobile: "09120000000",
				Text:   "hi",
			},
			want: "username and password are required",
		},
		{
			name: "missing mobile and text",
			req: send.SendByURLRequest{
				Username: "u",
				Password: "p",
			},
			want: "mobile and text are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("server should not be called for validation failure")
			})
			_, err := svc.SendByURL(context.Background(), tt.req)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}

func TestDeleteScheduled_Validation(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be called for validation failure")
	})
	_, err := svc.DeleteScheduled(context.Background(), "")
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !strings.Contains(err.Error(), "packID is required") {
		t.Errorf("error = %q, want substring %q", err.Error(), "packID is required")
	}
}
