package inventory_usecase

import (
	"context"
	"eav-intentory/internal/inventory/domain"
)

type WarehouseUsecase interface {
	CreateWarehouse(ctx context.Context, warehouse *domain.Warehouse) (int, error)
	GetWarehouses(ctx context.Context) ([]domain.Warehouse, error)
	GetWarehouseById(ctx context.Context, id int) (*domain.Warehouse, error)
	UpdateWarehouse(ctx context.Context, id int, warehouse domain.Warehouse)
	DeleteWarehouse(ctx context.Context, id int) error
}

type warehouseUsecase struct {
	warehouseRepository domain.WarehouseRepository
}

func (w warehouseUsecase) GetWarehouses(ctx context.Context) ([]domain.Warehouse, error) {
	//TODO implement me
	panic("implement me")
}

func (w warehouseUsecase) GetWarehouseById(ctx context.Context, id int) (*domain.Warehouse, error) {
	//TODO implement me
	panic("implement me")
}

func (w warehouseUsecase) UpdateWarehouse(ctx context.Context, id int, warehouse domain.Warehouse) {
	//TODO implement me
	panic("implement me")
}

func (w warehouseUsecase) DeleteWarehouse(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewWarehouseUsecase(repository domain.WarehouseRepository) WarehouseUsecase {
	return &warehouseUsecase{warehouseRepository: repository}
}

func (w warehouseUsecase) CreateWarehouse(ctx context.Context, warehouse *domain.Warehouse) (int, error) {

	id, err := w.warehouseRepository.Create(ctx, warehouse)
	if err != nil {
		return 0, err
	}

	return id, nil
}
