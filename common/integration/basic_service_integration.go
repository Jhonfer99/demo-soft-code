package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsangopanta/demo-soft-code/common/domains"
	"go.uber.org/zap"
)

// Constants
const (
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
	HeaderTransactionID = "X-Transaction-ID"
	HeaderRetryAfter    = "Retry-After"
	
	MediaTypeJSON = "application/json"
	
	MaxRetries     = 3
	DefaultTimeout = 30 * time.Second
)

type APIError domains.APIError
type ErrorResponse domains.ErrorResponse

func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Code, e.Message)
}



// Retry policy
type RetryPolicy struct {
	MaxRetries   int
	BaseDelay    time.Duration
	MaxDelay     time.Duration
	RetryableCodes []int
}

func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:   3,
		BaseDelay:    100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		RetryableCodes: []int{
			http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout,
		},
	}
}

// HTTP client wrapper
type HTTPClient struct {
	client      *http.Client
	baseURL     string
	logger      *zap.Logger
	retryPolicy *RetryPolicy
	headers     map[string]string
	mu          sync.RWMutex
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(baseURL string, logger *zap.Logger) *HTTPClient {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL:     strings.TrimSuffix(baseURL, "/"),
		logger:      logger,
		retryPolicy: DefaultRetryPolicy(),
		headers:     make(map[string]string),
	}
}

// SetHeaders sets default headers
func (c *HTTPClient) SetHeaders(headers map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for k, v := range headers {
		c.headers[k] = v
	}
}

// SetRetryPolicy sets retry policy
func (c *HTTPClient) SetRetryPolicy(policy *RetryPolicy) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.retryPolicy = policy
}

// SetTimeout sets request timeout
func (c *HTTPClient) SetTimeout(timeout time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.client.Timeout = timeout
}

// GET request
func (c *HTTPClient) GET(ctx context.Context, path string, response interface{}) error {
	return c.DoRequest(ctx, http.MethodGet, path, nil, response, nil)
}

// POST request
func (c *HTTPClient) POST(ctx context.Context, path string, body, response interface{}) error {
	return c.DoRequest(ctx, http.MethodPost, path, body, response, nil)
}

// PUT request
func (c *HTTPClient) PUT(ctx context.Context, path string, body, response interface{}) error {
	return c.DoRequest(ctx, http.MethodPut, path, body, response, nil)
}

// PATCH request
func (c *HTTPClient) PATCH(ctx context.Context, path string, body, response interface{}) error {
	return c.DoRequest(ctx, http.MethodPatch, path, body, response, nil)
}

// DELETE request
func (c *HTTPClient) DELETE(ctx context.Context, path string, response interface{}) error {
	return c.DoRequest(ctx, http.MethodDelete, path, nil, response, nil)
}

// DoRequest executes HTTP request with retry logic
func (c *HTTPClient) DoRequest(
	ctx context.Context,
	method string,
	path string,
	body interface{},
	response interface{},
	customHeaders map[string]string,
) error {
	// Build URL
	url := c.baseURL + path
	
	// Prepare request body
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	
	// Set headers
	c.setRequestHeaders(req, body, customHeaders)
	
	// Generate transaction ID for logging
	txID := generateTransactionID()
	req.Header.Set(HeaderTransactionID, txID)
	
	// Log request
	c.logRequest(txID, method, url, body)
	
	// Execute with retry
	var lastError error
	var lastResponse *http.Response
	
	for attempt := 0; attempt <= c.retryPolicy.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := c.calculateRetryDelay(attempt, lastResponse)
			c.logger.Info("Retrying request",
				zap.String("txID", txID),
				zap.Int("attempt", attempt),
				zap.Duration("delay", delay))
			
			select {
			case <-time.After(delay):
				// Continue with retry
			case <-ctx.Done():
				return ctx.Err()
			}
			
			// Reset body reader for retry
			if bodyBytes != nil {
				req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}
		
		// Execute request
		startTime := time.Now()
		resp, err := c.client.Do(req)
		duration := time.Since(startTime)
		
		if err != nil {
			lastError = fmt.Errorf("execute request: %w", err)
			c.logError(txID, method, url, err, duration)
			
			// Check if error is retryable
			if !c.isRetryableError(err) || attempt == c.retryPolicy.MaxRetries {
				return lastError
			}
			continue
		}
		
		defer resp.Body.Close()
		lastResponse = resp
		
		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastError = fmt.Errorf("read response: %w", err)
			continue
		}
		
		// Log response
		c.logResponse(txID, method, url, resp.StatusCode, respBody, duration)
		
		// Check if status code indicates retry
		if c.shouldRetry(resp.StatusCode) && attempt < c.retryPolicy.MaxRetries {
			c.logger.Warn("Retryable status code",
				zap.String("txID", txID),
				zap.Int("status", resp.StatusCode),
				zap.Int("attempt", attempt))
			continue
		}
		
		// Handle error responses
		if resp.StatusCode >= 400 {
			return c.handleHTTPError(txID, resp.StatusCode, respBody)
		}
		
		// Parse successful response
		if response != nil {
			if err := json.Unmarshal(respBody, response); err != nil {
				return fmt.Errorf("unmarshal response: %w", err)
			}
		}
		
		return nil
	}
	
	return lastError
}

