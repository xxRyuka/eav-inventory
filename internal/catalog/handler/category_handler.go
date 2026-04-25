package handler

import (
	catalog2 "eav-intentory/internal/catalog/domain"
	"eav-intentory/internal/catalog/handler/dto"
	"eav-intentory/internal/catalog/usecase"
	"eav-intentory/pkg/response"
	"fmt"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	categoryService usecase.CategoryUseCase
}

func NewCategoryHandler(useCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryService: useCase}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.ErrorJson(w, http.StatusMethodNotAllowed, "gecersiz method", fmt.Errorf("Gecersiz Method"))
		return
	}
	req := dto.CreateCategoryRequest{}
	err := response.ReadJson(w, r, &req)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "json bind edilemedi", err)
		return
	}
	domainAttributes := make([]catalog2.CategoryAttribute, 0, len(req.Attributes))

	for _, attributeDto := range req.Attributes {
		attr := catalog2.CategoryAttribute{
			AttributeID: attributeDto.AttributeID,
			IsRequired:  attributeDto.IsRequired,
		}

		domainAttributes = append(domainAttributes, attr)
	}

	category := catalog2.Category{
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

func (h *CategoryHandler) GetCategoryById(w http.ResponseWriter, r *http.Request) {

	//id := r.URL.Query().Get("id") // queryden değil patthen cekeceğiz

	id := r.PathValue("id")

	idINT, err := strconv.Atoi(id)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "0'dan büyük bir tam sayi id'si gonder", err)
		return
	}
	category, err := h.categoryService.GetCategoryById(r.Context(), idINT)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "Hata Olustu ", fmt.Errorf("error %w", err))
		return
	}

	var attrDtos []dto.CategoryAttributeResponse

	for _, k := range category.Attributes {
		attrDto := dto.CategoryAttributeResponse{
			AttributeID: k.AttributeID,
			IsRequired:  k.IsRequired,
			Code:        k.Attribute.Code,
			Name:        k.Attribute.Name,
			DataType:    string(k.Attribute.DataType),
		}

		attrDtos = append(attrDtos, attrDto)
	}

	resp := dto.CategoryResponse{
		ID:         category.ID,
		Name:       category.Name,
		ParentID:   category.ParentID,
		Attributes: attrDtos,
	}

	response.WriteJson(w, http.StatusOK, resp, "")
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	pageSize := r.URL.Query().Get("pageSize")
	page := r.URL.Query().Get("page")

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "query validation err", fmt.Errorf("err : %w", err))
		return
	}
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "query validation err", fmt.Errorf("err : %w", err))
		return
	}

	categories, totalCount, err := h.categoryService.GetCategories(r.Context(), pageSizeInt, pageInt)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "usecase and repo err", fmt.Errorf("err : %w", err))
		return
	}

	resp := response.CalculatedPagedResponse(categories, totalCount, pageSizeInt, pageInt)

	response.WriteJson(w, 200, resp, "basarili")
}

func (h *CategoryHandler) UpdateBaseCategory(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, fmt.Sprintf("gecersiz id (ID:%V)", id), err)
		return
	}
	var request dto.UpdateCategoryRequest
	err = response.ReadJson(w, r, &request)

	category := catalog2.Category{
		ID:       idInt,
		Name:     request.Name,
		ParentID: request.ParentID,
	}
	err = h.categoryService.UpdateCategory(r.Context(), &category)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "alt katmanlarda hata ", err)
		return
	}

	response.WriteJson(w, http.StatusOK, true, "update succsess")
}

