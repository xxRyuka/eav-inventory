package domain

import "context"

type AttributeRepository interface {
	Create(ctx context.Context, attribute *Attribute) error
	GetAll(ctx context.Context, limit int, offset int) ([]Attribute, int, error)
	GetById(ctx context.Context, id int) (*Attribute, error)
}
