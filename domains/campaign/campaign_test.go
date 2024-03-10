package campaign_test

import (
	"emailn/constants"
	"emailn/domains/campaign"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var (
	name    = "test"
	content = "content"
	emails  = []string{"test@email.com", "test2@email.com"}
)

func TestNewCampaign(t *testing.T) {

	assert := assert.New(t)

	now := time.Now().Add(-time.Minute)

	campaign, err := campaign.NewCampaign(name, content, emails)

	assert.NotNil(campaign.ID)
	assert.Nil(err)
	assert.Equal(name, campaign.Name)
	assert.Equal(content, campaign.Content)
	assert.Greater(campaign.CreatedAt, now)
	assert.Equal("PENDING", campaign.Status)
	assert.Len(campaign.Contacts, 2)
	assert.Equal(emails[0], campaign.Contacts[0].Email)
	assert.Equal(emails[1], campaign.Contacts[1].Email)

}

func TestNewCampaignEmailError(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, []string{"invalidEmail"})

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrInvalidEmail)

}

func TestNewCampaignNameContentError(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign("1", "1", emails)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)

}

func TestNewCampaignMinNameContentError2(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign("a", "", emails)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)

}

func TestNewCampaignMinNameContentError3(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign("a", "a", emails)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)
}

func TestNewCampaignMaxNameContentError(t *testing.T) {

	assert := assert.New(t)
	fake := faker.New()

	campaign, err := campaign.NewCampaign(fake.Lorem().Text(102), fake.Lorem().Text(501), emails)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMaxLength)
}

func TestNewCampaignMaxNameContentError2(t *testing.T) {

	assert := assert.New(t)
	fake := faker.New()

	campaign, err := campaign.NewCampaign(name, fake.Lorem().Text(503), emails)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMaxLength)
}

func TestNewCampaignMaxNameContentError3(t *testing.T) {

	assert := assert.New(t)
	fake := faker.New()

	campaign, err := campaign.NewCampaign(fake.Lorem().Text(102), content, emails)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMaxLength)
}

func TestNewCampaignInvalidEmailError(t *testing.T) {
	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, []string{"invalidEmail"})

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrInvalidEmail)
}

func TestNewCampaignInvalidEmailError2(t *testing.T) {
	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, []string{})

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)
}

func TestShouldCancelCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(name, content, emails)

	campaign.Cancel()

	assert.Equal("CANCELLED", campaign.Status)
}

func TestShouldDeleteCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(name, content, emails)

	campaign.Delete()

	assert.Equal("DELETED", campaign.Status)
}
