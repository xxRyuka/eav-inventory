package domain

import "context"

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetById(ctx context.Context, id int) (*Category, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, category *Category) error
	// GetAll Tekil nesne dönerken (örn: GetById) pointer dönülür (*Product), liste dönerken value dönülür ([]Product).
	GetAll(ctx context.Context, limit, offset int) ([]Category, error)
	// GetPaged veya Filtreleme burda mı olmalı ? veya getAll'e mi koymaliyiz
}
