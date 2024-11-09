package managers

import (
	"OzonHW1/client/internal/commands"
	order_service "OzonHW1/pkg/order-service/v1"
	"bufio"
	"context"
	"errors"
	"time"
)

var CommandNotFoundError = errors.New("command not found")

type Command interface {
	Execute(context.Context, []string)
}
type CommandManager struct {
	commandManager map[string]Command
}

func NewCommandManager(facade order_service.OrderServiceClient, reader *bufio.Reader, quit chan bool) *CommandManager {
	commandManager := make(map[string]Command)
	commandManager["accept-order"] = commands.NewAcceptOrderCommand(facade)
	commandManager["return-order"] = commands.NewReturnOrderCommand(facade)
	commandManager["issue-order"] = commands.NewIssueOrderCommand(facade)
	commandManager["list-orders"] = commands.NewListOrdersCommand(facade)
	commandManager["accept-return"] = commands.NewAcceptReturnCommand(facade)
	commandManager["list-returns"] = commands.NewListReturnsCommand(facade, reader)
	commandManager["exit"] = commands.NewExitCommand(quit)
	commandManager["help"] = commands.NewHelpCommand()
	return &CommandManager{commandManager: commandManager}
}

func (c *CommandManager) ExecuteCommand(task Task) error {
	ctx, cancel := context.WithTimeout(task.Ctx, 10*time.Second)
	defer cancel()

	if _, exists := c.commandManager[task.Command]; !exists {
		return CommandNotFoundError
	}
	c.commandManager[task.Command].Execute(ctx, task.Args)
	return nil
}
