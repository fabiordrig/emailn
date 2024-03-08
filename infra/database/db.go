package database

import (
	"emailn/domains/campaign"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb() *gorm.DB {

	dsn := "host=localhost user=root password=root dbname=emailn port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&campaign.Campaign{}, &campaign.Contact{})
	return db
}
