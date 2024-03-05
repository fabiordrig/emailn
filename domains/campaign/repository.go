package campaign

type Repository interface {
	Save(campaign *Campaign) error
	FindAll() ([]Campaign, error)
	FindByID(id string) (*Campaign, error)
}
