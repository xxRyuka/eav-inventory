package inventory_repository

import (
	"context"
	"eav-intentory/internal/inventory/domain"
	"eav-intentory/internal/shared/postgres_tx_manager"
	"fmt"
)

type StockMovementRepository struct {
	db postgres_tx_manager.DbExecutor
}

func NewStockMovementRepository(db postgres_tx_manager.DbExecutor) domain.StockMovementRepository {
	return &StockMovementRepository{db: db}
}

// movementin uc katmanından domain layerda belirlediğim const enumlardan purhcase_in olarak gelmesi gerekiyor burda kontrol etmeyeceğim buranın gorevi değil !
func (s StockMovementRepository) Create(ctx context.Context, movement domain.StockMovement) error {

	querySM := `insert into stock_movements (
					warehouse_id,
					product_id,
					quantity,
					movement_type)	values ($1,$2,$3,$4)`

	exec, err := s.db.Exec(ctx, querySM, movement.WarehouseID, movement.ProductID, movement.Quantity, movement.MovementType)
	if err != nil {
		return err
	}

	if exec.RowsAffected() == 0 {
		return fmt.Errorf("0 rows affected ")
	}

	return nil
}
