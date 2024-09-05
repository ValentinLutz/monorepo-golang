package testutil

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	client  *http.Client
	baseUrl string
}

func NewHttpClient(baseUrl string) *HttpClient {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: httpTransport,
	}

	return &HttpClient{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (hc *HttpClient) Get(url string, correlationId uuid.UUID, queryParams map[string]string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("GET", hc.baseUrl+url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range queryParams {
		request.URL.Query().Add(key, value)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Correlation-Id", correlationId.String())
	request.SetBasicAuth("test", "test")

	return hc.client.Do(request)
}

func (hc *HttpClient) Post(url string, correlationId uuid.UUID, queryParams map[string]string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", hc.baseUrl+url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range queryParams {
		request.URL.Query().Add(key, value)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Correlation-Id", correlationId.String())
	request.SetBasicAuth("test", "test")

	return hc.client.Do(request)
}

func ParseResponseBody(response *http.Response) (map[string]any, error) {
	var responseBody map[string]any
	err := json.NewDecoder(response.Body).Decode(&responseBody)

	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return responseBody, nil
}

func ParseResponseBodyList(response *http.Response) ([]any, error) {
	var responseBody []any
	err := json.NewDecoder(response.Body).Decode(&responseBody)

	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return responseBody, nil
}
