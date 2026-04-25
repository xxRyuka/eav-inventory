package domain

type MovementType string

const (
	PurchaseIn MovementType = "PURCHASE_IN" // Satin alma işlemi - tedarikçiden gelen mal
	TransferIn MovementType = "TRANSFER_IN" // Depoya transfer girişi

	OrderOut    MovementType = "ORDER_OUT"    // Musteriye mal gitti
	TransferOut MovementType = "TRANSFER_OUT" // Depodan mal çıktı

	Adjustment MovementType = "ADJUSTMENT" // Stok düzeltme (sayım sonrası)
)

// StockMovements 1 işlem için 2 db rows yazılır ornegin bi depoda mal girdiyse öbür depodan cıkmıstır
type StockMovements struct {
	WarehouseID  int // hangi depoda gerceklesti
	ProductID    int // Hangi ürün
	Quantity     int // Kaç Tane
	MovementType MovementType

	//ReferenceID int tetikleyen işlemin id'sini tutacağım burda
}
