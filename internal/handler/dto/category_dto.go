package dto

// Request's

type CreateCategoryRequest struct {
	Name       string                           `json:"name"`
	ParentID   *int                             `json:"parentID,omitempty"`
	Attributes []CreateCategoryAttributeRequest `json:"attributes"`
}

type CreateCategoryAttributeRequest struct {
	AttributeID int  `json:"attributeID"`
	IsRequired  bool `json:"isRequired"`
}

type UpdateCategoryRequest struct {
	//ID       int id'yi path üzerinden almam lazım
	Name     string `json:"name"`
	ParentID *int   `json:"parentID"`
}

// Response's

type CategoryAttributeResponse struct {
	AttributeID int  `json:"attributeID"`
	IsRequired  bool `json:"isRequired"`

	// Attributeleri inner join ile cektiğimiz
	Code     string `json:"code"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type CategoryResponse struct {
	ID         int                         `json:"ID"`
	Name       string                      `json:"name"`
	ParentID   *int                        `json:"parentID"`
	Attributes []CategoryAttributeResponse `json:"attributes"`
}
