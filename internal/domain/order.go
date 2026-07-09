package domain

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "PENDING"
	OrderStatusPaid     OrderStatus = "PAID"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

type OrderItem struct {
	Product  *Product
	Quantity int
	Price    float64
}

func NewOrderItem(product *Product, quantity int, price float64) *OrderItem {
	return &OrderItem{
		Product:  product,
		Quantity: quantity,
		Price:    price,
	}
}

type Order struct {
	ID      string
	Cliente string
	Items   []*OrderItem
	status  OrderStatus
}

func NewOrder(id string, cliente string) *Order {
	return &Order{
		ID:      id,
		Cliente: cliente,
		Items:   []*OrderItem{},
		status:  OrderStatusPending,
	}
}

func TotalSum(orders []*Order) float64 {
	var total float64
	for _, order := range orders {
		for _, item := range order.Items {
			total += item.Price * float64(item.Quantity)
		}
	}
	return total
}

func PayOrder(order *Order) {
	order.status = OrderStatusPaid
}

func CancelOrder(order *Order) {
	order.status = OrderStatusCanceled
}
