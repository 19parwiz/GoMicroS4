package usecase

import (
	"context"
	"github.com/19parwiz/order-service/internal/domain"
)

type AutoIncRepo interface {
	Next(ctx context.Context, coll string) (uint64, error)
}

// OrderRepository defines the contract for order data operations
type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) error
	Update(ctx context.Context, filter domain.OrderFilter, update domain.OrderUpdateData) error
	GetWithFilter(ctx context.Context, filter domain.OrderFilter) (domain.Order, error)
	GetAllWithFilter(ctx context.Context, filter domain.OrderFilter, page, limit int64) ([]domain.Order, int64, error)
	Delete(ctx context.Context, filter domain.OrderFilter) error
}

type InventoryClient interface {
	GetProduct(ctx context.Context, productID uint64) (domain.Product, error)
}

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, event domain.Order) error
}
