package http

import (
	"errors"
	"go-hexagonal/internal/core/domain"
	"net/http"
)

var errorMapping = map[error]int{
	domain.ErrorInsufficientFunds:    http.StatusBadRequest,
	domain.ErrorInvalidAmount:        http.StatusBadRequest,
	domain.ErrConcurrentModification: http.StatusBadRequest,
}

func mapErrorToHTTP(err error) int {

	// Map'i tek tek döner, sarmalanmış hatalrı kontrol eder..
	for domainErr, httpCode := range errorMapping {
		if errors.Is(err, domainErr) {
			return httpCode
		}
	}

	return http.StatusInternalServerError
}
