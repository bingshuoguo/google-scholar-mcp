package scholar

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bingshuoguo/google-scholar-mcp/internal/config"
	"golang.org/x/time/rate"
)

type Response struct {
	StatusCode int
	URL        string
	Body       []byte
}

type Client struct {
	baseURL        string
	httpClient     *http.Client
	limiter        *rate.Limiter
	logger         *slog.Logger
	userAgent      string
	acceptLanguage string
}

func NewClient(cfg config.Config, logger *slog.Logger) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 20
	transport.MaxIdleConnsPerHost = 10
	transport.MaxConnsPerHost = 10
	transport.IdleConnTimeout = 90 * time.Second

	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		limiter:        rate.NewLimiter(rate.Limit(cfg.RateLimitRPS), 1),
		logger:         logger,
		userAgent:      cfg.UserAgent,
		acceptLanguage: cfg.AcceptLanguage,
	}
}

func (c *Client) Get(ctx context.Context, path string, query url.Values) (*Response, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, wrap(ErrTimeout, "rate limiter wait failed", err)
	}

	fullURL := c.baseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", c.userAgent)
		req.Header.Set("Accept-Language", c.acceptLanguage)

		res, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if isTimeout(err) {
				if attempt < 2 {
					sleepWithJitter(ctx, attempt)
					continue
				}
				return nil, wrap(ErrTimeout, "request to Google Scholar timed out", err)
			}
			if attempt < 2 {
				sleepWithJitter(ctx, attempt)
				continue
			}
			return nil, wrap(ErrUpstreamUnavailable, "request to Google Scholar failed", err)
		}

		body, readErr := io.ReadAll(io.LimitReader(res.Body, 4<<20))
		closeErr := res.Body.Close()
		if readErr != nil {
			return nil, wrap(ErrUpstreamUnavailable, "failed to read Google Scholar response", readErr)
		}
		if closeErr != nil {
			c.logger.Debug("close response body", "error", closeErr)
		}

		response := &Response{
			StatusCode: res.StatusCode,
			URL:        req.URL.String(),
			Body:       body,
		}

		if isRetryableStatus(res.StatusCode) && attempt < 2 {
			lastErr = fmt.Errorf("retryable status: %d", res.StatusCode)
			sleepWithJitter(ctx, attempt)
			continue
		}

		if res.StatusCode == http.StatusTooManyRequests || looksBlocked(body) {
			return nil, wrap(ErrUpstreamBlocked, "Google Scholar blocked the request", nil)
		}
		if res.StatusCode != http.StatusOK {
			return nil, wrap(ErrUpstreamUnavailable, fmt.Sprintf("unexpected Google Scholar status %d", res.StatusCode), nil)
		}

		return response, nil
	}

	if lastErr != nil {
		return nil, wrap(ErrUpstreamUnavailable, "request retries exhausted", lastErr)
	}
	return nil, wrap(ErrUpstreamUnavailable, "request retries exhausted", nil)
}

func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func sleepWithJitter(ctx context.Context, attempt int) {
	base := time.Duration(300*(attempt+1)) * time.Millisecond
	jitter := time.Duration(rand.IntN(250)) * time.Millisecond
	timer := time.NewTimer(base + jitter)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func looksBlocked(body []byte) bool {
	lower := strings.ToLower(string(body))
	return strings.Contains(lower, "unusual traffic from your computer network") ||
		strings.Contains(lower, "our systems have detected unusual traffic") ||
		strings.Contains(lower, "captcha")
}
