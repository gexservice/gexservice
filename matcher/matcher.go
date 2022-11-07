package matcher

import (
	"context"
	"fmt"

	"github.com/centny/orderbook"
	"github.com/codingeasygo/util/xprop"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

var Quote string = "USDT"

type FeeCache interface {
	LoadFee(ctx context.Context, userID int64, symbol string) (fee decimal.Decimal, err error)
}

type ConstFee float64

func (c ConstFee) LoadFee(ctx context.Context, userID int64, symbol string) (fee decimal.Decimal, err error) {
	fee = decimal.NewFromFloat(float64(c))
	return
}

type MatcherEvent struct {
	Symbol       string
	Orders       []*gexdb.Order
	DoneOrder    map[string]bool
	PartOrder    map[string]bool
	CancelOrder  map[string]bool
	Balances     map[string]*gexdb.Balance
	Holdings     map[string]*gexdb.Holding
	Blowups      map[string]*gexdb.Holding
	DoneOrderIDs map[int64][]int64
	Depth        *orderbook.Depth
}

type MatcherMonitor interface {
	OnMatched(ctx context.Context, event *MatcherEvent)
}

type MatcherMonitorF func(ctx context.Context, event *MatcherEvent)

func (m MatcherMonitorF) OnMatched(ctx context.Context, event *MatcherEvent) { m(ctx, event) }

func NewMatcherEvent(symbol string) (event *MatcherEvent) {
	event = &MatcherEvent{
		Symbol:       symbol,
		DoneOrder:    map[string]bool{},
		PartOrder:    map[string]bool{},
		CancelOrder:  map[string]bool{},
		Balances:     map[string]*gexdb.Balance{},
		Holdings:     map[string]*gexdb.Holding{},
		Blowups:      map[string]*gexdb.Holding{},
		DoneOrderIDs: map[int64][]int64{},
	}
	return
}

func (m *MatcherEvent) AddOrder(orders ...*gexdb.Order) {
	m.Orders = append(m.Orders, orders...)
}

func (m *MatcherEvent) AddMatched(doneOrder []*orderbook.Order, partOrder, cancelOrder *orderbook.Order) {
	for _, order := range doneOrder {
		m.DoneOrder[order.ID()] = true
	}
	if partOrder != nil && !m.DoneOrder[partOrder.ID()] {
		m.PartOrder[partOrder.ID()] = true
	}
	if cancelOrder != nil {
		m.CancelOrder[cancelOrder.ID()] = true
	}
}

func (m *MatcherEvent) AddBalance(balances ...*gexdb.Balance) {
	for _, balance := range balances {
		m.Balances[BalanceKey(balance)] = balance
	}
}

func (m *MatcherEvent) AddHolding(holdings ...*gexdb.Holding) {
	for _, holding := range holdings {
		m.Holdings[HoldingKey(holding)] = holding
	}
}

func (m *MatcherEvent) AddBlowup(holdings ...*gexdb.Holding) {
	for _, holding := range holdings {
		m.Blowups[HoldingKey(holding)] = holding
	}
}

func BalanceKey(balance *gexdb.Balance) (key string) {
	key = fmt.Sprintf("%v-%v-%v", balance.UserID, balance.Area, balance.Asset)
	return
}

func HoldingKey(holding *gexdb.Holding) (key string) {
	key = fmt.Sprintf("%v-%v", holding.UserID, holding.Symbol)
	return
}

type ErrNotCancelable string

func (e ErrNotCancelable) Error() string { return string(e) }

type ErrStackable interface {
	error
	Stack() string
	IsBalanceNotEnought() bool
	IsBalanceNotFound() bool
	IsNotCancelable() bool
}

type ErrMatcher struct {
	Base  error
	Debug string
}

func NewErrMatcher(base error, format string, args ...interface{}) (err *ErrMatcher) {
	err = &ErrMatcher{
		Base:  base,
		Debug: fmt.Sprintf(format, args...),
	}
	return
}

func (e *ErrMatcher) Error() string {
	return e.Base.Error()
}

func (e *ErrMatcher) String() string {
	return e.Base.Error()
}

func (e *ErrMatcher) Stack() string {
	info := ""
	if v, ok := e.Base.(ErrStackable); ok {
		info = v.Stack()
	} else {
		info = "  " + e.Base.Error()
	}
	info += "\n  " + e.Debug
	return info
}

func (e *ErrMatcher) IsBalanceNotEnought() bool {
	return IsErrBalanceNotEnought(e.Base)
}

func (e *ErrMatcher) IsBalanceNotFound() bool {
	return IsErrBalanceNotFound(e.Base)
}

func (e *ErrMatcher) IsNotCancelable() bool {
	return IsErrNotCancelable(e.Base)
}

func ErrStack(err error) string {
	if v, ok := err.(ErrStackable); ok {
		return v.Stack()
	} else if err != nil {
		return err.Error()
	} else {
		return fmt.Sprintf("%v", err)
	}
}

func IsErrBalanceNotEnought(err error) bool {
	if v, ok := err.(ErrStackable); ok {
		return v.IsBalanceNotEnought()
	} else {
		return gexdb.IsErrBalanceNotEnought(err)
	}
}

func IsErrBalanceNotFound(err error) bool {
	if v, ok := err.(ErrStackable); ok {
		return v.IsBalanceNotFound()
	} else {
		return gexdb.IsErrBalanceNotFound(err)
	}
}

func IsErrNotCancelable(err error) bool {
	if v, ok := err.(ErrStackable); ok {
		return v.IsNotCancelable()
	} else {
		_, ok := err.(ErrNotCancelable)
		return ok
	}
}

type Matcher interface {
	Bootstrap(ctx context.Context) (changed *MatcherEvent, err error)
	ProcessCancel(ctx context.Context, userID int64, orderID string) (order *gexdb.Order, err error)
	ProcessMarket(ctx context.Context, userID int64, side gexdb.OrderSide, total, quantity decimal.Decimal) (order *gexdb.Order, err error)
	ProcessLimit(ctx context.Context, userID int64, side gexdb.OrderSide, quantity, price decimal.Decimal) (order *gexdb.Order, err error)
	ProcessOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error)
	ChangeLever(ctx context.Context, userID int64, lever int) (err error)
	Depth(max int) (depth *orderbook.Depth)
}

