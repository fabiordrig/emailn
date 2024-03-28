package campaign_test

import (
	"emailn/constants"
	"emailn/contracts"
	"emailn/domains"
	"emailn/domains/campaign"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(campaign *campaign.Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

func (m *MockRepository) FindAll() ([]campaign.Campaign, error) {
	args := m.Called()
	return args.Get(0).([]campaign.Campaign), args.Error(1)
}

func (m *MockRepository) FindByID(id string) (*campaign.Campaign, error) {
	args := m.Called(id)
	return args.Get(0).(*campaign.Campaign), args.Error(1)
}

func (m *MockRepository) Update(campaign *campaign.Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

func (m *MockRepository) Delete(campaign *campaign.Campaign) error {
	args := m.Called(campaign)
	return args.Error(0)
}

var (
	newCampaign = contracts.NewCampaign{
		Name:      "test",
		Content:   "content",
		Emails:    []string{"test@t.com"},
		CreatedBy: "test@t.com",
	}
	repoMock = new(MockRepository)
	service  = campaign.NewService(repoMock)
)

func TestNewCreateCampaign(t *testing.T) {
	assert := assert.New(t)

	repoMock.On("Create", mock.Anything).Return(nil)

	campaign, err := service.Create(newCampaign)

	assert.Nil(err)
	assert.NotNil(campaign.ID)
	assert.Equal(newCampaign.Name, campaign.Name)
	assert.Equal(newCampaign.Content, campaign.Content)
	assert.Len(campaign.Contacts, 1)
	assert.Equal(newCampaign.Emails[0], campaign.Contacts[0].Email)

}

func TestCreateCampaign(t *testing.T) {
	assert := assert.New(t)

	repoMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {

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

func TestCreateCampaignError(t *testing.T) {
	assert := assert.New(t)

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("Create", mock.Anything).Return(errors.New("error"))

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

	mock := contracts.NewCampaign{
		Name:      fake.Lorem().Text(10),
		Content:   fake.Lorem().Text(400),
		Emails:    []string{fake.Internet().Email()},
		CreatedBy: fake.Internet().Email(),
	}

	_, err := service.Create(mock)

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

func TestFindAllCampaignCorrect(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaigns := []campaign.Campaign{
		{
			ID:      uuid.New(),
			Name:    fake.Lorem().Text(10),
			Content: fake.Lorem().Text(400),
			Contacts: []campaign.Contact{
				{
					Email: fake.Internet().Email(),
				},
			},
		},
	}

	repoMock.On("FindAll").Return(mockedCampaigns, nil)

	campaigns, err := service.FindAll()

	assert.Nil(err)

	assert.Len(campaigns, 1)

}

func TestFindAllCampaignError(t *testing.T) {
	assert := assert.New(t)

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindAll").Return([]campaign.Campaign{}, errors.New("error"))

	errorService := campaign.NewService(errorRepoMock)

	_, err := errorService.FindAll()

	assert.Equal(err, errors.New("error"))
}

func TestFindCampaignByIDCorrect(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	repoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)

	campaign, err := service.FindByID(uuid.New().String())

	assert.Nil(err)

	assert.Equal(mockedCampaign, *campaign)

}

func TestFindCampaignByIDError(t *testing.T) {
	assert := assert.New(t)

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindByID", mock.Anything).Return(&campaign.Campaign{}, errors.New("error"))

	errorService := campaign.NewService(errorRepoMock)

	_, err := errorService.FindByID(uuid.New().String())

	assert.Equal(err, errors.New("error"))
}

func TestCancelCampaignCorrect(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}
	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)

	errorRepoMock.On("Update", mock.Anything).Return(nil)

	errorService := campaign.NewService(errorRepoMock)

	err := errorService.Cancel("id")

	assert.Nil(err)
}

func TestCancelCampaignFindByIdError(t *testing.T) {
	assert := assert.New(t)

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindByID", mock.Anything).Return(&campaign.Campaign{}, errors.New("error"))

	errorService := campaign.NewService(errorRepoMock)

	err := errorService.Cancel(uuid.New().String())

	assert.Equal(err, constants.ErrNotFound)
}

func TestCancelCampaignStatusError(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.IN_PROGRESS,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	repoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)
	repoMock.On("Update", mock.Anything).Return(errors.New("error"))

	err := service.Cancel(uuid.New().String())

	assert.Equal(err, constants.ErrUnprocessableEntity)
}

func TestCancelCampaignError(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)
	errorRepoMock.On("Update", mock.Anything).Return(errors.New("error"))

	errorService := campaign.NewService(errorRepoMock)

	err := errorService.Cancel(uuid.New().String())

	assert.Equal(err, constants.ErrInternalServer)
}

func TestDeleteCampaignCorrect(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	newRepoMock := new(MockRepository)
	newRepoMock.On("FindByID", mockedCampaign.ID.String()).Return(&mockedCampaign, nil)
	newRepoMock.On("Delete", &mockedCampaign).Return(nil)

	newService := campaign.NewService(newRepoMock)

	err := newService.Delete(mockedCampaign.ID.String())

	assert.Nil(err)
}

func TestDeleteCampaignFindByIdError(t *testing.T) {
	assert := assert.New(t)

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindByID", mock.Anything).Return(&campaign.Campaign{}, errors.New("error"))

	errorService := campaign.NewService(errorRepoMock)

	err := errorService.Delete(uuid.New().String())

	assert.Equal(err, constants.ErrNotFound)
}

func TestDeleteCampaignStatusWhenCampaignIsInProgress(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.IN_PROGRESS,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)

	errorService := campaign.NewService(errorRepoMock)

	err := errorService.Delete(uuid.New().String())

	assert.Equal(err, constants.ErrUnprocessableEntity)
}

