package campaign_test

import (
	"emailn/constants"
	"emailn/contracts"
	"emailn/domains"
	"emailn/domains/campaign"
	"errors"
	"testing"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(campaign *campaign.Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

func (m *MockRepository) FindAll() ([]campaign.Campaign, error) {
	args := m.Called()
	return args.Get(0).([]campaign.Campaign), nil
}

var (
	newCampaign = contracts.NewCampaign{
		Name:    "test",
		Content: "content",
		Emails:  []string{"test@t.com"},
	}
	repoMock = new(MockRepository)
	service  = campaign.NewService(repoMock)
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

	repoMock.On("Save", mock.MatchedBy(func(campaign *campaign.Campaign) bool {

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

	errorService := campaign.NewService(errorRepoMock)

	_, err := errorService.Create(newCampaign)

	assert.Equal(err, errors.New("error"))
}

func TestCreateCampaignEmailError(t *testing.T) {
	assert := assert.New(t)
	newCampaign := contracts.NewCampaign{
		Name:    "test",
		Content: "content",
		Emails:  []string{"t@.com"},
	}

	_, err := service.Create(newCampaign)

	assert.Equal(err, constants.ErrInvalidEmail)

}

func TestCreateCampaignNameContentError(t *testing.T) {
	assert := assert.New(t)
	newCampaign := contracts.NewCampaign{
		Name:    "1",
		Content: "1",
		Emails:  []string{"t@t.com"},
	}

	fake := faker.New()

	_, err := service.Create(newCampaign)

	assert.Equal(err, constants.ErrStringMinLength)

	newCampaign = contracts.NewCampaign{
		Name:    fake.Lorem().Text(10),
		Content: "1",
		Emails:  []string{"@t.com"},
	}

	_, err = service.Create(newCampaign)

	assert.Equal(err, constants.ErrStringMinLength)

	newCampaign = contracts.NewCampaign{
		Name:    "1",
		Content: fake.Lorem().Text(10),
		Emails:  []string{"@t.com"},
	}

	_, err = service.Create(newCampaign)

	assert.Equal(err, constants.ErrStringMinLength)
}

func TestCreateCampaignMaxNameContentError(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()
	newCampaign := contracts.NewCampaign{
		Name:    fake.Lorem().Text(102),
		Content: fake.Lorem().Text(501),
		Emails:  []string{"test@test.com"},
	}

	_, err := service.Create(newCampaign)

	assert.Equal(err, constants.ErrStringMaxLength)

	newCampaign = contracts.NewCampaign{
		Name:    fake.Lorem().Text(101),
		Content: fake.Lorem().Text(400),
		Emails:  []string{"test@test.com"},
	}

	_, err = service.Create(newCampaign)

	assert.Equal(err, constants.ErrStringMaxLength)

	newCampaign = contracts.NewCampaign{
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(503),
		Emails:  []string{"test@test.com"},
	}

	_, err = service.Create(newCampaign)

	assert.Equal(err, constants.ErrStringMaxLength)

}

func TestCreateCampaignEmailsError(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	newCampaign := contracts.NewCampaign{
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Emails:  []string{},
	}

	_, err := service.Create(newCampaign)

	assert.Equal(err, constants.ErrStringMinLength)

}

func TestCreateCampaignRequiredField(t *testing.T) {

	assert := assert.New(t)

	newCampaign := contracts.NewCampaign{

		Emails: []string{"test@test.com"},
	}

	_, err := service.Create(newCampaign)

	assert.Equal(err, constants.ErrRequiredField)
}

func TestCreateCampaignValidateCorrect(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	newCampaign := contracts.NewCampaign{
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Emails:  []string{fake.Internet().Email()},
	}

	_, err := service.Create(newCampaign)

	assert.Nil(err)
}

func TestValidateStructUnknownTagError(t *testing.T) {
	assert := assert.New(t)

	type TestData struct {
		Field1 string `validate:"oneof=1 2 3"`
	}

	obj := TestData{
		Field1: "4",
	}

	err := domains.ValidateStruct(obj)

	assert.Equal(constants.ErrUnknown, err)
}
