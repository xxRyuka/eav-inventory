package inventory_dto

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Location string `json:"location"`
}
