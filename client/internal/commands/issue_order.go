package commands

import (
	order_service "OzonHW1/pkg/order-service/v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math"
	"strconv"
)

type IssueOrderCommand struct {
	Name         string
	Description  string
	orderService IssueOrderService
}
type IssueOrderService interface {
	UpdateIssued(ctx context.Context, in *order_service.UpdateIssuedRequest, opts ...grpc.CallOption) (*order_service.UpdateIssuedResponse, error)
}

func NewIssueOrderCommand(orderService IssueOrderService) *IssueOrderCommand {
	return &IssueOrderCommand{Name: "issue-order",
		Description:  "Issue the order to the user",
		orderService: orderService,
	}
}

func (cmd *IssueOrderCommand) GetName() string {
	return cmd.Name
}

func (cmd *IssueOrderCommand) GetDescription() string {
	return cmd.Description
}

func (cmd *IssueOrderCommand) Execute(ctx context.Context, args []string) {
	err := cmd.validateAndPrepareArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	request := &order_service.UpdateIssuedRequest{
		OrderIds: args,
	}
	if _, err := cmd.orderService.UpdateIssued(ctx, request); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("orders have been successfully issued")

}
func (cmd *IssueOrderCommand) validateAndPrepareArgs(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("command must have at least one argument <orderID>")
	}
	for i := range args {
		orderID, err := strconv.Atoi(args[i])
		if err != nil || orderID < 1 {
			return fmt.Errorf("orderID is entered incorrectly, it must be a number greater than 0 and less than %d", math.MaxInt)
		}
	}

	return nil
}
