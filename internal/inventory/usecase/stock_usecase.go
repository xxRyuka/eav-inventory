package inventory_usecase

import (
	"context"
	"eav-intentory/internal/inventory/domain"
	inventory_repository "eav-intentory/internal/inventory/repository"
	"eav-intentory/internal/inventory/usecase/command"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockUsecase interface {
	PurchaseIn(ctx context.Context, command command.PurchaseInCommand) error
}

type StockUseCase struct {
	db *pgxpool.Pool // sadece tx acmak icin kullanacağız !
}

func NewStockUseCase(db *pgxpool.Pool) StockUsecase {
	return &StockUseCase{
		db: db,
	}
}

func (c *StockUseCase) PurchaseIn(ctx context.Context, command command.PurchaseInCommand) error {
	// 1. Önce usecase seviyesinde temel iş kuralı kontrol edilir.
	//
	// Quantity 0 veya negatif olamaz.
	// Çünkü stok girişi pozitif bir miktar olmalıdır.
	if command.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	movement := domain.StockMovement{
		WarehouseID:  command.WarehouseID,
		ProductID:    command.ProductID,
		Quantity:     command.Quantity,
		MovementType: domain.PurchaseIn,
	}

	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	stockMovementRepo := inventory_repository.NewStockMovementRepository(tx)
	err = stockMovementRepo.Create(ctx, movement)
	if err != nil {
		return fmt.Errorf("stock movement create : %w", err)
	}

	stockRepo := inventory_repository.NewStockRepository(tx)

	err = stockRepo.CreateOrIncrease(ctx, movement.WarehouseID, movement.ProductID, movement.Quantity)
	if err != nil {
		return fmt.Errorf("stock upsert : %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit err : %w", err)
	}
	return nil
}
