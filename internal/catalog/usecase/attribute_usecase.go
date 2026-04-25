package usecase

import (
	"context"
	catalog "eav-intentory/internal/catalog/domain"
	"eav-intentory/internal/catalog/usecase/command"
	"fmt"
)

type AttributeUsecase interface {
	CreateAttribute(ctx context.Context, attribute *catalog.Attribute) error
	GetAttributeByID(ctx context.Context, id int) (*catalog.Attribute, error)
	GetAttributes(ctx context.Context, page, pageSize int) ([]catalog.Attribute, int, error)

	DeleteAttribute(ctx context.Context, id int) error
	UpdateAttribute(ctx context.Context, attribute *command.UpdateAttributeCommand) error
}

type attributeUsecase struct {
	attributeRepo catalog.AttributeRepository
}

func (a *attributeUsecase) UpdateAttribute(ctx context.Context, attribute *command.UpdateAttributeCommand) error {

	attr := catalog.Attribute{
		ID:       attribute.ID,
		Code:     attribute.Code,
		Name:     attribute.Name,
		DataType: catalog.DataType(attribute.DataType),
	}
	err := a.attributeRepo.Update(ctx, &attr)
	if err != nil {
		return err
	}
	return nil
}

func (a *attributeUsecase) DeleteAttribute(ctx context.Context, id int) error {
	if id < 0 {
		return fmt.Errorf("Gecersiz id")
	}

	err := a.attributeRepo.Remove(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func NewAttributeService(repository catalog.AttributeRepository) AttributeUsecase {
	return &attributeUsecase{attributeRepo: repository}
}

func (a *attributeUsecase) CreateAttribute(ctx context.Context, attribute *catalog.Attribute) error {
	err := attribute.Validate()
	if err != nil {
		return fmt.Errorf("Validasyon Hatasi %w", err)
	}
	err = a.attributeRepo.Create(ctx, attribute)
	if err != nil {
		return err
	}
	return nil

}

func (a *attributeUsecase) GetAttributeByID(ctx context.Context, id int) (*catalog.Attribute, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id 0dan büyük olmalıdır")
	}

	attr, err := a.attributeRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

func (a *attributeUsecase) GetAttributes(ctx context.Context, page, limit int) ([]catalog.Attribute, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}

	if limit <= 0 {
		limit = 10
	}
	// offset kac tane sayfa atlayacağımızdıe yani 2. sayfayadaysak ve limit 10 ise (2-1)*10 =10 yani 10 tane atlayacak ve 11. veriden basalyacak 2. sayfa
	offset := (page - 1) * limit
	attributes, i, err := a.attributeRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return attributes, i, nil
}
