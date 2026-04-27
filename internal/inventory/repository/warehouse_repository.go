package inventory_repository

import (
	"context"
	"eav-intentory/internal/inventory/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WarehouseRepository struct {
	db *pgxpool.Pool
}

func (w WarehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) (int, error) {

	var id int
	query := `insert into warehouses (location,name,code) values ($1,$2,$3,$4) returning id`

	row := w.db.QueryRow(ctx, query, warehouse.Location, warehouse.Name, warehouse.Code)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (w WarehouseRepository) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	//TODO implement me
	panic("implement me")
}

func (w WarehouseRepository) GetById(ctx context.Context, id int) (*domain.Warehouse, error) {
	//TODO implement me
	panic("implement me")
}

func (w WarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	//TODO implement me
	panic("implement me")
}

func (w WarehouseRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
