package response

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
type PagedRespondeData struct {
	Data       any `json:"pagedData,omitempty"`
	TotalCount int `json:"totalCount,omitempty"`
	PageSize   int `json:"pageSize,omitempty"`
	Page       int `json:"page,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

func CalculatedPagedResponse(data any, totalCount, pageSize, page int) *PagedRespondeData {
	return &PagedRespondeData{
		Data:       data,
		TotalCount: totalCount,
		PageSize:   pageSize,
		Page:       page,
		TotalPages: int(math.Ceil(float64(totalCount) / float64(pageSize))),
	}
}

// WriteJson Basarili durum
func WriteJson(w http.ResponseWriter, status int, data any, message string) {
	envelope := ApiResponse{
		Success: true,
		Data:    data,
		Message: message,
		//Error:   "", // basarili durum donuyoruz jsonda bulunmasına gerek yok
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Yanıtı doğrudan ağ soketine (w) yaz (RAM'de fazladan []byte tahsis etmeden)
	err := json.NewEncoder(w).Encode(envelope)
	if err != nil {
		fmt.Println("Json Yanıtı Olusturulurken Hata Olustu")
		return
	}
}

func ErrorJson(w http.ResponseWriter, status int, message string, err error) {
	envelope := ApiResponse{
		Success: false,
		//Data:    nil,
		Message: message,
		Error:   err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Yanıtı doğrudan ağ soketine (w) yaz (RAM'de fazladan []byte tahsis etmeden)

	ok := json.NewEncoder(w).Encode(envelope)
	if ok != nil {
		fmt.Println("Json Yanıtı Olusturulurken Hata Olustu")
		return
	}
}

// ReadJson gelen HTTP Request body'sini okuyup hedef DTO struct'ına dönüştürür (Bind).
// Go'da any (eski adıyla interface{}) her türlü tipi tutabilen bir kutudur.
// Biz bu ReadJson fonksiyonunu Handler (Controller) içinden çağırırken, zaten nesnenin bellek adresini (pointer) göndereceğiz.
// err := response.ReadJson(w, r, &requestDTO) <= örneğin
func ReadJson(w http.ResponseWriter, r *http.Request, destination any) error {
	// Gelen isteğin boyutunu sınırla (Güvenlik için: maks 1MB)
	r.Body = http.MaxBytesReader(w, r.Body, int64(1048576))

	decoder := json.NewDecoder(r.Body)
	// Gelen JSON'da bizim DTO'da olmayan bir alan varsa hata fırlat (Strict parsing)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(destination) // Burda neden & kullanmadın interface oldugu için mi ?
	if err != nil {
		return err
	}
	return nil
}

//Dışarıdan zaten &requestDTO (pointer) geldiği için, destination değişkeninin içi doğrudan o nesnenin bellek adresini tutuyor.
//Eğer içeride bir daha &destination yapsaydık:
//"Pointer'ı tutan Interface'in Pointer'ı" gibi saçma bir çift referans (double pointer) yaratmış olurduk
//ve json.Decode hata verirdi.
