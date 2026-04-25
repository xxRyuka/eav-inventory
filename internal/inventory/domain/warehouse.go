package domain

type Warehouse struct {
	ID       int
	Location string
	Name     string
	Code     string // should be unique
	// ilerde is Active Alanları eklenecek ve proje geneline audit log entegreasyonu yapılacak
}
