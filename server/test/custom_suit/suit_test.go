package custom_suit

import (
	"OzonHW1/server/internal/storage"
	"OzonHW1/server/internal/storage/postgres"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"log"
)

type ItemSuite struct {
	suite.Suite
	service postgres.OrderService
}

func initSQLRows() string {
	return `
INSERT INTO orders (id, user_id, expiry_date, packaging, weigh, cost, accept_return_order_date, returned_from_client_at)  
	VALUES (1, 2, current_timestamp + interval '5 days' , 'box', 2.5, 120, null,null),
	(3, 2, current_timestamp + interval '5 days' , 'box', 2.5, 120, null, null),
	(21, 2, current_timestamp + interval '5 days' , 'box', 2.5, 120, null, null),
	(22, 2, current_timestamp + interval '5 days' , 'box', 2.5, 120, null, null),
(23, 2, current_timestamp + interval '5 days' , 'box', 2.5, 120, null, null),
(24, 3, current_timestamp + interval '5 days' , 'box', 2.5, 120, null, null),
(25, 3, current_timestamp + interval '5 days' , 'box', 2.5, 120, current_timestamp + interval '2 days', null);
INSERT INTO orders(id, user_id, expiry_date, returned_to_courier_at, packaging, weigh, cost)
VALUES (10, 2, current_timestamp + interval '5 days',  current_timestamp + interval '4 days', 'box', 2.5, 120);
INSERT INTO orders (id, user_id, expiry_date, packaging, weigh, cost, accept_return_order_date, returned_from_client_at)VALUES 
(31, 3, current_timestamp + interval '5 days' , 'box', 2.5, 120, current_timestamp + interval '1 days',current_timestamp + interval '2 days');`
}
func newStorageFacade(pool *pgxpool.Pool) postgres.OrderService {
	txManager := postgres.NewTxManager(pool)

	pgRepository := postgres.NewPgRepository(txManager, 0)

	return storage.NewStorageFacade(txManager, pgRepository)
}
func (s *ItemSuite) SetupSuite() {
	const psqlDSN = "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, psqlDSN)
	if err != nil {
		s.T().Fatal(err)
	}
	s.service = newStorageFacade(pool)
	exec, err := pool.Exec(ctx, initSQLRows())
	if err != nil || exec.RowsAffected() == 0 {
		log.Fatal(err)
	}
}
func (s *ItemSuite) TestAddOrderSuccess() {
	err := s.service.AddOrder(context.Background(), 2, 2, 5, "box", 2.5, 120.0)
	s.Require().NoError(err)
}

func (s *ItemSuite) TestAddOrderFailed() {
	err := s.service.AddOrder(context.Background(), 1, 2, 5, "box", 2.5, 120.0)
	s.Require().Error(err)
}

func (s *ItemSuite) TestReturnOrderToCourierSuccess() {
	err := s.service.ReturnOrderToCourier(context.Background(), 3)
	s.Require().NoError(err)
}
func (s *ItemSuite) TestReturnOrderToCourierFailed() {
	err := s.service.ReturnOrderToCourier(context.Background(), 10)
	s.Require().Error(err)
}
func (s *ItemSuite) TestIssueOrderSuccess() {
	err := s.service.UpdateIssued(context.Background(), []string{"21", "22"})
	s.Require().NoError(err)
}
func (s *ItemSuite) TestIssueOrderFailed() {
	err := s.service.UpdateIssued(context.Background(), []string{"23", "24"})
	s.Require().Error(err)

}
func (s *ItemSuite) TestIssueOrderFailedDoesntExist() {
	err := s.service.UpdateIssued(context.Background(), []string{"30"})
	s.Require().Error(err)

}
func (s *ItemSuite) TestIssueOrderFailedNotInPickUp() {
	err := s.service.UpdateIssued(context.Background(), []string{"25"})
	s.Require().Error(err)
}
func (s *ItemSuite) TestAcceptReturnSuccess() {
	err := s.service.AcceptReturn(context.Background(), 3, 31)
	s.Require().NoError(err)
}
