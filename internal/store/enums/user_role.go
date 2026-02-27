package enums

import (
	"database/sql/driver"
	"encoding/json/jsontext"
	"fmt"
)

type UserRole struct {
	slug string
}

func NewUserRole(s string) (UserRole, error) {
	switch s {
	case UserRoleUser.slug:
		return UserRoleUser, nil
	case UserRoleManager.slug:
		return UserRoleManager, nil
	case UserRoleModerator.slug:
		return UserRoleModerator, nil
	case UserRoleAdmin.slug:
		return UserRoleAdmin, nil
	default:
		return UserRole{}, fmt.Errorf("unknown user role: %s", s)
	}
}

var (
	UserRoleUser      = UserRole{slug: "user"}
	UserRoleManager   = UserRole{slug: "manager"}
	UserRoleModerator = UserRole{slug: "moderator"}
	UserRoleAdmin     = UserRole{slug: "admin"}
)

func (u *UserRole) String() string {
	return u.slug
}

func (u *UserRole) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("can not assert user role to string")
	}
	r, err := NewUserRole(s)
	if err != nil {
		return err
	}
	*u = r
	return nil
}

func (u UserRole) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u UserRole) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteToken(jsontext.String(u.slug))
}

func (u *UserRole) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != '"' {
		return fmt.Errorf("user role status must be a JSON string")
	}
	e, err := NewUserRole(tok.String())
	if err != nil {
		return err
	}
	*u = e
	return nil
}
