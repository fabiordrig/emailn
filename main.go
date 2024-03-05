package main

import (
	"emailn/domains/campaign"
	"emailn/domains/routes"
	"emailn/infra/database"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	service := campaign.NewService(&database.CampaignRepository{})
	handler := routes.Handler{
		CampaignService: *service,
	}

	r.Post("/campaigns", handler.HandlerError(handler.CreateCampaign))
	r.Get("/campaigns", handler.HandlerError(handler.FindALlCampaigns))

	http.ListenAndServe(":8000", r)

}
