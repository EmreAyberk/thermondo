package response

import (
	"movie-rating-service/config"
)

type SuccessResponse struct {
	Status  string      `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"success"`
	Message string `json:"message"`
	Cause   string `json:"cause,omitempty"`
}

func Success(data interface{}) SuccessResponse {
	return SuccessResponse{Status: "success", Data: data}
}

func Error(message, cause string) ErrorResponse {
	res := ErrorResponse{Status: "error", Message: message}
	if config.Cfg.DebugMode {
		res.Cause = cause
	}
	return res
}
