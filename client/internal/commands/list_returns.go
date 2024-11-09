package commands

import (
	orderWriter "OzonHW1/client/internal/api/order"
	order_service "OzonHW1/pkg/order-service/v1"
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"strings"
)

type ListReturnsCommand struct {
	Name         string
	Description  string
	orderService ListReturnsService
	reader       *bufio.Reader
}
type ListReturnsService interface {
	GetReturned(ctx context.Context, in *order_service.GetReturnedRequest, opts ...grpc.CallOption) (*order_service.GetReturnedResponse, error)
}

func NewListReturnsCommand(orderService ListReturnsService, reader *bufio.Reader) *ListReturnsCommand {
	return &ListReturnsCommand{Name: "list-returns",
		Description:  "Get a list of returns",
		orderService: orderService,
		reader:       reader,
	}
}

func (cmd *ListReturnsCommand) Execute(ctx context.Context, args []string) {
	startIndex := int64(0)
	pageSize := int64(5)
	fmt.Println(`
		q				close the table
		d				page forward
		b				page back
	`)
	for {
		req := &order_service.GetReturnedRequest{
			Offset: startIndex,
			Limit:  pageSize,
		}
		orders, err := cmd.orderService.GetReturned(ctx, req)
		for i := startIndex; i < int64(len(orders.Orders)) && i-startIndex < pageSize; i++ {
			orderWriter.Print(orders.Orders[i])
		}
		fmt.Print("> ")
		line, err := cmd.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		line = strings.TrimSpace(line)
		if len(line) != 1 {
			fmt.Println("enter only one character")
			continue
		}
		if strings.ToLower(line) == "q" {
			return
		} else if strings.ToLower(line) == "d" {
			if int64(len(orders.Orders))/(startIndex+pageSize) != 0 {
				startIndex += pageSize
			}
		} else if strings.ToLower(line) == "b" {
			if startIndex != 0 {
				startIndex -= pageSize
			}
		}
	}
}
