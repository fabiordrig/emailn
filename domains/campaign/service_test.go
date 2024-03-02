package campaign

import (
	"emailn/contracts"
	"emailn/domains"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(campaign *Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

var (
	newCampaign = contracts.NewCampaign{
		Name:    "test",
		Content: "content",
		Emails:  []string{"test@t.com"},
	}
	repoMock = new(MockRepository)
	service  = NewService(repoMock)
)

func TestCreateCampaign(t *testing.T) {
	assert := assert.New(t)

	repoMock.On("Save", mock.Anything).Return(nil)

	campaign, err := service.Create(newCampaign)

	assert.Nil(err)
	assert.NotNil(campaign.ID)
	assert.Equal(newCampaign.Name, campaign.Name)
	assert.Equal(newCampaign.Content, campaign.Content)
	assert.Len(campaign.Contacts, 1)
	assert.Equal(newCampaign.Emails[0], campaign.Contacts[0].Email)

}

func TestSaveCampaign(t *testing.T) {
	assert := assert.New(t)

	repoMock.On("Save", mock.MatchedBy(func(campaign *Campaign) bool {

		if campaign.Name != newCampaign.Name {
			return false
		} else if campaign.Content != newCampaign.Content {
			return false
		} else if len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}

		return true
	})).Return(nil)

	_, err := service.Create(newCampaign)

	assert.Nil(err)

	repoMock.AssertExpectations(t)

}

func TestSaveCampaignError(t *testing.T) {
	assert := assert.New(t)

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("Save", mock.Anything).Return(errors.New("error"))

	errorService := NewService(errorRepoMock)

	_, err := errorService.Create(newCampaign)

	assert.Equal(err, errors.New("error"))
}

func TestCreateCampaignError(t *testing.T) {
	assert := assert.New(t)
	newCampaign := contracts.NewCampaign{
		Name:    "test",
		Content: "content",
		Emails:  []string{"t@.com"},
	}

	_, err := service.Create(newCampaign)

	assert.Equal(err, domains.ErrInvalidEmail)

}
