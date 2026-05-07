package inventory_repository

import (
	"context"
	"eav-intentory/internal/inventory/domain"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WarehouseRepository struct {
	db *pgxpool.Pool
}

func NewWarehouseRepository(db *pgxpool.Pool) domain.WarehouseRepository {
	return &WarehouseRepository{db: db}
}
func (w WarehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) (int, error) {

	var id int
	query := `insert into warehouses (location,name,code) values ($1,$2,$3) returning id`

	row := w.db.QueryRow(ctx, query, warehouse.Location, warehouse.Name, warehouse.Code)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (w WarehouseRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Warehouse, int, error) {

	countQuery := "select count(*) from warehouses"
	count := 0
	err := w.db.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `select id ,name, code, location from warehouses order by id limit $1 offset $2`

	rows, err := w.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var warehouses []domain.Warehouse
	for rows.Next() {
		var warehouse domain.Warehouse // burda pointer olarak versem ne olur ki ?
		err = rows.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Code, &warehouse.Location)
		if err != nil {
			return nil, 0, err
		}
		warehouses = append(warehouses, warehouse)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return warehouses, count, nil
}

func (w WarehouseRepository) GetById(ctx context.Context, id int) (*domain.Warehouse, error) {

	query := `select id ,name, code, location from warehouses where id = $1`
	var warehouse domain.Warehouse // burda pointer olarak versem ne olur ki ?
	err := w.db.QueryRow(ctx, query, id).Scan(&warehouse.ID, &warehouse.Name, &warehouse.Code, &warehouse.Location)
	if err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (w WarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `update warehouses set name =$1,code = $2,location=  $3 where id = $4`
	// todo: handleri id'yi pathden warehouseyi bodyden cekmeli !
	exec, err := w.db.Exec(ctx, query, warehouse.Name, warehouse.Code, warehouse.Location, warehouse.ID)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return fmt.Errorf("0 rows affected")
	}

	return nil
}

func (w WarehouseRepository) Delete(ctx context.Context, id int) error {

	query := `delete from warehouses where id = $1`
	exec, err := w.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return fmt.Errorf("0 rows affected")
	}

	return nil
}
