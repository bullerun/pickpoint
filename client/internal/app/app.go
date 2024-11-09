package app

import (
	"OzonHW1/client/internal/managers"
	"bufio"
	"context"
	"fmt"
	"strings"
)

type App struct {
	reader *bufio.Reader
	cm     *managers.CommandManager
}

func MustLoad(reader *bufio.Reader, manager *managers.CommandManager) *App {
	return &App{reader, manager}
}

func (c *App) Run() {
	ctx := context.Background()
	for {
		fmt.Print("> ")
		line, _ := c.reader.ReadString('\n')
		line = strings.TrimSpace(line)
		args := strings.Split(line, " ")
		if len(args) < 1 {
			fmt.Println("Invalid command")
			continue
		}
		if err := c.cm.ExecuteCommand(managers.Task{
			Ctx:     ctx,
			Command: args[0],
			Args:    args[1:],
		}); err != nil {
			fmt.Println(err)
		}
	}
}
