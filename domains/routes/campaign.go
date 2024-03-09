package routes

import (
	"emailn/constants"
	"emailn/contracts"
	"errors"
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

func (h *Handler) FindCampaignByID(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	id := r.PathValue("id")

	campaign, err := h.CampaignService.FindByID(id)

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return campaign, http.StatusOK, nil

}

func (h *Handler) CancelCampaign(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	id := r.PathValue("id")

	err := h.CampaignService.Cancel(id)

	if err == nil {
		return map[string]string{"message": "campaign cancelled"}, http.StatusNoContent, nil
	}

	var status int
	switch {
	case errors.Is(err, constants.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, constants.ErrUnprocessableEntity):
		status = http.StatusUnprocessableEntity
	default:
		status = http.StatusInternalServerError
	}

	return nil, status, err

}
