package configopts

import (
	"http_utils/pkg/config"
	"http_utils/pkg/errors"
	"time"

	"go.uber.org/zap"

	"net/url"
	"strings"
)

type ClientConfigOpt func(*config.ClientConfig)

func WithBaseUrl(baseUrl string) ClientConfigOpt {
	return func(config *config.ClientConfig) {
		if len(baseUrl) == 0 {
			config.AppendError(errors.ErrNoBaseUrl)
			return
		}

		u, err := NormalizeBaseUrl(baseUrl)
		if err != nil {
			config.AppendError(err)
			return
		}

		config.BaseUrl = &u
	}
}

func NormalizeBaseUrl(baseUrl string) (url.URL, error) {
	// We need to support DNS names, but url.Parse will fail if address doesn't contain any protocol.
	if !strings.HasPrefix(baseUrl, "https://") && !strings.HasPrefix(baseUrl, "http://") {
		// Prefer https over http.
		baseUrl = "https://" + baseUrl
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		return url.URL{}, err
	}

	if u.Port() == "" {
		if u.Scheme == "http" {
			u.Host += ":80"
		} else {
			u.Host += ":443"
		}
	}

	return *u, nil
}

func WithHTTPRequestTimeout(timeout time.Duration) ClientConfigOpt {
	return func(clientConfig *config.ClientConfig) {
		clientConfig.HTTPRequestTimeout = timeout
	}
}

func WithHTTPTLSHandshakeTimeout(timeout time.Duration) ClientConfigOpt {
	return func(clientConfig *config.ClientConfig) {
		clientConfig.HTTPTLSHandshakeTimeout = timeout
	}
}

func WithHTTPRequestProxyUrl(proxyUrl *url.URL) ClientConfigOpt {
	return func(clientConfig *config.ClientConfig) {
		clientConfig.HTTPRequestProxyUrl = proxyUrl
	}
}

func WithRetry() ClientConfigOpt {
	return func(clientConfig *config.ClientConfig) {
		clientConfig.RetryOnFailure = true
	}
}

func WithMaxRetryAttempts(attempts int) ClientConfigOpt {
	return func(clientConfig *config.ClientConfig) {
		clientConfig.MaxRetryAttempts = attempts
	}
}

func WithLogger(logger *zap.Logger) ClientConfigOpt {
	return func(clientConfig *config.ClientConfig) {
		clientConfig.Logger = logger
	}
}
