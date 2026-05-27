package model

import (
	"fmt"
	"regexp"
)

var vatRegexp = regexp.MustCompile(`^DE\d{9}$`)

// VATNumber represents a German VAT number (DE + 9 digits).
type VATNumber string

// NewVATNumber validates and creates a VATNumber value object.
func NewVATNumber(raw string) (VATNumber, error) {
	if !vatRegexp.MatchString(raw) {
		return "", fmt.Errorf("invalid VAT number format %q: must be DE followed by 9 digits", raw)
	}
	return VATNumber(raw), nil
}

// String returns the string representation.
func (v VATNumber) String() string {
	return string(v)
}
