package postgres

import (
	"context"
	"eav-intentory/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: pool}
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	productQuery := `insert into products (category_id,name,sku) values ($1,$2,$3) returning id`
	err = tx.QueryRow(ctx, productQuery, p.CategoryId, p.Name, p.SKU).Scan(&p.ID)
	if err != nil {
		return err
	}

	attrributeQuery := `insert into product_attribute_values (product_id,category_attribute_id,value) values ($1,$2,$3)`

	for _, value := range p.AttributeValues {
		exec, err := tx.Exec(ctx, attrributeQuery, p.ID, value.AttributeID, value.Value)
		if err != nil {
			return fmt.Errorf("nitelik degeri (%v) eklenirken hata: %w", value.Value, err)
		}
		fmt.Println("Etkilenen satir : ", exec.RowsAffected(), " <- ")

	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("transaction commit edilemedi: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetById(ctx context.Context, id int) (*domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) GetAll(ctx context.Context, limit int, offset int) ([]domain.Product, int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) UpdateAttributes(ctx context.Context, productId int, values []domain.ProductAttributeValue) error {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) SearchByAttribute(ctx context.Context, filters map[int]string) ([]domain.Product, error) {
	//TODO implement me
	panic("implement me")
}
