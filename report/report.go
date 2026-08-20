package report

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

func (s *Service) GetSentMessageStatus(ctx context.Context, messageId string) (*GetSentMessageStatusResponseData, error) {
	if messageId == "" {
		return nil, fmt.Errorf("sms_ir: message id is required")
	}

	path := "/v1/send/" + messageId
	var resp GetSentMessageStatusResponseData
	if err := transport.Do[GetSentMessageStatusResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *Service) ListTodaySendPacks(ctx context.Context, params ListTodaySendPacksRequestParams) (*ListTodaySendPacksResponseData, error) {
	q := url.Values{}
	if params.PageNumber != 0 {
		q.Set("pageNumber", strconv.Itoa(params.PageNumber))
	}
	if params.PageSize != 0 {
		q.Set("pageSize", strconv.Itoa(params.PageSize))
	}

	path := "/v1/pack?" + q.Encode()

	var resp ListTodaySendPacksResponseData
	if err := transport.Do[ListTodaySendPacksResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *Service) GetSendPackReport(ctx context.Context, params GetSendPackReportRequestParams) (*GetSendPackReportResponseData, error) {
	if params.PackID == "" {
		return nil, fmt.Errorf("sms_ir: PackID is required")
	}
	path := "/v1/pack/" + params.PackID
	var resp GetSendPackReportResponseData
	if err := transport.Do[GetSendPackReportResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *Service) ListTodaySentMessages(ctx context.Context, params ListTodaySentMessagesRequestParams) (*ListTodaySentMessagesResponseData, error) {
	q := url.Values{}
	if params.PageSize != 0 {
		q.Set("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.PageNumber != 0 {
		q.Set("pageNumber", strconv.Itoa(params.PageNumber))
	}

	path := "/v1/send/live?" + q.Encode()

	var resp ListTodaySentMessagesResponseData

	if err := transport.Do[ListTodaySentMessagesResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *Service) ListArchivedSentMessages(ctx context.Context, params ListArchivedSentMessagesRequestParams) (*ListArchivedSentMessagesResponseData, error) {
	if params.FromDate == 0 {
		return nil, fmt.Errorf("fromDate is required")
	}
	if params.ToDate == 0 {
		return nil, fmt.Errorf("toDate is required")
	}
	q := url.Values{}
	q.Set("fromDate", strconv.FormatInt(params.FromDate, 10))
	q.Set("toDate", strconv.FormatInt(params.ToDate, 10))

	if params.PageSize != 0 {
		q.Set("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.PageNumber != 0 {
		q.Set("pageNumber", strconv.Itoa(params.PageNumber))
	}

	path := "/v1/send/archive?" + q.Encode()

	var resp ListArchivedSentMessagesResponseData
	if err := transport.Do[ListArchivedSentMessagesResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *Service) ListLatestReceivedMessages(ctx context.Context, params GetLatestReceivedMessagesRequestParams) (*GetLatestReceivedMessagesResponseData, error) {
	q := url.Values{}
	if params.Count != 0 {
		q.Set("count", strconv.Itoa(params.Count))
	}

	path := "/v1/receive/latest?" + q.Encode()
	var resp GetLatestReceivedMessagesResponseData
	if err := transport.Do[GetLatestReceivedMessagesResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *Service) ListTodayReceivedMessages(ctx context.Context, params ListTodayReceivedMessagesRequestParams) (*ListTodayReceivedMessagesResponseData, error) {
	q := url.Values{}
	q.Set("sortByNewest", strconv.FormatBool(params.SortByNewest))
	if params.PageSize != 0 {
		q.Set("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.PageNumber != 0 {
		q.Set("pageNumber", strconv.Itoa(params.PageNumber))
	}

	path := "/v1/receive/live?" + q.Encode()
	var resp ListTodayReceivedMessagesResponseData
	if err := transport.Do[ListTodayReceivedMessagesResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (s *Service) ListArchivedReceivedMessages(ctx context.Context, params ListArchivedReceivedMessagesRequestParams) (*ListArchivedReceivedMessagesResponseData, error) {
	if params.FromDate == 0 {
		return nil, fmt.Errorf("fromDate is required")
	}
	if params.ToDate == 0 {
		return nil, fmt.Errorf("toDate is required")
	}
	q := url.Values{}
	if params.PageSize != 0 {
		q.Set("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.PageNumber != 0 {
		q.Set("pageNumber", strconv.Itoa(params.PageNumber))
	}

	path := "/v1/receive/archive?" + q.Encode()

	var resp ListArchivedReceivedMessagesResponseData
	if err := transport.Do[ListArchivedReceivedMessagesResponseData](ctx, s.tc, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
