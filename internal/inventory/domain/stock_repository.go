package domain

import "context"

type StockRepository interface {
	CreateOrIncrease(ctx context.Context, warehouseId int, productId int, quantity int) error
}
