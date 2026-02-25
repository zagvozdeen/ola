package enums

import (
	"database/sql/driver"
	"encoding/json/jsontext"
	"fmt"
)

type RequestStatus struct {
	slug  string
	emoji string
	label string
}

func NewRequestStatus(s string) (RequestStatus, error) {
	switch s {
	case RequestStatusCreated.slug:
		return RequestStatusCreated, nil
	case RequestStatusInProgress.slug:
		return RequestStatusInProgress, nil
	case RequestStatusReviewed.slug:
		return RequestStatusReviewed, nil
	default:
		return RequestStatus{}, fmt.Errorf("unknown request status: %s", s)
	}
}

var (
	RequestStatusCreated    = RequestStatus{slug: "created", emoji: "🆕", label: "Новый"}
	RequestStatusInProgress = RequestStatus{slug: "in_progress", emoji: "💼", label: "В работе"}
	RequestStatusReviewed   = RequestStatus{slug: "reviewed", emoji: "✅", label: "Завершён"}
)

func (s RequestStatus) String() string {
	return s.slug
}

func (s RequestStatus) Emoji() string {
	return s.emoji
}

func (s RequestStatus) Label() string {
	return s.label
}

func (s *RequestStatus) Scan(src any) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("can not assert request status to string")
	}

	status, err := NewRequestStatus(str)
	if err != nil {
		return err
	}
	*s = status
	return nil
}

func (s RequestStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s RequestStatus) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(s.slug))
}

func (s *RequestStatus) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '"' {
		return fmt.Errorf("request status must be a JSON string")
	}
	status, err := NewRequestStatus(tok.String())
	if err != nil {
		return err
	}
	*s = status
	return nil
}
