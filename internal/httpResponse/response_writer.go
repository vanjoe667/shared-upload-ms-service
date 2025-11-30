package httpResponse

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, message string, err error) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    res := ErrorResponse{
        StatusCode: status,
        Message:    message,
    }

    if err != nil {
        res.Error = err.Error()
    }

    json.NewEncoder(w).Encode(res)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    res := SuccessResponse{
        StatusCode: status,
        Message:    "success",
        Data:       data,
    }

    json.NewEncoder(w).Encode(res)
}