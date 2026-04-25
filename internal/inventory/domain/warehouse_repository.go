package domain

import "context"

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *Warehouse) (int, error)
	// GetAll ilerde pagination ve filter eklenecek
	GetAll(ctx context.Context) ([]Warehouse, error)
	GetById(ctx context.Context, id int) (*Warehouse, error)
	Update(ctx context.Context, warehouse *Warehouse) error
	Delete(ctx context.Context, id int) error
}
