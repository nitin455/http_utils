package http_client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	client := NewClient("www.google.com")

	assert.NotNil(t, client, "Client should not be nil")
}

func prepareTestHttpServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
}

func TestGet(t *testing.T) {
	server := prepareTestHttpServer()
	defer server.Close()
	logger := zap.NewExample()
	client := NewClient(server.URL,
		WithPath("health"),
		WithLogger(logger))

	assert.NotNil(t, client, "Client should not be nil")

	out := make(map[string]interface{})
	_, err := client.Get(&out)

	assert.Nil(t, err, "Error should be nil for Get call.")
	assert.NotNil(t, out, "Output for Get call should not be nil")
}

func TestGet_Negative(t *testing.T) {
	logger := zap.NewExample()
	client := NewClient("http://localhost:10080",
		WithPath("health"),
		WithLogger(logger))

	assert.NotNil(t, client, "Client should not be nil")

	out := make(map[string]interface{})
	_, err := client.Get(&out)

	assert.NotNil(t, err, "Error should not be nil for Get call.")
	assert.NotNil(t, out, "Output for Get call should not be nil")
}

func payload() map[string]interface{} {
	return map[string]interface{}{
		"name":      "Test",
		"tenant_id": "f3efc0f5-d502-49ac-b748-71906e92364d",
		"device_id": "f3efc0f5-d502-49ac-b748-71906e92364e",
		"checksum":  "AAASDRRX",
		"size":      1024,
	}
}
func TestPost(t *testing.T) {
	server := prepareTestHttpServer()
	defer server.Close()
	logger := zap.NewExample()
	client := NewClient(server.URL,
		WithPath("/api/v1/post"),
		WithLogger(logger),
		WithPayload(payload()))

	assert.NotNil(t, client, "Client should not be nil")

	out := make(map[string]interface{})
	_, err := client.Post(&out)

	assert.Nil(t, err, "Error should be nil for Post call.")
	assert.NotNil(t, out, "Output for Post call should not be nil")
}

func TestDelete(t *testing.T) {
	server := prepareTestHttpServer()
	defer server.Close()
	logger := zap.NewExample()
	client := NewClient(server.URL,
		WithPath("/api/v1/delete"),
		WithLogger(logger),
		WithPayload(payload()))

	assert.NotNil(t, client, "Client should not be nil")

	_, err := client.Delete(nil)
	assert.Nil(t, err, "Error should be nil for delete call.")
}