func TestDeleteCampaignError(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	errorRepoMock := new(MockRepository)
	errorRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)
	errorRepoMock.On("Delete", mock.Anything).Return(errors.New("error"))

	errorService := campaign.NewService(errorRepoMock)

	err := errorService.Delete(uuid.New().String())

	assert.Equal(err, constants.ErrInternalServer)
}

func TestStartCampaign(t *testing.T) {
	assert := assert.New(t)

	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	newRepoMock := new(MockRepository)
	newRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)
	newRepoMock.On("Update", mock.Anything).Return(nil)

	newService := campaign.NewService(newRepoMock)

	err := newService.Start(mockedCampaign.ID.String())

	assert.Nil(err)
}

func TestStartCampaignError(t *testing.T) {
	assert := assert.New(t)

	newRepoMock := new(MockRepository)
	newRepoMock.On("FindByID", mock.Anything).Return(&campaign.Campaign{}, errors.New("error"))

	newService := campaign.NewService(newRepoMock)

	err := newService.Start("1")

	assert.Equal(err, constants.ErrNotFound)
}

func TestStartCampaignErrorStatus(t *testing.T) {
	assert := assert.New(t)

	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.IN_PROGRESS,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	newRepoMock := new(MockRepository)
	newRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)

	newService := campaign.NewService(newRepoMock)

	err := newService.Start(mockedCampaign.ID.String())

	assert.Equal(err, constants.ErrUnprocessableEntity)
}

func TestStartCampaignErrorUpdate(t *testing.T) {
	assert := assert.New(t)

	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	newRepoMock := new(MockRepository)
	newRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)
	newRepoMock.On("Update", mock.Anything).Return(errors.New("error"))

	newService := campaign.NewService(newRepoMock)

	err := newService.Start(mockedCampaign.ID.String())

	assert.Equal(err, errors.New("error"))
}

func TestShouldSendEmail(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	newRepoMock := new(MockRepository)
	newRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)

	newService := campaign.NewService(newRepoMock)

	err := newService.SendEmail(&mockedCampaign)

	assert.Nil(err)
}

func TestShouldSendEmailError(t *testing.T) {
	assert := assert.New(t)
	fake := faker.New()

	mockedCampaign := campaign.Campaign{
		ID:      uuid.New(),
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(400),
		Status:  campaign.PENDING,
		Contacts: []campaign.Contact{
			{
				Email: fake.Internet().Email(),
			},
		},
	}

	newRepoMock := new(MockRepository)
	newRepoMock.On("FindByID", mock.Anything).Return(&mockedCampaign, nil)

	newService := campaign.NewService(newRepoMock)
	newService.SendEmail = func(campaign *campaign.Campaign) error {
		return errors.New("error")
	}

	err := newService.SendEmail(&mockedCampaign)

	assert.NotNil(err)
}
