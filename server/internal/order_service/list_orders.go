package order_service

import (
	desc "OzonHW1/pkg/order-service/v1"
	entity "OzonHW1/pkg/order_entity"
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) ListOrders(ctx context.Context, req *desc.ListOrdersRequest) (*desc.ListOrdersResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Order.ListOrders")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	orders, err := s.storage.ListOrders(ctx, req.UserId, &entity.Flag{Latest: req.Latest, InTheDeliveryPoint: req.InTheDeliveryPoint})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &desc.ListOrdersResponse{
		Orders: make([]*desc.Order, 0, len(orders)),
	}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, convertPostgresOrderToProtoOrder(o))
	}
	return resp, nil
}
