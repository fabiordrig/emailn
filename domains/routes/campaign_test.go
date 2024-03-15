package routes_test

import (
	"bytes"
	"context"
	"emailn/constants"
	"emailn/contracts"
	"emailn/domains/campaign"
	"emailn/domains/routes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) Create(newCampaign contracts.NewCampaign) (*campaign.Campaign, error) {
	args := m.Called(newCampaign)
	return args.Get(0).(*campaign.Campaign), args.Error(1)
}

func (m *mockService) FindAll() ([]campaign.Campaign, error) {
	args := m.Called()
	return args.Get(0).([]campaign.Campaign), args.Error(1)
}

func (m *mockService) FindByID(id string) (*campaign.Campaign, error) {
	args := m.Called(id)
	return args.Get(0).(*campaign.Campaign), args.Error(1)
}

func (m *mockService) Cancel(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

var (
	fake      = faker.New()
	createdBy = fake.Internet().Email()
)

func TestCreateCampaignShouldCreateACampaign(t *testing.T) {
	assert := assert.New(t)

	newCampaign := contracts.NewCampaign{
		Name:      fake.Beer().Name(),
		Content:   fake.Beer().Alcohol(),
		Emails:    []string{fake.Internet().Email()},
		CreatedBy: createdBy,
	}

	serviceMock := new(mockService)
	body := &campaign.Campaign{
		ID:      uuid.New(),
		Name:    newCampaign.Name,
		Content: newCampaign.Content,
		Contacts: []campaign.Contact{
			{
				Email: newCampaign.Emails[0],
			},
		},
	}
	serviceMock.On("Create", newCampaign).Return(body, nil)
	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(newCampaign)

	req, _ := http.NewRequest("POST", "/campaign", &buffer)
	ctx := context.WithValue(req.Context(), "email", createdBy)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	response, status, err := handler.CreateCampaign(rr, req)

	assert.Nil(err)
	assert.Equal(http.StatusCreated, status)
	assert.Equal(body.ID.String(), response.(map[string]string)["id"])

}

func TestCreateCampaignShouldReturnBadRequestWhenServiceReturnError(t *testing.T) {
	assert := assert.New(t)

	newCampaign := contracts.NewCampaign{
		Name:    fake.Beer().Name(),
		Content: fake.Beer().Alcohol(),
		Emails:  []string{fake.Internet().Email()},
	}

	serviceMock := new(mockService)
	serviceMock.On("Create", mock.Anything).Return(
		&campaign.Campaign{},
		errors.New(constants.ErrInvalidEmail.Error()),
	)
	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(newCampaign)

	req, _ := http.NewRequest("POST", "/campaign", &buffer)
	ctx := context.WithValue(req.Context(), "email", createdBy)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	response, status, err := handler.CreateCampaign(rr, req)

	assert.Equal(http.StatusBadRequest, status)
	assert.NotNil(err)
	assert.Nil(response)

}

func TestFindAllCampaignsShouldReturnAllCampaigns(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	serviceMock := new(mockService)
	campaigns := []campaign.Campaign{
		{
			ID:      uuid.New(),
			Name:    fake.Beer().Name(),
			Content: fake.Beer().Alcohol(),
			Contacts: []campaign.Contact{
				{
					Email: fake.Internet().Email(),
				},
			},
		},
	}

	serviceMock.On("FindAll").Return(campaigns, nil)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("GET", "/campaign", nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.FindALlCampaigns(rr, req)

	assert.Nil(err)
	assert.Equal(http.StatusOK, status)
	assert.Equal(campaigns, response)

}

func TestFindAllCampaignsShouldReturnNotFoundWhenServiceReturnError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("FindAll").Return(
		[]campaign.Campaign{},
		errors.New("error"),
	)
	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("GET", "/campaign", nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.FindALlCampaigns(rr, req)

	assert.Equal(http.StatusNotFound, status)
	assert.NotNil(err)
	assert.Nil(response)

}

func TestFindCampaignByIDShouldReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	serviceMock := new(mockService)
	campaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Beer().Name(),
		Content: fake.Beer().Alcohol(),
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	serviceMock.On("FindByID", mock.Anything).Return(&campaign, nil)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("GET", "/campaigns/"+campaign.ID.String(), nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.FindCampaignByID(rr, req)

	assert.Nil(err)
	assert.Equal(http.StatusOK, status)
	assert.Equal(&campaign, response)

}

func TestFindCampaignByIDShouldReturnNotFoundWhenServiceReturnError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("FindByID", mock.Anything).Return(
		&campaign.Campaign{},
		errors.New("error"),
	)
	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("GET", "/campaigns/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.FindCampaignByID(rr, req)

	assert.Equal(http.StatusNotFound, status)
	assert.NotNil(err)
	assert.Nil(response)

}

func TestCancelCampaignShouldReturnNoContent(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Cancel", mock.Anything).Return(nil)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("POST", "/campaigns/"+uuid.New().String()+"/cancel", nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.CancelCampaign(rr, req)

	assert.Nil(err)
	assert.Equal(http.StatusNoContent, status)
	assert.Equal(response, map[string]string{"message": "campaign cancelled"})

}

func TestCancelCampaignShouldReturnNotFoundWhenServiceReturnNotFoundError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Cancel", mock.Anything).Return(constants.ErrNotFound)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("POST", "/campaigns/"+uuid.New().String()+"/cancel", nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.CancelCampaign(rr, req)

	assert.Equal(http.StatusNotFound, status)
	assert.NotNil(err)
	assert.Nil(response)
}

func TestCancelCampaignShouldReturnUnprocessableEntityWhenServiceReturnUnprocessableEntityError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Cancel", mock.Anything).Return(constants.ErrUnprocessableEntity)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("POST", "/campaigns/"+uuid.New().String()+"/cancel", nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.CancelCampaign(rr, req)

	assert.Equal(http.StatusUnprocessableEntity, status)
	assert.NotNil(err)
	assert.Nil(response)
}

