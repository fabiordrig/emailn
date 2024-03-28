package campaign

import (
	"emailn/constants"
	"emailn/contracts"
)

type Service interface {
	Create(newCampaign contracts.NewCampaign) (*Campaign, error)
	FindAll() ([]Campaign, error)
	FindByID(id string) (*Campaign, error)
	Cancel(id string) error
	Delete(id string) error
}

type EmailSender interface {
	SendEmail(campaign *Campaign) error
}
type ServiceImp struct {
	Repository  Repository
	EmailSender EmailSender
}

func NewService(repo Repository, emailSender EmailSender) *ServiceImp {
	return &ServiceImp{Repository: repo, EmailSender: emailSender}
}

func (s *ServiceImp) Create(newCampaign contracts.NewCampaign) (*Campaign, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)

	if err != nil {
		return nil, err
	}

	err = s.Repository.Create(campaign)

	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *ServiceImp) FindAll() ([]Campaign, error) {
	return s.Repository.FindAll()
}

func (s *ServiceImp) FindByID(id string) (*Campaign, error) {
	return s.Repository.FindByID(id)
}

func (s *ServiceImp) Cancel(id string) error {

	existentCampaign, err := s.Repository.FindByID(id)

	if err != nil {
		return constants.ErrNotFound
	}

	if existentCampaign.Status != PENDING {
		return constants.ErrUnprocessableEntity
	}
	existentCampaign.Cancel()
	err = s.Repository.Update(existentCampaign)

	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

func (s *ServiceImp) Delete(id string) error {
	campaign, err := s.Repository.FindByID(id)

	if err != nil {
		return constants.ErrNotFound
	}

	if campaign.Status != PENDING {
		return constants.ErrUnprocessableEntity
	}
	campaign.Delete()
	err = s.Repository.Delete(campaign)

	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}

func (s *ServiceImp) Start(id string) error {

	campaign, err := s.Repository.FindByID(id)

	if err != nil {
		return constants.ErrNotFound
	}
	if campaign.Status != PENDING {
		return constants.ErrUnprocessableEntity
	}

	campaign.Status = DONE

	err = s.EmailSender.SendEmail(campaign)

	if err != nil {
		return constants.ErrInternalServer
	}

	err = s.Repository.Update(campaign)

	return err
}
