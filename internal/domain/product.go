package domain

type Product struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}

func (p *Product) UpdateName(newName string) {
	p.Name = newName
}

func (p *Product) UpdatePrice(newPrice float64) {
	p.Price = newPrice
}

func (p *Product) RestoreQuantity(quantity int) {
	p.Quantity += quantity
}

func (p *Product) ReduceQuantity(quantity int) {
	p.Quantity -= quantity
}
