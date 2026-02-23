package enums

import (
	"database/sql/driver"
	"encoding/json/jsontext"
	"fmt"
)

type ProductType struct {
	slug string
}

func NewProductType(s string) (ProductType, error) {
	switch s {
	case ProductTypeProduct.slug:
		return ProductTypeProduct, nil
	case ProductTypeService.slug:
		return ProductTypeService, nil
	default:
		return ProductType{}, fmt.Errorf("unknown product type: %s", s)
	}
}

var (
	ProductTypeProduct = ProductType{slug: "product"}
	ProductTypeService = ProductType{slug: "service"}
)

func (u *ProductType) String() string {
	return u.slug
}

func (u *ProductType) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("can not assert product type to string")
	}
	r, err := NewProductType(s)
	if err != nil {
		return err
	}
	*u = r
	return nil
}

func (u ProductType) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u ProductType) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(u.slug))
}

func (u *ProductType) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '"' {
		return fmt.Errorf("product type must be a JSON string")
	}
	e, err := NewProductType(tok.String())
	if err != nil {
		return err
	}
	*u = e
	return nil
}
