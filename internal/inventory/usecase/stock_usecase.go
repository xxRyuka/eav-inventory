package inventory_usecase

import "context"

type StockUsecase interface {
	PurchaseIn(ctx context.Context)
}

type StockUseCase struct {
}

func NewStockUseCase() {

}

func (c *StockUseCase) PurchaseIn(ctx context.Context) {

}
