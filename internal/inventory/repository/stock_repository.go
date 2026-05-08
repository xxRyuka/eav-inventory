package inventory_repository

import (
	"context"
	"eav-intentory/internal/inventory/domain"
	"eav-intentory/internal/shared/postgres_tx_manager"
	"fmt"
)

type StockRepository struct {
	db postgres_tx_manager.DbExecutor // Neden * vermiyoruz burda ? pgx verirken pointer olarak veriyorduk onu neden pointer verdikte bunu direk veriyoruz ?
}

func NewStockRepository(db postgres_tx_manager.DbExecutor) domain.StockRepository {
	return &StockRepository{db: db}
}

// upsert aslında tek satır kodla olan bireşy değil ekleme işleminde nerde hata olacaksa ona göre koşul yazıyoruz !
func (s StockRepository) CreateOrIncrease(ctx context.Context, warehouseId int, productId int, quantity int) error {

	query := `insert into stocks (warehouse_id,product_id,available_quantity) values ($1,$2,$3) 
				on conflict (warehouse_id, product_id)
				do update set available_quantity = excluded.available_quantity + stocks.available_quantity `

	exec, err := s.db.Exec(ctx, query, warehouseId, productId, quantity)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return fmt.Errorf("no rowms affected ")
	}

	return nil
}
