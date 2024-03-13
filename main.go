package main

import (
	"emailn/domains/campaign"
	"emailn/domains/routes"
	"emailn/domains/routes/middlewares"
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
	r.Use(middlewares.AuthMiddleware)

	db := database.NewDb()
	service := campaign.NewService(&database.CampaignRepository{
		Db: db,
	})
	handler := routes.Handler{
		CampaignService: service,
	}

	r.Post("/campaigns", routes.HandlerError(handler.CreateCampaign))
	r.Get("/campaigns", routes.HandlerError(handler.FindALlCampaigns))
	r.Get("/campaigns/{id}", routes.HandlerError(handler.FindCampaignByID))
	r.Post("/campaigns/{id}/cancel", routes.HandlerError(handler.CancelCampaign))
	r.Delete("/campaigns/{id}", routes.HandlerError(handler.DeleteCampaign))

	http.ListenAndServe(":8000", r)

}