// todo: no pagination for now
func (h *CategoryHandler) GetCategoriesWithAttributes(w http.ResponseWriter, r *http.Request) {

	categoriesWithAttributes, err := h.categoryService.GetCategoriesWithAttributes(r.Context())
	if err != nil {
		response.ErrorJson(w, 500, "sunucu kaynaklı problem", fmt.Errorf("err : %w", err))
		return
	}

	var responseData []dto.CategoryResponse

	for _, category := range categoriesWithAttributes {
		categoryResp := dto.CategoryResponse{
			ID:       category.ID,
			Name:     category.Name,
			ParentID: category.ParentID,
		}

		var catAttrResp []dto.CategoryAttributeResponse
		for _, attributeResponse := range category.Attributes {
			catAttrResp = append(catAttrResp, dto.CategoryAttributeResponse{
				AttributeID: attributeResponse.AttributeID,
				IsRequired:  attributeResponse.IsRequired,
				Code:        attributeResponse.Attribute.Code,
				Name:        attributeResponse.Attribute.Name,
				DataType:    string(attributeResponse.Attribute.DataType),
			})
		}

		// Burda attributeleri dtoya atiyoruz
		categoryResp.Attributes = catAttrResp

		// burdada responseye yazıyoruz
		responseData = append(responseData, categoryResp)
	}

	response.WriteJson(w, 200, responseData, "succsess")

}

func (h *CategoryHandler) AssignAttributeToCategory(w http.ResponseWriter, r *http.Request) {

	var request dto.AssignAttributeToCategoryRequest
	err := response.ReadJson(w, r, &request)
	if err != nil {
		response.ErrorJson(w, 400, "json okunurken hata meydana geldi ", fmt.Errorf("err :%w", err))
		return
	}
	err = h.categoryService.AddAttributeToCategory(r.Context(), request.CategoryID, request.AttributeID, request.IsRequired)
	if err != nil {
		response.ErrorJson(w, 500, "kategoriye attribute eklenirken sunucu kaynaklı hata meydana geldi", fmt.Errorf("err :%w", err))
		return
	}

	response.WriteJson(w, 201, true, "basarili")
}

func (h *CategoryHandler) RemoveAttributeFromCategory(w http.ResponseWriter, r *http.Request) {

	//var req dto.RemoveAttributeFromCategoryRequest -> rest mimarisinde delete isteklerinde body olmamalı !
	//err := response.ReadJson(w, r, &req) // üzerine işlem yapacağı ve değiştireceği için pointerini veriyorum
	//if err != nil {
	//	response.ErrorJson(w, http.StatusBadRequest, "json okunurken hata meydana geldi", err)
	//	return
	//}

	categoryId := r.PathValue("categoryId")
	categoryID, err := strconv.Atoi(categoryId)
	if err != nil {
		return
	}
	attributeId := r.PathValue("attributeId")
	attributeID, err := strconv.Atoi(attributeId)
	if err != nil {
		return
	}
	err = h.categoryService.RemoveAttributeFromCategory(r.Context(), categoryID, attributeID)
	if err != nil {
		response.ErrorJson(w, 500, "sunucu kaynaklı problem", fmt.Errorf("err : %w", err))
		return
	}

	response.WriteJson(w, 200, true, "kaldirma islemi basarili")
}

func (h *CategoryHandler) UpdateAttributeFromCategory(w http.ResponseWriter, r *http.Request) {

	categoryID := r.PathValue("categoryId")
	categoryId, err := strconv.Atoi(categoryID)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "path degerleri  okunurken hata meydana geldi ", err)

		return
	}
	attributeID := r.PathValue("attributeId")
	attributeId, err := strconv.Atoi(attributeID)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "path degerleri  okunurken hata meydana geldi ", err)
		return
	}

	var req dto.UpdateAttributeFromCategoryRequest
	err = response.ReadJson(w, r, &req) // üzerine işlem yapacağı ve değiştireceği için pointerini veriyorum
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "json okunurken hata meydana geldi", err)
		return
	}
	err = h.categoryService.UpdateAttributeToCategory(r.Context(), req.IsRequired, attributeId, categoryId)

	if err != nil {
		response.ErrorJson(w, 500, "sunucu kaynaklı problem", fmt.Errorf("err : %w", err))
		return
	}

	response.WriteJson(w, 200, true, "güncelleme islemi basarili")
}
