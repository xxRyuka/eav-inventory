package domain

import (
	"context"
	"eav-intentory/internal/shared/postgres_tx_manager"
)

type StockMovementRepository interface {
	Create(ctx context.Context, tx postgres_tx_manager.DbExecutor, movement StockMovement) error
}