// Calculate retry delay with exponential backoff
func (c *HTTPClient) calculateRetryDelay(attempt int, resp *http.Response) time.Duration {
	// Check for Retry-After header
	if resp != nil {
		if retryAfter := resp.Header.Get(HeaderRetryAfter); retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				return time.Duration(seconds) * time.Second
			}
		}
	}
	
	// Exponential backoff with jitter
	delay := float64(c.retryPolicy.BaseDelay) * math.Pow(2, float64(attempt))
	delay = math.Min(delay, float64(c.retryPolicy.MaxDelay))
	
	// Add jitter (±20%)
	jitter := 0.8 + 0.4*float64(time.Now().UnixNano()%100)/100
	delay *= jitter
	
	return time.Duration(delay)
}

// Check if error is retryable
func (c *HTTPClient) isRetryableError(err error) bool {
	errStr := strings.ToLower(err.Error())
	
	retryableErrors := []string{
		"timeout",
		"deadline exceeded",
		"connection reset",
		"connection refused",
		"eof",
		"network",
		"tls",
	}
	
	for _, retryableErr := range retryableErrors {
		if strings.Contains(errStr, retryableErr) {
			return true
		}
	}
	
	return false
}

// Check if status code is retryable
func (c *HTTPClient) shouldRetry(statusCode int) bool {
	for _, code := range c.retryPolicy.RetryableCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

// Set request headers
func (c *HTTPClient) setRequestHeaders(req *http.Request, body interface{}, customHeaders map[string]string) {
	// Set default headers
	c.mu.RLock()
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	c.mu.RUnlock()
	
	// Set custom headers
	for k, v := range customHeaders {
		req.Header.Set(k, v)
	}
	
	// Set content type
	if body != nil {
		req.Header.Set(HeaderContentType, MediaTypeJSON)
	}
	
	// Always accept JSON
	req.Header.Set("Accept", MediaTypeJSON)
}

// Handle HTTP error responses
func (c *HTTPClient) handleHTTPError(txID string, statusCode int, body []byte) error {
	// Try to parse as structured error
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != nil {
		errResp.Error.StatusCode = statusCode
		return errResp.Error
	}
	
	// Try to parse as simple error
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Code != "" {
		apiErr.StatusCode = statusCode
		return &apiErr
	}
	
	// Generic error
	return &APIError{
		Code:       fmt.Sprintf("HTTP_%d", statusCode),
		Message:    string(body),
		StatusCode: statusCode,
	}
}

// Logging methods
func (c *HTTPClient) logRequest(txID, method, url string, body interface{}) {
	logFields := []zap.Field{
		zap.String("txID", txID),
		zap.String("method", method),
		zap.String("url", url),
		zap.String("direction", "outbound"),
	}
	
	if body != nil {
		// Mask sensitive data before logging
		maskedBody := maskSensitiveData(body)
		bodyJSON, _ := json.Marshal(maskedBody)
		logFields = append(logFields, zap.String("body", string(bodyJSON)))
	}
	
	c.logger.Info("HTTP request", logFields...)
}

func (c *HTTPClient) logResponse(txID, method, url string, statusCode int, body []byte, duration time.Duration) {
	logFields := []zap.Field{
		zap.String("txID", txID),
		zap.String("method", method),
		zap.String("url", url),
		zap.Int("status", statusCode),
		zap.Duration("duration", duration),
		zap.String("direction", "inbound"),
	}
	
	// Try to parse and mask response
	if len(body) > 0 {
		var response interface{}
		if err := json.Unmarshal(body, &response); err == nil {
			maskedResponse := maskSensitiveData(response)
			responseJSON, _ := json.Marshal(maskedResponse)
			logFields = append(logFields, zap.String("body", string(responseJSON)))
		} else {
			// Not JSON, log truncated
			bodyStr := string(body)
			if len(bodyStr) > 500 {
				bodyStr = bodyStr[:500] + "..."
			}
			logFields = append(logFields, zap.String("body", bodyStr))
		}
	}
	
	if statusCode >= 400 {
		c.logger.Error("HTTP error response", logFields...)
	} else {
		c.logger.Info("HTTP response", logFields...)
	}
}

func (c *HTTPClient) logError(txID, method, url string, err error, duration time.Duration) {
	c.logger.Error("HTTP request failed",
		zap.String("txID", txID),
		zap.String("method", method),
		zap.String("url", url),
		zap.Duration("duration", duration),
		zap.Error(err))
}

// Mask sensitive data in logs
func maskSensitiveData(data interface{}) interface{} {
	sensitiveFields := []string{
		"password",
		"token",
		"secret",
		"creditCard",
		"cvv",
		"ssn",
		"authorization",
	}
	
	return recursiveMask(data, sensitiveFields)
}

func recursiveMask(data interface{}, sensitiveFields []string) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		masked := make(map[string]interface{})
		for key, value := range v {
			lowerKey := strings.ToLower(key)
			isSensitive := false
			
			for _, field := range sensitiveFields {
				if strings.Contains(lowerKey, field) {
					isSensitive = true
					break
				}
			}
			
			if isSensitive {
				masked[key] = "***MASKED***"
			} else {
				masked[key] = recursiveMask(value, sensitiveFields)
			}
		}
		return masked
		
	case []interface{}:
		masked := make([]interface{}, len(v))
		for i, item := range v {
			masked[i] = recursiveMask(item, sensitiveFields)
		}
		return masked
		
	default:
		return data
	}
}

