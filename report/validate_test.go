package report_test

import (
	"context"
	"strings"
	"testing"

	"github.com/1tzArad/sms_ir/report"
)

func TestGetSentMessageStatus_Validation(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be called for validation failure")
	})

	_, err := svc.GetSentMessageStatus(context.Background(), "")
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !strings.Contains(err.Error(), "message id is required") {
		t.Errorf("error = %q, want substring %q", err.Error(), "message id is required")
	}
}

func TestGetSendPackReport_Validation(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be called for validation failure")
	})

	_, err := svc.GetSendPackReport(context.Background(), report.GetSendPackReportRequestParams{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !strings.Contains(err.Error(), "PackID is required") {
		t.Errorf("error = %q, want substring %q", err.Error(), "PackID is required")
	}
}

func TestListArchivedSentMessages_Validation(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be called for validation failure")
	})

	tests := []struct {
		name   string
		params report.ListArchivedSentMessagesRequestParams
		want   string
	}{
		{
			name:   "missing from date",
			params: report.ListArchivedSentMessagesRequestParams{ToDate: 2},
			want:   "fromDate is required",
		},
		{
			name:   "missing to date",
			params: report.ListArchivedSentMessagesRequestParams{FromDate: 1},
			want:   "toDate is required",
		},
		{
			name:   "missing both dates",
			params: report.ListArchivedSentMessagesRequestParams{},
			want:   "fromDate is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ListArchivedSentMessages(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}

func TestListArchivedReceivedMessages_Validation(t *testing.T) {
	svc := newService(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be called for validation failure")
	})

	tests := []struct {
		name   string
		params report.ListArchivedReceivedMessagesRequestParams
		want   string
	}{
		{
			name:   "missing from date",
			params: report.ListArchivedReceivedMessagesRequestParams{ToDate: 2},
			want:   "fromDate is required",
		},
		{
			name:   "missing to date",
			params: report.ListArchivedReceivedMessagesRequestParams{FromDate: 1},
			want:   "toDate is required",
		},
		{
			name:   "missing both dates",
			params: report.ListArchivedReceivedMessagesRequestParams{},
			want:   "fromDate is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ListArchivedReceivedMessages(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}
