# A lightweight wrapper over Go’s HTTP client (net/http)
It is a wrapper on net/http package to simplify usage of all HTTP methods. It helps users to unmarshal their individual responses into their own structures easily. It exposes multiple configuration for http client while initialization.

This wrapper is supposed to provide a clean interface for setting various parameters in your request (query params, body, content type etc.) and handle various redundant tasks like the marshalling/unmarshalling of request/response, closing response body after it has been read etc.