// Generate transaction ID
func generateTransactionID() string {
	return fmt.Sprintf("tx-%d-%s", 
		time.Now().UnixNano(),
		randomString(8))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// // Example usage
// type User struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

// type BeneficiaryService struct {
// 	client *HTTPClient
// }

// func NewBeneficiaryService(baseURL string, logger *zap.Logger) *BeneficiaryService {
// 	client := NewHTTPClient(baseURL, logger)
	
// 	// Set default headers
// 	client.SetHeaders(map[string]string{
// 		"X-Service-Name": "beneficiary-service",
// 		"User-Agent":     "go-client/1.0",
// 	})
	
// 	return &BeneficiaryService{
// 		client: client,
// 	}
// }

// func (s *BeneficiaryService) GetBeneficiaries(ctx context.Context, userID string) ([]User, error) {
// 	var response struct {
// 		Data []User      `json:"data"`
// 		Error *APIError `json:"error,omitempty"`
// 	}
	
// 	path := fmt.Sprintf("/users/%s/beneficiaries", userID)
	
// 	if err := s.client.GET(ctx, path, &response); err != nil {
// 		return nil, err
// 	}
	
// 	if response.Error != nil {
// 		return nil, response.Error
// 	}
	
// 	return response.Data, nil
// }

// func (s *BeneficiaryService) CreateBeneficiary(ctx context.Context, userID string, beneficiary User) error {
// 	var response struct {
// 		Success bool      `json:"success"`
// 		Error   *APIError `json:"error,omitempty"`
// 	}
	
// 	path := fmt.Sprintf("/users/%s/beneficiaries", userID)
	
// 	if err := s.client.POST(ctx, path, beneficiary, &response); err != nil {
// 		return err
// 	}
	
// 	if !response.Success {
// 		return response.Error
// 	}
	
// 	return nil
// }

// func (s *BeneficiaryService) UpdateBeneficiary(ctx context.Context, beneficiaryID string, updates User) error {
// 	var response struct {
// 		Success bool      `json:"success"`
// 		Error   *APIError `json:"error,omitempty"`
// 	}
	
// 	path := fmt.Sprintf("/beneficiaries/%s", beneficiaryID)
	
// 	if err := s.client.PUT(ctx, path, updates, &response); err != nil {
// 		return err
// 	}
	
// 	if !response.Success {
// 		return response.Error
// 	}
	
// 	return nil
// }

// func (s *BeneficiaryService) DeleteBeneficiary(ctx context.Context, beneficiaryID string) error {
// 	var response struct {
// 		Success bool      `json:"success"`
// 		Error   *APIError `json:"error,omitempty"`
// 	}
	
// 	path := fmt.Sprintf("/beneficiaries/%s", beneficiaryID)
	
// 	if err := s.client.DELETE(ctx, path, &response); err != nil {
// 		return err
// 	}
	
// 	if !response.Success {
// 		return response.Error
// 	}
	
// 	return nil
// }

// // Middleware for authentication
// func (s *BeneficiaryService) WithAuth(token string) *BeneficiaryService {
// 	s.client.SetHeaders(map[string]string{
// 		HeaderAuthorization: "Bearer " + token,
// 	})
// 	return s
// }

// // Main example
// func main() {
// 	// Create logger
// 	logger, _ := zap.NewDevelopment()
// 	defer logger.Sync()
	
// 	// Create service
// 	service := NewBeneficiaryService("https://api.example.com/v1", logger)
	
// 	// Set custom retry policy
// 	service.client.SetRetryPolicy(&RetryPolicy{
// 		MaxRetries:   5,
// 		BaseDelay:    500 * time.Millisecond,
// 		MaxDelay:     10 * time.Second,
// 		RetryableCodes: []int{429, 500, 502, 503, 504},
// 	})
	
// 	// Use service with auth
// 	ctx := context.Background()
// 	service.WithAuth("your-jwt-token")
	
// 	// Get beneficiaries
// 	beneficiaries, err := service.GetBeneficiaries(ctx, "user-123")
// 	if err != nil {
// 		logger.Error("Failed to get beneficiaries", zap.Error(err))
// 		return
// 	}
	
// 	logger.Info("Retrieved beneficiaries",
// 		zap.Int("count", len(beneficiaries)))
	
// 	// Create beneficiary
// 	newBeneficiary := User{
// 		ID:   "ben-456",
// 		Name: "John Doe",
// 	}
	
// 	if err := service.CreateBeneficiary(ctx, "user-123", newBeneficiary); err != nil {
// 		logger.Error("Failed to create beneficiary", zap.Error(err))
// 		return
// 	}
	
// 	logger.Info("Beneficiary created successfully")
// }