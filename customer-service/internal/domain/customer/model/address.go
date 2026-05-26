package model

import "fmt"

// Address is a value object — immutable, no identity.
type Address struct {
	Line1      string
	Line2      string
	City       string
	PostalCode string
	Country    string
}

// NewAddress validates and creates an Address value object.
func NewAddress(line1, line2, city, postalCode, country string) (Address, error) {
	if line1 == "" {
		return Address{}, fmt.Errorf("line1 is required")
	}
	if city == "" {
		return Address{}, fmt.Errorf("city is required")
	}
	if country == "" {
		return Address{}, fmt.Errorf("country is required")
	}
	return Address{
		Line1:      line1,
		Line2:      line2,
		City:       city,
		PostalCode: postalCode,
		Country:    country,
	}, nil
}
