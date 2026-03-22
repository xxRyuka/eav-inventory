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

type ProductHandler struct {
	service usecase.ProductUseCase
}

func NewProductHandler(useCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{service: useCase}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.ErrorJson(w, http.StatusMethodNotAllowed, "gecersiz method", fmt.Errorf("Gecersiz Method"))
		return
	}
	req := dto.CreateProductRequest{}
	err := response.ReadJson(w, r, &req)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "bind edilemedi", fmt.Errorf("json okunamadı %w", err))
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

func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {

	idPath := r.PathValue("id")

	id, err := strconv.Atoi(strings.TrimSpace(idPath))
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "0'dan büyük tam sayi bir id giriniz (path)", fmt.Errorf("%w", err))
		return
	}

	product, err := h.service.GetProductById(r.Context(), id)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "Veritabaninda hata meydana geldi", fmt.Errorf("olusan hata : %w", err))
		return
	}

	resp := dto.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		CategoryID: product.CategoryId,
		SKU:        product.SKU,
		//Attributes: nil,
	}
	var attrs []dto.ProductAttributeResponse
	for i, value := range product.AttributeValues {
		attr := dto.ProductAttributeResponse{
			AttributeID: product.AttributeValues[i].AttributeID,
			Code:        value.Attribute.Code,
			Name:        value.Attribute.Name,
			DataType:    string(value.Attribute.DataType),
			Value:       value.Value,
		}
		attrs = append(attrs, attr)

	}
	resp.Attributes = attrs

	response.WriteJson(w, http.StatusOK, resp, "")
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")
	//ram_code := r.URL.Query().Get("ram_code") // "50,v,a"  dizisini string olarak donmus burda. ama biz 2 tane ram_code gondermistik ve sadece ilkini verdi
	// yani burda gercekten map[string][]string kullanmamız lazım diger degerleride yakalamak için

	// attributeMap["key"] bize key'e ait value's stringini donecek queryden cekiyoruz
	attributeMap := make(map[string][]string)

	// Queryi Cekerken ilk parametre query Adı(key) ikinci parametre query'nin valuesleri
	for key, values := range r.URL.Query() {
		if key == "page" || key == "pageSize" {
			continue
		}
		attributeMap[key] = values
	}

	pageVal, err := strconv.Atoi(page)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "Query parametreleri yanlıs", err)
		return
	}
	pageSizeVal, err := strconv.Atoi(pageSize)
	if err != nil {
		response.ErrorJson(w, http.StatusBadRequest, "Query parametreleri yanlıs", err)

		return
	}

	products, total, err := h.service.GetProducts(r.Context(), pageSizeVal, pageVal, attributeMap)
	if err != nil {
		response.ErrorJson(w, http.StatusInternalServerError, "ürünler veritabanından cekilirken hata olustu", fmt.Errorf("Hata : %w", err))
		return
	}

	var resp []dto.ProductResponse

	for _, product := range products {
		prodDto := dto.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			CategoryID: product.CategoryId,
			SKU:        product.SKU,
			//Attributes: product.AttributeValues,
		}
		var pavs []dto.ProductAttributeResponse
		for _, pav := range product.AttributeValues {
			pavDto := dto.ProductAttributeResponse{
				Code:        pav.Attribute.Code,
				Name:        pav.Attribute.Name,
				DataType:    string(pav.Attribute.DataType),
				AttributeID: pav.AttributeID,
				Value:       pav.Value,
			}
			pavs = append(pavs, pavDto)

		}
		prodDto.Attributes = pavs
		resp = append(resp, prodDto)
	}

	x := response.CalculatedPagedResponse(resp, total, pageSizeVal, pageVal)
	response.WriteJson(w, http.StatusOK, x, "")
}
