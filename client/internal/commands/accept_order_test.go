package commands

import (
	"OzonHW1/client/internal/commands/mock"
	order_service "OzonHW1/pkg/order-service/v1"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
	"time"

	"math"
	"strconv"

	"github.com/stretchr/testify/assert"
)

const (
	orderIDIndex     = 0
	userIDIndex      = 1
	daysToStoreIndex = 2
	packingIndex     = 3
	weightIndex      = 4
	costIndex        = 5
)

var validArgs = []string{"1", "2", "5", "Box", "2.5", "100.0"}
var cmd = &AcceptOrderCommand{}

func copyArgs() []string {
	argsCopy := make([]string, len(validArgs))
	copy(argsCopy, validArgs)
	return argsCopy
}
func TestValidateArgsAndCreateOrder_Success(t *testing.T) {
	t.Parallel()
	args := copyArgs()
	order, err := cmd.validateArgsAndCreateOrder(args)
	require.NoError(t, err)
	require.NotNil(t, order)
	assert.Equal(t, int64(1), order.ID)
	assert.Equal(t, int64(2), order.UserID)
	assert.Equal(t, "box", order.Packaging)
	assert.Equal(t, float32(2.5), order.Weigh)
	assert.Equal(t, float32(120.0), order.Cost) // 100 (cost) + 20 (box cost)
}

func TestValidateArgsAndCreateOrder_InvalidArgsCount(t *testing.T) {
	t.Parallel()
	args := []string{"1", "2", "10", "box", "2.5"}

	_, err := cmd.validateArgsAndCreateOrder(args)
	assert.ErrorContains(t, err, "incorrect number of arguments. Expecting 6. [orderID] [userID] [daysToStore] [packing] [weight] [cost]")
}

func TestValidateArgsAndCreateOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()
	maxValid := strconv.Itoa(math.MaxInt)
	// MaxInt32 и MaxInt64 заканчиваются на 7
	badArgs := []string{"-2", "-1", "0", maxValid[:len(maxValid)-1] + "8"}
	args := copyArgs()
	for _, el := range badArgs {
		args[orderIDIndex] = el
		_, err := cmd.validateArgsAndCreateOrder(args)
		assert.ErrorContains(t, err, "orderID is entered incorrectly, it must be a number greater than 0 and less than "+strconv.Itoa(math.MaxInt))
	}
}

func TestValidateArgsAndCreateOrder_InvalidUserID(t *testing.T) {
	t.Parallel()
	maxValid := strconv.Itoa(math.MaxInt)
	// MaxInt32 и MaxInt64 заканчиваются на 7
	badArgs := []string{"-2", "-1", "0", maxValid[:len(maxValid)-1] + "8"}
	args := copyArgs()
	for _, el := range badArgs {
		args[userIDIndex] = el
		_, err := cmd.validateArgsAndCreateOrder(args)
		assert.ErrorContains(t, err, "userID is entered incorrectly, it must be a number greater than 0 and less than "+strconv.Itoa(math.MaxInt))
	}
}

func TestValidateArgsAndCreateOrder_InvalidDaysToStore(t *testing.T) {
	t.Parallel()
	maxValid := strconv.Itoa(math.MaxInt)
	// MaxInt32 и MaxInt64 заканчиваются на 7
	badArgs := []string{"-2", "-1", "0", maxValid[:len(maxValid)-1] + "8"}
	args := copyArgs()
	for _, el := range badArgs {
		args[daysToStoreIndex] = el
		_, err := cmd.validateArgsAndCreateOrder(args)
		assert.ErrorContains(t, err, "daysToStore is entered incorrectly, it must be a number greater than 0 and less than "+strconv.Itoa(math.MaxInt))
	}
}

func TestValidateArgsAndCreateOrder_UnknownPackagingType(t *testing.T) {
	t.Parallel()
	args := copyArgs()
	args[packingIndex] = "unknown"
	_, err := cmd.validateArgsAndCreateOrder(args)
	assert.ErrorContains(t, err, "unknown packaging type:")
}

