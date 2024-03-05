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
	Email string `validate:"required,email"`
}

// Campaign is a struct that represents a campaign
type Campaign struct {
	ID        uuid.UUID `validate:"required"`
	Name      string    `validate:"required,min=2,max=50"`
	Status    string    `validate:"required,oneof=PENDING IN_PROGRESS CANCELED DONE"`
	CreatedAt time.Time `validate:"required"`
	Content   string    `validate:"required,min=2,max=500"`
	Contacts  []Contact `validate:"min=1,dive"`
}

func NewCampaign(name, content string, emails []string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))

	for i, email := range emails {
		contacts[i] = Contact{Email: email}
	}

	campaign := &Campaign{
		ID:        uuid.New(),
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
