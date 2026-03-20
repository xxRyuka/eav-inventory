package dto

// Request's

type CreateAttributeRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	DataType string `json:"dataType"`
}

// Response

type AttributeResponse struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}
