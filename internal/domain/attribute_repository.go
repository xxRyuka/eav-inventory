package domain

import "context"

type AttributeRepository interface {
	Create(ctx context.Context, attribute *Attribute) error
	GetAll(ctx context.Context, limit int, offset int) ([]Attribute, int, error)
	GetById(ctx context.Context, id int) (*Attribute, error)
	// TODO: Update ve Remove Kodlarını yazarken projeyi hatirlicaz !
	Update(ctx context.Context, attribute *Attribute) error
	Remove(ctx context.Context, id int) error
}
