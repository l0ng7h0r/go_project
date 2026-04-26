package usecase

import (
	"errors"

	"github.com/l0ng7h0r/golang/internal/domain"
	"github.com/l0ng7h0r/golang/internal/repository"
)

type OrderUsecase struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewOrderUsecase(orderRepo *repository.OrderRepository, cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *OrderUsecase {
	return &OrderUsecase{orderRepo: orderRepo, cartRepo: cartRepo, productRepo: productRepo}
}

// CreateOrderFromCart creates an order from the user's cart and clears the cart
func (u *OrderUsecase) CreateOrderFromCart(
	userID, receiverName, phone, province, district, logistic, logisticBranch string,
) (string, error) {
	cart, err := u.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		return "", err
	}
	if len(cart.Items) == 0 {
		return "", errors.New("cart is empty")
	}

	var totalPrice float64
	var orderItems []domain.OrderItem

	for _, item := range cart.Items {
		product, err := u.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return "", errors.New("product not found: " + item.ProductID)
		}
		if product.Stock < item.Quantity {
			return "", errors.New("insufficient stock for product: " + product.Name)
		}
		totalPrice += product.Price * float64(item.Quantity)
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	order := &domain.Order{
		UserID:         userID,
		TotalPrice:     totalPrice,
		Status:         "pending",
		ReceiverName:   receiverName,
		Phone:          phone,
		Province:       province,
		District:       district,
		Logistic:       logistic,
		LogisticBranch: logisticBranch,
		OrderItems:     orderItems,
	}

	orderID, err := u.orderRepo.CreateOrder(order)
	if err != nil {
		return "", err
	}

	// Clear cart after successful order
	_ = u.cartRepo.ClearCart(cart.ID)

	return orderID, nil
}

func (u *OrderUsecase) GetOrderByID(id string) (*domain.Order, error) {
	return u.orderRepo.GetOrderByID(id)
}

func (u *OrderUsecase) GetMyOrders(userID string) ([]domain.Order, error) {
	return u.orderRepo.GetOrdersByUserID(userID)
}

func (u *OrderUsecase) GetAllOrders() ([]domain.Order, error) {
	return u.orderRepo.GetAllOrders()
}

func (u *OrderUsecase) UpdateOrderStatus(id string, status string) error {
	validStatuses := map[string]bool{
		"pending": true, "confirmed": true, "shipped": true,
		"delivered": true, "cancelled": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// If cancelled, restore stock
	if status == "cancelled" {
		order, err := u.orderRepo.GetOrderByID(id)
		if err == nil && order.Status != "cancelled" {
			for _, item := range order.OrderItems {
				_ = u.productRepo.RestoreStock(item.ProductID, item.Quantity)
			}
		}
	}

	return u.orderRepo.UpdateOrderStatus(id, status)
}
