package inventory_usecase

import (
	"context"
	"eav-intentory/internal/inventory/domain"
)

type WarehouseUsecase interface {
	CreateWarehouse(ctx context.Context, warehouse *domain.Warehouse) (int, error)
}

type warehouseUsecase struct {
	warehouseRepository domain.WarehouseRepository
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
