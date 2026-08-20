package send

import (
	"context"
	"fmt"
	"github.com/1tzArad/sms_ir/internal/transport"
	"net/http"
	"net/url"
	"strconv"
)

type Service struct {
	tc *transport.Client
}

func NewService(tc *transport.Client) *Service {
	return &Service{tc: tc}
}

func (s *Service) Bulk(ctx context.Context, req BulkRequest) (*BulkResponseData, error) {
	if err := validateBulk(req); err != nil {
		return nil, err
	}

	var resp BulkResponseData
	err := transport.Do(ctx, s.tc, http.MethodPost, "/v1/send/bulk", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) LikeToLike(ctx context.Context, req LikeToLikeRequest) (*LikeToLikeResponseData, error) {
	if err := validateLikeToLike(req); err != nil {
		return nil, err
	}

	var resp LikeToLikeResponseData
	err := transport.Do(ctx, s.tc, http.MethodPost, "/v1/send/likeToLike", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) Verify(ctx context.Context, req VerifyRequest) (*VerifyResponseData, error) {
	if req.Mobile == "" {
		return nil, fmt.Errorf("sms_ir: mobile is required")
	}
	if req.TemplateID == 0 {
		return nil, fmt.Errorf("sms_ir: templateId is required")
	}

	var resp VerifyResponseData
	err := transport.Do(ctx, s.tc, http.MethodPost, "/v1/send/verify", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) SendByURL(ctx context.Context, req SendByURLRequest) (*SendByURLResponseData, error) {
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("sms_ir: username and password are required")
	}
	if req.Mobile == "" || req.Text == "" {
		return nil, fmt.Errorf("sms_ir: mobile and text are required")
	}

	q := url.Values{}
	q.Set("username", req.Username)
	q.Set("password", req.Password)
	q.Set("line", strconv.FormatInt(req.Line, 10))
	q.Set("mobile", req.Mobile)
	q.Set("text", req.Text)

	path := "/v1/send?" + q.Encode()

	var resp SendByURLResponseData
	err := transport.Do[SendByURLResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *Service) DeleteScheduled(ctx context.Context, packID string) (*DeleteScheduledResponseData, error) {
	if packID == "" {
		return nil, fmt.Errorf("sms_ir: packID is required")
	}

	path := "/v1/send/scheduled/" + url.PathEscape(packID)

	var resp DeleteScheduledResponseData
	err := transport.Do[DeleteScheduledResponseData](ctx, s.tc, http.MethodDelete, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
