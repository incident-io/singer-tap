package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	incident "github.com/incident-io/sdk-go"
	"github.com/pkg/errors"
)

// New returns a client for the incident.io API that retries transient failures
// and returns an error on any non-2xx response.
func New(ctx context.Context, apiKey, apiEndpoint, version string, opts ...incident.ClientOption) (*incident.ClientWithResponses, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Logger = nil

	base := retryClient.StandardClient()

	// The SDK signals a failed request with a nil JSON200 rather than an error, and
	// the streams dereference JSON200 directly, so turn it into an error here.
	base.Transport = Wrap(cleanhttp.DefaultTransport(), func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		resp, err := next.RoundTrip(req)
		if err == nil && resp.StatusCode > 299 {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("status %d: no response body", resp.StatusCode)
			}

			return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
		}

		return resp, err
	})

	clientOpts := append([]incident.ClientOption{
		incident.WithBaseURL(apiEndpoint),
		incident.WithHTTPClient(base),
		// Add a user-agent so we can tell which version these requests came from.
		incident.WithUserAgent(fmt.Sprintf("tap-incident/%s", version)),
	}, opts...)

	cl, err := incident.New(apiKey, clientOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "creating client")
	}

	return cl, nil
}

// RoundTripperFunc wraps a function to implement the RoundTripper interface, allowing
// easy wrapping of existing round-trippers.
type RoundTripperFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Wrap allows easy wrapping of an existing RoundTripper with a function that can
// optionally call the original, or do its own thing.
func Wrap(next http.RoundTripper, apply func(req *http.Request, next http.RoundTripper) (*http.Response, error)) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return apply(req, next)
	})
}
