package order_service

import (
	"OzonHW1/pkg/kafka_entity"
	desc "OzonHW1/pkg/order-service/v1"
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *Implementation) AddOrder(ctx context.Context, req *desc.AddOrderRequest) (*desc.AddOrderResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Order.AddOrder")
	defer span.Finish()
	span.LogEvent("AddOrder")
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := s.storage.AddOrder(ctx, req.Id, req.UserId, req.ShelfLife, req.Packaging, req.Weigh, req.Cost)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	event := kafka_entity.AddOrderEventMessage{
		OrderID:   req.Id,
		UserID:    req.UserId,
		ShelfLife: req.ShelfLife,
		Packaging: req.Packaging,
		Weigh:     req.Weigh,
		Cost:      req.Cost,
	}

	if err := s.sendToKafka(event, acceptOrder, req.Id); err != nil {
		log.Printf("send to kafka failed, err:%v", err)
	}

	return &desc.AddOrderResponse{}, nil
}
