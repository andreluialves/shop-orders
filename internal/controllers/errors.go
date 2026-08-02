package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/andreluialves/shop-orders/internal/domain"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, domain.ErrOrderNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, domain.ErrInsufficientQuantity):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, domain.ErrOrderAlreadyPaid):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, domain.ErrOrderAlreadyCanceled):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, domain.ErrChangeStatusInvalid):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, domain.ErrCustomerNameRequired),
		errors.Is(err, domain.ErrCustomerNameTooShort),
		errors.Is(err, domain.ErrCustomerNameTooLong),
		errors.Is(err, domain.ErrProductNameRequired),
		errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrInvalidQuantity),
		errors.Is(err, domain.ErrInvalidCustomer),
		errors.Is(err, domain.ErrProductNameInvalid),
		errors.Is(err, domain.ErrProductPriceInvalid),
		errors.Is(err, domain.ErrEmptyOrder):

		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		log.Printf("internal error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
