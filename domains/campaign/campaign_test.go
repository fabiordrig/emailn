package campaign_test

import (
	"emailn/domains/campaign"
	"testing"
	"time"

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

	campaign, _ := campaign.NewCampaign(name, content, emails)

	assert.NotNil(campaign.ID)
	assert.Equal(name, campaign.Name)
	assert.Equal(content, campaign.Content)
	assert.Greater(campaign.CreatedAt, now)
	assert.Len(campaign.Contacts, 2)
	assert.Equal(emails[0], campaign.Contacts[0].Email)
	assert.Equal(emails[1], campaign.Contacts[1].Email)

}

func TestNewCampaignEmailError(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, []string{})

	assert.Nil(campaign)
	assert.Equal(err.Error(), "at least one email is required")

}

func TestNewCampaignNameContentError(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign("", "", emails)

	assert.Nil(campaign)
	assert.Equal(err.Error(), "name and content are required")

}

func TestNewCampaignInvalidEmailError(t *testing.T) {

	assert := assert.New(t)

	campaign, err := campaign.NewCampaign(name, content, []string{"invalid"})

	assert.Nil(campaign)
	assert.Equal(err.Error(), "invalid email")

}
