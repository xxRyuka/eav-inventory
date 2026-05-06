package inventory_usecase

import (
	"context"
	"eav-intentory/internal/inventory/domain"
	"fmt"
)

type WarehouseUsecase interface {
	CreateWarehouse(ctx context.Context, warehouse *domain.Warehouse) (int, error)
	GetWarehouses(ctx context.Context, pagesize, page int) ([]domain.Warehouse, int, error)
	GetWarehouseById(ctx context.Context, id int) (*domain.Warehouse, error)
	UpdateWarehouse(ctx context.Context, id int, warehouse *domain.Warehouse) error
	DeleteWarehouse(ctx context.Context, id int) error
}

type warehouseUsecase struct {
	warehouseRepository domain.WarehouseRepository
}

func NewWarehouseUsecase(repository domain.WarehouseRepository) WarehouseUsecase {
	return &warehouseUsecase{warehouseRepository: repository}
}

func (w warehouseUsecase) GetWarehouses(ctx context.Context, pagesize, page int) ([]domain.Warehouse, int, error) {
	offset := (page - 1) * pagesize

	warehouses, i, err := w.warehouseRepository.GetAll(ctx, pagesize, offset)
	if err != nil {
		return nil, 0, err
	}

	return warehouses, i, nil
}

func (w warehouseUsecase) GetWarehouseById(ctx context.Context, id int) (*domain.Warehouse, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	warehouse, err := w.warehouseRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

// todo : id path ile warehouse bodyden dto paketli olarak alınacak !
func (w warehouseUsecase) UpdateWarehouse(ctx context.Context, id int, warehouse *domain.Warehouse) error {
	if id <= 0 {
		return fmt.Errorf("invalid id %d", id)
	}
	// 2. İş mantığı: Güncellenecek kayıt var mı?
	existing, err := w.warehouseRepository.GetById(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("Gecerli depo bulunamadı ! (id : %d) ", id)
	}

	warehouse.ID = id
	err = w.warehouseRepository.Update(ctx, warehouse)
	if err != nil {
		return err
	}

	return nil

}

func (w warehouseUsecase) DeleteWarehouse(ctx context.Context, id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid id %d", id)
	}
	// 2. İş mantığı: Güncellenecek kayıt var mı?
	existing, err := w.warehouseRepository.GetById(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("Gecerli depo bulunamadı ! (id : %d) ", id)
	}

	err = w.warehouseRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (w warehouseUsecase) CreateWarehouse(ctx context.Context, warehouse *domain.Warehouse) (int, error) {

	id, err := w.warehouseRepository.Create(ctx, warehouse)
	if err != nil {
		return 0, err
	}

	return id, nil
}
