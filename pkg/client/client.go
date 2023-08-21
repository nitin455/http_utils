package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"http_utils/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

// Client - Holds http client to make http calls.
type Client struct {
	data         []byte
	headers      http.Header
	client       *http.Client
	path         string
	query        url.Values
	url          *url.URL
	logger       *zap.Logger
	useProto     bool
	signedUrl    string
	useSignedUrl bool
}

// NewClient returns new client
func NewClient(targetUrl string, options ...ClientOpt) (client *Client) {
	client = &Client{
		headers:   make(http.Header),
		url:       prepareURL(targetUrl),
		query:     make(url.Values, 0),
		client:    &http.Client{},
		signedUrl: targetUrl,
	}

	// apply options
	for _, option := range options {
		option(client)
	}

	// check if logger is updated. if not, use default logger.
	if client.logger == nil {
		logger.InitLogger(logger.DefaultLogLevel)
		client.logger = logger.Logger
	}

	// Check if we have urlencoded header present. If yes, add content length.
	if client.headers.Get("Content-Type") == "application/x-www-form-urlencoded" {
		client.headers.Add("Content-Length", strconv.Itoa(len(client.query.Encode())))
	}

	// Add default headers.
	if client.headers.Get("Content-Type") == "" && !client.useSignedUrl {
		client.headers.Add("Content-Type", "application/json")
	}

	client.headers.Add("cache-control", "no-cache")

	return
}

// prepareURL - Returns an instance of *url.URL.
func prepareURL(baseURL string) *url.URL {
	u := new(url.URL)
	u.Scheme = "http"

	if strings.HasPrefix(baseURL, "https") {
		u.Scheme = "https"
		baseURL = strings.Replace(baseURL, "https://", "", 1)
	} else {
		baseURL = strings.Replace(baseURL, "http://", "", 1)
	}

	u.Host = baseURL

	return u
}

// Do perform actual request on client
// Also does error checking, so everything thats not 2xx family will be returned as http error
func (c *Client) Do(r *http.Request, response interface{}) (string, error) {
	c.logger.Debug("New http request",
		zap.String("URL", r.URL.String()), zap.Binary("Payload", c.data),
		zap.String("Method", r.Method))

	resp, err := c.client.Do(r)
	if err != nil {
		c.logger.Error("Failed to execute http request",
			zap.String("URL", r.URL.String()), zap.Binary("Payload", c.data),
			zap.String("Method", r.Method), zap.Error(err))
		return "", err
	}

	body := ""
	if resp.Body != nil {
		defer resp.Body.Close()
		byteResp, err := io.ReadAll(resp.Body)
		if err != nil {
			c.logger.Error("Failed to unmarshal response",
				zap.String("URL", r.URL.String()), zap.Binary("Payload", c.data),
				zap.String("Method", r.Method), zap.Error(err))
			return "", err
		}
		body = string(byteResp)
	}

	if resp.StatusCode/100 != 2 {
		c.logger.Error("Different status received than 2xx",
			zap.String("URL", r.URL.String()), zap.String("Method", r.Method),
			zap.Int("Status", resp.StatusCode), zap.String("Response", body))
		return body, newHttpError(resp.StatusCode, resp.Status)
	}

	c.logger.Debug("Http response",
		zap.String("URL", r.URL.String()), zap.String("Method", r.Method),
		zap.Int("Status", resp.StatusCode), zap.String("Response", body))

	if response != nil && body != "" {
		if c.useProto {
			unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
			err := unmarshaler.Unmarshal(strings.NewReader(body), response.(proto.Message))
			if err != nil {
				c.logger.Error("Failed to decode proto response",
					zap.String("URL", r.URL.String()), zap.String("Method", r.Method),
					zap.String("Response", body))
				return "", err
			}
		}

		if err := json.NewDecoder(strings.NewReader(body)).Decode(&response); err != nil {
			c.logger.Error("Failed to decode response",
				zap.String("URL", r.URL.String()), zap.String("Method", r.Method),
				zap.String("Response", body))
			return "", err
		}
	}

	return "", nil
}

// Get does http request and unmarshalls response into `into`.
// if non 2xx response is returned, it returns HttpError
func (c *Client) Get(response interface{}) (string, error) {
	req, err := c.NewRequest(http.MethodGet)
	if err != nil {
		c.logger.Error("Cant create new http GET request.", zap.Error(err))
		return "", err
	}

	return c.Do(req, response)
}

// Post does http request with `data` from client obj and unmarshalls response into `into`.
// if non 2xx response is returned, it returns HttpError.
func (c *Client) Post(into interface{}) (string, error) {
	req, err := c.NewRequest(http.MethodPost)
	if err != nil {
		c.logger.Error("Cant create new http POST request.", zap.Error(err))
		return "", err
	}

	return c.Do(req, into)
}

// Put does http request with `data` from client and unmarshalls response into `into`.
// if non 2xx response is returned, it returns HttpError
func (c *Client) Put(into interface{}) (string, error) {
	req, err := c.NewRequest(http.MethodPut)
	if err != nil {
		c.logger.Error("Cant create new http PUT request.", zap.Error(err))
		return "", err
	}

	return c.Do(req, into)
}

// Delete does http request with `data` from client obj and unmarshalls response into `into`.
// if non 2xx response is returned, it returns HttpError.
func (c *Client) Delete(into interface{}) (string, error) {
	req, err := c.NewRequest(http.MethodDelete)
	if err != nil {
		c.logger.Error("Cant create new http DELETE request.", zap.Error(err))
		return "", err
	}

	return c.Do(req, into)
}

// NewRequest creates http request and returns it
func (c *Client) NewRequest(method string) (*http.Request, error) {
	finalUrl := c.signedUrl
	if !c.useSignedUrl {
		u := c.url.ResolveReference(&url.URL{Path: c.path})
		u.RawQuery = c.query.Encode()
		finalUrl = u.String()
	}

	// when body is provided
	var reqBody io.Reader
	if c.data != nil {
		reqBody = bytes.NewBuffer(c.data)
	}

	// prepare request
	req, err := http.NewRequestWithContext(context.Background(), method, finalUrl, reqBody)
	if err != nil {
		return nil, err
	}

	// add headers
	for k, v := range c.headers {
		req.Header[k] = v
	}

	return req, nil
}

func newHttpError(statusCode int, status string) error {
	return fmt.Errorf("%d: %s", statusCode, status)
}
