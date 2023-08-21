package config

import (
	"http_utils/pkg/errors"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// ClientConfig - Holds client configurations.
type ClientConfig struct {
	BaseUrl                 *url.URL
	HTTPRequestTimeout      time.Duration
	HTTPTLSHandshakeTimeout time.Duration
	HTTPRequestProxyUrl     *url.URL
	RetryOnFailure          bool
	MaxRetryAttempts        int
	errors                  error
	Logger                  *zap.Logger
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
