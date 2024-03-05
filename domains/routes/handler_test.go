package routes_test

import (
	"emailn/constants"
	"emailn/domains/routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerErrorWhenRouteReturnInternalError(t *testing.T) {
	assert := assert.New(t)

	route := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, http.StatusInternalServerError, constants.ErrInternalServer
	}

	handlerFunc := routes.HandlerError(route)

	req, _ := http.NewRequest("GET", "/", nil)

	rec := httptest.NewRecorder()

	handlerFunc.ServeHTTP(rec, req)

	assert.Equal(http.StatusInternalServerError, rec.Code)

}

func TestHandlerErrorWhenRouteReturnBadRequest(t *testing.T) {
	assert := assert.New(t)

	route := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, http.StatusBadRequest, constants.ErrInvalidEmail
	}

	handlerFunc := routes.HandlerError(route)

	req, _ := http.NewRequest("GET", "/", nil)

	rec := httptest.NewRecorder()

	handlerFunc.ServeHTTP(rec, req)

	assert.Equal(http.StatusBadRequest, rec.Code)

}

func TestHandlerErrorWhenRouteReturnSuccess(t *testing.T) {
	assert := assert.New(t)

	route := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return map[string]string{"message": "success"}, http.StatusOK, nil
	}

	handlerFunc := routes.HandlerError(route)

	req, _ := http.NewRequest("GET", "/", nil)

	rec := httptest.NewRecorder()

	handlerFunc.ServeHTTP(rec, req)

	assert.Equal(http.StatusOK, rec.Code)
	assert.JSONEq(`{"message": "success"}`, rec.Body.String())

}
