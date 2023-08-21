# A lightweight wrapper over Go’s HTTP client (net/http)
It is a wrapper on net/http package to simplify usage of all HTTP methods. It helps users to unmarshal their individual responses into their own structures easily. It exposes multiple configuration for http client while initialization.

This wrapper is supposed to provide a clean interface for setting various parameters in your request (query params, body, content type etc.) and handle various redundant tasks like the marshalling/unmarshalling of request/response, closing response body after it has been read etc.

## Usage:
```
package main

import (
	http_client "http_utils/pkg/client"
	"http_utils/pkg/logger"

	"go.uber.org/zap"
)

func main() {
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
```
## Sample logs:
```
$ go run examples/main.go
{"level":"debug","ts":1692619273.4722042,"caller":"client/client.go:90","msg":"New http request","URL":"http://127.0.0.1:58100/api/v1/post","Payload":"eyJmaWVsZDEiOiJ2YWx1ZTEiLCJmaWVsZDIiOiJ2YWx1ZTIifQ==","Method":"POST"}     
{"level":"debug","ts":1692619273.4743419,"caller":"client/client.go:122","msg":"Http response","URL":"http://127.0.0.1:58100/api/v1/post","Method":"POST","Status":200,"Response":"{\"status\":\"success\"}"}
```