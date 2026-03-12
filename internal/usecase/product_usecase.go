package usecase

import (
	"context"
	"eav-intentory/internal/domain"
	"errors"
	"fmt"
)

// Burda usecaaselerde bişey kafamı karıstırdı neden interface ve struct kullanıyoruz ki sadece struct işimizi gormez mi ?
// sonucta usecase db gibi değil değişmez birşey lütfen haksizsam haksızsın de gercekten gerekliliğini ogrenmek istiyorum?
// ve niye bu struct private ?
type ProductUseCase interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
}

type productUseCase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
}

func (p *productUseCase) CreateProduct(ctx context.Context, product *domain.Product) error {
	err := product.Validate()
	if err != nil {
		return err
	}

	category, err := p.categoryRepo.GetById(ctx, product.CategoryId)
	if err != nil {
		return errors.New("kategori bulunamadi")
	}

	attributes := category.Attributes
	for _, attribute := range attributes {
		if attribute.IsRequired == true {
			found := false
			for _, pav := range product.AttributeValues {
				if pav.AttributeID == attribute.ID {
					found = true
					break
				}
				if !found {

					return fmt.Errorf("zorunlu nitelik eksik: %s (ID: %d)", attribute.Name, attribute.ID)
				}
			}
		}
	}

	err = p.productRepo.Create(ctx, product)
	if err != nil {
		return err
	}
	return nil

}

// Burdan neden bir interface donuyoz ya ?? bide neden basına * koymuyoz
func NewProductUseCase(productRepository domain.ProductRepository, categoryRepository domain.CategoryRepository) ProductUseCase {
	return &productUseCase{
		productRepo:  productRepository,
		categoryRepo: categoryRepository,
	}
}
