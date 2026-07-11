package main

import (
	"fmt"

	"github.com/andreluialves/shop-orders/internal/domain"
	"github.com/andreluialves/shop-orders/internal/repository"
	"github.com/andreluialves/shop-orders/internal/service"
)

func main() {

	// Repositories
	productRepository := repository.NewMemoryProductRepository()
	orderRepository := repository.NewMemoryOrderRepository()

	// Service
	orderService := service.NewOrderService(
		productRepository,
		orderRepository,
	)

	// Produtos
	products := []*domain.Product{}

	notebook, err := domain.NewProduct("P001", "Notebook", 3500.00, 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	mouse, err := domain.NewProduct("P002", "Mouse", 80.00, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	teclado, err := domain.NewProduct("P003", "Teclado", 180.00, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	products = append(products, notebook, mouse, teclado)

	for _, product := range products {
		if err := productRepository.Save(product); err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("=== PRODUTOS CADASTRADOS ===")

	products, err = productRepository.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, product := range products {
		fmt.Printf(
			"%s | %s | R$ %.2f | Estoque: %d\n",
			product.ID,
			product.Name,
			product.Price,
			product.Quantity,
		)
	}

	fmt.Println()

	//--------------------------------------------------
	// Pedido válido
	//--------------------------------------------------

	fmt.Println("=== CRIANDO PEDIDO ===")

	order, err := orderService.CreateOrder(
		"Ana",
		[]service.CreateOrderItem{
			{
				ProductID: "P001",
				Quantity:  1,
			},
			{
				ProductID: "P002",
				Quantity:  2,
			},
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Pedido: %s\n", order.ID)
	fmt.Printf("Cliente: %s\n", order.Customer)
	fmt.Printf("Status: %s\n", order.Status())
	fmt.Printf("Total: R$ %.2f\n", order.TotalSum())

	fmt.Println()

	//--------------------------------------------------
	// Estoque atualizado
	//--------------------------------------------------

	fmt.Println("=== ESTOQUE ATUAL ===")

	for _, product := range products {
		fmt.Printf(
			"%s | Estoque: %d\n",
			product.Name,
			product.Quantity,
		)
	}

	fmt.Println()

	//--------------------------------------------------
	// Pagamento
	//--------------------------------------------------

	fmt.Println("=== PAGANDO PEDIDO ===")

	if err := orderService.PayOrder(order.ID); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Pedido pago com sucesso.")
	}

	fmt.Println()

	//--------------------------------------------------
	// Buscar pedido
	//--------------------------------------------------

	fmt.Println("=== BUSCANDO PEDIDO ===")

	foundOrder, err := orderService.FindOrderByID(order.ID)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(
			"Pedido %s encontrado com status %s\n",
			foundOrder.ID,
			foundOrder.Status(),
		)
	}

	fmt.Println()

	//--------------------------------------------------
	// Cancelar pedido pago
	//--------------------------------------------------

	fmt.Println("=== CANCELANDO PEDIDO PAGO ===")

	if err := orderService.CancelOrder(order.ID); err != nil {
		fmt.Println(err)
	}

	fmt.Println()

	//--------------------------------------------------
	// Estoque insuficiente
	//--------------------------------------------------

	fmt.Println("=== ESTOQUE INSUFICIENTE ===")

	_, err = orderService.CreateOrder(
		"Maria",
		[]service.CreateOrderItem{
			{
				ProductID: "P001",
				Quantity:  100,
			},
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()

	//--------------------------------------------------
	// Pedido inválido
	//--------------------------------------------------

	fmt.Println("=== PEDIDO INVÁLIDO ===")

	_, err = orderService.CreateOrder(
		"",
		[]service.CreateOrderItem{},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()

	//--------------------------------------------------
	// Filtro: Pedidos pagos
	//--------------------------------------------------

	fmt.Println("=== FILTRO: PEDIDOS PAGOS ===")

	paidOrders, err := orderService.FilterOrders(
		service.PaidOrders(),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, order := range paidOrders {
		fmt.Printf(
			"Pedido: %s | Cliente: %s | Status: %s | Total: %.2f\n",
			order.ID,
			order.Customer,
			order.Status(),
			order.TotalSum(),
		)
	}

	fmt.Println()

	//--------------------------------------------------
	// Filtro: Pedidos pendentes
	//--------------------------------------------------

	fmt.Println("=== FILTRO: PEDIDOS PENDENTES ===")

	pendingOrders, err := orderService.FilterOrders(
		service.PendingOrders(),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, order := range pendingOrders {
		fmt.Printf(
			"Pedido: %s | Cliente: %s | Status: %s | Total: %.2f\n",
			order.ID,
			order.Customer,
			order.Status(),
			order.TotalSum(),
		)
	}
}
