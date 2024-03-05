package campaign

type Repository interface {
	Save(campaign *Campaign) error
	FindAll() ([]Campaign, error)
}
