package domain

import "context"

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetById(ctx context.Context, id int) (*Product, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, product *Product) error
	// GetAll Tekil nesne dönerken (örn: GetById) pointer dönülür (*Product), liste dönerken value dönülür ([]Product).
	GetAll(ctx context.Context, limit int, offset int) ([]Product, int, error)
	UpdateAttributes(ctx context.Context, productId int, values []ProductAttributeValue) error // yeni bir attribute listesi alıyor içindekiler yoksa eklicez id mevcutsa güncelleme yapcaz dimi i?
	SearchByAttribute(ctx context.Context, filters map[int]string) ([]Product, error)
}
