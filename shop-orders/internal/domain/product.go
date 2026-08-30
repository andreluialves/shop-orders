package domain

type Product struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}

func NewProduct(Name string, Price float64, Quantity int) (*Product, error) {
	product := Product{
		Name:     Name,
		Price:    Price,
		Quantity: Quantity,
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	return &product, nil
}

func (p Product) Validate() error {
	// if p.ID == "" {
	// 	return ErrProductIDInvalid
	// }

	if p.Name == "" {
		return ErrProductNameInvalid
	}

	if p.Price <= 0 {
		return ErrProductPriceInvalid
	}

	if p.Quantity < 0 {
		return ErrInvalidQuantity
	}

	return nil
}

func RestoreProduct(id string, name string, price float64, quantity int) *Product {
	return &Product{
		ID:       id,
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}
}

func (p *Product) RestoreQuantity(quantity int) error {
	if err := p.ValidateQuantity(quantity); err != nil {
		return err
	}

	p.Quantity += quantity
	return nil
}

func (p *Product) ReduceQuantity(quantity int) error {
	if err := p.ValidateQuantity(quantity); err != nil {
		return err
	}

	p.Quantity -= quantity
	return nil
}

func (p Product) ValidateQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if p.Quantity < quantity {
		return ErrInsufficientQuantity
	}

	return nil
}
