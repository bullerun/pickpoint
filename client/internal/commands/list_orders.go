package commands

import (
	orderWriter "OzonHW1/client/internal/api/order"
	order_service "OzonHW1/pkg/order-service/v1"
	entity "OzonHW1/pkg/order_entity"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"strconv"
)

type ListOrdersCommand struct {
	Name         string
	Description  string
	orderService ListOrdersService
}

type ListOrdersService interface {
	ListOrders(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption) (*order_service.ListOrdersResponse, error)
}

func NewListOrdersCommand(orderService ListOrdersService) *ListOrdersCommand {
	return &ListOrdersCommand{Name: "list-orders",
		Description:  "Get a list of orders",
		orderService: orderService,
	}
}
func (cmd *ListOrdersCommand) Execute(ctx context.Context, args []string) {
	userID, flag, err := cmd.validateArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	request := &order_service.ListOrdersRequest{
		UserId:             userID,
		InTheDeliveryPoint: flag.InTheDeliveryPoint,
		Latest:             flag.Latest,
	}
	orders, err := cmd.orderService.ListOrders(ctx, request)
	if err != nil {
		fmt.Println(err)
		return
	}
	if orders == nil || len(orders.Orders) == 0 {
		fmt.Println("there are no orders that meet the conditions")
		return
	}
	for _, order := range orders.Orders {
		orderWriter.Print(order)
	}
}
func (cmd *ListOrdersCommand) validateArgs(args []string) (int64, *entity.Flag, error) {
	if !(len(args) == 3 || len(args) == 5 || len(args) == 1) {
		return 0, nil, fmt.Errorf("the number of arguments must be either 1 [UserID], 3 [UserID] flag flagAtr, or 5 [UserID] flag flagAtr flag flagAtr")
	}
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil || userID < 1 {
		return 0, nil, fmt.Errorf("userID is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
	}
	flag := &entity.Flag{}
	for i := 1; i <= len(args)-2 && i <= 5; i += 2 {
		if err := cmd.checkFlag(args, i, flag); err != nil {
			return 0, nil, err
		}
	}
	return userID, flag, nil
}
func (cmd *ListOrdersCommand) checkFlag(args []string, start int, flag *entity.Flag) error {
	if "--last" == args[start] {
		n, err := strconv.ParseInt(args[start+1], 10, 64)
		if err != nil || n < 1 {
			return fmt.Errorf("the --last flag must be followed by a positive number")
		}
		flag.Latest = n
	} else if "--in-the" == args[start] {
		if args[start+1] == "true" {
			flag.InTheDeliveryPoint = true
		} else if args[start+1] == "false" {
			flag.InTheDeliveryPoint = false
		} else {
			return fmt.Errorf("the --in-the flag should be followed by \"true\" or \"false\"")
		}
	} else {
		return errors.New("unknown flag")
	}
	return nil
}
