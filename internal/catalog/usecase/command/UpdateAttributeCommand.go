package command

type UpdateAttributeCommand struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}
