package utils

import (
	"consultant-management/backend/internal/logger"
	"net/http"
)

// HandleError logs the error and sends an appropriate HTTP response
func HandleError(w http.ResponseWriter, err error, message string, code int) {
	logger.ErrorLogger.Printf("%s: %v", message, err)
	http.Error(w, message, code)
}
