package campaign

type Repository interface {
	Create(campaign *Campaign) error
	FindAll() ([]Campaign, error)
	FindByID(id string) (*Campaign, error)
	Delete(campaign *Campaign) error
	Update(campaign *Campaign) error
}
