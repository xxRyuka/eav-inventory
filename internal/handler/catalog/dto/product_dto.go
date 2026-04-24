package dto

// Request's

type CreateProductRequest struct {
	Name       string                    `json:"name"`
	CategoryID int                       `json:"categoryID"`
	SKU        string                    `json:"sku"`
	Attributes []ProductAttributeRequest `json:"attributes"` // burda append yapabilmek için böyle direk domain nesnesi kullandım alanları aynı diye ama yaptıgım yanlıs mı ? best praticesi nasıl ?
}
type ProductAttributeRequest struct {
	AttributeID int    `json:"attributeID"`
	Value       string `json:"value"`
}

//Response's

type ProductResponse struct {
	ID         int                        `json:"id"`
	Name       string                     `json:"name"`
	CategoryID int                        `json:"categoryID"`
	SKU        string                     `json:"sku"`
	Attributes []ProductAttributeResponse `json:"attributes"`
}

type ProductAttributeResponse struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	AttributeID int    `json:"attributeID"`
	Value       string `json:"value"`
}
