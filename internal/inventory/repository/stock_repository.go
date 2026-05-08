package inventory_repository

import (
	"context"
	"eav-intentory/internal/inventory/domain"
	"eav-intentory/internal/shared/postgres_tx_manager"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockRepository struct {
	db *pgxpool.Pool
}

func NewStockRepository(db *pgxpool.Pool) domain.StockRepository {
	return &StockRepository{db: db}
}

// upsert aslında tek satır kodla olan bireşy değil ekleme işleminde nerde hata olacaksa ona göre koşul yazıyoruz !
func (s StockRepository) CreateOrIncrease(ctx context.Context, tx postgres_tx_manager.DbExecutor, warehouseId int, productId int, quantity int) error {

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
