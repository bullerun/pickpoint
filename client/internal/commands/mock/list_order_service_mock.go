// Code generated by http://github.com/gojuno/minimock (v3.4.0). DO NOT EDIT.

package mock

//go:generate minimock -i OzonHW1/client/internal/commands.ListOrdersService -o list_order_service_mock.go -n ListOrderServiceMock -p mock

import (
	order_service "OzonHW1/pkg/order-service/v1"
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"google.golang.org/grpc"
)

// ListOrderServiceMock implements mm_commands.ListOrdersService
type ListOrderServiceMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcListOrders          func(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption) (lp1 *order_service.ListOrdersResponse, err error)
	funcListOrdersOrigin    string
	inspectFuncListOrders   func(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption)
	afterListOrdersCounter  uint64
	beforeListOrdersCounter uint64
	ListOrdersMock          mListOrderServiceMockListOrders
}

// NewListOrderServiceMock returns a mock for mm_commands.ListOrdersService
func NewListOrderServiceMock(t minimock.Tester) *ListOrderServiceMock {
	m := &ListOrderServiceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ListOrdersMock = mListOrderServiceMockListOrders{mock: m}
	m.ListOrdersMock.callArgs = []*ListOrderServiceMockListOrdersParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mListOrderServiceMockListOrders struct {
	optional           bool
	mock               *ListOrderServiceMock
	defaultExpectation *ListOrderServiceMockListOrdersExpectation
	expectations       []*ListOrderServiceMockListOrdersExpectation

	callArgs []*ListOrderServiceMockListOrdersParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// ListOrderServiceMockListOrdersExpectation specifies expectation struct of the ListOrdersService.ListOrders
type ListOrderServiceMockListOrdersExpectation struct {
	mock               *ListOrderServiceMock
	params             *ListOrderServiceMockListOrdersParams
	paramPtrs          *ListOrderServiceMockListOrdersParamPtrs
	expectationOrigins ListOrderServiceMockListOrdersExpectationOrigins
	results            *ListOrderServiceMockListOrdersResults
	returnOrigin       string
	Counter            uint64
}

// ListOrderServiceMockListOrdersParams contains parameters of the ListOrdersService.ListOrders
type ListOrderServiceMockListOrdersParams struct {
	ctx  context.Context
	in   *order_service.ListOrdersRequest
	opts []grpc.CallOption
}

// ListOrderServiceMockListOrdersParamPtrs contains pointers to parameters of the ListOrdersService.ListOrders
type ListOrderServiceMockListOrdersParamPtrs struct {
	ctx  *context.Context
	in   **order_service.ListOrdersRequest
	opts *[]grpc.CallOption
}

// ListOrderServiceMockListOrdersResults contains results of the ListOrdersService.ListOrders
type ListOrderServiceMockListOrdersResults struct {
	lp1 *order_service.ListOrdersResponse
	err error
}

// ListOrderServiceMockListOrdersOrigins contains origins of expectations of the ListOrdersService.ListOrders
type ListOrderServiceMockListOrdersExpectationOrigins struct {
	origin     string
	originCtx  string
	originIn   string
	originOpts string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmListOrders *mListOrderServiceMockListOrders) Optional() *mListOrderServiceMockListOrders {
	mmListOrders.optional = true
	return mmListOrders
}

// Expect sets up expected params for ListOrdersService.ListOrders
func (mmListOrders *mListOrderServiceMockListOrders) Expect(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption) *mListOrderServiceMockListOrders {
	if mmListOrders.mock.funcListOrders != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Set")
	}

	if mmListOrders.defaultExpectation == nil {
		mmListOrders.defaultExpectation = &ListOrderServiceMockListOrdersExpectation{}
	}

	if mmListOrders.defaultExpectation.paramPtrs != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by ExpectParams functions")
	}

	mmListOrders.defaultExpectation.params = &ListOrderServiceMockListOrdersParams{ctx, in, opts}
	mmListOrders.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmListOrders.expectations {
		if minimock.Equal(e.params, mmListOrders.defaultExpectation.params) {
			mmListOrders.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmListOrders.defaultExpectation.params)
		}
	}

	return mmListOrders
}

