package bitnut

import (
    "context"
    "net/http"
)

// CreateOrderService create order
type CreateOrderService struct {
    c                *Client
    symbol           string
    side             SideType
    orderType        OrderType
    quantity         *string
    quoteOrderQty    *string
    price            *string
    newClientOrderID *string
}

// Symbol set symbol
func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
    s.symbol = symbol
    return s
}

// Side set side
func (s *CreateOrderService) Side(side SideType) *CreateOrderService {
    s.side = side
    return s
}

// Type set type
func (s *CreateOrderService) Type(orderType OrderType) *CreateOrderService {
    s.orderType = orderType
    return s
}

// Quantity set quantity
func (s *CreateOrderService) Quantity(quantity string) *CreateOrderService {
    s.quantity = &quantity
    return s
}

// QuoteOrderQty set quoteOrderQty
func (s *CreateOrderService) QuoteOrderQty(quoteOrderQty string) *CreateOrderService {
    s.quoteOrderQty = &quoteOrderQty
    return s
}

// Price set price
func (s *CreateOrderService) Price(price string) *CreateOrderService {
    s.price = &price
    return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateOrderService) NewClientOrderID(newClientOrderID string) *CreateOrderService {
    s.newClientOrderID = &newClientOrderID
    return s
}

func (s *CreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, err error) {
    r := &request{
        method:   http.MethodPost,
        endpoint: endpoint,
        secType:  secTypeSigned,
    }
    m := params{
        "symbol": s.symbol,
        "side":   s.side,
        "type":   s.orderType,
    }
    if s.quantity != nil {
        m["quantity"] = *s.quantity
    }
    if s.quoteOrderQty != nil {
        m["quoteOrderQty"] = *s.quoteOrderQty
    }

    if s.price != nil {
        m["price"] = *s.price
    }
    if s.newClientOrderID != nil {
        m["newClientOrderId"] = *s.newClientOrderID
    }

    r.setFormParams(m)
    data, err = s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return []byte{}, err
    }
    return data, nil
}

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateOrderResponse, err error) {
    data, err := s.createOrder(ctx, "/v1/trade/order", opts...)
    if err != nil {
        return nil, err
    }
    res = new(CreateOrderResponse)
    err = json.Unmarshal(data, res)
    if err != nil {
        return nil, err
    }
    return res, nil
}

// CreateOrderResponse define create order response
type CreateOrderResponse struct {
    Symbol           string `json:"symbol"`
    OrderID          int64  `json:"orderId"`
    ClientOrderID    string `json:"clientOrderId"`
    TransactTime     int64  `json:"transactTime"`
    Price            string `json:"price"`
    OrigQuantity     string `json:"origQty"`
    ExecutedQuantity string `json:"executedQty"`

    Status OrderStatusType `json:"status"`
    Type   OrderType       `json:"type"`
    Side   SideType        `json:"side"`
}

// GetOrderService get an order
type GetOrderService struct {
    c                 *Client
    symbol            string
    orderID           *int64
    origClientOrderID *string
}

// Symbol set symbol
func (s *GetOrderService) Symbol(symbol string) *GetOrderService {
    s.symbol = symbol
    return s
}

// OrderID set orderID
func (s *GetOrderService) OrderID(orderID int64) *GetOrderService {
    s.orderID = &orderID
    return s
}

// OrigClientOrderID set origClientOrderID
func (s *GetOrderService) OrigClientOrderID(origClientOrderID string) *GetOrderService {
    s.origClientOrderID = &origClientOrderID
    return s
}

// Do send request
func (s *GetOrderService) Do(ctx context.Context, opts ...RequestOption) (res *Order, err error) {
    r := &request{
        method:   http.MethodPost,
        endpoint: "/api/v3/order",
        secType:  secTypeSigned,
    }
    r.setParam("symbol", s.symbol)
    if s.orderID != nil {
        r.setParam("orderId", *s.orderID)
    }
    if s.origClientOrderID != nil {
        r.setParam("origClientOrderId", *s.origClientOrderID)
    }
    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return nil, err
    }
    res = new(Order)
    err = json.Unmarshal(data, res)
    if err != nil {
        return nil, err
    }
    return res, nil
}