var Shared *MatcherCenter

func Bootstrap(conf *xprop.Config) (err error) {
	Shared, err = BootstrapMatcherCenterByConfig(conf)
	if err == nil {
		Shared.Start()
	}
	return
}

func ProcessCancel(ctx context.Context, userID int64, symbol string, orderID string) (order *gexdb.Order, err error) {
	order, err = Shared.ProcessCancel(ctx, userID, symbol, orderID)
	return
}

func ProcessMarket(ctx context.Context, userID int64, symbol string, side gexdb.OrderSide, total, quantity decimal.Decimal) (order *gexdb.Order, err error) {
	order, err = Shared.ProcessMarket(ctx, userID, symbol, side, total, quantity)
	return
}

func ProcessLimit(ctx context.Context, userID int64, symbol string, side gexdb.OrderSide, quantity, price decimal.Decimal) (order *gexdb.Order, err error) {
	order, err = Shared.ProcessLimit(ctx, userID, symbol, side, quantity, price)
	return
}

func ProcessOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	order, err = Shared.ProcessOrder(ctx, args)
	return
}

func ChangeLever(ctx context.Context, userID int64, symbol string, lever int) (err error) {
	err = Shared.ChangeLever(ctx, userID, symbol, lever)
	return
}

func bestPrice(depth *orderbook.Depth) (ask, bid []decimal.Decimal) {
	if depth == nil {
		return
	}
	if len(depth.Asks) > 0 {
		ask = depth.Asks[0]
	} else {
		ask = nil
	}
	if len(depth.Bids) > 0 {
		bid = depth.Bids[0]
	} else {
		bid = nil
	}
	return
}
