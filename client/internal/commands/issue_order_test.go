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

func TestIssueOrderCommand_validateAndPrepareArgs(t *testing.T) {

	cmd := &IssueOrderCommand{}

	type args struct {
		args []string
	}
	tests := []struct {
		name        string
		args        args
		orderId     []int64
		wantErr     assert.ErrorAssertionFunc
		wantErrDisc string
	}{
		{
			name:    "Success issue order 1",
			args:    args{args: []string{"1", "2"}},
			orderId: []int64{1, 2},
			wantErr: assert.NoError,
		},
		{
			name:    "Success issue order 2",
			args:    args{args: []string{"56", "323"}},
			orderId: []int64{56, 323},
			wantErr: assert.NoError,
		},
		{
			name:        "Failed issue order",
			args:        args{args: []string{"0", "1"}},
			wantErr:     assert.Error,
			wantErrDisc: "orderID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:    "Failed issue order",
			args:    args{args: []string{"5", "0"}},
			wantErr: assert.Error, wantErrDisc: "orderID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:        "Failed issue order",
			args:        args{args: []string{"-1", "-1436352"}},
			wantErr:     assert.Error,
			wantErrDisc: "orderID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:        "Failed issue order",
			args:        args{args: nil},
			wantErr:     assert.Error,
			wantErrDisc: "command must have at least one argument <orderID>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := cmd.validateAndPrepareArgs(tt.args.args)
			if !tt.wantErr(t, err) {
				return
			}
			if err != nil {
				assert.ErrorContains(t, err, tt.wantErrDisc)
			}
		})
	}
}

func TestIssueOrderCommand_Execute(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)
	issueOrderServiceMock := mock.NewIssueOrderServiceMock(mc)
	cmd := NewIssueOrderCommand(issueOrderServiceMock)
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
			args: args{args: []string{"1", "2", "6", "9"}},
			setup: func() {
				issueOrderServiceMock.UpdateIssuedMock.Expect(ctx, &order_service.UpdateIssuedRequest{OrderIds: []string{"1", "2", "6", "9"}}).Return(nil, nil)
			},
			want: "orders have been successfully issued",
		},
		{
			name: "Failed issue order",
			args: args{args: []string{"1", "2", "6", "9"}},
			setup: func() {
				issueOrderServiceMock.UpdateIssuedMock.Expect(ctx, &order_service.UpdateIssuedRequest{OrderIds: []string{"1", "2", "6", "9"}}).Return(nil, errors.New("service error"))
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
