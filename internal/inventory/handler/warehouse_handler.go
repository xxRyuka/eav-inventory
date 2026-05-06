package inventory_handler

import (
	"eav-intentory/internal/inventory/domain"
	inventory_dto "eav-intentory/internal/inventory/handler/dto"
	inventory_usecase "eav-intentory/internal/inventory/usecase"
	"eav-intentory/pkg/response"
	"fmt"
	"net/http"
	"strconv"
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

func (h WarehouseHandler) GetWarehouses(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageINT, _ := strconv.Atoi(page)

	pageSize := r.URL.Query().Get("pageSize")
	pageSizeINT, _ := strconv.Atoi(pageSize)

	if pageSizeINT <= 0 {
		pageSizeINT = 10
	} else if pageSizeINT > 100 {
		pageSizeINT = 100
	}

	if pageINT <= 0 {
		pageINT = 1
	}

	warehouses, i, err := h.service.GetWarehouses(r.Context(), pageSizeINT, pageINT)
	if err != nil {
		response.ErrorJson(w, 500, "sayfalanmıs depolar getirilirken hata meydana geldi ! ", fmt.Errorf("err : %w", err))
		return
	}

	warehouseResponses := make([]inventory_dto.WarehouseResponse, 0, len(warehouses))

	for k, _ := range warehouses {

		warehouseDto := inventory_dto.WarehouseResponse{
			ID:       warehouses[k].ID,
			Name:     warehouses[k].Name,
			Code:     warehouses[k].Code,
			Location: warehouses[k].Location,
		}
		warehouseResponses = append(warehouseResponses, warehouseDto)
	}

	resp := response.CalculatedPagedResponse(warehouseResponses, i, pageSizeINT, pageINT) // *prd

	response.WriteJson(w, 200, resp, "")
}
func (h WarehouseHandler) GetWarehouseById(w http.ResponseWriter, r *http.Request) {

	idSTR := r.PathValue("id")
	idINT, err := strconv.Atoi(idSTR)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "0'dan büyük bir tam sayi id'si gonder (id : %d)", err)
		return
	}
	warehouse, err := h.service.GetWarehouseById(r.Context(), idINT)
	if err != nil {
		response.ErrorJson(w, 500, "depo sorgulanırken hata meydana geldi", fmt.Errorf("err : %w", err))
		return
	}

	respDto := inventory_dto.WarehouseResponse{
		ID:       warehouse.ID,
		Name:     warehouse.Name,
		Code:     warehouse.Code,
		Location: warehouse.Location,
	}

	response.WriteJson(w, 200, respDto, "")
}

// id'yi path'den cekiyorum dto'yu bodyden sonra id'yi dtoya mapleyip uc layera gonderiyor olacağım
func (h WarehouseHandler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	idSTR := r.PathValue("id")
	idINT, _ := strconv.Atoi(idSTR)

	if idINT == 0 {
		response.ErrorJson(w, http.StatusBadRequest, "Gecersiz ID ", fmt.Errorf("invalid id"))
	}

	var warehouseDto inventory_dto.WarehouseResponse
	err := response.ReadJson(w, r, &warehouseDto)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "invalid wh", fmt.Errorf("invalid wh"))
		return
	}
	warehouseDto.ID = idINT
	wh := domain.Warehouse{
		ID:       idINT,
		Location: warehouseDto.Location,
		Name:     warehouseDto.Name,
		Code:     warehouseDto.Code,
	}
	err = h.service.UpdateWarehouse(r.Context(), wh.ID, &wh)
	if err != nil {
		response.ErrorJson(w, 500, "an error occured while updating warehouse", fmt.Errorf("err : %w", err))
		return
	}

	response.WriteJson(w, 200, true, "")
}

func (h WarehouseHandler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	idSTR := r.PathValue("id")
	idINT, _ := strconv.Atoi(idSTR)

	if idINT == 0 {
		response.ErrorJson(w, http.StatusBadRequest, "Gecersiz ID ", fmt.Errorf("invalid id"))
	}

	err := h.service.DeleteWarehouse(r.Context(), idINT)
	if err != nil {
		response.ErrorJson(w, 500, "an error occured while deleting warehouse", fmt.Errorf("err : %w", err))
		return
	}
	response.WriteJson(w, 200, true, "")

}
