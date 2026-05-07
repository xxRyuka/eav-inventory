package inventory_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockRepository struct {
	db *pgxpool.Pool
}

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
