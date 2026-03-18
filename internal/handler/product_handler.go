package handler

import (
	"eav-intentory/internal/domain"
	"eav-intentory/internal/usecase"
	"eav-intentory/pkg/response"
	"fmt"
	"net/http"
)

type ProductHandler struct {
	service usecase.ProductUseCase
}

func NewProductHandler(useCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{service: useCase}
}

type createProductRequest struct {
	Name       string         `json:"name"`
	CategoryID int            `json:"categoryID"`
	SKU        string         `json:"SKU"`
	Attributes []attributeDto `json:"attributes"` // burda append yapabilmek için böyle direk domain nesnesi kullandım alanları aynı diye ama yaptıgım yanlıs mı ? best praticesi nasıl ?
}

type attributeDto struct {
	AttributeID int    `json:"attributeID"`
	Value       string `json:"value"`
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.ErrorJson(w, http.StatusMethodNotAllowed, "gecersiz method", fmt.Errorf("Gecersiz Method"))
		return
	}
	req := createProductRequest{}
	err := response.ReadJson(w, r, &req)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "bind edilemedi", fmt.Errorf("json okunamadı"))
		return
	}
	product := domain.Product{
		CategoryId: req.CategoryID,
		Name:       req.Name,
		SKU:        req.SKU,
		//AttributeValues: nil,
	}

	for _, attribute := range req.Attributes {
		product.AttributeValues = append(product.AttributeValues, domain.ProductAttributeValue{
			AttributeID: attribute.AttributeID,
			Value:       attribute.Value,
		})
	}
	err = h.service.CreateProduct(r.Context(), &product)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.WriteJson(w, http.StatusCreated, product.ID, "Basarili")
}
