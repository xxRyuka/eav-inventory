package usecase

import (
	"context"
	catalog "eav-intentory/internal/catalog/domain"
	"fmt"
)

// todo: usecase errors girilecek

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category *catalog.Category) error
	GetCategoryById(ctx context.Context, id int) (*catalog.Category, error)
	GetCategories(ctx context.Context, pageSize, page int) ([]catalog.Category, int, error)
	UpdateCategory(ctx context.Context, category *catalog.Category) error
	GetCategoriesWithAttributes(ctx context.Context) ([]catalog.Category, error)
	RemoveAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error
	AddAttributeToCategory(ctx context.Context, categoryID, attributeID int, isRequired bool) error
	UpdateAttributeToCategory(ctx context.Context, isRequired bool, attributeID, categoryID int) error
}

type categoryUseCase struct {
	categoryRepository catalog.CategoryRepository
}

func (c *categoryUseCase) GetCategoriesWithAttributes(ctx context.Context) ([]catalog.Category, error) {

	// ekstra isleme gerek yok burda suanlık

	categoriesWithAttirbutes, err := c.categoryRepository.GetCategoriesWithAttirbutes(ctx)
	if err != nil {
		return nil, err
	}

	return categoriesWithAttirbutes, nil
}

func (c *categoryUseCase) RemoveAttributeFromCategory(ctx context.Context, categoryID, attributeID int) error {
	if categoryID <= 0 || attributeID <= 0 {
		return fmt.Errorf("Gecersiz id parametresi, id 0'dan kucuk olamaz")
	}

	err := c.categoryRepository.RemoveAttributeToCategory(ctx, categoryID, attributeID)
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryUseCase) AddAttributeToCategory(ctx context.Context, categoryID, attributeID int, isRequired bool) error {
	if categoryID <= 0 || attributeID <= 0 {
		return fmt.Errorf("Gecersiz id parametresi, id 0'dan kucuk olamaz")
	}

	err := c.categoryRepository.AddAttributeToCategory(ctx, categoryID, attributeID, isRequired)
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryUseCase) UpdateAttributeToCategory(ctx context.Context, isRequired bool, attributeID, categoryID int) error {
	if categoryID <= 0 || attributeID <= 0 {
		return fmt.Errorf("Gecersiz id parametresi, id 0'dan kucuk olamaz")
	}

	err := c.categoryRepository.UpdateAttributeToCategory(ctx, isRequired, attributeID, categoryID)
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryUseCase) UpdateCategory(ctx context.Context, category *catalog.Category) error {
	if category.ID <= 0 {
		return fmt.Errorf(" gecersiz id %v", category.ID)
	}
	err := c.categoryRepository.Update(ctx, category)
	if err != nil {
		return err
	}
	return nil
}

func (c *categoryUseCase) GetCategories(ctx context.Context, pageSize, page int) ([]catalog.Category, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize < 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	categories, totalCount, err := c.categoryRepository.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return categories, totalCount, nil
}

func NewCategoryUseCase(repository catalog.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{categoryRepository: repository}
}

func (c *categoryUseCase) CreateCategory(ctx context.Context, category *catalog.Category) error {

	err := category.Validate()
	if err != nil {
		return err
	}
	err = c.categoryRepository.Create(ctx, category)
	if err != nil {
		return err
	}
	return nil
}

func (c *categoryUseCase) GetCategoryById(ctx context.Context, id int) (*catalog.Category, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Lütfen 0'dan büyük geçerli bir id giriniz")
	}
	category, err := c.categoryRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}
