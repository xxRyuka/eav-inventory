package inventory_handler

import (
	inventory_dto "eav-intentory/internal/inventory/handler/dto"
	"eav-intentory/pkg/response"
	"net/http"
)

type WarehouseHandler struct {
	service WarehouseUsecase
}

func NewWarehouseHandler(warehouseUsecase WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{service: warehouseUsecase}
}

func (h WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {

	var req inventory_dto.CreateCategoryRequest
	response.ReadJson(w, r, &req)
}
