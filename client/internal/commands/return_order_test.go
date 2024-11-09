package commands

import (
	"OzonHW1/client/internal/commands/mock"
	order_service "OzonHW1/pkg/order-service/v1"
	"bytes"
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestReturnOrderCommand_validateArgs(t *testing.T) {
	cmd := &ReturnOrderCommand{}
	type args struct {
		args []string
	}
	type answer struct {
		orderID int64
	}
	tests := []struct {
		name        string
		args        args
		wantAnswer  answer
		wantErr     assert.ErrorAssertionFunc
		wantErrDisc string
	}{
		{
			name:       "Success return order",
			args:       args{[]string{"1"}},
			wantAnswer: answer{1},
			wantErr:    assert.NoError,
		},
		{
			name:        "Failed return order",
			args:        args{[]string{"1", "2"}},
			wantErr:     assert.Error,
			wantErrDisc: "command accepts only one orderID argument",
		},
		{
			name:        "Failed return order",
			args:        args{[]string{}},
			wantErr:     assert.Error,
			wantErrDisc: "command accepts only one orderID argument",
		},
		{
			name:        "Failed return order",
			args:        args{[]string{"0"}},
			wantErr:     assert.Error,
			wantErrDisc: "orderID is entered incorrectly, it must be a number greater than 0 and less than",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderID, err := cmd.validateArgs(tt.args.args)
			if !tt.wantErr(t, err) {
				return
			}
			if err == nil {
				assert.Equal(t, tt.wantAnswer.orderID, orderID)

			} else {
				assert.ErrorContains(t, err, tt.wantErrDisc)
			}
		})
	}
}

func TestReturnOrderCommand_Execute(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)
	issueOrderServiceMock := mock.NewReturnOrderServiceMock(mc)
	cmd := NewReturnOrderCommand(issueOrderServiceMock)
	type args struct {
		args []string
	}
	tests := []struct {
		name  string
		args  args
		setup func()
		want  string
	}{
		{
			name: "Success return order",
			args: args{args: []string{"1"}},
			setup: func() {
				issueOrderServiceMock.ReturnOrderToCourierMock.Expect(ctx, &order_service.ReturnOrderToCourierRequest{OrderId: 1}).Return(nil, nil)
			},
			want: "The order was returned to the courier",
		},
		{
			name: "Failed return order",
			args: args{args: []string{"1"}},
			setup: func() {
				issueOrderServiceMock.ReturnOrderToCourierMock.Expect(ctx, &order_service.ReturnOrderToCourierRequest{OrderId: 1}).Return(nil, errors.New("service error"))
			},
			want: "service error",
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
			tt.setup()
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
