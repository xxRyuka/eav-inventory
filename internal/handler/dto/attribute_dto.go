package dto

type CreateAttributeRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	DataType string `json:"dataType"`
}
