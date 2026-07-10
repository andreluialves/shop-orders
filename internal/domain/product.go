package domain

type Product struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}

func NewProduct(ID string, Name string, Price float64, Quantity int) (*Product, error) {
	product := Product{
		ID:       ID,
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
	if p.ID == "" {
		return ErrProductIDInvalid
	}

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

func (p *Product) RestoreQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if p.Quantity < quantity {
		return ErrInsufficientQuantity
	}

	p.Quantity += quantity
	return nil
}

func (p *Product) ReduceQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if p.Quantity < quantity {
		return ErrInsufficientQuantity
	}

	p.Quantity -= quantity
	return nil
}
