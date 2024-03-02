package campaign

import (
	"emailn/domains"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Email string `validate:"email,required"`
}

// Campaign is a struct that represents a campaign
type Campaign struct {
	ID        uuid.UUID `validate:"required"`
	Name      string    `validate:"min=2,max=50,required"`
	CreatedAt time.Time `validate:"required"`
	Content   string    `validate:"min=2,max=500,required"`
	Contacts  []Contact `validate:"min=1,required"`
}

func NewCampaign(name, content string, emails []string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))

	if len(emails) == 0 {
		return nil, domains.ErrAtLeastOneEmailIsRequired
	}

	if name == "" || content == "" {
		return nil, domains.ErrNameAndContentAreRequired
	}

	emailList := strings.Join(emails, ", ")

	if _, err := mail.ParseAddressList(emailList); err != nil {
		return nil, domains.ErrInvalidEmail
	}

	for i, email := range emails {
		contacts[i] = Contact{Email: email}
	}

	return &Campaign{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		Content:   content,
		Contacts:  contacts,
	}, nil
}
