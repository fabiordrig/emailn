package routes

import (
	"emailn/contracts"
	"net/http"

	"github.com/go-chi/render"
)

func (h *Handler) CreateCampaign(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	var payload contracts.NewCampaign
	render.DecodeJSON(r.Body, &payload)

	campaign, err := h.CampaignService.Create(payload)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return map[string]string{
		"id":        campaign.ID.String(),
		"name":      campaign.Name,
		"createdAt": campaign.CreatedAt.String(),
		"content":   campaign.Content,
	}, http.StatusCreated, err

}

func (h *Handler) FindALlCampaigns(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	campaigns, err := h.CampaignService.FindAll()

	if err != nil {

		return nil, http.StatusNotFound, err
	}

	return campaigns, http.StatusOK, nil

}
