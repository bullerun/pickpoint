package postgres

import (
	entity "OzonHW1/pkg/order_entity"
	"OzonHW1/server/internal/imdb"
	"OzonHW1/server/internal/metrics"
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"sync"
)

const returnedCacheKey = "Returned"

type PgRepository struct {
	txManager TransactionManager

	returnedOrderLock *sync.RWMutex
	returnedCache     *imdb.LRUCache[string, []*entity.Order]
	orderLock         *sync.RWMutex
	orderCache        *imdb.LRUCache[int64, int64]
}

func NewPgRepository(txManager TransactionManager, capacity int) *PgRepository {
	return &PgRepository{
		txManager:         txManager,
		returnedOrderLock: &sync.RWMutex{},
		returnedCache:     imdb.NewLRUCache[string, []*entity.Order](capacity),
		orderLock:         &sync.RWMutex{},
		orderCache:        imdb.NewLRUCache[int64, int64](capacity),
	}
}

func (r *PgRepository) AddOrder(ctx context.Context, ID, userID, shelfLife int64, packaging string, weigh float32, cost float32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Order.AddOrder")
	defer span.Finish()
	r.orderLock.RLock()
	_, exist := r.orderCache.Get(ID)
	r.orderLock.RUnlock()
	if exist {
		return ErrOrderAlreadyExists
	}
	tx := r.txManager.GetQueryEngine(ctx)
	result, err := tx.Exec(ctx, `
	INSERT INTO orders (id, user_id, expiry_date, packaging, weigh, cost)  
	select $1, $2, current_timestamp + interval '1 days' * $3 , $4, $5, $6
	on conflict (id) DO NOTHING`,
		ID, userID, shelfLife, packaging, weigh, cost)
	if err != nil {
		return err
	}
	r.orderLock.Lock()
	r.orderCache.Put(ID, userID)
	r.orderLock.Unlock()
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrOrderAlreadyExists
	}
	return nil
}

func (r *PgRepository) UpdateIssued(ctx context.Context, orderIDs []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Order.UpdateIssued")
	defer span.Finish()
	var found bool
	var referenceUserID int64
	var userID int64

	i := 0
	r.orderLock.RLock()
	for ; i < len(orderIDs); i++ {
		parseInt, _ := strconv.ParseInt(orderIDs[i], 10, 64)
		referenceUserID, found = r.orderCache.Get(parseInt)
		if found {
			break
		}
	}
	for ; i < len(orderIDs); i++ {
		parseInt, _ := strconv.ParseInt(orderIDs[i], 10, 64)
		userID, found = r.orderCache.Get(parseInt)
		if found && (userID != referenceUserID || userID == 0) {
			r.orderLock.RUnlock()
			return ErrDifferentUserIDs
		}
	}

	r.orderLock.RUnlock()

	tx := r.txManager.GetQueryEngine(ctx)
	countOrders := len(orderIDs)
	strOrderIDs := strings.Join(orderIDs, ", ")
	query := fmt.Sprintf(`
	UPDATE orders
	SET  accept_return_order_date = current_timestamp + interval '2 days'
	WHERE orders.id IN (  %s  ) AND EXISTS (
		SELECT 1
		FROM orders
		WHERE id IN ( %s ) AND accept_return_order_date is null AND expiry_date >= current_timestamp
		GROUP BY user_id
		HAVING COUNT(DISTINCT id) = $1
	);
	`, strOrderIDs, strOrderIDs)

	result, err := tx.Exec(ctx, query, countOrders)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.Wrap(ErrNothingHappened, "orders could not be issued")
	}
	metrics.IssueOrder(float64(rowsAffected))

	for i := 0; i < len(orderIDs); i++ {
		parseInt, _ := strconv.ParseInt(orderIDs[i], 10, 64)
		r.orderCache.Put(parseInt, referenceUserID)
	}

	return nil
}

func (r *PgRepository) ReturnOrderToCourier(ctx context.Context, orderID int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Order.ReturnOrderToCourier")
	defer span.Finish()
	tx := r.txManager.GetQueryEngine(ctx)
	result, err := tx.Exec(ctx, `
	UPDATE orders SET returned_to_courier_at = current_timestamp where id = $1 and returned_to_courier_at is null and (accept_return_order_date is null or returned_from_client_at is not null);`, orderID)
	if err != nil {
		return err
	}
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return errors.Wrap(ErrNothingHappened, "it was not possible to return the order to the courier")
	}

	r.returnedOrderLock.Lock()
	r.returnedCache.DeleteAll()
	r.returnedOrderLock.Unlock()
	return nil
}

func (r *PgRepository) ListOrders(ctx context.Context, userID int64, flag *entity.Flag) ([]*entity.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Order.ListOrders")
	defer span.Finish()
	var isIn string
	if flag.InTheDeliveryPoint {
		isIn = "and (returned_from_client_at is not null or accept_return_order_date is null)"
	}
	var orders []*entity.Order
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, fmt.Sprintf("SELECT * FROM orders WHERE user_id = $1 %s ORDER BY created_at DESC LIMIT $2;", isIn), userID, flag.Latest)
	if err == nil {
		r.orderLock.Lock()
		for _, el := range orders {
			r.orderCache.Put(el.ID, el.UserID)
		}
		r.orderLock.Unlock()
	}
	return orders, err
}
func (r *PgRepository) AcceptReturn(ctx context.Context, UserID, OrderID int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Order.AcceptReturn")
	defer span.Finish()
	tx := r.txManager.GetQueryEngine(ctx)
	result, err := tx.Exec(ctx, "UPDATE orders SET returned_from_client_at = current_timestamp WHERE id = $1 AND user_id = $2 AND accept_return_order_date >= current_timestamp;", OrderID, UserID)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.Wrap(ErrNothingHappened, "it was not possible to return the order")
	}

	r.returnedOrderLock.Lock()
	r.returnedCache.DeleteAll()
	r.returnedOrderLock.Unlock()
	return nil
}

func (r *PgRepository) GetReturned(ctx context.Context, offset int64, limit int64) []*entity.Order {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository.Order.GetReturned")
	defer span.Finish()
	key := fmt.Sprintf("%s:%d:%d", returnedCacheKey)
	r.returnedOrderLock.RLock()
	orders, found := r.returnedCache.Get(key)
	r.returnedOrderLock.RUnlock()
	if found {
		return orders
	}

	query := `SELECT * FROM orders WHERE returned_from_client_at IS NOT null or returned_to_courier_at IS NOT NULL  ORDER BY created_at DESC limit $1 OFFSET $2;`
	tx := r.txManager.GetQueryEngine(ctx)
	_ = pgxscan.Select(ctx, tx, &orders, query, limit, offset)

	r.returnedOrderLock.Lock()
	r.returnedCache.Put(key, orders)
	r.returnedOrderLock.Unlock()
	return orders
}