// ExpectCtxParam1 sets up expected param ctx for ListOrdersService.ListOrders
func (mmListOrders *mListOrderServiceMockListOrders) ExpectCtxParam1(ctx context.Context) *mListOrderServiceMockListOrders {
	if mmListOrders.mock.funcListOrders != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Set")
	}

	if mmListOrders.defaultExpectation == nil {
		mmListOrders.defaultExpectation = &ListOrderServiceMockListOrdersExpectation{}
	}

	if mmListOrders.defaultExpectation.params != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Expect")
	}

	if mmListOrders.defaultExpectation.paramPtrs == nil {
		mmListOrders.defaultExpectation.paramPtrs = &ListOrderServiceMockListOrdersParamPtrs{}
	}
	mmListOrders.defaultExpectation.paramPtrs.ctx = &ctx
	mmListOrders.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmListOrders
}

// ExpectInParam2 sets up expected param in for ListOrdersService.ListOrders
func (mmListOrders *mListOrderServiceMockListOrders) ExpectInParam2(in *order_service.ListOrdersRequest) *mListOrderServiceMockListOrders {
	if mmListOrders.mock.funcListOrders != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Set")
	}

	if mmListOrders.defaultExpectation == nil {
		mmListOrders.defaultExpectation = &ListOrderServiceMockListOrdersExpectation{}
	}

	if mmListOrders.defaultExpectation.params != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Expect")
	}

	if mmListOrders.defaultExpectation.paramPtrs == nil {
		mmListOrders.defaultExpectation.paramPtrs = &ListOrderServiceMockListOrdersParamPtrs{}
	}
	mmListOrders.defaultExpectation.paramPtrs.in = &in
	mmListOrders.defaultExpectation.expectationOrigins.originIn = minimock.CallerInfo(1)

	return mmListOrders
}

// ExpectOptsParam3 sets up expected param opts for ListOrdersService.ListOrders
func (mmListOrders *mListOrderServiceMockListOrders) ExpectOptsParam3(opts ...grpc.CallOption) *mListOrderServiceMockListOrders {
	if mmListOrders.mock.funcListOrders != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Set")
	}

	if mmListOrders.defaultExpectation == nil {
		mmListOrders.defaultExpectation = &ListOrderServiceMockListOrdersExpectation{}
	}

	if mmListOrders.defaultExpectation.params != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Expect")
	}

	if mmListOrders.defaultExpectation.paramPtrs == nil {
		mmListOrders.defaultExpectation.paramPtrs = &ListOrderServiceMockListOrdersParamPtrs{}
	}
	mmListOrders.defaultExpectation.paramPtrs.opts = &opts
	mmListOrders.defaultExpectation.expectationOrigins.originOpts = minimock.CallerInfo(1)

	return mmListOrders
}

// Inspect accepts an inspector function that has same arguments as the ListOrdersService.ListOrders
func (mmListOrders *mListOrderServiceMockListOrders) Inspect(f func(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption)) *mListOrderServiceMockListOrders {
	if mmListOrders.mock.inspectFuncListOrders != nil {
		mmListOrders.mock.t.Fatalf("Inspect function is already set for ListOrderServiceMock.ListOrders")
	}

	mmListOrders.mock.inspectFuncListOrders = f

	return mmListOrders
}

// Return sets up results that will be returned by ListOrdersService.ListOrders
func (mmListOrders *mListOrderServiceMockListOrders) Return(lp1 *order_service.ListOrdersResponse, err error) *ListOrderServiceMock {
	if mmListOrders.mock.funcListOrders != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Set")
	}

	if mmListOrders.defaultExpectation == nil {
		mmListOrders.defaultExpectation = &ListOrderServiceMockListOrdersExpectation{mock: mmListOrders.mock}
	}
	mmListOrders.defaultExpectation.results = &ListOrderServiceMockListOrdersResults{lp1, err}
	mmListOrders.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmListOrders.mock
}

