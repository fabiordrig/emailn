package campaign

import "emailn/contracts"

type Service struct {
	Repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repository: repo}
}

func (s *Service) Create(newCampaign contracts.NewCampaign) (*Campaign, error) {
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

func (s *Service) FindAll() ([]Campaign, error) {
	return s.Repository.FindAll()
}