// Order define order info
type Order struct {
    Symbol                   string          `json:"symbol"`
    OrderID                  int64           `json:"orderId"`
    OrderListId              int64           `json:"orderListId"`
    ClientOrderID            string          `json:"clientOrderId"`
    Price                    string          `json:"price"`
    OrigQuantity             string          `json:"origQty"`
    ExecutedQuantity         string          `json:"executedQty"`
    CummulativeQuoteQuantity string          `json:"cummulativeQuoteQty"`
    Status                   OrderStatusType `json:"status"`
    TimeInForce              TimeInForceType `json:"timeInForce"`
    Type                     OrderType       `json:"type"`
    Side                     SideType        `json:"side"`

    Time       int64 `json:"time"`
    UpdateTime int64 `json:"updateTime"`
}

// ListOrdersService all account orders; active, canceled, or filled
type ListOrdersService struct {
    c         *Client
    symbol    string
    orderID   *int64
    startTime *int64
    endTime   *int64
    limit     *int
}

// Symbol set symbol
func (s *ListOrdersService) Symbol(symbol string) *ListOrdersService {
    s.symbol = symbol
    return s
}

// OrderID set orderID
func (s *ListOrdersService) OrderID(orderID int64) *ListOrdersService {
    s.orderID = &orderID
    return s
}

// StartTime set starttime
func (s *ListOrdersService) StartTime(startTime int64) *ListOrdersService {
    s.startTime = &startTime
    return s
}

// EndTime set endtime
func (s *ListOrdersService) EndTime(endTime int64) *ListOrdersService {
    s.endTime = &endTime
    return s
}

// Limit set limit
func (s *ListOrdersService) Limit(limit int) *ListOrdersService {
    s.limit = &limit
    return s
}

// Do send request
func (s *ListOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*Order, err error) {
    r := &request{
        method:   http.MethodGet,
        endpoint: "/api/v3/allOrders",
        secType:  secTypeSigned,
    }
    r.setParam("symbol", s.symbol)
    if s.orderID != nil {
        r.setParam("orderId", *s.orderID)
    }
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
        return []*Order{}, err
    }
    res = make([]*Order, 0)
    err = json.Unmarshal(data, &res)
    if err != nil {
        return []*Order{}, err
    }
    return res, nil
}

// CancelOrderService cancel an order
type CancelOrderService struct {
    c                 *Client
    symbol            string
    orderID           *int64
    origClientOrderID *string
}

// Symbol set symbol
func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
    s.symbol = symbol
    return s
}

// OrderID set orderID
func (s *CancelOrderService) OrderID(orderID int64) *CancelOrderService {
    s.orderID = &orderID
    return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelOrderService) OrigClientOrderID(origClientOrderID string) *CancelOrderService {
    s.origClientOrderID = &origClientOrderID
    return s
}

// Do send request
func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOrderResponse, err error) {
    r := &request{
        method:   http.MethodDelete,
        endpoint: "/api/v3/order",
        secType:  secTypeSigned,
    }
    r.setFormParam("symbol", s.symbol)
    if s.orderID != nil {
        r.setFormParam("orderId", *s.orderID)
    }
    if s.origClientOrderID != nil {
        r.setFormParam("origClientOrderId", *s.origClientOrderID)
    }

    data, err := s.c.callAPI(ctx, r, opts...)
    if err != nil {
        return nil, err
    }
    res = new(CancelOrderResponse)
    err = json.Unmarshal(data, res)
    if err != nil {
        return nil, err
    }
    return res, nil
}

// CancelOpenOrdersService cancel all active orders on a symbol.
type CancelOpenOrdersService struct {
    c      *Client
    symbol string
}

// Symbol set symbol
func (s *CancelOpenOrdersService) Symbol(symbol string) *CancelOpenOrdersService {
    s.symbol = symbol
    return s
}

// CancelOrderResponse may be returned included in a CancelOpenOrdersResponse.
type CancelOrderResponse struct {
    Symbol            string          `json:"symbol"`
    OrigClientOrderID string          `json:"origClientOrderId"`
    OrderID           int64           `json:"orderId"`
    OrderListID       int64           `json:"orderListId"`
    ClientOrderID     string          `json:"clientOrderId"`
    TransactTime      int64           `json:"transactTime"`
    Price             string          `json:"price"`
    OrigQuantity      string          `json:"origQty"`
    ExecutedQuantity  string          `json:"executedQty"`
    Status            OrderStatusType `json:"status"`
    TimeInForce       TimeInForceType `json:"timeInForce"`
    Type              OrderType       `json:"type"`
    Side              SideType        `json:"side"`
}
