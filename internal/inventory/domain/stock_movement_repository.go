package domain

import (
	"context"
)

type StockMovementRepository interface {
	Create(ctx context.Context, movement StockMovement) error
}
