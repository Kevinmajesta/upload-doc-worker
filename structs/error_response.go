package structs

type ErrorResponse struct {
	Code    int    `json:"code"`
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}
