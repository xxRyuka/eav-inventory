package domain

// Stock depodaki fiziksel toplam mal = availableQuantity + reservedQuantity
type Stock struct {
	WarehouseID       int // hangi depoda
	ProductID         int // hangi ürün
	AvaliableQuantity int
	ReservedQuantity  int
}