func TestCancelCampaignShouldReturnInternalServerErrorWhenServiceReturnError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Cancel", mock.Anything).Return(errors.New("error"))

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("POST", "/campaigns/"+uuid.New().String()+"/cancel", nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.CancelCampaign(rr, req)

	assert.Equal(http.StatusInternalServerError, status)
	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteCampaignShouldReturnNoContent(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Delete", mock.Anything).Return(nil)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("DELETE", "/campaigns/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.DeleteCampaign(rr, req)

	assert.Nil(err)
	assert.Equal(http.StatusNoContent, status)
	assert.Equal(response, map[string]string{"message": "campaign deleted"})

}

func TestDeleteCampaignShouldReturnNotFoundWhenServiceReturnNotFoundError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Delete", mock.Anything).Return(constants.ErrNotFound)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("DELETE", "/campaigns/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.DeleteCampaign(rr, req)

	assert.Equal(http.StatusNotFound, status)
	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteCampaignShouldReturnUnprocessableEntityWhenServiceReturnUnprocessableEntityError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Delete", mock.Anything).Return(constants.ErrUnprocessableEntity)

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("DELETE", "/campaigns/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.DeleteCampaign(rr, req)

	assert.Equal(http.StatusUnprocessableEntity, status)
	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteCampaignShouldReturnInternalServerErrorWhenServiceReturnError(t *testing.T) {
	assert := assert.New(t)

	serviceMock := new(mockService)
	serviceMock.On("Delete", mock.Anything).Return(errors.New("error"))

	handler := routes.Handler{
		CampaignService: serviceMock,
	}

	req, _ := http.NewRequest("DELETE", "/campaigns/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()

	response, status, err := handler.DeleteCampaign(rr, req)

	assert.Equal(http.StatusInternalServerError, status)
	assert.NotNil(err)
	assert.Nil(response)
}
