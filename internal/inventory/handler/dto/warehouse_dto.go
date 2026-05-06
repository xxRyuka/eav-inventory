package inventory_dto

// Request's

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Location string `json:"location"`
}

// Response's
type WarehouseResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Location string `json:"location"`
}
