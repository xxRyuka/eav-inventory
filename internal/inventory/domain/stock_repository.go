package domain

import (
	"context"
	"eav-intentory/internal/shared/postgres_tx_manager"
)

type StockRepository interface {
	CreateOrIncrease(ctx context.Context, tx postgres_tx_manager.DbExecutor, warehouseId int, productId int, quantity int) error
}
