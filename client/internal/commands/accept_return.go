package commands

import (
	order_service "OzonHW1/pkg/order-service/v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"strconv"
)

type AcceptReturnCommand struct {
	Name         string
	Description  string
	orderService AcceptReturnService
}
type AcceptReturnService interface {
	AcceptReturn(ctx context.Context, in *order_service.AcceptReturnRequest, opts ...grpc.CallOption) (*order_service.AcceptReturnResponse, error)
}

func NewAcceptReturnCommand(orderService AcceptReturnService) *AcceptReturnCommand {
	return &AcceptReturnCommand{Name: "accept-return",
		Description:  "Accept a return from a client",
		orderService: orderService,
	}
}
func (cmd *AcceptReturnCommand) Execute(ctx context.Context, args []string) {
	userID, orderID, err := cmd.validateArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	request := &order_service.AcceptReturnRequest{
		UserId:  userID,
		OrderId: orderID,
	}
	if _, err := cmd.orderService.AcceptReturn(ctx, request); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Accept return")
}
func (cmd *AcceptReturnCommand) validateArgs(args []string) (int64, int64, error) {
	if len(args) != 2 {
		return 0, 0, fmt.Errorf("incorrect number of arguments. Expecting 2. [userID] [orderID]")
	}
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil || userID < 1 {
		return 0, 0, fmt.Errorf("userID is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
	}
	orderID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || orderID < 1 {
		return 0, 0, fmt.Errorf("orderID is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
	}
	return userID, orderID, nil
}
