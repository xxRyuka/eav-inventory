package postgres

import (
	"context"
	"eav-intentory/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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
		attributeQuery := `insert into category_attributes (category_id,name,data_type,is_required) values ($1,$2,$3,$4)`
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

	category := &domain.Category{}
	categoryQuery := `select id, name, parent_id from categories where id=$1`
	// parent_id veritabanında NULL olabileceği için onu bir pointer ile karşılıyoruz.
	var parentID *int
	row := c.db.QueryRow(ctx, categoryQuery, id)
	err := row.Scan(&category.ID, &category.Name, &parentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Kategori bulunamadı (ID :%v)", id)
		}
		return nil, fmt.Errorf("Kategori sorgulanırken olusan hata : %w", err)
	}
	//if parentID != nil {

	category.ParentID = parentID
	//}

	//attributelerini cekeceğiz simdi

	attributesQuery := `select id ,name, data_type ,is_required from category_attributes where category_id=$1`

	rows, err := c.db.Query(ctx, attributesQuery, category.ID)
	if err != nil {
		return nil, fmt.Errorf("Kategori attributeleri sorgularınrken cıkan hata : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var attr domain.CategoryAttribute
		err := rows.Scan(&attr.ID, &attr.Name, &attr.DataType, &attr.IsRequired)
		if err != nil {
			return nil, fmt.Errorf("Kategori attributeleri scan ile yerleştirilirken cıkan hata : %w", err)
		}

		category.Attributes = append(category.Attributes, attr)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("nitelik dongusunde hata %w", err)
	}
	return category, nil
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
