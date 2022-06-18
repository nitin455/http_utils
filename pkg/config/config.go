package config

import (
	"http_utils/pkg/errors"
	"net/url"
	"time"
)

// ClientConfig - Holds client configurations.
type ClientConfig struct {
	BaseUrl                 *url.URL
	HTTPRequestTimeout      time.Duration
	HTTPTLSHandshakeTimeout time.Duration
	HTTPRequestProxyUrl     *url.URL
	RetryOnFailure          bool
	errors                  error
}

func (c *ClientConfig) AppendError(err error) {
	c.errors = errors.WrapError{
		Err:  err,
		Next: c.errors,
	}
}

func (c *ClientConfig) Validate() error {
	return c.errors
}
