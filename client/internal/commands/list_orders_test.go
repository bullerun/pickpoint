package commands

import (
	"OzonHW1/client/internal/commands/mock"
	order_service "OzonHW1/pkg/order-service/v1"
	entity "OzonHW1/pkg/order_entity"
	"bytes"
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"os"
	"testing"
)

func TestListOrdersCommand_validateArgs(t *testing.T) {
	cmd := &ListOrdersCommand{}
	type args struct {
		args []string
	}
	type answer struct {
		userId int64
		flag   *entity.Flag
	}
	tests := []struct {
		name        string
		args        args
		wantAnswer  answer
		wantErr     assert.ErrorAssertionFunc
		wantErrDisc string
	}{
		{
			name:       "Success list order",
			args:       args{[]string{"1"}},
			wantAnswer: answer{userId: 1, flag: &entity.Flag{}},
			wantErr:    assert.NoError,
		},
		{
			name:       "Success list order",
			args:       args{[]string{"1", "--last", "6"}},
			wantAnswer: answer{userId: 1, flag: &entity.Flag{Latest: 6}},
			wantErr:    assert.NoError,
		},
		{
			name:       "Success list order",
			args:       args{[]string{"1", "--in-the", "true"}},
			wantAnswer: answer{userId: 1, flag: &entity.Flag{InTheDeliveryPoint: true}},
			wantErr:    assert.NoError,
		},
		{
			name:       "Success list order",
			args:       args{[]string{"1", "--in-the", "false"}},
			wantAnswer: answer{userId: 1, flag: &entity.Flag{InTheDeliveryPoint: false}},
			wantErr:    assert.NoError,
		},
		{
			name:       "Success list order",
			args:       args{[]string{"1", "--last", "3", "--in-the", "true"}},
			wantAnswer: answer{userId: 1, flag: &entity.Flag{Latest: 3, InTheDeliveryPoint: true}},
			wantErr:    assert.NoError,
		},
		{
			name:        "Failed list order",
			args:        args{[]string{"1", "--last"}},
			wantErr:     assert.Error,
			wantErrDisc: "the number of arguments must be either 1 [UserID], 3 [UserID] flag flagAtr, or 5 [UserID] flag flagAtr flag flagAtr",
		},
		{
			name:        "Failed list order",
			args:        args{[]string{"1", "--last", "3", "true"}},
			wantErr:     assert.Error,
			wantErrDisc: "the number of arguments must be either 1 [UserID], 3 [UserID] flag flagAtr, or 5 [UserID] flag flagAtr flag flagAtr",
		},
		{
			name:        "Failed list order",
			args:        args{[]string{"0"}},
			wantErr:     assert.Error,
			wantErrDisc: "userID is entered incorrectly, it must be a number greater than 0 and less than",
		},
		{
			name:        "Failed list order",
			args:        args{[]string{"1", "--last", "0"}},
			wantErr:     assert.Error,
			wantErrDisc: "the --last flag must be followed by a positive number",
		},
		{
			name:        "Failed list order",
			args:        args{[]string{"1", "--in-the", "1"}},
			wantErr:     assert.Error,
			wantErrDisc: "the --in-the flag should be followed by \"true\" or \"false\"",
		},
		{
			name:        "Failed list order",
			args:        args{[]string{"1", "unknown flag", "1"}},
			wantErr:     assert.Error,
			wantErrDisc: "unknown flag",
		},
		{
			name:        "Failed list order",
			args:        args{[]string{"1", "--in-the", "true", "unknown flag", "6"}},
			wantErr:     assert.Error,
			wantErrDisc: "unknown flag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userId, flag, err := cmd.validateArgs(tt.args.args)
			if !tt.wantErr(t, err) {
				return
			}
			if err == nil {
				assert.Equal(t, tt.wantAnswer.userId, userId)
				assert.Equal(t, tt.wantAnswer.flag, flag)
			} else {
				assert.ErrorContains(t, err, tt.wantErrDisc)
			}
		})
	}
}

func TestListOrdersCommand_Execute(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)
	issueOrderServiceMock := mock.NewListOrderServiceMock(mc)
	cmd := NewListOrdersCommand(issueOrderServiceMock)
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
			name: "Success list order",
			args: args{args: []string{"1", "--last", "3", "--in-the", "true"}},
			setup: func() {
				issueOrderServiceMock.ListOrdersMock.Expect(ctx, &order_service.ListOrdersRequest{UserId: 1, InTheDeliveryPoint: true, Latest: 3}).Return(&order_service.ListOrdersResponse{Orders: []*order_service.Order{{
					Id:              1,
					UserId:          2,
					OrderCreateDate: timestamppb.New(now()),
					ExpiryDate:      timestamppb.New(now().AddDate(0, 0, 5)),
					Packaging:       "box",
					Weigh:           2.5,
					Cost:            float32(120)}}}, nil)
			},
			want: "Order Details:\nID: 1\nUser ID: 2",
		},
		{
			name: "Failed list order",
			args: args{args: []string{"1", "--last", "3", "--in-the", "true"}},
			setup: func() {
				issueOrderServiceMock.ListOrdersMock.Expect(ctx, &order_service.ListOrdersRequest{UserId: 1, InTheDeliveryPoint: true, Latest: 3}).Return(nil, errors.New("service error"))
			},
			want: "service error",
		},
		{
			name: "Failed list order",
			args: args{args: []string{"1", "--last", "3", "--in-the", "true"}},
			setup: func() {
				issueOrderServiceMock.ListOrdersMock.Expect(ctx, &order_service.ListOrdersRequest{UserId: 1, InTheDeliveryPoint: true, Latest: 3}).Return(nil, nil)
			},
			want: "there are no orders that meet the conditions",
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
