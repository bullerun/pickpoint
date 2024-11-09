package order_service

import (
	desc "OzonHW1/pkg/order-service/v1"
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetReturned(ctx context.Context, req *desc.GetReturnedRequest) (*desc.GetReturnedResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Order.GetReturned")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.Limit < 5 {
		req.Limit = 5
	}
	orders := s.storage.GetReturned(ctx, req.Offset, req.Limit)
	resp := &desc.GetReturnedResponse{
		Orders: make([]*desc.Order, 0, len(orders)),
	}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, convertPostgresOrderToProtoOrder(o))
	}
	return resp, nil
}
