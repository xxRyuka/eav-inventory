package domain

import (
	"errors"
	"strings"
)

type DataType string

// Domain Hataları (Dış dünyaya fırlatılacak standart sözlük)
var (
	ErrCategoryNameEmpty  = errors.New("kategori adi bos olamaz")
	ErrAttributeNameEmpty = errors.New("nitelik adi bos olamaz")
	ErrInvalidDataType    = errors.New("gecersiz veri tipi tanimlandi")
)

const (
	TypeString DataType = "string"
	TypeInt    DataType = "int"
	TypeFloat  DataType = "float"
	TypeBool   DataType = "bool"
)

type CategoryAttribute struct {
	ID         int
	Name       string
	DataType   DataType
	IsRequired bool
}

// Category Go'da int tipinin varsayılan değeri 0'dır.
// Eğer bir kategori en üst seviyedeyse (örn: Elektronik), veritabanında NULL tutmamız gerekir.
// Go'da NULL mantığını sadece pointer'lar (nil) ile kurabiliriz.
type Category struct {
	ID         int
	Name       string
	ParentID   *int
	Attributes []CategoryAttribute
}

func (c *Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrCategoryNameEmpty
	}
	for _, attr := range c.Attributes {
		if strings.TrimSpace(attr.Name) == "" {
			return ErrAttributeNameEmpty
		}
		switch attr.DataType {
		case TypeString, TypeInt, TypeBool, TypeFloat:
			// Geçerli
		default:
			return ErrInvalidDataType
		}
	}
	return nil
}
