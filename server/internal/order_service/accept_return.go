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

func (s *Implementation) AcceptReturn(ctx context.Context, req *desc.AcceptReturnRequest) (*desc.AcceptReturnResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Order.AcceptReturn")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.storage.AcceptReturn(ctx, req.UserId, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	event := kafka_entity.AcceptReturnEventMessage{
		OrderID: req.OrderId,
		UserID:  req.UserId,
	}

	if err := s.sendToKafka(event, acceptReturn, req.OrderId); err != nil {
		log.Printf("send to kafka failed, err:%v", err)
	}
	return &desc.AcceptReturnResponse{}, nil
}
