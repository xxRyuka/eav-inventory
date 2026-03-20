package handler

import (
	"eav-intentory/internal/domain"
	"eav-intentory/internal/handler/dto"
	"eav-intentory/internal/usecase"
	"eav-intentory/pkg/response"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type AttributeHandler struct {
	attrUsecase usecase.AttributeUsecase
}

func NewAttributeHandler(attributeUsecase usecase.AttributeUsecase) *AttributeHandler {
	return &AttributeHandler{attrUsecase: attributeUsecase}
}

func (h *AttributeHandler) CreateAttribute(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		response.ErrorJson(w, http.StatusMethodNotAllowed, "yanlis method bu endpoint post ile calısıyor ", fmt.Errorf("%v istegi gelmis post isteği atmalisin", r.Method))
		return
	}
	var req dto.CreateAttributeRequest
	err := response.ReadJson(w, r, &req)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "json okunurken hata meydana geldi", fmt.Errorf("Json okunup structa bind edilirken hata gerçekleşti %w", err))
		return
	}
	attr := domain.Attribute{
		Code:     req.Code,
		Name:     req.Name,
		DataType: domain.DataType(req.DataType),
	}
	err = h.attrUsecase.CreateAttribute(r.Context(), &attr)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "attribute olusturulurken hata ", fmt.Errorf("attribute veritabanında olusutrulurken hata olustu hata : %w ", err))
		return
	}

	response.WriteJson(w, http.StatusCreated, attr.ID, "")

}

func (h *AttributeHandler) GetAttributeByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.ErrorJson(w, http.StatusMethodNotAllowed, "yanlis method, bu endpoint get ile calısıyor", fmt.Errorf("%v istegi gelmis get isteği atmalisin", r.Method))
		return
	}

	id := r.PathValue("id")
	atoi, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "json okunurken hata meydana geldi gecerli bir tam sayi id'si gerekiyor", fmt.Errorf("gecersiz id lütfen tam sayı giriniz %w", err))
		return
	}

	if atoi < 0 {
		response.ErrorJson(w, http.StatusBadRequest, "gecerli bir tam sayi id'si gerekiyor", fmt.Errorf("gecersiz id lütfen tam sayı giriniz %w", err))
		return
	}

	attribute, err := h.attrUsecase.GetAttributeByID(r.Context(), atoi)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "attribute sorgulanırken ", fmt.Errorf("attribute id veritabanında sorgulanırken hata olustu, hata : %w ", err))
		return
	}

	attrDto := dto.AttributeResponse{
		ID:       attribute.ID,
		Code:     attribute.Code,
		Name:     attribute.Name,
		DataType: string(attribute.DataType),
	}

	response.WriteJson(w, http.StatusOK, attrDto, "")
}

func (h *AttributeHandler) GetAttributes(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")

	limitQuery, err := strconv.Atoi(limit)
	pageQuery, err := strconv.Atoi(page)

	attributes, i, err := h.attrUsecase.GetAttributes(r.Context(), pageQuery, limitQuery)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "usecaseden hata geldi", err)
		return
	}

	pagedcalculated := response.CalculatedPagedResponse(attributes, i, limitQuery, pageQuery)
	response.WriteJson(w, http.StatusOK, pagedcalculated, "")
}
