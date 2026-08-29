package domain

import "net/mail"

type Customer struct {
	ID      string
	Name    string
	Email   string
	Address string
	Phone   string
}

func NewCustomer(Name, Email, Address, Phone string) (*Customer, error) {
	customer := Customer{
		Name:    Name,
		Email:   Email,
		Address: Address,
		Phone:   Phone,
	}

	if err := customer.Validate(); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (c Customer) Validate() error {
	if c.Name == "" {
		return ErrCustomerNameRequired
	}

	if c.Email == "" {
		return ErrCustomerEmailRequired
	}

	if _, err := mail.ParseAddress(c.Email); err != nil {
		return ErrCustomerEmailInvalid
	}

	if c.Address == "" {
		return ErrCustomerAddressRequired
	}

	// if c.Phone == "" {
	// 	return ErrCustomerPhoneRequired
	// }

	return nil
}
