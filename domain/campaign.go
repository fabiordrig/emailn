package domain

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Email string
}

// Campaign is a struct that represents a campaign
type Campaign struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	Content   string
	Contacts  []Contact
}

func NewCampaign(name, content string, emails []string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))

	if len(emails) == 0 {
		return nil, errors.New("at least one email is required")
	}

	if name == "" || content == "" {
		return nil, errors.New("name and content are required")
	}

	emailList := strings.Join(emails, ", ")

	if _, err := mail.ParseAddressList(emailList); err != nil {
		return nil, errors.New("invalid email")
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
