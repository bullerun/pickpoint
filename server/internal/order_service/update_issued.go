package order_service

import (
	"OzonHW1/pkg/kafka_entity"
	desc "OzonHW1/pkg/order-service/v1"
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

func (s *Implementation) UpdateIssued(ctx context.Context, req *desc.UpdateIssuedRequest) (*desc.UpdateIssuedResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.Order.UpdateIssued")
	defer span.Finish()
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := s.storage.UpdateIssued(ctx, req.OrderIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	event := kafka_entity.UpdateIssuedEventMessage{
		OrderIDs: req.OrderIds,
	}
	for _, orderID := range req.OrderIds {
		id, err := strconv.ParseInt(orderID, 10, 64)
		if err != nil {
			log.Println("Error converting order id to int")
		}
		if err := s.sendToKafka(event, updateIssued, id); err != nil {
			log.Printf("send to kafka failed, err:%v", err)
		}
	}
	return &desc.UpdateIssuedResponse{}, nil
}