// Set uses given function f to mock the ListOrdersService.ListOrders method
func (mmListOrders *mListOrderServiceMockListOrders) Set(f func(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption) (lp1 *order_service.ListOrdersResponse, err error)) *ListOrderServiceMock {
	if mmListOrders.defaultExpectation != nil {
		mmListOrders.mock.t.Fatalf("Default expectation is already set for the ListOrdersService.ListOrders method")
	}

	if len(mmListOrders.expectations) > 0 {
		mmListOrders.mock.t.Fatalf("Some expectations are already set for the ListOrdersService.ListOrders method")
	}

	mmListOrders.mock.funcListOrders = f
	mmListOrders.mock.funcListOrdersOrigin = minimock.CallerInfo(1)
	return mmListOrders.mock
}

// When sets expectation for the ListOrdersService.ListOrders which will trigger the result defined by the following
// Then helper
func (mmListOrders *mListOrderServiceMockListOrders) When(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption) *ListOrderServiceMockListOrdersExpectation {
	if mmListOrders.mock.funcListOrders != nil {
		mmListOrders.mock.t.Fatalf("ListOrderServiceMock.ListOrders mock is already set by Set")
	}

	expectation := &ListOrderServiceMockListOrdersExpectation{
		mock:               mmListOrders.mock,
		params:             &ListOrderServiceMockListOrdersParams{ctx, in, opts},
		expectationOrigins: ListOrderServiceMockListOrdersExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmListOrders.expectations = append(mmListOrders.expectations, expectation)
	return expectation
}

// Then sets up ListOrdersService.ListOrders return parameters for the expectation previously defined by the When method
func (e *ListOrderServiceMockListOrdersExpectation) Then(lp1 *order_service.ListOrdersResponse, err error) *ListOrderServiceMock {
	e.results = &ListOrderServiceMockListOrdersResults{lp1, err}
	return e.mock
}

// Times sets number of times ListOrdersService.ListOrders should be invoked
func (mmListOrders *mListOrderServiceMockListOrders) Times(n uint64) *mListOrderServiceMockListOrders {
	if n == 0 {
		mmListOrders.mock.t.Fatalf("Times of ListOrderServiceMock.ListOrders mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmListOrders.expectedInvocations, n)
	mmListOrders.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmListOrders
}

func (mmListOrders *mListOrderServiceMockListOrders) invocationsDone() bool {
	if len(mmListOrders.expectations) == 0 && mmListOrders.defaultExpectation == nil && mmListOrders.mock.funcListOrders == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmListOrders.mock.afterListOrdersCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmListOrders.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// ListOrders implements mm_commands.ListOrdersService
func (mmListOrders *ListOrderServiceMock) ListOrders(ctx context.Context, in *order_service.ListOrdersRequest, opts ...grpc.CallOption) (lp1 *order_service.ListOrdersResponse, err error) {
	mm_atomic.AddUint64(&mmListOrders.beforeListOrdersCounter, 1)
	defer mm_atomic.AddUint64(&mmListOrders.afterListOrdersCounter, 1)

	mmListOrders.t.Helper()

	if mmListOrders.inspectFuncListOrders != nil {
		mmListOrders.inspectFuncListOrders(ctx, in, opts...)
	}

	mm_params := ListOrderServiceMockListOrdersParams{ctx, in, opts}

	// Record call args
	mmListOrders.ListOrdersMock.mutex.Lock()
	mmListOrders.ListOrdersMock.callArgs = append(mmListOrders.ListOrdersMock.callArgs, &mm_params)
	mmListOrders.ListOrdersMock.mutex.Unlock()

	for _, e := range mmListOrders.ListOrdersMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.lp1, e.results.err
		}
	}

	if mmListOrders.ListOrdersMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmListOrders.ListOrdersMock.defaultExpectation.Counter, 1)
		mm_want := mmListOrders.ListOrdersMock.defaultExpectation.params
		mm_want_ptrs := mmListOrders.ListOrdersMock.defaultExpectation.paramPtrs

		mm_got := ListOrderServiceMockListOrdersParams{ctx, in, opts}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmListOrders.t.Errorf("ListOrderServiceMock.ListOrders got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListOrders.ListOrdersMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.in != nil && !minimock.Equal(*mm_want_ptrs.in, mm_got.in) {
				mmListOrders.t.Errorf("ListOrderServiceMock.ListOrders got unexpected parameter in, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListOrders.ListOrdersMock.defaultExpectation.expectationOrigins.originIn, *mm_want_ptrs.in, mm_got.in, minimock.Diff(*mm_want_ptrs.in, mm_got.in))
			}

			if mm_want_ptrs.opts != nil && !minimock.Equal(*mm_want_ptrs.opts, mm_got.opts) {
				mmListOrders.t.Errorf("ListOrderServiceMock.ListOrders got unexpected parameter opts, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmListOrders.ListOrdersMock.defaultExpectation.expectationOrigins.originOpts, *mm_want_ptrs.opts, mm_got.opts, minimock.Diff(*mm_want_ptrs.opts, mm_got.opts))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmListOrders.t.Errorf("ListOrderServiceMock.ListOrders got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmListOrders.ListOrdersMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmListOrders.ListOrdersMock.defaultExpectation.results
		if mm_results == nil {
			mmListOrders.t.Fatal("No results are set for the ListOrderServiceMock.ListOrders")
		}
		return (*mm_results).lp1, (*mm_results).err
	}
	if mmListOrders.funcListOrders != nil {
		return mmListOrders.funcListOrders(ctx, in, opts...)
	}
	mmListOrders.t.Fatalf("Unexpected call to ListOrderServiceMock.ListOrders. %v %v %v", ctx, in, opts)
	return
}

// ListOrdersAfterCounter returns a count of finished ListOrderServiceMock.ListOrders invocations
func (mmListOrders *ListOrderServiceMock) ListOrdersAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListOrders.afterListOrdersCounter)
}

// ListOrdersBeforeCounter returns a count of ListOrderServiceMock.ListOrders invocations
func (mmListOrders *ListOrderServiceMock) ListOrdersBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmListOrders.beforeListOrdersCounter)
}

