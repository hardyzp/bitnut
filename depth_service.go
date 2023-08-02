package bitnut

import (
    "context"
    "net/http"
)

// DepthService show depth info
type DepthService struct {
    c      *Client
    symbol string
    limit  *int
}

// Symbol set symbol
func (s *DepthService) Symbol(symbol string) *DepthService {
    s.symbol = symbol
    return s
}

// Limit set limit
func (s *DepthService) Limit(limit int) *DepthService {
    s.limit = &limit
    return s
}

// Do send request
func (s *DepthService) Do(ctx context.Context, opts ...RequestOption) (depth *Depth, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/v1/tick/depth",
    }
    r.setParam("symbol", s.symbol)
    if s.limit != nil {
        r.setParam("limit", *s.limit)
    }
    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return nil, err
    }
    var res DepthResponse
    err = json.Unmarshal(data, &res)
    if err != nil {
        return nil, err
    }
    return &res.Data, nil
}

// DepthResponse define depth info with bids and asks
type Depth struct {
    Bids [][2]string `json:"bids"`
    Asks [][2]string `json:"asks"`
}

type DepthResponse struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data Depth  `json:"data"`
}
