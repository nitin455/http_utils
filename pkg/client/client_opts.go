package http_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	thirty = 30
)

// ClientOpt - Custom function type.
type ClientOpt func(*Client)

// WithHeader - Adds new header to http request.
func WithHeader(key, value string) ClientOpt {
	return func(cln *Client) {
		cln.headers[key] = []string{value}
	}
}

// WithQueryParam - Adds query param to http request.
func WithQueryParam(key, value string) ClientOpt {
	return func(cln *Client) {
		cln.query.Add(key, value)
	}
}

// WithPath - Adds url path to http request.
func WithPath(path string) ClientOpt {
	return func(cln *Client) {
		cln.path = path
	}
}

// WithTransport - Adds transport to http client.
func WithTransport(transport *http.Transport) ClientOpt {
	return func(cln *Client) {
		if transport == nil {
			transport = http.DefaultTransport.(*http.Transport)
		}

		cln.client.Transport = transport
	}
}

// WithTimeout - Adds timeout to http client.
func WithTimeout(timeout time.Duration) ClientOpt {
	return func(cln *Client) {
		if timeout <= 0 {
			timeout = thirty * time.Second
		}

		cln.client.Timeout = timeout
	}
}

// WithPayload - Adds request payload to http client wrapper.
func WithPayload(payload interface{}) ClientOpt {
	return func(cln *Client) {
		switch data := payload.(type) {
		case string:
			cln.data = []byte(data)
		case []byte:
			cln.data = data
		default:
			// convert to json string.
			bytes, err := json.Marshal(payload)
			if err != nil {
				fmt.Printf("Failed to convert payload into bytes. error - %v", err)
			}

			cln.data = bytes
		}
	}
}

// WithLogger - Used for logging details while making client calls.
func WithLogger(logger *zap.Logger) ClientOpt {
	return func(cln *Client) {
		cln.logger = logger
	}
}

// WithUseProto - Sets if client should use proto for unmarshalling response or not.
func WithUseProto(useProto bool) ClientOpt {
	return func(cln *Client) {
		cln.useProto = useProto
	}
}

// WithSignedUrl - Sets if client should use signed url for making http call or not.
func WithSignedUrl(useSignedUrl bool) ClientOpt {
	return func(cln *Client) {
		cln.useSignedUrl = useSignedUrl
	}
}
