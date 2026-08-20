package transport

import (
	"fmt"
	"github.com/1tzArad/sms_ir/statuscode"
)

type APIError struct {
	Code    statuscode.Code
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("sms_ir: request failed [code=%d]: %s", e.Code, e.Message)
}
