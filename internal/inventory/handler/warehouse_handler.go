package inventory_handler

import (
	"eav-intentory/internal/inventory/domain"
	inventory_dto "eav-intentory/internal/inventory/handler/dto"
	inventory_usecase "eav-intentory/internal/inventory/usecase"
	"eav-intentory/pkg/response"
	"fmt"
	"net/http"
)

type WarehouseHandler struct {
	service inventory_usecase.WarehouseUsecase
}

func NewWarehouseHandler(warehouseUsecase inventory_usecase.WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{service: warehouseUsecase}
}

func (h WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {

	var req inventory_dto.CreateCategoryRequest
	err := response.ReadJson(w, r, &req)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "bind error ", fmt.Errorf("err : %w", err))
		return
	}

	warehouse := domain.Warehouse{
		Location: req.Location,
		Name:     req.Name,
		Code:     req.Code,
	}
	id, err := h.service.CreateWarehouse(r.Context(), &warehouse)
	if err != nil {
		response.ErrorJson(w, 500, "service layer err ", fmt.Errorf("err : %w", err))
		return
	}
	response.WriteJson(w, 201, id, "succsess")

}
