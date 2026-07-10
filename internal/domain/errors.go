package domain

import "errors"

var ErrProductNotFound = errors.New("produto não encontrado")

var ErrProductIDInvalid = errors.New("ID do produto inválido")

var ErrProductNameInvalid = errors.New("nome do produto inválido")

var ErrProductPriceInvalid = errors.New("preço do produto inválido")

var ErrInvalidQuantity = errors.New("quantidade inválida")

var ErrInsufficientQuantity = errors.New("quantidade insuficiente")

var ErrOrderNotFound = errors.New("pedido não encontrado")

var ErrEmptyOrder = errors.New("pedido vazio")

var ErrChangeStatusInvalid = errors.New("mudança de status de pedido inválida")

var ErrInvalidCustomer = errors.New("cliente inválido")
