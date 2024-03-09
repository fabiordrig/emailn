package database

import (
	"emailn/domains/campaign"

	"gorm.io/gorm"
)

type CampaignRepository struct {
	Db *gorm.DB
}

func (c *CampaignRepository) Save(campaign *campaign.Campaign) error {
	tx := c.Db.Save(campaign)
	return tx.Error
}

func (c *CampaignRepository) FindAll() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign
	tx := c.Db.Preload("Contacts").Find(&campaigns)

	return campaigns, tx.Error
}

func (c *CampaignRepository) FindByID(id string) (*campaign.Campaign, error) {

	var campaign campaign.Campaign

	tx := c.Db.Preload("Contacts").First(&campaign, "id = ?", id)

	return &campaign, tx.Error

}

func (c *CampaignRepository) Update(campaign *campaign.Campaign) error {
	tx := c.Db.Save(campaign)
	return tx.Error
}
