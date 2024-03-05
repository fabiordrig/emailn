package campaign

import "emailn/contracts"

type Service interface {
	Create(newCampaign contracts.NewCampaign) (*Campaign, error)
	FindAll() ([]Campaign, error)
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
