package models

import "errors"

var (
	ErrNotFound        = errors.New("entity not found")
	ErrUniqueViolation = errors.New("entity already exists")
	ErrCartEmpty       = errors.New("cart is empty")
)
