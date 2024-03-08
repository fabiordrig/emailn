package campaign

import (
	"emailn/domains"
	"time"

	"github.com/google/uuid"
)

const (
	PENDING     = "PENDING"
	IN_PROGRESS = "IN_PROGRESS"
	CANCELED    = "CANCELED"
	DONE        = "DONE"
)

type Contact struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" validate:"required"`
	Email      string    `gorm:"size:60" validate:"required,email"`
	CampaignID uuid.UUID
}

// Campaign is a struct that represents a campaign
type Campaign struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" validate:"required"`
	Name      string    `gorm:"size:60" validate:"required,min=2,max=50"`
	Status    string    `gorm:"size:20" validate:"required,oneof=PENDING IN_PROGRESS CANCELED DONE"`
	CreatedAt time.Time `validate:"required"`
	Content   string    `gorm:"size:2000" validate:"required,min=2,max=500"`
	Contacts  []Contact `validate:"min=1,dive"`
}

func NewCampaign(name, content string, emails []string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))
	id := uuid.New()

	for i, email := range emails {
		contacts[i] = Contact{Email: email, ID: uuid.New(), CampaignID: id}

	}

	campaign := &Campaign{
		ID:        id,
		Name:      name,
		Status:    PENDING,
		CreatedAt: time.Now(),
		Content:   content,
		Contacts:  contacts,
	}

	err := domains.ValidateStruct(campaign)

	if err != nil {
		return nil, err
	}

	return campaign, nil
}
