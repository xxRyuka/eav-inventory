package usecase

import (
	"context"
	"eav-intentory/internal/domain"
)

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
}

type categoryUseCase struct {
	categoryRepository domain.CategoryRepository
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

func NewCategoryUseCase(repository domain.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{categoryRepository: repository}
}
