package client

import (
	"http_utils/pkg/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Client - Holds client details.
type Client struct {
	cnt *http.Client
}

func NewClient(cfg *config.ClientConfig) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	cnt := &http.Client{}

	if cfg.HTTPRequestProxyUrl != nil {
		log.Info("Setting Proxy URL manually")
		transport.Proxy = http.ProxyURL(cfg.HTTPRequestProxyUrl)
	}

	if cfg.HTTPTLSHandshakeTimeout.Seconds() > 0 {
		log.Info("Overriding TLS Handshake timeout")
		transport.TLSHandshakeTimeout = cfg.HTTPTLSHandshakeTimeout
	}

	if cfg.HTTPRequestTimeout.Seconds() > 0 {
		log.Info("Overriding request timeout")
		cnt.Timeout = cfg.HTTPRequestTimeout
	}

	cnt.Transport = transport

	return &Client{
		cnt: cnt,
	}
}
