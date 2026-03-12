package postgres

import (
	"context"
	"eav-intentory/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool // sql.db vermemiz gerkmiyormuydu ya ben burdaki pool muhabbetini falan tam anlamadım
}

func NewCategoryRepository(pool *pgxpool.Pool) domain.CategoryRepository {
	return &CategoryRepository{db: pool}
}

func (c *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "insert into categories (name,parent_id) values ($1,$2) returning id",
		category.Name,
		category.ParentID).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("kategori eklenirken hata: %w", err)
	}

	for _, attribute := range category.Attributes {
		attributeQuery := `insert into category_attributes (category_id,name,datatype,isrequired) values ($1,$2,$3,$4)`
		_, err := tx.Exec(ctx, attributeQuery, category.ID, attribute.Name, attribute.DataType, attribute.IsRequired)
		if err != nil {
			return fmt.Errorf("nitelik (%s) eklenirken hata: %w", attribute.Name, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("transaction commit edilemedi: %w", err)
	}
	return nil
}

func (c *CategoryRepository) GetById(ctx context.Context, id int) (*domain.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CategoryRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	//TODO implement me
	panic("implement me")
}

func (c *CategoryRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Category, error) {
	//TODO implement me
	panic("implement me")
}
