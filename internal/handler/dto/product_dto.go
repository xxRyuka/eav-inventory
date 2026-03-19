package dto

type CreateProductRequest struct {
	Name       string         `json:"name"`
	CategoryID int            `json:"categoryID"`
	SKU        string         `json:"SKU"`
	Attributes []AttributeDto `json:"attributes"` // burda append yapabilmek için böyle direk domain nesnesi kullandım alanları aynı diye ama yaptıgım yanlıs mı ? best praticesi nasıl ?
}
type AttributeDto struct {
	AttributeID int    `json:"attributeID"`
	Value       string `json:"value"`
}
