package httpResponse

type ErrorResponse struct {
    StatusCode int         `json:"status_code"`
    Message    string      `json:"message"`
    Data       any         `json:"data,omitempty"`
    Error      string      `json:"error,omitempty"`
}

type SuccessResponse struct {
    StatusCode int         `json:"status_code"`
    Message    string      `json:"message"`
    Data       any         `json:"data"`
}