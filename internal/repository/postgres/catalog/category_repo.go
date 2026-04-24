package catalog

import (
	"context"
	"eav-intentory/internal/domain/catalog"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) catalog.CategoryRepository {
	return &CategoryRepository{db: pool}
}
func (c *CategoryRepository) Create(ctx context.Context, category *catalog.Category) error {
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

func (c *CategoryRepository) GetById(ctx context.Context, id int) (*catalog.Category, error) {

	category := &catalog.Category{}
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

	var catAttrs []catalog.CategoryAttribute
	for rows.Next() {
		var cat_attr catalog.CategoryAttribute
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

func (c *CategoryRepository) GetAll(ctx context.Context, limit, offset int) ([]catalog.Category, int, error) {

	totalCountQuery := `select count(*) from categories `
	var totalCount int
	err := c.db.QueryRow(ctx, totalCountQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}
	if totalCount == 0 {
		return nil, 0, fmt.Errorf("no category row ")
	}

	query := `select id, name, parent_id from categories limit $1 offset $2`

	rows, err := c.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var categoryList []catalog.Category
	for rows.Next() {
		var category catalog.Category
		err = rows.Scan(&category.ID, &category.Name, &category.ParentID)
		if err != nil {
			return nil, 0, err
		}
		categoryList = append(categoryList, category)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("rows error : %w", err)
	}
	return categoryList, totalCount, nil
}

func (c *CategoryRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *CategoryRepository) Update(ctx context.Context, category *catalog.Category) error {
	query := `update categories set name= $1, parent_id=$2 where id =$3`
	exec, err := c.db.Exec(ctx, query, category.Name, category.ParentID, category.ID)
	if err != nil {
		// db.Exec() fonksiyonu asla pgx.ErrNoRows hatası fırlatmaz !!!!
		//if errors.Is(err, pgx.ErrNoRows) {
		//	return fmt.Errorf("0 rows affected")
		//}
		return err
	}
	if exec.RowsAffected() == 0 {
		return fmt.Errorf("güncellenecek kategori bulunamadı (ID: %d)", category.ID)
	}
	return nil
}

// TODO: usecase ve handler kodları yazılmadı !
func (c *CategoryRepository) GetCategoriesWithAttirbutes(ctx context.Context) ([]catalog.Category, error) {
	//once kategorileri sonra kategorilere baglı, attributeleri çekeceğim plan olarak xxx yanlıs yaklasım
	// Error handlingi en son detayli yapcam

	//burda n+1 problemi dogdugu için tek sql ile bütün verileri cekip method içinde maplicez

	//query := `select c.name ,c.id ,c.parent_id  from categories c
	//					join `

	query := `SELECT 
            c.id, c.name, c.parent_id,
            ca.attribute_id, ca.is_required,
            a.code, a.name, a.data_type
        FROM categories c
        LEFT JOIN category_attributes ca ON c.id = ca.category_id
        LEFT JOIN attributes a ON ca.attribute_id = a.id`

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Aynısı gelen kategorileri üst üste katlayacağımız (gruplayacağımız) Sözlük (Map)
	categoryMap := make(map[int]*catalog.Category) // mapler her zaman make ile olusturulur
	for rows.Next() {
		var category catalog.Category

		//var attribute domain.Attribute
		// bunlaro pointer olarak almamızın sebebi null gelebilecek olması
		var attributeID *int
		var attributeCode *string
		var attributeName *string
		var attributeDataType *string

		var isRequired *bool

		//var catAttr domain.CategoryAttribute

		// Burda catAttr üzerinedn &catAttr.xx , &catAttr.Attribute.xy seklinde almakla tek tek field acmanin farkı ne ?
		err = rows.Scan(&category.ID, &category.Name, &category.ParentID, &attributeID, &isRequired, &attributeCode, &attributeName, &attributeDataType)
		if err != nil {
			return nil, err
		}

		if _, ok := categoryMap[category.ID]; !ok {
			categoryMap[category.ID] = &catalog.Category{
				ID:         category.ID,
				Name:       category.Name,
				ParentID:   category.ParentID,
				Attributes: make([]catalog.CategoryAttribute, 0),
			}
		}

		if attributeID != nil {
			attr := catalog.CategoryAttribute{
				AttributeID: *attributeID,
				Attribute: catalog.Attribute{
					ID:       *attributeID,
					Code:     *attributeCode,
					Name:     *attributeName,
					DataType: catalog.DataType(*attributeDataType),
				},
				IsRequired: *isRequired,
			}

			categoryMap[category.ID].Attributes = append(categoryMap[category.ID].Attributes, attr) // kategoriye eristim sonrada onun attributelerine erisip ekledim
		}

	}
	var categories []catalog.Category

	for _, cat := range categoryMap {
		categories = append(categories, *cat)
	}
	return categories, nil
}

func (c *CategoryRepository) AddAttributeToCategory(ctx context.Context, categoryID, attributeID int, isRequired bool) error {

	query := `insert into category_attributes (attribute_id ,category_id,is_required) values ($1 , $2, $3)`

	exec, err := c.db.Exec(ctx, query, attributeID, categoryID, isRequired)
	if err != nil {
		return fmt.Errorf("err %w", err)
	}

	if exec.RowsAffected() == 0 {
		return fmt.Errorf("Attribute assingment failure ")
	}
	return nil
}

func (c *CategoryRepository) RemoveAttributeToCategory(ctx context.Context, categoryID, attributeID int) error {
	query := `delete from category_attributes where attribute_id=$1 and category_id =$2`

	exec, err := c.db.Exec(ctx, query, attributeID, categoryID)
	if err != nil {
		return err
	}

	if exec.RowsAffected() == 0 {
		return fmt.Errorf("remove process failed 0 row affected")
	}
	return nil
}

func (c *CategoryRepository) UpdateAttributeToCategory(ctx context.Context, isRequired bool, attributeID, categoryID int) error {

	query := `update from category_attributes , set is_required =$1 where attribute_id =$2 and category_id=$3`

	exec, err := c.db.Exec(ctx, query, isRequired, attributeID, categoryID)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return fmt.Errorf("update process failed 0 row affected")
	}
	return nil

}
