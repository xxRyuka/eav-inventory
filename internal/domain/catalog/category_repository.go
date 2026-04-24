package catalog

import "context"

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetById(ctx context.Context, id int) (*Category, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, category *Category) error
	// GetAll Tekil nesne dönerken (örn: GetById) pointer dönülür (*Product), liste dönerken value dönülür ([]Product).
	GetAll(ctx context.Context, limit, offset int) ([]Category, int, error)
	// GetPaged veya Filtreleme burda mı olmalı ? veya getAll'e mi koymaliyiz : getCategories içinde query ve path parametrelerinden erişmek en mantıklısı diye düsündüm
	// atomic bir endpoint ile ihtiyaca gore polimorfizm sagladım

	//TODO: usecase and handler not implemented
	GetCategoriesWithAttirbutes(ctx context.Context) ([]Category, error)
	AddAttributeToCategory(ctx context.Context, categoryID, attributeID int, isRequired bool) error

	UpdateAttributeToCategory(ctx context.Context, isRequired bool, attributeID, categoryID int) error
	RemoveAttributeToCategory(ctx context.Context, categoryID, attributeID int) error
}
