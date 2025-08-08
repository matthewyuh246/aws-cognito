package httpclient

import (
	"context"
	"crypto/rand"
	"errors"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/matthewyuh246/aws-cognito/pkg/logger"
)

type Config struct {
	Timeout     time.Duration
	MaxRetries  int
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
	JitterMax   time.Duration
}

type Client struct {
	httpClient *http.Client
	config     Config
	logger     *logger.Logger
}

func NewClient(config Config, logger *logger.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
		logger: logger,
	}
}

func (c *Client) DoWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := c.calculateBackoffWithJitter(attempt)
			
			c.logger.Info("HTTPリトライ実行", map[string]interface{}{
				"attempt":    attempt,
				"max_retry":  c.config.MaxRetries,
				"backoff_ms": backoff.Milliseconds(),
				"last_error": lastErr.Error(),
			})

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, err := c.httpClient.Do(req)
		if err == nil && !c.shouldRetry(resp.StatusCode) {
			return resp, nil
		}

		if err != nil {
			lastErr = err
		} else {
			resp.Body.Close()
			lastErr = errors.New("HTTP error: " + resp.Status)
		}

		if !c.isRetriableError(err, resp) {
			break
		}
	}

	return nil, lastErr
}

func (c *Client) calculateBackoffWithJitter(attempt int) time.Duration {
	exponential := c.config.BaseBackoff * time.Duration(math.Pow(2, float64(attempt-1)))
	backoff := time.Duration(math.Min(float64(exponential), float64(c.config.MaxBackoff)))
	
	jitter := c.generateJitter()
	return backoff + jitter
}

func (c *Client) generateJitter() time.Duration {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return time.Duration(time.Now().UnixNano()) % c.config.JitterMax
	}
	
	jitter := int(buf[0])<<24 | int(buf[1])<<16 | int(buf[2])<<8 | int(buf[3])
	if jitter < 0 {
		jitter = -jitter
	}
	
	return time.Duration(jitter) % c.config.JitterMax
}

func (c *Client) shouldRetry(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

func (c *Client) isRetriableError(err error, resp *http.Response) bool {
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return true
		}
		
		var netErr net.Error
		if errors.As(err, &netErr) {
			return netErr.Temporary() || netErr.Timeout()
		}
		
		errStr := strings.ToLower(err.Error())
		networkErrors := []string{
			"connection refused", "connection reset", "no such host",
			"timeout", "network unreachable",
		}
		
		for _, pattern := range networkErrors {
			if strings.Contains(errStr, pattern) {
				return true
			}
		}
	}
	
	if resp != nil {
		return c.shouldRetry(resp.StatusCode)
	}
	
	return false
}