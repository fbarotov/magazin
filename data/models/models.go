package models

import (
	"fmt"
	"strings"
)

type User struct {
	ID string `json:"-"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type Product struct {
	ID string `json:"-"`
	Name string `json:"name"`
	Price float64 `json:"price"`
}

func (p *Product) String() string {
	return fmt.Sprintf("ID: %s, Product: %s, %v", p.ID, p.Name, p.Price)
}

type Quantity struct {
	ProductID string
	Quantity int
}

type Purchase struct {
	ID string `json:"-"`
	Products []Product `json:"products"`
	Price float64 `json:"price"`
}

func (p *Purchase) String() string  {
	var sb strings.Builder
	sb.WriteString("Order Summary: \n")
	for _, prod := range p.Products {
		sb.WriteString(prod.String() + "\n")
	}
	sb.WriteString(fmt.Sprintf("Purchase Price: %v", p.Price))
	return sb.String()
}
