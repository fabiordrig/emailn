package middlewares_test

import (
	"emailn/domains/routes/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_NoAuthorizationHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	handler := middlewares.AuthMiddleware(testHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "{\"error\":\"Unauthorized\"}\n", rr.Body.String())
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	handler := middlewares.AuthMiddleware(testHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "{\"error\":\"Invalid token\"}\n", rr.Body.String())
}
