package inventory_dto

type PurchaseInRequest struct {
	WarehouseID int `json:"warehouse_id,omitempty"`
	ProductID   int `json:"product_id,omitempty"`
	Quantity    int `json:"quantity,omitempty"`
}
