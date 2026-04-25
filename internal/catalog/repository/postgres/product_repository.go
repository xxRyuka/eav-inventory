package postgres

import (
	"context"
	catalog "eav-intentory/internal/catalog/domain"
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

func (r *ProductRepository) Create(ctx context.Context, p *catalog.Product) error {
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

func (r *ProductRepository) GetById(ctx context.Context, id int) (*catalog.Product, error) {

	var product catalog.Product
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
	var pavs []catalog.ProductAttributeValue

	for rows.Next() {
		var pav catalog.ProductAttributeValue
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

func (r *ProductRepository) Update(ctx context.Context, product *catalog.Product) error {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) GetAll(ctx context.Context, limit int, offset int, filters map[string][]string) ([]catalog.Product, int, error) {
	var args []any // parametreleri tutacağım kutu
	argsID := 1    //parametreleri $1 $2 diye yerlestirirken artık ($%d),argsID olarak verip her işlemde ++ yapacağımm

	totalCountQuery := `select count(*) from products p where 1=1 `
	productsQuery := `select p.id ,p."name" ,p.sku ,p.category_id  from products p where 1=1` //burda attributeleri cekmiyoruz

	for key, valueArr := range filters {
		attrFilterQuery := fmt.Sprintf(`	and  exists(
		select 1 from product_attribute_values pav 
		join "attributes" a on pav.attribute_id =a.id 
		where pav.product_id =p.id 
		and a.code = '%s' -- bunu ekliyoruz ''  yoksa parametreyi kolon sanıyor
		and pav.value=any($%d) )`,
			key, argsID)
		args = append(args, valueArr)
		argsID++

		productsQuery += attrFilterQuery
		totalCountQuery += attrFilterQuery
	}
	// todo :     "error": "Hata : Total Count Hesaplanırken olusan hata : expected 1 arguments, got 0"
	totalCount := 0
	err := r.db.QueryRow(ctx, totalCountQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("Total Count Hesaplanırken olusan hata : %w", err)
	}
	args = append(args, limit, offset)
	paginationQuery := fmt.Sprintf(` order by p.name asc limit $%d offset $%d `, argsID, argsID+1)
	argsID++

	productsQuery += paginationQuery
	productRows, err := r.db.Query(ctx, productsQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Urun Querysinde hata %w", err)
	}
	defer productRows.Close()
	var products []catalog.Product
	var productIds []int // az sorna attributeleri çekerken sadece elimizdeki ürünlerle ilişkili attributeleri çekmek için idlerini bir liste içine alıyorum
	for productRows.Next() {
		var product catalog.Product
		err = productRows.Scan(&product.ID, &product.Name, &product.SKU, &product.CategoryId)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
		productIds = append(productIds, product.ID)
	}

	err = productRows.Err()
	if err != nil {
		return nil, 0, err
	}
	// Bu noktada elimizde products var şimdi attributelerini eklememiz gerekiyor

	attributesQuery := ` select pav.attribute_id ,pav.product_id ,pav.value ,a.id ,a."name" ,a.code ,a.data_type 
								from product_attribute_values pav 
								join "attributes" a on a.id =pav.attribute_id
								where pav.product_id = any($1) -- sadece ilislili ürünleri getiriyorum 
`

	attributeRows, err := r.db.Query(ctx, attributesQuery, productIds)
	if err != nil {
		return nil, 0, err
	}
	defer attributeRows.Close()
	//var pavs []domain.ProductAttributeValue // mapi kurdugum için buna gerek yokmus
	productPavMap := make(map[int][]catalog.ProductAttributeValue) // idye göre eşleştieme yapacağım aşağıya oyüzden bir map olusturuyorum
	for attributeRows.Next() {
		var pav catalog.ProductAttributeValue
		var productId int
		err = attributeRows.Scan(&pav.AttributeID, &productId, &pav.Value, &pav.Attribute.ID, &pav.Attribute.Name, &pav.Attribute.Code, &pav.Attribute.DataType)
		if err != nil {
			return nil, 0, err
		}

		//pavs = append(pavs, pav) // bunu yapmaya gerek yokmus
		productPavMap[productId] = append(productPavMap[productId], pav) // appende değilde direk pav'e eşitleseydik her ürün için sadece 1 tane pav eklenecekti
	}
	err = attributeRows.Err()
	if err != nil {
		return nil, 0, err
	}

	for i, _ := range products {

		products[i].AttributeValues = productPavMap[products[i].ID]
	}
	return products, totalCount, nil
}

func (r *ProductRepository) UpdateAttributes(ctx context.Context, productId int, values []catalog.ProductAttributeValue) error {
	//TODO implement me
	panic("implement me")
}

func (r *ProductRepository) SearchByAttribute(ctx context.Context, filters map[int]string) ([]catalog.Product, error) {
	//TODO implement me
	panic("implement me")
}
