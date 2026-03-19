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
	db *pgxpool.Pool
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
		attributeQuery := `insert into category_attributes (attribute_id,category_id,is_required) values ($1,$2,$3)`
		_, err := tx.Exec(ctx, attributeQuery, attribute.AttributeID, category.ID, attribute.IsRequired)
		if err != nil {
			return fmt.Errorf("nitelik (%v) eklenirken hata: %w", attribute.AttributeID, err)
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

	attributesQuery := `select ca.attribute_id ,ca.is_required ,a.id ,a."name" ,a.code ,a.data_type 
								from category_attributes ca
								join "attributes" a  on a.id  = ca.attribute_id 
								where ca.category_id =$1`

	rows, err := c.db.Query(ctx, attributesQuery, category.ID)
	if err != nil {
		return nil, fmt.Errorf("kategori %v attributeleri sorgulanırken hata olustu %w", category.Name, err)
	}
	defer rows.Close()

	var catAttrs []domain.CategoryAttribute
	for rows.Next() {
		var cat_attr domain.CategoryAttribute
		err = rows.Scan(
			&cat_attr.AttributeID,
			&cat_attr.IsRequired,
			&cat_attr.Attribute.ID,
			&cat_attr.Attribute.Name,
			&cat_attr.Attribute.Code,
			&cat_attr.Attribute.DataType)
		if err != nil {
			return nil, fmt.Errorf("%v idli attribute %v kateogorisine eklenirken olusan hata %w", cat_attr.AttributeID, category.Name, err)
		}
		catAttrs = append(catAttrs, cat_attr)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("nitelik dongusunde hata %w", err)
	}
	category.Attributes = catAttrs
	return category, nil
}

func (c *CategoryRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Category, error) {
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
