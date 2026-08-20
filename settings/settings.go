package settings

import (
	"context"
	"github.com/1tzArad/sms_ir/internal/transport"
	"net/http"
)

type GetCreditResponseData = float64

type ListLinesResponseData []int64

// Service

type Service struct {
	tc *transport.Client
}

func NewService(tc *transport.Client) *Service {
	return &Service{tc: tc}
}

func (s *Service) GetCredit(ctx context.Context) (*GetCreditResponseData, error) {
	var resp GetCreditResponseData
	if err := transport.Do[GetCreditResponseData](ctx, s.tc, http.MethodGet, "v1/credit", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) ListLines(ctx context.Context) (*ListLinesResponseData, error) {
	var resp ListLinesResponseData
	if err := transport.Do[ListLinesResponseData](ctx, s.tc, http.MethodGet, "/v1/line", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
