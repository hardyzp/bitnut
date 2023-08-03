package bitnut

import (
    "bytes"
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "crypto/tls"
    "encoding/base64"
    "fmt"
    "github.com/bitly/go-simplejson"
    "github.com/hardyzp/bitnut/common"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"

    jsoniter "github.com/json-iterator/go"
)

// SideType define side type of order
type SideType string

// OrderType define order type
type OrderType string

// TimeInForceType define time in force type of order
type TimeInForceType string

// NewOrderRespType define response JSON verbosity
type NewOrderRespType string

// OrderStatusType define order status type
type OrderStatusType string

// SymbolType define symbol type
type SymbolType string

// SymbolStatusType define symbol status type
type SymbolStatusType string

// SymbolFilterType define symbol filter type
type SymbolFilterType string

// UserDataEventType define spot user data event type
type UserDataEventType string

// TransactionType define transaction type
type TransactionType string

// AccountType define the account types
type AccountType string

// Endpoints
const (
    baseAPIMainURL    = "https://api.binance.com"
    baseAPITestnetURL = "https://testnet.binance.vision"
)

// UseTestnet switch all the API endpoints from production to the testnet
var UseTestnet = false

// Redefining the standard package
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Global enums
const (
    SideTypeBuy  SideType = "BUY"
    SideTypeSell SideType = "SELL"

    OrderTypeLimit  OrderType = "LIMIT"
    OrderTypeMarket OrderType = "MARKET"

    TimeInForceTypeGTC TimeInForceType = "GTC"
    TimeInForceTypeIOC TimeInForceType = "IOC"
    TimeInForceTypeFOK TimeInForceType = "FOK"

    NewOrderRespTypeACK    NewOrderRespType = "ACK"
    NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
    NewOrderRespTypeFULL   NewOrderRespType = "FULL"

    OrderStatusTypeNew      OrderStatusType = "NEW"
    OrderStatusTypeFilled   OrderStatusType = "FILLED"
    OrderStatusTypeCanceled OrderStatusType = "CANCELED"

    SymbolTypeSpot SymbolType = "SPOT"

    timestampKey = "timestamp"
    signatureKey = "signature"

    AccountTypeSpot AccountType = "SPOT"
)

