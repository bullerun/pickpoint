package commands

import (
	"context"
	"fmt"
)

type HelpCommand struct {
	Name        string
	Description string
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{Name: "help", Description: "Help Command"}
}
func (cmd *HelpCommand) Execute(ctx context.Context, args []string) {
	fmt.Println(`
	accept-order <orderID> <userID> <daysToStore> <packing>	<weight> <cost>		Accept order from courier
 	return-order <orderID>											Return the order to the courier
	issue-order <orderId> ...										Issue the order to the user
	list-orders <userID> --last=<number> --in-the=<bool>		Get a list of orders
	accept-return <userID> <orderID>							Accept the return from the user
	list-returns 													Get a list of returns
	exit 															Exiting the application
	help															Show help message 						
    `)
}
