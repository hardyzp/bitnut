package bitnut

import (
    "context"
    "net/http"
)

// ServerTimeService get server time
type ServerTimeService struct {
    c *Client
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...RequestOption) (serverTime int64, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/v1/time",
    }
    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return 0, err
    }
    j, err := newJSON(data)
    if err != nil {
        return 0, err
    }
    serverTime = j.Get("data").Get("ts").MustInt64()
    return serverTime, nil
}

// SetServerTimeService set server time
type SetServerTimeService struct {
    c *Client
}

// Do send request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (timeOffset int64, err error) {
    serverTime, err := s.c.NewServerTimeService().Do(ctx)
    if err != nil {
        return 0, err
    }
    timeOffset = currentTimestamp() - serverTime
    s.c.TimeOffset = timeOffset
    return timeOffset, nil
}
