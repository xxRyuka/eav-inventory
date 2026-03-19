package usecase

import (
	"context"
	"eav-intentory/internal/domain"
	"fmt"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategoryById(ctx context.Context, id int) (*domain.Category, error)
}

type categoryUseCase struct {
	categoryRepository domain.CategoryRepository
}

func NewCategoryUseCase(repository domain.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{categoryRepository: repository}
}

func (c *categoryUseCase) CreateCategory(ctx context.Context, category *domain.Category) error {

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

func (c *categoryUseCase) GetCategoryById(ctx context.Context, id int) (*domain.Category, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Lütfen 0'dan büyük geçerli bir id giriniz")
	}
	category, err := c.categoryRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}
