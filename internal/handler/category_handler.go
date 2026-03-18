package handler

import (
	"eav-intentory/internal/domain"
	"eav-intentory/internal/usecase"
	"eav-intentory/pkg/response"
	"fmt"
	"net/http"
)

type CategoryHandler struct {
	categoryService usecase.CategoryUseCase
}

func NewCategoryHandler(useCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryService: useCase}
}

// Dto's
type CategoryAttributeDto struct {
	Name       string `json:"name"`
	DataType   string `json:"dataType"` // JSON'dan saf string olarak gelir
	IsRequired bool   `json:"isRequired"`
}

type createCategoryRequest struct {
	Name       string                 `json:"name"`
	ParentID   *int                   `json:"parentID,omitempty"`
	Attributes []CategoryAttributeDto `json:"attributes"`
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.ErrorJson(w, http.StatusMethodNotAllowed, "gecersiz method", fmt.Errorf("Gecersiz Method"))
		return
	}
	req := createCategoryRequest{}
	err := response.ReadJson(w, r, &req)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "json bind edilemedi", err)
		return
	}
	domainAttributes := make([]domain.CategoryAttribute, 0, len(req.Attributes))

	for _, attributeDto := range req.Attributes {
		attr := domain.CategoryAttribute{
			//ID:         0,
			Name:       attributeDto.Name,
			DataType:   domain.DataType(attributeDto.DataType),
			IsRequired: attributeDto.IsRequired,
		}

		domainAttributes = append(domainAttributes, attr)
	}

	category := domain.Category{
		Name:       req.Name,
		ParentID:   req.ParentID,
		Attributes: domainAttributes,
	}

	err = h.categoryService.CreateCategory(r.Context(), &category)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	response.WriteJson(w, http.StatusCreated, category.ID, "islem basarili")
}
