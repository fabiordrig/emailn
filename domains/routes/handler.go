package routes

import (
	"emailn/constants"
	"emailn/domains/campaign"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type Handler struct {
	CampaignService campaign.Service
}

type RouterFunc func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)

func HandlerError(handler RouterFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj, status, err := handler(w, r)

		if err != nil {
			if errors.Is(err, constants.ErrInternalServer) {
				render.Status(r, http.StatusInternalServerError)
			} else {
				render.Status(r, status)
			}
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
		}

		render.Status(r, status)

		if obj != nil {
			render.JSON(w, r, obj)
		}
	})
}