// Calls returns a list of arguments used in each call to ListOrderServiceMock.ListOrders.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmListOrders *mListOrderServiceMockListOrders) Calls() []*ListOrderServiceMockListOrdersParams {
	mmListOrders.mutex.RLock()

	argCopy := make([]*ListOrderServiceMockListOrdersParams, len(mmListOrders.callArgs))
	copy(argCopy, mmListOrders.callArgs)

	mmListOrders.mutex.RUnlock()

	return argCopy
}

// MinimockListOrdersDone returns true if the count of the ListOrders invocations corresponds
// the number of defined expectations
func (m *ListOrderServiceMock) MinimockListOrdersDone() bool {
	if m.ListOrdersMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.ListOrdersMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.ListOrdersMock.invocationsDone()
}

// MinimockListOrdersInspect logs each unmet expectation
func (m *ListOrderServiceMock) MinimockListOrdersInspect() {
	for _, e := range m.ListOrdersMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ListOrderServiceMock.ListOrders at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterListOrdersCounter := mm_atomic.LoadUint64(&m.afterListOrdersCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.ListOrdersMock.defaultExpectation != nil && afterListOrdersCounter < 1 {
		if m.ListOrdersMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to ListOrderServiceMock.ListOrders at\n%s", m.ListOrdersMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to ListOrderServiceMock.ListOrders at\n%s with params: %#v", m.ListOrdersMock.defaultExpectation.expectationOrigins.origin, *m.ListOrdersMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcListOrders != nil && afterListOrdersCounter < 1 {
		m.t.Errorf("Expected call to ListOrderServiceMock.ListOrders at\n%s", m.funcListOrdersOrigin)
	}

	if !m.ListOrdersMock.invocationsDone() && afterListOrdersCounter > 0 {
		m.t.Errorf("Expected %d calls to ListOrderServiceMock.ListOrders at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.ListOrdersMock.expectedInvocations), m.ListOrdersMock.expectedInvocationsOrigin, afterListOrdersCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ListOrderServiceMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockListOrdersInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ListOrderServiceMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ListOrderServiceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockListOrdersDone()
}