package postgres

import (
	"context"
	"eav-intentory/internal/domain"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

	attrributeQuery := `insert into product_attribute_values (attribute_id,product_id,value) values ($1,$2,$3)`

	for _, value := range p.AttributeValues {
		exec, err := tx.Exec(ctx, attrributeQuery, value.AttributeID, p.ID, value.Value)
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

	var product domain.Product
	productQuery := `select p.id ,p."name" ,p.sku ,p.category_id from products p where p.id=$1`

	err := r.db.QueryRow(ctx, productQuery, id).Scan(&product.ID, &product.Name, &product.SKU, &product.CategoryId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // ürün bulunamaz ise !
			return nil, fmt.Errorf("%v idsinine sahip ürün bulunamadı", id)
		}
		return nil, fmt.Errorf("Urun veritabaninda  sorgulanırken hata olustu : %w", err)
	}

	attrQuery := `select pav.attribute_id ,pav.value ,a."name" ,a.code ,a.data_type  from product_attribute_values pav 
join "attributes" a on pav.attribute_id =a.id 
where pav.product_id =$1`

	rows, err := r.db.Query(ctx, attrQuery, product.ID)
	if err != nil {
		return nil, fmt.Errorf("Nitelikler Veritabaninda sorgulanırken hata olustu : %w", err)
	}
	defer rows.Close()
	var pavs []domain.ProductAttributeValue

	for rows.Next() {
		var pav domain.ProductAttributeValue
		err := rows.Scan(&pav.AttributeID, &pav.Value, &pav.Attribute.Name, &pav.Attribute.Code, &pav.Attribute.DataType)
		if err != nil {
			return nil, fmt.Errorf("pav's sorgusunda hata : %w", err)
		}
		pavs = append(pavs, pav)
	}

	product.AttributeValues = pavs

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Satirlar Sorgulanırken olusan hata : %w", err)
	}
	return &product, nil
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

	// Tek sorgu atmak yerine önce ürünleri çekip sonra baglı oldugu ürünleri cekeceğiz.
	//cünkü tek sorgu yaparsak ürün*ürününNiteliği kadar satır doner

	// işe önce ürünleri çekmek ile başlayacağız ve limit ve offseti burda uyguluyor olacağız
	productsQuery := `select p.id ,p."name" ,p.sku ,p.category_id  from products p order by name asc limit $1 offset $2`

	productRows, err := r.db.Query(ctx, productsQuery, limit, offset)

	defer productRows.Close()
	if err != nil {
		return nil, 0, err
	}

	var products []domain.Product
	var productIds []int
	for productRows.Next() {
		var product domain.Product

		err = productRows.Scan(&product.ID, &product.Name, &product.SKU, &product.CategoryId)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
		productIds = append(productIds, product.ID)
	}
	// Bu noktada elimizde sayfalanmış şekilde ürünler var sirada bu ürünlerin attributelerini koymak var

	attrsQuery := `select pav.product_id, pav.attribute_id ,pav.value,a."name" ,a.code ,a.data_type ,a.id   
							from product_attribute_values pav 
							join "attributes" a on pav.attribute_id =a.id 
							where pav.product_id  = any($1)`
	attrRows, err := r.db.Query(ctx, attrsQuery, productIds)
	if err != nil {
		return nil, 0, err
	}

	// idlerini ve pavleri içeren bir map olusturuyorum az sonra eslestircem bunu kullanmazsam for dongusunde sistmei cok yoracak
	attrsMap := make(map[int][]domain.ProductAttributeValue)
	for attrRows.Next() {
		var pav domain.ProductAttributeValue
		var productID int
		err = attrRows.Scan(&productID, &pav.AttributeID, &pav.Value, &pav.Attribute.Name, &pav.Attribute.Code, &pav.Attribute.DataType, &pav.Attribute.ID)

		if err != nil {
			return nil, 0, err
		}

		attrsMap[productID] = append(attrsMap[productID], pav)
	}

	for i, _ := range products {
		products[i].AttributeValues = attrsMap[products[i].ID]

	}
	totalCountQuery := `select count(*) from products`
	var totalCount int
	r.db.QueryRow(ctx, totalCountQuery).Scan(&totalCount)

	return products, totalCount, nil
}

func (r *ProductRepository) UpdateAttributes(ctx context.Context, productId int, values []domain.ProductAttributeValue) error {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) SearchByAttribute(ctx context.Context, filters map[int]string) ([]domain.Product, error) {
	//TODO implement me
	panic("implement me")
}
