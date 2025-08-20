package connector

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/haydenhigg/clover/candle"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	KRAKEN_HOST         = "api.kraken.com"
	KRAKEN_CONTENT_TYPE = "application/x-www-form-urlencoded; charset=utf-8"
)

type Kraken struct {
	Key    []byte
	Secret []byte
}

// request helpers
func (c *Kraken) buildURL(path string, query url.Values) string {
	u := url.URL{
		Scheme: "https",
		Host:   KRAKEN_HOST,
		Path:   path,
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u.String()
}

func (c *Kraken) buildSignature(path string, body url.Values) string {
	hasher := sha256.New()
	hasher.Write([]byte(body.Get("nonce")))
	hasher.Write([]byte(body.Encode()))

	h := hmac.New(sha512.New, c.Secret)
	h.Write([]byte(path))
	h.Write(hasher.Sum(nil))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// basic requests
type Payload struct {
	Query url.Values
	Body  url.Values
}

func doRequest(request *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func (c *Kraken) public(
	method,
	path string,
	payload *Payload,
) ([]byte, error) {
	if payload == nil {
		payload = new(Payload)
	}

	// set up the request
	fullPath := "/0/public" + path
	u := c.buildURL(fullPath, payload.Query)

	var bodyReader io.Reader
	if payload.Body != nil {
		bodyReader = strings.NewReader(payload.Body.Encode())
	}

	// create the *http.Request
	request, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", KRAKEN_CONTENT_TYPE)

	return doRequest(request)
}

func (c *Kraken) private(
	method,
	path string,
	payload *Payload,
) ([]byte, error) {
	if payload == nil {
		payload = new(Payload)
	}

	// set up the request
	fullPath := "/0/private" + path
	u := c.buildURL(fullPath, payload.Query)

	if payload.Body == nil {
		payload.Body = url.Values{}
	}

	payload.Body.Set("nonce", strconv.FormatInt(time.Now().UnixMilli(), 10))
	bodyReader := strings.NewReader(payload.Body.Encode())

	// create the *http.Request
	request, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", KRAKEN_CONTENT_TYPE)

	// these are added directly to request.Header to sidestep canonicalization
	request.Header["API-Key"] = []string{string(c.Key)}
	request.Header["API-Sign"] = []string{c.buildSignature(
		fullPath,
		payload.Body,
	)}

	return doRequest(request)
}

// connector functions
func (c *Kraken) FetchCandlesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*candle.Candle, error) {
	// make request
	sinceTimestamp := since.Truncate(interval).Unix() - 1
	rawResponse, err := c.public("GET", "/OHLC", &Payload{
		Query: url.Values{
			"pair":     {pair},
			"interval": {strconv.Itoa(int(interval.Minutes()))},
			"since":    {strconv.FormatInt(sinceTimestamp, 10)},
		},
	})
	if err != nil {
		return nil, err
	}

	// unmarshal raw response
	var response struct {
		Errors []string           `json:"error"`
		Result map[string][][]any `json:"result"`
	}
	json.Unmarshal(rawResponse, &response)

	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0])
	}

	rawCandles, ok := response.Result[pair]
	if !ok {
		return nil, errors.New("no candles retrieved for pair")
	}

	// process returned candles
	candles := []*candle.Candle{}

	for _, rawCandle := range rawCandles[:len(rawCandles)-1] {
		open, _ := strconv.ParseFloat(rawCandle[1].(string), 64)
		high, _ := strconv.ParseFloat(rawCandle[2].(string), 64)
		low, _ := strconv.ParseFloat(rawCandle[3].(string), 64)
		close, _ := strconv.ParseFloat(rawCandle[4].(string), 64)
		volume, _ := strconv.ParseFloat(rawCandle[6].(string), 64)

		candles = append(candles, &candle.Candle{
			Time:   time.Unix(int64(rawCandle[0].(float64)), 0),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	return candles, nil
}

func (c *Kraken) FetchBalances() (map[string]float64, error) {
	// make request
	rawResponse, err := c.private("POST", "/Balance", nil)
	if err != nil {
		return nil, err
	}

	// unmarshal raw response
	var response struct {
		Errors []string          `json:"error"`
		Result map[string]string `json:"result"`
	}
	json.Unmarshal(rawResponse, &response)

	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0])
	}

	// process returned quantities
	store := map[string]float64{}

	for c, v := range response.Result {
		balance, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}

		// don't include balances of 0
		if balance == 0 {
			continue
		}

		store[c] = balance
	}

	return store, nil
}

func (c *Kraken) PlaceMarketOrder(side, pair string, quantity float64) error {
	// make request
	rawResponse, err := c.private("POST", "/AddOrder", &Payload{
		Body: url.Values{
			"ordertype": {"market"},
			"type":      {side},
			"volume":    {strconv.FormatFloat(quantity, 'f', 8, 64)},
			"pair":      {pair},
		},
	})
	if err != nil {
		return err
	}

	// unmarshal raw response
	var response struct {
		Errors []string       `json:"error"`
		Result map[string]any `json:"result"`
	}
	json.Unmarshal(rawResponse, &response)

	if len(response.Errors) > 0 {
		return errors.New(response.Errors[0])
	}

	return nil
}
