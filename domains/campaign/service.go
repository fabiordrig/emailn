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
}

type ServiceImp struct {
	Repository Repository
}

func NewService(repo Repository) *ServiceImp {
	return &ServiceImp{Repository: repo}
}

func (s *ServiceImp) Create(newCampaign contracts.NewCampaign) (*Campaign, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return nil, err
	}

	err = s.Repository.Save(campaign)

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
	err = s.Repository.Save(existentCampaign)

	if err != nil {
		return constants.ErrInternalServer
	}

	return nil
}
