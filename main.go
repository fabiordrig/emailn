package main

import (
	"emailn/constants"
	"emailn/contracts"
	"emailn/domains/campaign"
	"emailn/infra/database"
	"errors"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	service := campaign.NewService(&database.CampaignRepository{})

	r.Post("/campaigns", func(w http.ResponseWriter, r *http.Request) {

		var payload contracts.NewCampaign
		render.DecodeJSON(r.Body, &payload)

		campaign, err := service.Create(payload)

		if err != nil {
			if errors.Is(err, constants.ErrInternalServer) {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{
					"error": err.Error(),
				})
				return
			} else {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, map[string]string{
					"error": err.Error(),
				})
				return
			}
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]interface{}{
			"id":        campaign.ID,
			"name":      campaign.Name,
			"createdAt": campaign.CreatedAt,
			"content":   campaign.Content,
		})
	})

	http.ListenAndServe(":8000", r)

}
