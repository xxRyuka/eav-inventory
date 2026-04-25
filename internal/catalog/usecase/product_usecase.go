package usecase

import (
	"context"
	catalog "eav-intentory/internal/catalog/domain"
	"fmt"
	"strconv"
	"strings"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *catalog.Product) error
	GetProductById(ctx context.Context, id int) (*catalog.Product, error)
	GetProducts(ctx context.Context, pageSize, page int, filters map[string][]string) ([]catalog.Product, int, error)
}

type productUseCase struct {
	productRepo  catalog.ProductRepository
	categoryRepo catalog.CategoryRepository
}

func NewProductUseCase(productRepository catalog.ProductRepository, categoryRepository catalog.CategoryRepository) ProductUseCase {
	return &productUseCase{
		productRepo:  productRepository,
		categoryRepo: categoryRepository,
	}
}

func (p *productUseCase) CreateProduct(ctx context.Context, product *catalog.Product) error {

	err := product.Validate()
	if err != nil {
		return err
	}
	// Gelen CategoryID gerçekten veritabanında var mı?
	category, err := p.categoryRepo.GetById(ctx, product.CategoryId)
	if err != nil {
		return err
	}
	attributes := category.Attributes
	//Kullanıcının gönderdiği nitelikler, bu kategorinin şablonunda mevcut mu?
	//Şablon "Bu alan zorunludur" diyorsa, kullanıcı boş geçmiş mi?

	for _, attribute := range attributes {

		var userValue *catalog.ProductAttributeValue // burda neden kodda pointer kullandın == Kodda kendim gordum yazarken pointer kullanmaz isek nil kontrolu yapamayiz
		for i, value := range product.AttributeValues {
			if value.AttributeID == attribute.AttributeID {
				userValue = &product.AttributeValues[i] // Buraya direk value diyemez miydik ? => diyemezdik üzerinde değişiklik yaptıgımız tek sey for dongusu için gecici olarak acılan değişken olurdu bu durumda !
				// Artık user value değişkenini değiştirmek gercek değeride değiştirecek !
				break
			}
		}
		if attribute.IsRequired && userValue == nil {
			return fmt.Errorf("Zorunlu Nitelik Alanı Boş %v", attribute.Attribute.Name)
		}
		// KURAL 2: Tip Kontrolü (Eğer kullanıcı değer gönderdiyse)
		if userValue != nil {
			cleanVal := strings.TrimSpace(userValue.Value)

			switch attribute.Attribute.DataType {
			case catalog.TypeInt:
				_, err = strconv.Atoi(cleanVal)
				if err != nil {
					return fmt.Errorf("'%s' alani tam sayi olmalidir, girilen gecersiz deger: %s", attribute.Attribute.Name, userValue.Value)
				}
			case catalog.TypeBool:
				if _, err = strconv.ParseBool(cleanVal); err != nil {
					return fmt.Errorf("%s alani boolean ifade olmalidir gidirlen gecersiz deger %s", attribute.Attribute.Name, userValue.Value)
				}
			case catalog.TypeString:
				// String ise zaten string'dir, ekstra bir kontrole gerek yok ama boş mu diye bakabiliriz.
				if cleanVal == "" {
					return fmt.Errorf("'%s' alani bos birakilamaz", attribute.Attribute.Name)
				}
			}
			userValue.Value = cleanVal

		}

	}

	err = p.productRepo.Create(ctx, product)
	if err != nil {
		return fmt.Errorf("Veritabanına Kaydolurken Hata Oldu %w", err)
	}

	return nil
	// ?? Şablon "Bu alan tam sayıdır" diyorsa, kullanıcı harf girmiş mi? bu kontrolu nasıl yaparım bilemedim
}

func (p *productUseCase) GetProductById(ctx context.Context, id int) (*catalog.Product, error) {
	if id <= 0 {
		return nil, fmt.Errorf("0'dan büyük bir tam sayı değer giriniz id için (id:%v)", id)
	}

	product, err := p.productRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("db err : %w", err)
	}

	return product, nil
}
func (p *productUseCase) GetProducts(ctx context.Context, pageSize, page int, filters map[string][]string) ([]catalog.Product, int, error) {
	// optimizations

	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize

	// filters mapi ile ilgili burda ekstra validasyon yapmama gerek var mı bilemedim

	products, total, err := p.productRepo.GetAll(ctx, pageSize, offset, filters)
	if err != nil {
		return nil, 0, err
	}
	// totalide hesaplatmamız lazim
	return products, total, nil
}
