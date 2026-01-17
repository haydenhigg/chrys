package driver

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/haydenhigg/chrys/frame"
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

type KrakenDriver struct {
	Key    []byte
	Secret []byte
}

func NewKraken(key, secret string) (*KrakenDriver, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	d := &KrakenDriver{
		Key:    []byte(key),
		Secret: decodedSecret,
	}

	return d, nil
}

// request helpers
func (d *KrakenDriver) buildURL(path string, query url.Values) string {
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

func (d *KrakenDriver) buildSignature(path string, body url.Values) string {
	hasher := sha256.New()
	hasher.Write([]byte(body.Get("nonce")))
	hasher.Write([]byte(body.Encode()))

	h := hmac.New(sha512.New, d.Secret)
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

func (d *KrakenDriver) public(
	method,
	path string,
	payload *Payload,
) ([]byte, error) {
	if payload == nil {
		payload = new(Payload)
	}

	// set up the request
	fullPath := "/0/public" + path
	u := d.buildURL(fullPath, payload.Query)

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

func (d *KrakenDriver) private(
	method,
	path string,
	payload *Payload,
) ([]byte, error) {
	if payload == nil {
		payload = new(Payload)
	}

	// set up the request
	fullPath := "/0/private" + path
	u := d.buildURL(fullPath, payload.Query)

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
	request.Header["API-Key"] = []string{string(d.Key)}
	request.Header["API-Sign"] = []string{d.buildSignature(
		fullPath,
		payload.Body,
	)}

	return doRequest(request)
}

// driver functions
func (d *KrakenDriver) FetchFramesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*frame.Frame, error) {
	// make request
	rawResponse, err := d.public("GET", "/OHLC", &Payload{
		Query: url.Values{
			"pair":     {pair},
			"interval": {strconv.Itoa(int(interval.Minutes()))},
			"since":    {strconv.FormatInt(since.Unix() - 1, 10)},
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

	rawFrames, ok := response.Result[pair]
	if !ok {
		return nil, errors.New("no frames retrieved for pair")
	}

	// process returned frames
	frames := []*frame.Frame{}

	for _, rawFrame := range rawFrames[:len(rawFrames)-1] {
		open, _ := strconv.ParseFloat(rawFrame[1].(string), 64)
		high, _ := strconv.ParseFloat(rawFrame[2].(string), 64)
		low, _ := strconv.ParseFloat(rawFrame[3].(string), 64)
		close, _ := strconv.ParseFloat(rawFrame[4].(string), 64)
		volume, _ := strconv.ParseFloat(rawFrame[6].(string), 64)

		frames = append(frames, &frame.Frame{
			Time:   time.Unix(int64(rawFrame[0].(float64)), 0),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	return frames, nil
}

func (d *KrakenDriver) FetchBalances() (map[string]float64, error) {
	// make request
	rawResponse, err := d.private("POST", "/Balance", nil)
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

	for d, v := range response.Result {
		balance, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}

		// don't include balances of 0
		if balance == 0 {
			continue
		}

		store[d] = balance
	}

	return store, nil
}

func (d *KrakenDriver) MarketOrder(side, pair string, quantity float64) error {
	// make request
	rawResponse, err := d.private("POST", "/AddOrder", &Payload{
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
