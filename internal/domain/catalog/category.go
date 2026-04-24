package catalog

import (
	"errors"
	"strings"
)

type DataType string

// Domain Hataları (Dış dünyaya fırlatılacak standart sözlük)
var (
	ErrCategoryNameEmpty = errors.New("kategori adi bos olamaz")
	ErrInvalidDataType   = errors.New("gecersiz veri tipi tanimlandi")
)

const (
	TypeString DataType = "string"
	TypeInt    DataType = "int"
	TypeFloat  DataType = "float"
	TypeBool   DataType = "bool"
)

// refactoring !

// value object
type CategoryAttribute struct {
	AttributeID int
	Attribute   Attribute
	IsRequired  bool
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

// burayıda refactor ettim denemek lazım
func (c *Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrCategoryNameEmpty
	}
	for _, attr := range c.Attributes {
		if (attr.AttributeID) <= 0 {
			return errors.New("attribute id 0'dan kucuk olamaz ")
		}
	}
	return nil
}
