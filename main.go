package main

import (
	"emailn/domains/campaign"
	"emailn/domains/routes"
	"emailn/domains/routes/middlewares"
	"emailn/infra/database"
	"emailn/infra/mail"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db := database.NewDb()
	emailSender := mail.NewSMTPSender()
	service := campaign.NewService(&database.CampaignRepository{
		Db: db,
	},
		emailSender,
	)
	handler := routes.Handler{
		CampaignService: service,
	}

	r.Route("/api", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Post("/campaigns", routes.HandlerError(handler.CreateCampaign))
		r.Get("/campaigns", routes.HandlerError(handler.FindALlCampaigns))
		r.Get("/campaigns/{id}", routes.HandlerError(handler.FindCampaignByID))
		r.Post("/campaigns/{id}/cancel", routes.HandlerError(handler.CancelCampaign))
		r.Delete("/campaigns/{id}", routes.HandlerError(handler.DeleteCampaign))
	})

	http.ListenAndServe(":8000", r)

}
