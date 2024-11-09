package orderWriter

import (
	order_service "OzonHW1/pkg/order-service/v1"
	"fmt"
)

func Print(order *order_service.Order) {
	fmt.Printf("Order Details:\n")
	fmt.Printf("ID: %d\n", order.Id)
	fmt.Printf("User ID: %d\n", order.UserId)
	fmt.Printf("Order Create Date: %s\n", order.OrderCreateDate)
	fmt.Printf("Expiry Date: %s\n", order.ExpiryDate)
	fmt.Printf("Accept Date: %s\n", order.AcceptDate)
	fmt.Printf("Is Issued: %s\n", order.AcceptDate)
	fmt.Printf("Is Returned from User: %s\n", order.ReturnFromClient)
	fmt.Printf("Packaging: %s\n", order.Packaging)
	fmt.Printf("Weigh: %.2f \n", order.Weigh)
	fmt.Printf("Cost: %.2f \n", order.Cost)
	fmt.Println("________________________________________________________________")
}
