package dto

type ProductPurchased struct {
	Product
	Purchase
}

type Purchase struct {
	ProductID string
	Quantity int
}

type Product struct {
	Name string
	Price float64
}