package command

type PurchaseInCommand struct {
	WarehouseID  int
	ProductID    int
	Quantity     int
	MovementType string // bunu constlardan alcaz !!
}
