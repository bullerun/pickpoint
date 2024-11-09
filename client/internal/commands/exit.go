package commands

import (
	"context"
	"fmt"
)

type ExitCommand struct {
	Name        string
	Description string
	quit        chan bool
}

func NewExitCommand(quit chan bool) *ExitCommand {
	return &ExitCommand{Name: "exit",
		Description: "Exit command",
		quit:        quit,
	}
}
func (cmd *ExitCommand) Execute(ctx context.Context, args []string) {
	fmt.Println("Exit...")
	cmd.quit <- true
}
