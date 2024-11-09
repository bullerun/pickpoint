package commands

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestHelpCommand_Execute(t *testing.T) {
	ctx := context.Background()
	cmd := NewHelpCommand()
	type args struct {
		args []string
	}
	tests := []struct {
		name string

		args args
		want string
	}{
		{
			name: "Success help",
			args: args{},
			want: `accept-order <orderID> <userID> <daysToStore> <packing>	<weight> <cost>		Accept order from courier
 	return-order <orderID>											Return the order to the courier
	issue-order <orderId> ...										Issue the order to the user
	list-orders <userID> --last=<number> --in-the=<bool>		Get a list of orders
	accept-return <userID> <orderID>							Accept the return from the user
	list-returns 													Get a list of returns
	exit 															Exiting the application
	help`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("can't complete the test")
			}
			oldStdOut := os.Stdout
			defer func() {
				os.Stdout = oldStdOut
			}()
			os.Stdout = w
			cmd.Execute(ctx, tt.args.args)
			w.Close()
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Errorf("%v", err)
			}
			output := buf.String()
			assert.Contains(t, output, tt.want)
		})
	}
}
