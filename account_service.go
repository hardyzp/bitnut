package bitnut

import (
    "context"
    "fmt"
    "net/http"
)

// GetBalanceService get account balance
type GetBalanceService struct {
    c    *Client
    coin string
}

func (s *GetBalanceService) SetCoin(coin string) *GetBalanceService {
    s.coin = coin
    return s
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (*Balance, error) {
    r := &request{
        method:   http.MethodPost,
        endpoint: "/v1/asset/balance",
        secType:  secTypeSigned,
    }
    m := params{}
    m["coin"] = s.coin
    r.setFormParams(m)
    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return &Balance{}, err
    }
    ret := new(BalanceResponse)
    err = json.Unmarshal(data, &ret)
    if err != nil {
        fmt.Println(err)
        return &Balance{}, err
    }
    fmt.Println("ret", ret)
    return &ret.Data, nil
}

// BalanceResponse define user balance of your account
type BalanceResponse struct {
    Code int     `json:"code"`
    Msg  string  `json:"msg"`
    Data Balance `json:"data"`
}

// Balance define user balance of your account
type Balance struct {
    Coin   string `json:"coin"`
    Free   string `json:"free"`
    Freeze string `json:"freeze"`
}
