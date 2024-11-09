package commands

import (
	order_service "OzonHW1/pkg/order-service/v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"strconv"
)

type ReturnOrderCommand struct {
	Name         string
	Description  string
	orderService ReturnOrderService
}
type ReturnOrderService interface {
	ReturnOrderToCourier(ctx context.Context, in *order_service.ReturnOrderToCourierRequest, opts ...grpc.CallOption) (*order_service.ReturnOrderToCourierResponse, error)
}

func NewReturnOrderCommand(orderService ReturnOrderService) *ReturnOrderCommand {
	return &ReturnOrderCommand{Name: "return-order",
		Description:  "Return an order to the courier",
		orderService: orderService,
	}
}

func (cmd *ReturnOrderCommand) Execute(ctx context.Context, args []string) {
	orderID, err := cmd.validateArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	req := &order_service.ReturnOrderToCourierRequest{
		OrderId: orderID,
	}
	if _, err := cmd.orderService.ReturnOrderToCourier(ctx, req); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("The order was returned to the courier")
}
func (cmd *ReturnOrderCommand) validateArgs(args []string) (int64, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("command accepts only one orderID argument")
	}
	orderID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil || orderID < 1 {
		return 0, fmt.Errorf("orderID is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
	}
	return orderID, nil
}
