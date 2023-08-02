package bitnut

import (
    "context"
    "net/http"

    "github.com/hardyzp/bitnut/common"
)

type ListSymbolTickerService struct {
    c       *Client
    symbol  *string
    symbols []string
}

type SymbolTicker struct {
    Symbol             string `json:"symbol"`
    PriceChange        string `json:"priceChange"`
    PriceChangePercent string `json:"priceChangePercent"`
    HighPrice          string `json:"highPrice"`
    LowPrice           string `json:"lowPrice"`
    LastPrice          string `json:"lastPrice"`
    Volume             string `json:"volume"`
    QuoteVolume        string `json:"quoteVolume"`
}

func (s *ListSymbolTickerService) Symbol(symbol string) *ListSymbolTickerService {
    s.symbol = &symbol
    return s
}

func (s *ListSymbolTickerService) Symbols(symbols []string) *ListSymbolTickerService {
    s.symbols = symbols
    return s
}

func (s *ListSymbolTickerService) Do(ctx context.Context, opts ...RequestOption) (res []*SymbolTicker, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/v1/tick/24info",
    }
    if s.symbol != nil {
        r.setParam("symbol", *s.symbol)
    } else if s.symbols != nil {
        s, _ := json.Marshal(s.symbols)
        r.setParam("symbols", string(s))
    }

    data, err := s.c.callAPI(ctx, r, opts...)
    data = common.ToJSONList(data)
    if err != nil {
        return []*SymbolTicker{}, err
    }
    res = make([]*SymbolTicker, 0)
    err = json.Unmarshal(data, &res)
    if err != nil {
        return []*SymbolTicker{}, err
    }
    return res, nil
}
