package inventory_handler

import (
	"eav-intentory/internal/inventory/domain"
	inventory_dto "eav-intentory/internal/inventory/handler/dto"
	inventory_usecase "eav-intentory/internal/inventory/usecase"
	"eav-intentory/internal/inventory/usecase/command"
	"eav-intentory/pkg/response"
	"fmt"
	"net/http"
)

type StockHandler struct {
	service inventory_usecase.StockUsecase
}

func NewStockHandler(service inventory_usecase.StockUsecase) *StockHandler { // neden handlerda pointer donuyoruz ?
	return &StockHandler{service: service}
}

func (h *StockHandler) PurchaseIn(w http.ResponseWriter, r *http.Request) {

	fmt.Println("in handler")
	var req inventory_dto.PurchaseInRequest
	err := response.ReadJson(w, r, &req)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "gecersiz json", fmt.Errorf("json : %w", err))
		return
	}

	cmd := command.PurchaseInCommand{
		WarehouseID:  req.WarehouseID,
		ProductID:    req.ProductID,
		Quantity:     req.Quantity,
		MovementType: string(domain.PurchaseIn),
	}

	err = h.service.PurchaseIn(r.Context(), cmd)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "hata mesajını okuyun", fmt.Errorf("server : %w", err))
		return
	}

	response.WriteJson(w, http.StatusCreated, true, "")
}
