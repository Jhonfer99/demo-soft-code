package domains

type APIError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"statusCode"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}