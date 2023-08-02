package bitnut

import (
    "context"
    "net/http"
)

// ListTradesService list trades
type ListTradesService struct {
    c         *Client
    symbol    string
    startTime *int64
    endTime   *int64
    limit     *int
    orderId   *int64
}

// Symbol set symbol
func (s *ListTradesService) Symbol(symbol string) *ListTradesService {
    s.symbol = symbol
    return s
}

// StartTime set starttime
func (s *ListTradesService) StartTime(startTime int64) *ListTradesService {
    s.startTime = &startTime
    return s
}

// EndTime set endtime
func (s *ListTradesService) EndTime(endTime int64) *ListTradesService {
    s.endTime = &endTime
    return s
}

// Limit set limit
func (s *ListTradesService) Limit(limit int) *ListTradesService {
    s.limit = &limit
    return s
}

// OrderId set OrderId
func (s *ListTradesService) OrderId(OrderId int64) *ListTradesService {
    s.orderId = &OrderId
    return s
}

// Do send request
func (s *ListTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*Trade, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/api/v3/myTrades",
        secType:  secTypeSigned,
    }
    r.setParam("symbol", s.symbol)
    if s.limit != nil {
        r.setParam("limit", *s.limit)
    }
    if s.startTime != nil {
        r.setParam("startTime", *s.startTime)
    }
    if s.endTime != nil {
        r.setParam("endTime", *s.endTime)
    }

    if s.orderId != nil {
        r.setParam("orderId", *s.orderId)
    }
    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return []*Trade{}, err
    }
    res = make([]*Trade, 0)
    err = json.Unmarshal(data, &res)
    if err != nil {
        return []*Trade{}, err
    }
    return res, nil
}

// HistoricalTradesService trades
type HistoricalTradesService struct {
    c      *Client
    symbol string
    limit  *int
}

// Symbol set symbol
func (s *HistoricalTradesService) Symbol(symbol string) *HistoricalTradesService {
    s.symbol = symbol
    return s
}

// Limit set limit
func (s *HistoricalTradesService) Limit(limit int) *HistoricalTradesService {
    s.limit = &limit
    return s
}

// Do send request
func (s *HistoricalTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*Trade, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/api/v3/historicalTrades",
        secType:  secTypeAPIKey,
    }
    r.setParam("symbol", s.symbol)
    if s.limit != nil {
        r.setParam("limit", *s.limit)
    }

    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return
    }
    res = make([]*Trade, 0)
    err = json.Unmarshal(data, &res)
    if err != nil {
        return
    }
    return
}

// Trade define trade info
type Trade struct {
    ID            int64  `json:"id"`
    Price         string `json:"price"`
    Quantity      string `json:"qty"`
    QuoteQuantity string `json:"quoteQty"`
    Time          int64  `json:"time"`
    IsBuyerMaker  bool   `json:"isBuyerMaker"`
}

// AggTradesService list aggregate trades
type AggTradesService struct {
    c         *Client
    symbol    string
    startTime *int64
    endTime   *int64
    limit     *int
}

// Symbol set symbol
func (s *AggTradesService) Symbol(symbol string) *AggTradesService {
    s.symbol = symbol
    return s
}

// StartTime set startTime
func (s *AggTradesService) StartTime(startTime int64) *AggTradesService {
    s.startTime = &startTime
    return s
}

// EndTime set endTime
func (s *AggTradesService) EndTime(endTime int64) *AggTradesService {
    s.endTime = &endTime
    return s
}

// Limit set limit
func (s *AggTradesService) Limit(limit int) *AggTradesService {
    s.limit = &limit
    return s
}

// Do send request
func (s *AggTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*AggTrade, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/api/v3/aggTrades",
    }
    r.setParam("symbol", s.symbol)
    if s.startTime != nil {
        r.setParam("startTime", *s.startTime)
    }
    if s.endTime != nil {
        r.setParam("endTime", *s.endTime)
    }
    if s.limit != nil {
        r.setParam("limit", *s.limit)
    }
    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return []*AggTrade{}, err
    }
    res = make([]*AggTrade, 0)
    err = json.Unmarshal(data, &res)
    if err != nil {
        return []*AggTrade{}, err
    }
    return res, nil
}

// AggTrade define aggregate trade info
type AggTrade struct {
    AggTradeID   int64  `json:"a"`
    Price        string `json:"p"`
    Quantity     string `json:"q"`
    FirstTradeID int64  `json:"f"`
    LastTradeID  int64  `json:"l"`
    Timestamp    int64  `json:"T"`
    IsBuyerMaker bool   `json:"m"`
}

// RecentTradesService list recent trades
type RecentTradesService struct {
    c      *Client
    symbol string
    limit  *int
}

// Symbol set symbol
func (s *RecentTradesService) Symbol(symbol string) *RecentTradesService {
    s.symbol = symbol
    return s
}

// Limit set limit
func (s *RecentTradesService) Limit(limit int) *RecentTradesService {
    s.limit = &limit
    return s
}

// Do send request
func (s *RecentTradesService) Do(ctx context.Context, opts ...RequestOption) (res []*Trade, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/api/v1/trades",
    }
    r.setParam("symbol", s.symbol)
    if s.limit != nil {
        r.setParam("limit", *s.limit)
    }
    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return []*Trade{}, err
    }
    res = make([]*Trade, 0)
    err = json.Unmarshal(data, &res)
    if err != nil {
        return []*Trade{}, err
    }
    return res, nil
}
