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

func TestAcceptReturnCommand_validateArgs(t *testing.T) {
	cmd := &AcceptReturnCommand{}
	type args struct {
		args []string
	}
	tests := []struct {
		name        string
		args        args
		userId      int64
		orderId     int64
		wantErr     assert.ErrorAssertionFunc
		wantErrDisc string
	}{
		{
			name:    "Success accept return 1",
			args:    args{args: []string{"1", "2"}},
			userId:  1,
			orderId: 2,
			wantErr: assert.NoError,
		},
		{
			name:    "Success accept return 2",
			args:    args{args: []string{"1", "1"}},
			userId:  1,
			orderId: 1,
			wantErr: assert.NoError,
		},
		{
			name:        "Failed accept return",
			args:        args{args: []string{"0", "1"}},
			wantErr:     assert.Error,
			wantErrDisc: "userID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:        "Failed accept return",
			args:        args{args: []string{"10", "0"}},
			wantErr:     assert.Error,
			wantErrDisc: "orderID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:        "Failed accept return",
			args:        args{args: []string{"100000000000000000000000000000000000000000", "1"}},
			wantErr:     assert.Error,
			wantErrDisc: "userID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:        "Failed accept return",
			args:        args{args: []string{"1", "100000000000000000000000000000000000000000"}},
			wantErr:     assert.Error,
			wantErrDisc: "orderID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:    "Failed accept return",
			args:    args{args: []string{"-100000000000000000000000000000000000000000", "-1"}},
			wantErr: assert.Error,
		},
		{
			name:    "Failed accept return",
			args:    args{args: []string{"-1", "-100000000000000000000000000000000000000000"}},
			wantErr: assert.Error,
		},
		{
			name:        "Failed accept return",
			args:        args{args: []string{"1", "2", "2"}},
			wantErr:     assert.Error,
			wantErrDisc: "incorrect number of arguments. Expecting 2. [userID] [orderID]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotOrderID, gotUserID, err := cmd.validateArgs(tt.args.args)
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.userId, gotOrderID)
				assert.Equal(t, tt.orderId, gotUserID)
			} else if tt.wantErrDisc != "" {
				assert.ErrorContains(t, err, tt.wantErrDisc)
			}
		})
	}
}
func TestAcceptReturnCommand_Execute(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)
	issueOrderServiceMock := mock.NewAcceptReturnServiceMock(mc)
	cmd := NewAcceptReturnCommand(issueOrderServiceMock)
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
			name: "Success issue order",
			args: args{args: []string{"1", "2"}},
			setup: func() {
				issueOrderServiceMock.AcceptReturnMock.Expect(ctx, &order_service.AcceptReturnRequest{OrderId: 2, UserId: 1}).Return(nil, nil)
			},
			want: "Accept return",
		},
		{
			name: "Failed issue order",
			args: args{args: []string{"1", "2"}},
			setup: func() {
				issueOrderServiceMock.AcceptReturnMock.Expect(ctx, &order_service.AcceptReturnRequest{OrderId: 2, UserId: 1}).Return(nil, errors.New("service error"))
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
