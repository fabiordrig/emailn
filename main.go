package main

import (
	"emailn/contracts"
	"emailn/domains/campaign"
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

	service := campaign.Service{}

	r.Post("/campaigns", func(w http.ResponseWriter, r *http.Request) {

		var payload contracts.NewCampaign
		render.DecodeJSON(r.Body, &payload)

		campaign, err := service.Create(payload)

		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
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
