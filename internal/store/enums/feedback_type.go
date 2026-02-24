package enums

import (
	"database/sql/driver"
	"encoding/json/jsontext"
	"fmt"
)

type FeedbackType struct {
	slug string
}

func NewFeedbackType(s string) (FeedbackType, error) {
	switch s {
	case FeedbackTypeManagerContact.slug:
		return FeedbackTypeManagerContact, nil
	case FeedbackTypePartnershipOffer.slug:
		return FeedbackTypePartnershipOffer, nil
	case FeedbackTypeFeedbackRequest.slug:
		return FeedbackTypeFeedbackRequest, nil
	default:
		return FeedbackType{}, fmt.Errorf("unknown feedback type: %s", s)
	}
}

var (
	FeedbackTypeManagerContact   = FeedbackType{slug: "manager_contact"}
	FeedbackTypePartnershipOffer = FeedbackType{slug: "partnership_offer"}
	FeedbackTypeFeedbackRequest  = FeedbackType{slug: "feedback_request"}
)

func (t *FeedbackType) String() string {
	return t.slug
}

func (t *FeedbackType) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("can not assert feedback type to string")
	}
	v, err := NewFeedbackType(s)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func (t FeedbackType) Value() (driver.Value, error) {
	return t.slug, nil
}

func (t FeedbackType) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(t.slug))
}

func (t *FeedbackType) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '"' {
		return fmt.Errorf("feedback type must be a JSON string")
	}
	v, err := NewFeedbackType(tok.String())
	if err != nil {
		return err
	}
	*t = v
	return nil
}
