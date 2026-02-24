package enums

import (
	"database/sql/driver"
	"encoding/json/jsontext"
	"fmt"
)

type OrderSource struct {
	slug string
}

func NewOrderSource(s string) (OrderSource, error) {
	switch s {
	case OrderSourceLanding.slug:
		return OrderSourceLanding, nil
	case OrderSourceSPA.slug:
		return OrderSourceSPA, nil
	case OrderSourceTMA.slug:
		return OrderSourceTMA, nil
	default:
		return OrderSource{}, fmt.Errorf("unknown order source: %s", s)
	}
}

var (
	OrderSourceLanding = OrderSource{slug: "landing"}
	OrderSourceSPA     = OrderSource{slug: "spa"}
	OrderSourceTMA     = OrderSource{slug: "tma"}
)

func (u *OrderSource) String() string {
	return u.slug
}

func (u *OrderSource) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("can not assert order source to string")
	}
	r, err := NewOrderSource(s)
	if err != nil {
		return err
	}
	*u = r
	return nil
}

func (u OrderSource) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u OrderSource) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(u.slug))
}

func (u *OrderSource) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '"' {
		return fmt.Errorf("order source must be a JSON string")
	}
	e, err := NewOrderSource(tok.String())
	if err != nil {
		return err
	}
	*u = e
	return nil
}
