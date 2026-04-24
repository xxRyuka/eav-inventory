package catalog

import (
	"errors"
	"strings"
)

var (
	ErrAttributeCodeEmpty     = errors.New("nitelik kodu (code) bos olamaz")
	ErrAttributeNameEmpty     = errors.New("nitelik adi (name) bos olamaz")
	ErrAttributeDataTypeEmpty = errors.New("nitelik veri tipi (data_type) bos olamaz")
)

type Attribute struct {
	ID       int
	Code     string
	Name     string
	DataType DataType
}

func (a *Attribute) Validate() error {

	if strings.TrimSpace(a.Code) == "" {
		return ErrAttributeCodeEmpty
	}
	if strings.TrimSpace(a.Name) == "" {
		return ErrAttributeNameEmpty
	}
	if strings.TrimSpace(string(a.DataType)) == "" {
		return ErrAttributeDataTypeEmpty
	}
	return nil
}
