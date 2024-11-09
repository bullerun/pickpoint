package order_service

import (
	desc "OzonHW1/pkg/order-service/v1"
	entity "OzonHW1/pkg/order_entity"
	"database/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertSqlTimeToProtoTime(time *sql.NullTime) *timestamppb.Timestamp {
	if time.Valid {
		return timestamppb.New(time.Time)
	}
	return nil
}
func convertPostgresOrderToProtoOrder(o *entity.Order) *desc.Order {
	return &desc.Order{
		Id:               o.ID,
		UserId:           o.UserID,
		ExpiryDate:       timestamppb.New(o.ExpiryDate),
		AcceptDate:       convertSqlTimeToProtoTime(&o.AcceptDate),
		ReturnFromClient: convertSqlTimeToProtoTime(&o.ReturnFromClient),
		ReturnToCourier:  convertSqlTimeToProtoTime(&o.ReturnToCourier),
		Packaging:        o.Packaging,
		Weigh:            o.Weigh,
		Cost:             o.Cost,
	}

}
