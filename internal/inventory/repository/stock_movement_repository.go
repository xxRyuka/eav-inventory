package inventory_repository

import (
	"context"
	"eav-intentory/internal/inventory/domain"
	"eav-intentory/internal/shared/postgres_tx_manager"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockMovementRepository struct {
	db *pgxpool.Pool
}

func NewStockMovementRepository(db *pgxpool.Pool) domain.StockMovementRepository {
	return StockMovementRepository{db: db}
}

// movementin uc katmanından domain layerda belirlediğim const enumlardan purhcase_in olarak gelmesi gerekiyor burda kontrol etmeyeceğim buranın gorevi değil !
func (s StockMovementRepository) Create(ctx context.Context, tx postgres_tx_manager.DbExecutor, movement domain.StockMovement) error {

	querySM := `insert into stock_movements (
					warehouse_id,
					product_id,
					quantity,
					movement_type)	values ($1,$2,$3,$4)`

	exec, err := tx.Exec(ctx, querySM, movement.WarehouseID, movement.ProductID, movement.Quantity, movement.MovementType)
	if err != nil {
		return err
	}

	if exec.RowsAffected() == 0 {
		return fmt.Errorf("0 rows affected ")
	}

	return nil
}
