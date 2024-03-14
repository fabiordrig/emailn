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
	name      = "test"
	content   = "content"
	emails    = []string{"test@email.com", "test2@email.com"}
	createdBy = "test@t.com"
)

func TestNewCampaign(t *testing.T) {

	assert := assert.New(t)

	now := time.Now().Add(-time.Minute)

	campaign, err := campaign.NewCampaign(name, content, emails, createdBy)

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

	campaign, err := campaign.NewCampaign(name, content, []string{"invalidEmail"}, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrInvalidEmail)

}

func TestNewCampaignNameContentError(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign("1", "1", emails, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)

}

func TestNewCampaignMinNameContentError2(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign("a", "", emails, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)

}

func TestNewCampaignMinNameContentError3(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign("a", "a", emails, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)
}

func TestNewCampaignMaxNameContentError(t *testing.T) {

	assert := assert.New(t)
	fake := faker.New()

	campaign, err := campaign.NewCampaign(fake.Lorem().Text(102), fake.Lorem().Text(501), emails, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMaxLength)
}

func TestNewCampaignMaxNameContentError2(t *testing.T) {

	assert := assert.New(t)
	fake := faker.New()

	campaign, err := campaign.NewCampaign(name, fake.Lorem().Text(503), emails, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMaxLength)
}

func TestNewCampaignMaxNameContentError3(t *testing.T) {

	assert := assert.New(t)
	fake := faker.New()

	campaign, err := campaign.NewCampaign(fake.Lorem().Text(102), content, emails, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMaxLength)
}

func TestNewCampaignInvalidEmailError(t *testing.T) {
	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, []string{"invalidEmail"}, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrInvalidEmail)
}

func TestNewCampaignInvalidEmailError2(t *testing.T) {
	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, []string{}, createdBy)

	assert.Nil(campaign)
	assert.Equal(err, constants.ErrStringMinLength)
}

func TestNewCampaignShouldCreate(t *testing.T) {
	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, emails, createdBy)

	assert.Nil(err)
	assert.NotNil(campaign)
	assert.Equal("PENDING", campaign.Status)
	assert.Equal(name, campaign.Name)
	assert.Equal(content, campaign.Content)
	assert.Len(campaign.Contacts, 2)
	assert.Equal(emails[0], campaign.Contacts[0].Email)
	assert.Equal(emails[1], campaign.Contacts[1].Email)
	assert.Equal(createdBy, campaign.CreatedBy)
}

func TestShouldCancelCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(name, content, emails, createdBy)

	campaign.Cancel()

	assert.Equal("CANCELLED", campaign.Status)
}

func TestShouldDeleteCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := campaign.NewCampaign(name, content, emails, createdBy)

	campaign.Delete()

	assert.Equal("DELETED", campaign.Status)
}
