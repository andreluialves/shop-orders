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

var ErrOrderAlreadyPaid = errors.New("pedido já pago")

var ErrOrderAlreadyCanceled = errors.New("pedido já cancelado")

var ErrChangeStatusInvalid = errors.New("mudança de status de pedido inválida")

var ErrInvalidCustomer = errors.New("cliente inválido")

var ErrInvalidCustomerID = errors.New("ID do cliente inválido")

var ErrProductIDAlreadyExists = errors.New("ID do produto já existe")

var ErrCustomerNameRequired = errors.New("nome do cliente é obrigatório")

var ErrCustomerNameTooShort = errors.New("nome do cliente é muito curto")

var ErrCustomerNameTooLong = errors.New("nome do cliente é muito longo")

var ErrProductNameRequired = errors.New("nome do produto é obrigatório")

var ErrInvalidPrice = errors.New("preço inválido")
