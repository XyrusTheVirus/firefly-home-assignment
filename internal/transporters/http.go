package transporters

import (
	"context"
	"firefly-home-assigment/configs"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// HTTP client with a timeout
var client = &http.Client{
	Timeout: time.Duration(configs.EnvInt("HTTP_LIMIT_TIMEOUT", "30")) * time.Second,
}

// Limiter to control the rate of HTTP requests
var limiter = rate.NewLimiter(
	rate.Every(time.Duration(configs.EnvInt("RATE_LIMIT_MILLISECONDS", "200"))*time.Millisecond),
	configs.EnvInt("RATE_LIMIT_BURST", "5"),
) // 5 requests per second

type Http struct {
	Body    io.Reader
	Headers map[string]string
	Method  string
	Url     string
}

// NewHttp creates a new Http transporter with the specified parameters
func NewHttp(method, url string, body io.Reader, headers map[string]string) Http {
	return Http{
		Method:  method,
		Url:     url,
		Body:    body,
		Headers: headers,
	}
}

// Transport performs the HTTP request and returns the response or an error
func (h Http) Transport() (interface{}, error) {
	var err error
	var resp *http.Response
	var req *http.Request

	req, err = http.NewRequest(h.Method, h.Url, h.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request for URL %s: %s\n", h.Url, err.Error())
	}
	h.setHeaders(req)

	// Wait for a token from the limiter (respecting cancellation)
	if err := limiter.Wait(context.Background()); err != nil {
		panic(fmt.Sprintf("rate limit wait failed: %s\n", err.Error()))
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch URL %s: %s\n", h.Url, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non-OK HTTP status for URL %s: %s\n", h.Url, resp.Status)
	}

	return resp, nil
}

// setHeaders sets the headers for the HTTP request
func (h Http) setHeaders(req *http.Request) {
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}
}
