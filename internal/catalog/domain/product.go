package domain

import (
	"errors"
	"strings"
)

var (
	ErrProductNameNull = errors.New("Urun Adı Boş Olamaz")
	ErrSkuNull         = errors.New("SKU Alanı boş olamaz")
)

type ProductAttributeValue struct {
	//ID          int
	AttributeID int
	Value       string
	Attribute   Attribute
}

type Product struct {
	ID              int
	CategoryId      int
	Name            string
	SKU             string
	AttributeValues []ProductAttributeValue
}

func (p *Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrProductNameNull
	}
	if strings.TrimSpace(p.SKU) == "" {
		return ErrSkuNull
	}

	return nil
}
func (p *Product) UpdateAttribute(attrId int, newValue string) {

	for i, value := range p.AttributeValues {
		if value.AttributeID == attrId {
			p.AttributeValues[i].Value = newValue
			//value.Value = newValue // Burda koypa olusacağı için güncelleme islemi gerceklesmez
			return
		}
	}

	// Yeni Olustur
	newAttr := ProductAttributeValue{
		AttributeID: attrId,
		Value:       newValue,
	}
	p.AttributeValues = append(p.AttributeValues, newAttr)
	return
}
