package campaign

import "emailn/contracts"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(newCampaign contracts.NewCampaign) (*Campaign, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return nil, err
	}

	err = s.repo.Save(campaign)

	if err != nil {
		return nil, err
	}

	return campaign, nil
}