func currentTimestamp() int64 {
    return FormatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func FormatTimestamp(t time.Time) int64 {
    return t.UnixNano() / int64(time.Millisecond)
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
    j, err = simplejson.NewJson(data)
    if err != nil {
        return nil, err
    }
    return j, nil
}

// getAPIEndpoint return the base endpoint of the Rest API according the UseTestnet flag
func getAPIEndpoint() string {
    return baseAPIMainURL
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(apiKey, secretKey string) *Client {
    return &Client{
        APIKey:     apiKey,
        SecretKey:  secretKey,
        BaseURL:    getAPIEndpoint(),
        UserAgent:  "Bitnut/golang",
        HTTPClient: http.DefaultClient,
        Logger:     log.New(os.Stderr, "Bitnut-golang ", log.LstdFlags),
    }
}

// NewProxiedClient passing a proxy url
func NewProxiedClient(apiKey, secretKey, proxyUrl string) *Client {
    proxy, err := url.Parse(proxyUrl)
    if err != nil {
        log.Fatal(err)
    }
    tr := &http.Transport{
        Proxy:           http.ProxyURL(proxy),
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    return &Client{
        APIKey:    apiKey,
        SecretKey: secretKey,
        BaseURL:   getAPIEndpoint(),
        UserAgent: "Bitnut/golang",
        HTTPClient: &http.Client{
            Transport: tr,
        },
        Logger: log.New(os.Stderr, "Bitnut-golang ", log.LstdFlags),
    }
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
    APIKey     string
    SecretKey  string
    BaseURL    string
    UserAgent  string
    HTTPClient *http.Client
    Debug      bool
    Logger     *log.Logger
    TimeOffset int64
    do         doFunc
}

func (c *Client) debug(format string, v ...interface{}) {
    if c.Debug {
        c.Logger.Printf(format, v...)
    }
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
    // set request options from user
    for _, opt := range opts {
        opt(r)
    }
    err = r.validate()
    if err != nil {
        return err
    }

    fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
    if r.secType == secTypeSigned {
        r.setParam(timestampKey, currentTimestamp()-c.TimeOffset)
    }
    queryString := r.query.Encode()
    fmt.Println("queryString:", queryString)
    body := &bytes.Buffer{}
    bodyString := r.form.Encode()
    fmt.Println("rform:", r.form)
    fmt.Println("bodystring:", bodyString)
    header := http.Header{}
    if r.header != nil {
        header = r.header.Clone()
    }
    if bodyString != "" {
        header.Set("Content-Type", "application/x-www-form-urlencoded")
        body = bytes.NewBufferString(bodyString)
    }
    if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
        header.Set("BU-ACCESS-KEY", c.APIKey)
    }
    if queryString != "" {
        fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
    }
    if r.secType == secTypeSigned {
        // raw := fmt.Sprintf("%s%s", queryString, bodyString)
        raw := bodyString
        mac := hmac.New(sha256.New, []byte(c.SecretKey))
        _, err = mac.Write([]byte(raw))
        if err != nil {
            return err
        }
        signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
        header.Set("BU-ACCESS-SIGN", signature)
    }
    c.debug("full url: %s, body: %s", fullURL, bodyString)

    r.fullURL = fullURL
    r.header = header
    r.body = body
    return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
    err = c.parseRequest(r, opts...)
    if err != nil {
        return []byte{}, err
    }
    req, err := http.NewRequest(r.method, r.fullURL, r.body)
    if err != nil {
        return []byte{}, err
    }
    req = req.WithContext(ctx)
    req.Header = r.header
    c.debug("request: %#v", req)
    f := c.do
    if f == nil {
        f = c.HTTPClient.Do
    }
    res, err := f(req)
    if err != nil {
        return []byte{}, err
    }
    data, err = ioutil.ReadAll(res.Body)
    if err != nil {
        return []byte{}, err
    }
    defer func() {
        cerr := res.Body.Close()
        // Only overwrite the retured error if the original error was nil and an
        // error occurred while closing the body.
        if err == nil && cerr != nil {
            err = cerr
        }
    }()
    c.debug("response: %#v", res)
    c.debug("response body: %s", string(data))
    c.debug("response status code: %d", res.StatusCode)

    if res.StatusCode >= http.StatusBadRequest {
        apiErr := new(common.APIError)
        e := json.Unmarshal(data, apiErr)
        if e != nil {
            c.debug("failed to unmarshal json: %s", e)
        }
        return nil, apiErr
    }
    return data, nil
}

// SetApiEndpoint set api Endpoint
func (c *Client) SetApiEndpoint(url string) *Client {
    c.BaseURL = url
    return c
}

// NewServerTimeService init server time service
func (c *Client) NewServerTimeService() *ServerTimeService {
    return &ServerTimeService{c: c}
}

// NewSetServerTimeService init set server time service
func (c *Client) NewSetServerTimeService() *SetServerTimeService {
    return &SetServerTimeService{c: c}
}

// NewDepthService init depth service
func (c *Client) NewDepthService() *DepthService {
    return &DepthService{c: c}
}

// NewKlinesService init klines service
func (c *Client) NewKlinesService() *KlinesService {
    return &KlinesService{c: c}
}

// NewListSymbolTickerService init listing symbols tickers
func (c *Client) NewListSymbolTickerService() *ListSymbolTickerService {
    return &ListSymbolTickerService{c: c}
}

// NewCreateOrderService init creating order service
func (c *Client) NewCreateOrderService() *CreateOrderService {
    return &CreateOrderService{c: c}
}

// NewGetOrderService init get order service
func (c *Client) NewGetOrderService() *GetOrderService {
    return &GetOrderService{c: c}
}

// NewCancelOrderService init cancel order service
func (c *Client) NewCancelOrderService() *CancelOrderService {
    return &CancelOrderService{c: c}
}

// NewCancelOpenOrdersService init cancel open orders service
func (c *Client) NewCancelOpenOrdersService() *CancelOpenOrdersService {
    return &CancelOpenOrdersService{c: c}
}

// NewListOrdersService init listing orders service
func (c *Client) NewListOrdersService() *ListOrdersService {
    return &ListOrdersService{c: c}
}

// NewGetBalanceService  init getting account service
func (c *Client) NewGetBalanceService() *GetBalanceService {
    return &GetBalanceService{c: c}
}

// NewListTradesService init listing trades service
func (c *Client) NewListTradesService() *ListTradesService {
    return &ListTradesService{c: c}
}

// NewHistoricalTradesService init listing trades service
func (c *Client) NewHistoricalTradesService() *HistoricalTradesService {
    return &HistoricalTradesService{c: c}
}
