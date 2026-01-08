package domains

type ErrorResponse struct {
	Error *APIError `json:"error"`
}