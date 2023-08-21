package main

import (
	http_client "http_utils/pkg/client"
	"http_utils/pkg/logger"
	"net/http"
	"net/http/httptest"

	"go.uber.org/zap"
)

func runTestHttpServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
}

func main() {
	server := runTestHttpServer()
	defer server.Close()

	logger.InitLogger(logger.DefaultLogLevel)

	// create new http client.
	client := http_client.NewClient(server.URL,
		http_client.WithPath("/api/v1/post"),
		http_client.WithHeader("Content-Type", "application/json"),
		http_client.WithPayload([]byte("{\"field1\":\"value1\",\"field2\":\"value2\"}")),
		http_client.WithLogger(logger.Logger))

	resp := make(map[string]interface{})
	if _, err := client.Post(&resp); err != nil {
		logger.Logger.Error(
			"error received while making post api.",
			zap.Error(err),
			zap.String("url", server.URL),
		)
	}
}
