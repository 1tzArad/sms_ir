package statuscode_test

import (
	"testing"

	"github.com/1tzArad/sms_ir/statuscode"
)

const deliveryUnknownFallback = "وضعیت نامشخص"

func TestDeliveryState_String_KnownStatesAreMapped(t *testing.T) {
	known := []statuscode.DeliveryState{
		statuscode.DeliveredToDevice,
		statuscode.NotDeliveredToDevice,
		statuscode.ProcessingAtCarrier,
		statuscode.NotReceivedByCarrier,
		statuscode.ReceivedByCarrier,
		statuscode.DeliveryFailed,
		statuscode.DeliveryBlacklisted,
	}

	for _, d := range known {
		s := d.String()
		if s == "" {
			t.Errorf("(%d).String() returned empty message", d)
		}
		if s == deliveryUnknownFallback {
			t.Errorf("(%d).String() returned the unknown fallback", d)
		}
	}
}

func TestDeliveryState_String_UnknownStateReturnsFallback(t *testing.T) {
	for _, d := range []statuscode.DeliveryState{99, 0x7fffffff, -1} {
		if got := d.String(); got != deliveryUnknownFallback {
			t.Errorf("(%d).String() = %q, want %q", d, got, deliveryUnknownFallback)
		}
	}
}

func TestDeliveryState_String_ZeroValueIsUnknown(t *testing.T) {
	// DeliveryState(0) is not a documented constant, so it maps to the fallback.
	if statuscode.DeliveryState(0).String() != deliveryUnknownFallback {
		t.Error("DeliveryState(0) should map to the unknown fallback")
	}
}
