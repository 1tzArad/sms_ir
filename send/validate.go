package send

import "fmt"

func validateBulk(req BulkRequest) error {
	if req.LineNumber == 0 {
		return fmt.Errorf("smsir: lineNumber is required")
	}
	if req.MessageText == "" {
		return fmt.Errorf("smsir: messageText is required")
	}
	if len(req.Mobiles) == 0 {
		return fmt.Errorf("smsir: mobiles list cannot be empty")
	}
	if len(req.Mobiles) > 100 {
		return fmt.Errorf("smsir: mobiles list cannot exceed 100 numbers, got %d", len(req.Mobiles))
	}
	return nil
}

func validateLikeToLike(req LikeToLikeRequest) error {
	if req.LineNumber == 0 {
		return fmt.Errorf("smsir: lineNumber is required")
	}
	if len(req.Mobiles) == 0 || len(req.MessageTexts) == 0 {
		return fmt.Errorf("smsir: mobiles and messageTexts cannot be empty")
	}
	if len(req.Mobiles) != len(req.MessageTexts) {
		return fmt.Errorf(
			"smsir: mobiles and messageTexts must have equal length, got %d and %d",
			len(req.Mobiles), len(req.MessageTexts),
		)
	}
	return nil
}
