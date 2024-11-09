package storage

import (
	entity "OzonHW1/pkg/order_entity"
	postgressql "OzonHW1/server/internal/storage/postgres"
	"context"
)

type storageFacade struct {
	txManager    postgressql.TransactionManager
	pgRepository *postgressql.PgRepository
}

func NewStorageFacade(
	txManager postgressql.TransactionManager,
	pgRepository *postgressql.PgRepository,
) *storageFacade {
	return &storageFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
	}
}

func (s *storageFacade) AddOrder(ctx context.Context, ID, userID, shelfLife int64, packaging string, weigh float32, cost float32) error {
	return s.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		err := s.pgRepository.AddOrder(ctx, ID, userID, shelfLife, packaging, weigh, cost)
		if err != nil {
			return err
		}
		return nil
	})
}
func (s *storageFacade) AcceptReturn(ctx context.Context, UserID, OrderID int64) error {
	return s.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		return s.pgRepository.AcceptReturn(ctx, UserID, OrderID)
	})
}

func (s *storageFacade) ReturnOrderToCourier(ctx context.Context, orderID int64) error {
	return s.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		return s.pgRepository.ReturnOrderToCourier(ctx, orderID)
	})
}
func (s *storageFacade) UpdateIssued(ctx context.Context, orderIDs []string) error {
	return s.txManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		return s.pgRepository.UpdateIssued(ctx, orderIDs)
	})
}
func (s *storageFacade) ListOrders(ctx context.Context, userID int64, flag *entity.Flag) ([]*entity.Order, error) {
	var orders []*entity.Order
	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		o, err := s.pgRepository.ListOrders(ctx, userID, flag)
		orders = o
		return err
	})
	if err != nil {
		return nil, err
	}
	return orders, err
}
func (s *storageFacade) GetReturned(ctx context.Context, offset int64, limit int64) []*entity.Order {
	var orders []*entity.Order
	err := s.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		o := s.pgRepository.GetReturned(ctx, offset, limit)
		orders = o
		return nil
	})
	if err != nil {
		return nil
	}
	return orders
}
