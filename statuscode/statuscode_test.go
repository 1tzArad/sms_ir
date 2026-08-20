package statuscode_test

import (
	"testing"

	"github.com/1tzArad/sms_ir/statuscode"
)

func TestCode_String_KnownCodesAreMapped(t *testing.T) {
	known := []statuscode.Code{
		statuscode.Success,
		statuscode.SystemError,
		statuscode.InvalidAPIKey,
		statuscode.InactiveAPIKey,
		statuscode.IPNotWhitelisted,
		statuscode.InactiveAccount,
		statuscode.SuspendedAccount,
		statuscode.RateLimitExceeded,
		statuscode.InvalidLineNumber,
		statuscode.InsufficientCredit,
		statuscode.EmptyMessageText,
		statuscode.InvalidMobileNumbers,
		statuscode.TooManyMobiles,
		statuscode.TooManyTexts,
		statuscode.EmptyMobilesList,
		statuscode.EmptyTextsList,
		statuscode.InvalidSendDateTime,
		statuscode.MobilesAndTextsLengthMismatch,
		statuscode.PackNotFound,
		statuscode.RecordNotFoundForDeletion,
		statuscode.TemplateNotFound,
		statuscode.ParameterValueTooLong,
		statuscode.MobilesBlacklisted,
		statuscode.EmptyParameterName,
		statuscode.MessageTextNotApproved,
		statuscode.TooManyMessages,
		statuscode.PlanUpgradeRequired,
		statuscode.LineNeedsActivation,
	}

	for _, c := range known {
		s := c.String()
		// Every documented code must map to a real (non-empty, non-fallback) message.
		if s == "" {
			t.Errorf("(%d).String() returned empty message", c)
		}
		if s == unknownFallback {
			t.Errorf("(%d).String() returned the unknown-status fallback", c)
		}
	}
}

func TestCode_String_UnknownCodeReturnsFallback(t *testing.T) {
	for _, c := range []statuscode.Code{999, 0x7fffffff, -1} {
		if got := c.String(); got != unknownFallback {
			t.Errorf("(%d).String() = %q, want %q", c, got, unknownFallback)
		}
	}
}

const unknownFallback = "کد وضعیت نامشخص"

func TestCode_String_ZeroValueIsSystemError(t *testing.T) {
	// Code(0) is the SystemError constant and must not be treated as unknown.
	if statuscode.Code(0) != statuscode.SystemError {
		t.Fatal("Code(0) is not SystemError")
	}
	if statuscode.Code(0).String() == unknownFallback {
		t.Error("SystemError mapped to unknown fallback")
	}
}

func TestCode_IsRateLimit(t *testing.T) {
	tests := []struct {
		name string
		code statuscode.Code
		want bool
	}{
		{"rate limit exceeded", statuscode.RateLimitExceeded, true},
		{"success", statuscode.Success, false},
		{"invalid api key", statuscode.InvalidAPIKey, false},
		{"insufficient credit", statuscode.InsufficientCredit, false},
		{"unknown", statuscode.Code(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.IsRateLimit(); got != tt.want {
				t.Errorf("(%d).IsRateLimit() = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestCode_IsAuthError(t *testing.T) {
	authCodes := []statuscode.Code{
		statuscode.InvalidAPIKey,
		statuscode.InactiveAPIKey,
		statuscode.IPNotWhitelisted,
	}
	for _, c := range authCodes {
		if !c.IsAuthError() {
			t.Errorf("(%d).IsAuthError() = false, want true", c)
		}
	}

	nonAuthCodes := []statuscode.Code{
		statuscode.Success,
		statuscode.InactiveAccount,
		statuscode.SuspendedAccount,
		statuscode.RateLimitExceeded,
		statuscode.InsufficientCredit,
		statuscode.InvalidMobileNumbers,
		statuscode.Code(999),
	}
	for _, c := range nonAuthCodes {
		if c.IsAuthError() {
			t.Errorf("(%d).IsAuthError() = true, want false", c)
		}
	}
}
