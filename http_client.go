package httputils

import (
	"context"
	"http_utils/pkg/client"
	"http_utils/pkg/config"
	"http_utils/pkg/configopts"
)

var httpclient *client.Client

func NewConfig(opts ...configopts.ClientConfigOpt) (*config.ClientConfig, error) {
	cfg := new(config.ClientConfig)

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, cfg.Validate()
}

// InitClient takes up configuration to build http client.
func InitClient(cfg *config.ClientConfig) error {
	httpclient = client.NewClient(cfg)

	return nil
}

func Get(ctx context.Context, url string, into interface{}) error {

	return nil
}