func TestValidateArgsAndCreateOrder_InvalidWeight(t *testing.T) {
	t.Parallel()

	args := copyArgs()
	args[weightIndex] = "-2.5"
	_, err := cmd.validateArgsAndCreateOrder(args)
	assert.ErrorContains(t, err, "the weight must be a non-negative number")
}

func TestValidateArgsAndCreateOrder_InvalidCost(t *testing.T) {
	t.Parallel()
	args := copyArgs()
	args[costIndex] = "-100"
	_, err := cmd.validateArgsAndCreateOrder(args)
	assert.ErrorContains(t, err, "the cost must be a non-negative number")
}

// TestExecute_ServesError тут проверяется вывод в консоль, это анрильно запараллелить (используется os.Pipe() и чтобы прочитать, что вывелось в STDOUT
// нужно закрыть файл, а io.Copy() копирует данные до EOF так что плаки плаки)
// вся валидация идет не тут, а в ValidateArgsAndCreateOrder
func TestExecute_Success(t *testing.T) {
	ctx := context.Background()
	// Перехватываем вывод в консоль
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	os.Stdout = w

	// Создаем minimock контроллер
	mc := minimock.NewController(t)

	// Создаем мок для orderStorage
	serviceMock := mock.NewAcceptOrderServiceMock(mc)

	// Настраиваем мок: успешное добавление заказа
	serviceMock.AddOrderMock.Return(nil, nil)

	// Создаем команду
	cmd := NewAcceptOrderCommand(serviceMock)

	// Вызов метода Execute с корректными аргументами
	args := []string{"1", "2", "10", "box", "2.5", "100.0"}
	cmd.Execute(ctx, args)
	w.Close()

	// Проверяем, что вывод содержит сообщение об успешном добавлении заказа
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return
	}

	output := buf.String()
	assert.Contains(t, output, "Order added successfully")
}

// TestExecute_ServesError тут проверяется вывод в консоль, это анрильно запараллелить (используется os.Pipe() и чтобы прочитать, что вывелось в STDOUT
// нужно закрыть файл, а io.Copy() копирует данные до EOF так что плаки плаки)
// вся валидация идет не тут, а в ValidateArgsAndCreateOrder
func TestExecute_ServesError(t *testing.T) {
	ctx := context.Background()
	// Перехватываем вывод в консоль
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	os.Stdout = w

	mc := minimock.NewController(t)

	serviceMock := mock.NewAcceptOrderServiceMock(mc)

	serviceMock.AddOrderMock.Return(nil, fmt.Errorf("service error"))

	cmd := NewAcceptOrderCommand(serviceMock)

	args := []string{"1", "2", "10", "box", "2.5", "100.0"}
	cmd.Execute(ctx, args)
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	assert.Contains(t, output, "Error adding order: service error")
}

func TestExecuteAndValidate(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)
	acceptOrderServiceMock := mock.NewAcceptOrderServiceMock(mc)
	cmd := NewAcceptOrderCommand(acceptOrderServiceMock)
	fixedTimeNow := time.Now()
	now = func() time.Time {
		return fixedTimeNow
	}
	defer func() { now = time.Now }()
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
			name: "Success accept order",
			args: args{args: []string{"1", "2", "5", "Box", "2.5", "100.0"}},
			setup: func() {
				acceptOrderServiceMock.AddOrderMock.Expect(ctx, &order_service.AddOrderRequest{Id: 1, UserId: 2, ShelfLife: 5, Packaging: "box", Weigh: 2.5, Cost: 120.0}).Return(nil, nil)
			},
			want: "Order added successfully",
		},

		{
			name: "failure accept order",
			args: args{args: []string{"1", "2", "5", "Box", "2.5", "100.0"}},
			setup: func() {
				acceptOrderServiceMock.AddOrderMock.Expect(ctx, &order_service.AddOrderRequest{Id: 1, UserId: 2, ShelfLife: 5, Packaging: "box", Weigh: 2.5, Cost: 120.0}).Return(nil, errors.New("service error"))
			},
			want: "Error adding order: service error",
		},
		{
			name: "failure accept order",
			args: args{args: []string{"0", "2", "5", "Box", "2.5", "100.0"}},
			setup: func() {

			},
			want: "orderID is entered incorrectly, it must be a number greater than 0 and less than",
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
