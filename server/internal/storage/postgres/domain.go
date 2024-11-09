package postgres

import (
	entity "OzonHW1/pkg/order_entity"
	"context"
)

type OrderService interface {
	AddOrder(ctx context.Context, ID, userID, shelfLife int64, packaging string, weigh float32, cost float32) error
	AcceptReturn(ctx context.Context, UserID, OrderID int64) error
	ReturnOrderToCourier(ctx context.Context, orderID int64) error
	UpdateIssued(ctx context.Context, orderIDs []string) error
	ListOrders(ctx context.Context, userID int64, flag *entity.Flag) ([]*entity.Order, error)
	GetReturned(ctx context.Context, index int64, size int64) []*entity.Order
}
