package tracing

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// HTTPClient wraps an http.Client with tracing instrumentation.
type HTTPClient struct {
	TracerProvider trace.TracerProvider
	Client         *http.Client
}

func NewHTTPClient(tp trace.TracerProvider) *HTTPClient {
	return &HTTPClient{
		TracerProvider: tp,
		Client: &http.Client{
			Transport: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithTracerProvider(tp),
			),
		},
	}
}

// GetJSON executes HTTP GET against specified url and tried to parse
// the response into out object.
func (c *HTTPClient) GetJSON(ctx context.Context, endpoint string, url string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}
