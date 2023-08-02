package bitnut

import (
    "context"
    "net/http"
)

// GetBalanceService get account balance
type GetBalanceService struct {
    c    *Client
    coin string
}

func (s *GetBalanceService) setCoin(coin string) {
    s.coin = coin
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (res *Balance, err error) {
    r := &request{
        method:   http.MethodPost,
        endpoint: "/v1/asset/balance",
        secType:  secTypeSigned,
    }

    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return &Balance{}, err
    }
    res = new(Balance)
    err = json.Unmarshal(data, &res)
    if err != nil {
        return &Balance{}, err
    }
    return res, nil
}

// Balance define user balance of your account
type Balance struct {
    Coin   string `json:"coin"`
    Free   string `json:"free"`
    Freeze string `json:"freeze"`
}
