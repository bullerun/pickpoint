package commands

import (
	order_service "OzonHW1/pkg/order-service/v1"
	"OzonHW1/pkg/order_entity/packaging"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"strconv"
	"strings"
	"time"
)

type AcceptOrderCommand struct {
	Name         string
	Description  string
	orderService AcceptOrderService
}

var packagingMap = map[string]packaging.Type{
	"bag":  &packaging.Bag{},
	"box":  &packaging.Box{},
	"film": &packaging.Film{},
}

type AcceptOrderService interface {
	AddOrder(ctx context.Context, in *order_service.AddOrderRequest, opts ...grpc.CallOption) (*order_service.AddOrderResponse, error)
}

func NewAcceptOrderCommand(orderService AcceptOrderService) *AcceptOrderCommand {
	return &AcceptOrderCommand{
		Name:         "accept-order",
		Description:  "Accepts orders on their own.",
		orderService: orderService,
	}
}

type addArgs struct {
	ID         int64
	UserID     int64
	ExpiryDate int64
	Packaging  string
	Weigh      float32
	Cost       float32
}

func (cmd *AcceptOrderCommand) Execute(ctx context.Context, args []string) {
	order, err := cmd.validateArgsAndCreateOrder(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	req := &order_service.AddOrderRequest{
		Id: order.ID, UserId: order.UserID, ShelfLife: order.ExpiryDate, Packaging: order.Packaging, Weigh: order.Weigh, Cost: order.Cost,
	}

	if _, err := cmd.orderService.AddOrder(ctx, req); err != nil {
		fmt.Println("Error adding order:", err)
		return
	}
	fmt.Println("Order added successfully")
}
func (cmd *AcceptOrderCommand) validateArgsAndCreateOrder(args []string) (*addArgs, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf("incorrect number of arguments. Expecting 6. [orderID] [userID] [daysToStore] [packing] [weight] [cost]")
	}
	orderID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil || orderID < 1 {
		return nil, fmt.Errorf("orderID is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
	}
	userID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || userID < 1 {
		return nil, fmt.Errorf("userID is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
	}
	daysToStore, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil || daysToStore < 1 {
		return nil, fmt.Errorf("daysToStore is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
	}
	packagingType, exist := packagingMap[strings.ToLower(args[3])]
	if !exist {
		return nil, fmt.Errorf("unknown packaging type: %s", args[3])
	}
	weight, err := strconv.ParseFloat(args[4], 32)
	if err != nil || weight < 0 {
		return nil, fmt.Errorf("the weight must be a non-negative number")
	}
	cost, err := strconv.ParseFloat(args[5], 32)
	if err != nil || cost < 0 {
		return nil, fmt.Errorf("the cost must be a non-negative number")
	}
	if err := cmd.validate(weight, packagingType); err != nil {
		return nil, err
	}
	return &addArgs{
		ID:         orderID,
		UserID:     userID,
		ExpiryDate: daysToStore,
		Packaging:  packagingType.Name(),
		Weigh:      float32(weight),
		Cost:       float32(cost) + packagingType.Cost(),
	}, nil
}
func (cmd *AcceptOrderCommand) validate(weight float64, p packaging.Type) error {
	if weight >= p.MaxWeight() {
		return fmt.Errorf("weight must be less than %.2f", p.MaxWeight())
	}
	return nil
}

var now = func() time.Time {
	return time.Now()
}
