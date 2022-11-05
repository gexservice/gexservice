package matcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/centny/orderbook"
	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/debug"
	"github.com/codingeasygo/util/xsort"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

type RollbackQueue []func()

func (r RollbackQueue) Call() {
	for i := len(r) - 1; i > -1; i-- {
		if r[i] != nil {
			r[i]()
		}
	}
}

type Rollbackable func() (func(), error)

type FuturesMatcher struct {
	Timeout           time.Duration
	PrecisionQuantity int32
	PrecisionPrice    int32
	Area              gexdb.BalanceArea
	Symbol            string
	Quote             string
	Fee               FeeCache
	MarginMax         decimal.Decimal
	MarginAdd         decimal.Decimal
	NewOrderID        func() string
	PrepareProcess    func(ctx context.Context, matcher *FuturesMatcher, userID int64) error
	Monitor           MatcherMonitor
	BestAsk           []decimal.Decimal
	BestBid           []decimal.Decimal
	bookUser          map[int64]map[int64]int
	bookVal           *orderbook.OrderBook
	bookLock          sync.RWMutex
}

func NewFuturesMatcher(symbol, quote string, monitor MatcherMonitor) (matcher *FuturesMatcher) {
	matcher = &FuturesMatcher{
		Timeout:           5 * time.Second,
		PrecisionQuantity: 2,
		PrecisionPrice:    2,
		Area:              gexdb.BalanceAreaFutures,
		Symbol:            symbol,
		Quote:             quote,
		Fee:               ConstFee(0.002),
		MarginMax:         decimal.NewFromFloat(0.99),
		MarginAdd:         decimal.NewFromFloat(0.05),
		NewOrderID:        gexdb.NewOrderID,
		PrepareProcess:    func(ctx context.Context, matcher *FuturesMatcher, userID int64) error { return nil },
		Monitor:           monitor,
		bookUser:          map[int64]map[int64]int{},
		bookVal:           orderbook.NewOrderBook(),
		bookLock:          sync.RWMutex{},
	}
	return
}

func (f *FuturesMatcher) Bootstrap(ctx context.Context) (changed *MatcherEvent, err error) {
	changed = NewMatcherEvent(f.Symbol)
	var tx *pgx.Tx
	f.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("FuturesMatcher bootstrap is panic with %v,\n%v", rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if tx != nil {
			if err == nil {
				err = tx.Commit(ctx)
			} else {
				tx.Rollback(ctx)
			}
		}
		f.bookLock.Unlock()
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] begin tx fail")
		return
	}
	var orders []*gexdb.Order
	err = gexdb.ScanOrderFilterWheref(ctx, "#all", "symbol=$%v,status=any($%v)", []interface{}{f.Symbol, gexdb.OrderStatusArray{gexdb.OrderStatusPending, gexdb.OrderStatusPartialled}}, "", &orders)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] query pending order by %v fail", converter.JSON([]interface{}{f.Symbol, gexdb.OrderStatusArray{gexdb.OrderStatusPending, gexdb.OrderStatusPartialled}}))
		return
	}
	for _, order := range orders {
		xlog.Infof("SpotMatcher bootstrap start cancel pending order %v", converter.JSON(order))
		if order.Filled.IsPositive() {
			order.Status = gexdb.OrderStatusPartCanceled
		} else {
			order.Status = gexdb.OrderStatusCanceled
		}

		//sync balance
		err = f.syncBalanceByOrderCancel(tx, ctx, changed, order)
		if err != nil {
			err = NewErrMatcher(err, "[ProcessCancel] sync balance by %v fail", converter.JSON(order))
			return
		}

		//change status
		err = order.UpdateFilter(tx, ctx, "status")
		if err != nil {
			err = NewErrMatcher(err, "[ProcessCancel] update order by %v fail", converter.JSON(order))
			return
		}

		xlog.Infof("SpotMatcher bootstrap cancel pending order %v is success", order.OrderID)
		changed.AddOrder(order)
	}
	return
}

func (f *FuturesMatcher) ProcessCancel(ctx context.Context, userID int64, orderID string) (order *gexdb.Order, err error) {
	args := &gexdb.Order{
		OrderID: orderID,
		UserID:  userID,
		Creator: userID,
		Status:  gexdb.OrderStatusCanceled,
	}
	order, err = f.ProcessOrder(ctx, args)
	return
}

func (f *FuturesMatcher) ProcessMarket(ctx context.Context, userID int64, side gexdb.OrderSide, total, quantity decimal.Decimal) (order *gexdb.Order, err error) {
	args := &gexdb.Order{
		OrderID:    f.NewOrderID(),
		Type:       gexdb.OrderTypeTrade,
		UserID:     userID,
		Creator:    userID,
		Symbol:     f.Symbol,
		Side:       side,
		Quantity:   quantity,
		TotalPrice: total,
	}
	order, err = f.ProcessOrder(ctx, args)
	return
}

func (f *FuturesMatcher) ProcessLimit(ctx context.Context, userID int64, side gexdb.OrderSide, quantity, price decimal.Decimal) (order *gexdb.Order, err error) {
	if price.Sign() <= 0 {
		err = fmt.Errorf("process limit userID/quantity/price is required or too small")
		err = NewErrMatcher(err, "[ProcessLimit] args invalid")
		return
	}
	args := &gexdb.Order{
		OrderID:  f.NewOrderID(),
		Type:     gexdb.OrderTypeTrade,
		UserID:   userID,
		Creator:  userID,
		Symbol:   f.Symbol,
		Side:     side,
		Quantity: quantity,
		Price:    price,
	}
	order, err = f.ProcessOrder(ctx, args)
	return
}

func (f *FuturesMatcher) ProcessOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	if args.Status == gexdb.OrderStatusCanceled {
		if args.UserID <= 0 {
			err = fmt.Errorf("process cancel userID is required")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if len(args.OrderID) < 1 {
			err = fmt.Errorf("process cancel orderID is required")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		order, err = f.processCancelOrder(ctx, args)
		return
	}
	err = f.PrepareProcess(ctx, f, args.UserID)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessOrder] prepare process fail")
		return
	}
	if args.Price.IsPositive() {
		//check args
		args.Quantity = args.Quantity.Round(f.PrecisionQuantity)
		args.Price = args.Price.Round(f.PrecisionPrice)
		if args.Side != gexdb.OrderSideBuy && args.Side != gexdb.OrderSideSell {
			err = fmt.Errorf("process limit side only supporte buy/sell")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if args.UserID <= 0 || !args.Quantity.IsPositive() || !args.Price.IsPositive() {
			err = fmt.Errorf("process limit userID/quantity/price is required or too small")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		args.FeeRate, err = f.Fee.LoadFee(ctx, args.UserID, f.Symbol)
		if err == nil {
			order, err = f.processLimitOrder(ctx, args)
		}
	} else {
		//check args
		args.Quantity = args.Quantity.Round(f.PrecisionQuantity)
		args.TotalPrice = args.TotalPrice.Round(f.PrecisionPrice)
		if args.Side != gexdb.OrderSideBuy && args.Side != gexdb.OrderSideSell {
			err = fmt.Errorf("process market side only supporte buy/sell")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if args.UserID <= 0 {
			err = fmt.Errorf("process market userID is required")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if args.Side == gexdb.OrderSideBuy && (!args.Quantity.IsPositive() && !args.TotalPrice.IsPositive()) {
			err = fmt.Errorf("process buy market quantity  or invest is required or too small")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if args.Side == gexdb.OrderSideSell && !args.Quantity.IsPositive() {
			err = fmt.Errorf("process sell market quantity is required or too small")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if args.Side == gexdb.OrderSideBuy && args.TotalPrice.IsPositive() && len(f.BestAsk) > 0 && args.TotalPrice.DivRound(f.BestAsk[0], f.PrecisionQuantity).Sign() == 0 {
			err = fmt.Errorf("process buy market invest is too small")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		args.FeeRate, err = f.Fee.LoadFee(ctx, args.UserID, f.Symbol)
		if err == nil {
			order, err = f.processMarketOrder(ctx, args)
		}
	}
	return
}

func (f *FuturesMatcher) processCancelOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	ctx, cancel := context.WithTimeout(ctx, f.Timeout)
	changed := NewMatcherEvent(f.Symbol)
	var tx *pgx.Tx
	var cancelOrder *orderbook.Order
	var rollback func()
	f.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("FuturesMatcher process cancel by %v,%v is panic with %v,\n%v", args.UserID, args.OrderID, rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if tx != nil {
			if err == nil {
				err = tx.Commit(ctx)
			} else {
				tx.Rollback(ctx)
			}
		}
		if err != nil && rollback != nil {
			rollback()
		}
		changed.AddOrder(order)
		changed.AddMatched(nil, nil, cancelOrder)
		if err == nil {
			f.syncUserOrder(changed)
		}
		cancel()
		changed.Depth = f.bookVal.Depth(30)
		f.BestAsk, f.BestBid = bestPrice(changed.Depth)
		f.bookLock.Unlock()

		//monitor
		if err == nil && f.Monitor != nil {
			f.Monitor.OnMatched(ctx, changed)
		}
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] begin tx fail")
		return
	}
	startDepth := f.bookVal.Depth(1)

	//find order
	order, err = gexdb.FindOrderByOrderIDCall(tx, ctx, args.OrderID, true)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] find order by %v fail", args.OrderID)
		return
	}
	if order.UserID != args.UserID {
		err = define.ErrNotAccess
		return
	}
	if order.Status != gexdb.OrderStatusPartialled && order.Status != gexdb.OrderStatusPending { //is done
		err = ErrNotCancelable(fmt.Sprintf("status is %v", order.Status))
		err = NewErrMatcher(err, "[ProcessCancel] cancel order by %v fail", args.OrderID)
		return
	}
	if order.Filled.IsPositive() {
		order.Status = gexdb.OrderStatusPartCanceled
	} else {
		order.Status = gexdb.OrderStatusCanceled
	}

	//sync balance
	err = f.syncBalanceByOrderCancel(tx, ctx, changed, order)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] sync balance by %v fail", converter.JSON(order))
		return
	}

	//change status
	err = order.UpdateFilter(tx, ctx, "status")
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] update order by %v fail", converter.JSON(order))
		return
	}

	//cancel order
	cancelOrder, rollback = f.bookVal.CancelOrder(order.OrderID)

	//check blowup and apply
	rb, err := f.checkBlowup(tx, ctx, changed, func() (func(), error) { return func() {}, nil })
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] process blowup by %v fail", converter.JSON(order))
		return
	}
	rollback = RollbackQueue{rollback, rb}.Call

	//free blowup
	err = f.freeBlowup(tx, ctx, changed, startDepth)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] free blowup by %v fail", converter.JSON(order))
		return
	}
	return
}

func (f *FuturesMatcher) processMarketOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	ctx, cancel := context.WithTimeout(ctx, f.Timeout)
	changed := NewMatcherEvent(f.Symbol)
	var tx *pgx.Tx
	var doneOrder []*orderbook.Order
	var partOrder *orderbook.Order
	var partFilled decimal.Decimal
	var rollback func()
	f.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("FuturesMatcher process market by %v,%v,%v,%v is panic with %v,\n%v", args.UserID, args.Side, args.TotalPrice, args.Quantity, rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if tx != nil {
			if err == nil {
				err = tx.Commit(ctx)
			} else {
				tx.Rollback(ctx)
			}
			if err != nil && order != nil && order.Status == 0 {
				order.Status = gexdb.OrderStatusCanceled
			}
		}
		if err != nil && rollback != nil {
			rollback()
		}
		changed.AddOrder(order)
		changed.AddMatched(doneOrder, partOrder, nil)
		if err == nil {
			f.syncUserOrder(changed)
		}
		cancel()
		changed.Depth = f.bookVal.Depth(30)
		f.BestAsk, f.BestBid = bestPrice(changed.Depth)
		f.bookLock.Unlock()

		//monitor
		if err == nil && f.Monitor != nil {
			f.Monitor.OnMatched(ctx, changed)
		}
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] begin tx fail")
		return
	}
	startDepth := f.bookVal.Depth(1)

	//process order
	if args.TID > 0 {
		order, err = gexdb.FindOrderWherefCall(tx, ctx, true, "tid=$%v", args.TID)
		if err != nil {
			err = NewErrMatcher(err, "[ProcessMarket] find order by %v fail", args.TID)
			return
		}
		if order.Status != gexdb.OrderStatusWaiting {
			err = NewErrMatcher(fmt.Errorf("status invalid"), "[ProcessLimit] order by %v is not waiting fail", args.TID)
			return
		}
		order.FeeRate = args.FeeRate
	} else {
		order = &gexdb.Order{
			OrderID: f.NewOrderID(),
			Type:    gexdb.OrderTypeTrade,
			UserID:  args.UserID,
			Creator: args.UserID,
			Symbol:  f.Symbol,
			Side:    args.Side,
			FeeRate: args.FeeRate,
		}
	}

	//check blowup and apply
	rollback, err = f.checkBlowup(tx, ctx, changed, func() (rb func(), _ error) {
		//may apply multi time
		if order.Side == gexdb.OrderSideBuy {
			if args.TotalPrice.IsPositive() {
				doneOrder, partOrder, partFilled, _, rb, _ = f.bookVal.ProcessMarketPriceBuy(args.TotalPrice, f.PrecisionPrice)
			} else {
				doneOrder, partOrder, partFilled, _, rb, _ = f.bookVal.ProcessMarketQuantityOrder(orderbook.Buy, args.Quantity)
			}
		} else {
			doneOrder, partOrder, partFilled, _, rb, _ = f.bookVal.ProcessMarketQuantityOrder(orderbook.Sell, args.Quantity)
		}
		return
	})
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] process blowup by %v fail", converter.JSON(order))
		return
	}

	//free blowup
	err = f.freeBlowup(tx, ctx, changed, startDepth)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] free blowup by %v fail", converter.JSON(order))
		return
	}

	totalQuantity := decimal.Zero
	totalPrice := decimal.Zero
	for _, order := range doneOrder {
		totalQuantity = totalQuantity.Add(order.Quantity())
		totalPrice = totalPrice.Add(order.Price().Mul(order.Quantity()))
	}

	if partFilled.Sign() > 0 {
		totalQuantity = totalQuantity.Add(partFilled)
		totalPrice = totalPrice.Add(partOrder.Price().Mul(partFilled))
	}

	if totalPrice.IsPositive() && totalQuantity.IsPositive() {
		order.AvgPrice = totalPrice.DivRound(totalQuantity, f.PrecisionPrice)
		if args.TotalPrice.IsPositive() {
			order.Quantity = args.TotalPrice.DivRound(order.AvgPrice, f.PrecisionPrice)
		} else {
			order.Quantity = args.Quantity
		}
	}
	order.Filled = totalQuantity
	order.TotalPrice = totalPrice
	order.FeeBalance = f.Quote
	order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
	if order.Side == gexdb.OrderSideBuy {
		order.Holding = order.Filled
	} else {
		order.Holding = decimal.Zero.Sub(order.Filled)
	}

	//sync balance
	err = f.syncBalanceByOrderAdd(tx, ctx, changed, order)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] sync balance by %v fail", converter.JSON(order))
		return
	}

	//sync holding
	order.Profit, err = f.syncHoldingByPartDone(tx, ctx, changed, order, order.Filled)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] sync holding by %v,%v fail", converter.JSON(order), order.Filled)
		return
	}

	//sync book order
	if totalPrice.IsPositive() && totalQuantity.IsPositive() {
		if len(doneOrder) > 0 {
			err = f.doneBookOrder(tx, ctx, changed, order.OrderID, doneOrder...)
		}
		if err == nil && partOrder != nil {
			err = f.partBookOrder(tx, ctx, changed, order.OrderID, partOrder, partFilled)
		}
		if err != nil {
			err = NewErrMatcher(err, "[ProcessMarket] sync order by %v fail", converter.JSON(order))
			return
		}
	}

	if totalPrice.IsPositive() && totalQuantity.IsPositive() {
		order.Transaction.Trans = f.allTrans(order, order.Price, doneOrder, partOrder, partFilled)
		if order.Quantity.Equal(order.Filled) {
			order.Status = gexdb.OrderStatusDone
		} else {
			order.Status = gexdb.OrderStatusPartCanceled
		}
	} else {
		order.Status = gexdb.OrderStatusCanceled
	}

	//save order
	if order.TID > 0 {
		err = order.UpdateFilter(tx, ctx, "")
	} else {
		err = gexdb.AddOrderCall(tx, ctx, order)
	}
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] add order by %v fail", converter.JSON(order))
		return
	}
	return
}

func (f *FuturesMatcher) processLimitOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	ctx, cancel := context.WithTimeout(ctx, f.Timeout)
	changed := NewMatcherEvent(f.Symbol)
	var tx *pgx.Tx
	var doneOrder []*orderbook.Order
	var partOrder *orderbook.Order
	var partFilled decimal.Decimal
	var rollback func()
	f.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("FuturesMatcher process limit by %v,%v,%v,%v is panic with %v,\n%v", args.UserID, args.Side, args.Quantity, args.Price, rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if tx != nil {
			if err == nil {
				err = tx.Commit(ctx)
			} else {
				tx.Rollback(ctx)
			}
			if err != nil && order != nil && order.Status == 0 {
				order.Status = gexdb.OrderStatusCanceled
			}
		}
		if err != nil && rollback != nil {
			rollback()
		}
		changed.AddOrder(order)
		changed.AddMatched(doneOrder, partOrder, nil)
		if err == nil {
			f.syncUserOrder(changed)
		}
		cancel()
		changed.Depth = f.bookVal.Depth(30)
		f.BestAsk, f.BestBid = bestPrice(changed.Depth)
		f.bookLock.Unlock()

		//montiro
		if err == nil && f.Monitor != nil {
			f.Monitor.OnMatched(ctx, changed)
		}
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] begin tx")
		return
	}
	startDepth := f.bookVal.Depth(1)

	//process order
	if args.TID > 0 {
		order, err = gexdb.FindOrderWherefCall(tx, ctx, true, "tid=$%v", args.TID)
		if err != nil {
			err = NewErrMatcher(err, "[ProcessLimit] find order by %v fail", args.TID)
			return
		}
		if order.Status != gexdb.OrderStatusWaiting {
			err = NewErrMatcher(fmt.Errorf("status invalid"), "[ProcessLimit] order by %v is not waiting fail", args.TID)
			return
		}
		order.FeeRate = args.FeeRate
	} else {
		order = &gexdb.Order{
			OrderID:  f.NewOrderID(),
			Type:     gexdb.OrderTypeTrade,
			UserID:   args.UserID,
			Creator:  args.UserID,
			Symbol:   f.Symbol,
			Side:     args.Side,
			Quantity: args.Quantity,
			Price:    args.Price,
			FeeRate:  args.FeeRate,
		}
	}

	//sync balance
	err = f.syncBalanceByOrderAdd(tx, ctx, changed, order)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] sync balance by %v fail", converter.JSON(order))
		return
	}

	//check blowup and apply
	var bookSide orderbook.Side
	if order.Side == gexdb.OrderSideBuy {
		bookSide = orderbook.Buy
	} else {
		bookSide = orderbook.Sell
	}
	rollback, err = f.checkBlowup(tx, ctx, changed, func() (rb func(), xerr error) {
		//may apply multi time
		doneOrder, partOrder, partFilled, rb, xerr = f.bookVal.ProcessLimitOrder(bookSide, order.OrderID, order.Quantity, order.Price)
		if xerr != nil {
			xerr = NewErrMatcher(xerr, "[ProcessLimit] process limit order by %v fail", converter.JSON(order))
		}
		return
	})
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] process blowup by %v fail", converter.JSON(order))
		return
	}

	//free blowup
	err = f.freeBlowup(tx, ctx, changed, startDepth)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] free blowup by %v fail", converter.JSON(order))
		return
	}

	//sync done partial order
	var refDoneOrder []*orderbook.Order
	if partOrder != nil && partOrder.ID() == order.OrderID { //current order partial
		refDoneOrder = doneOrder
		order.Filled = partFilled
		order.AvgPrice = order.Price
		order.TotalPrice = partFilled.Mul(order.Price)
		order.Status = gexdb.OrderStatusPartialled
	} else if len(doneOrder) > 0 { //current order is done
		curOrder := doneOrder[len(doneOrder)-1]
		refDoneOrder = doneOrder[:len(doneOrder)-1]
		order.Filled = curOrder.Quantity()
		order.AvgPrice = order.Price
		order.TotalPrice = curOrder.Quantity().Mul(order.Price)
		order.Status = gexdb.OrderStatusDone
	} else { //current order is pending
		order.Filled = decimal.Zero
		order.AvgPrice = order.Price
		order.TotalPrice = decimal.Zero
		order.Status = gexdb.OrderStatusPending
	}
	order.FeeBalance = f.Quote
	order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
	if order.Filled.IsPositive() {
		order.Profit, err = f.syncHoldingByPartDone(tx, ctx, changed, order, order.Filled)
	}
	if err == nil && len(refDoneOrder) > 0 {
		err = f.doneBookOrder(tx, ctx, changed, order.OrderID, refDoneOrder...)
	}
	if err == nil && partOrder != nil && partOrder.ID() != order.OrderID {
		err = f.partBookOrder(tx, ctx, changed, order.OrderID, partOrder, partFilled)
	}
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] sync order by %v fail", converter.JSON(order))
		return
	}

	//save order
	order.Transaction.Trans = f.allTrans(order, order.Price, refDoneOrder, partOrder, partFilled)
	if order.TID > 0 {
		err = order.UpdateFilter(tx, ctx, "")
	} else {
		err = gexdb.AddOrderCall(tx, ctx, order)
	}
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] add order by %v fail", converter.JSON(order))
		return
	}
	return
}

func (f *FuturesMatcher) checkBlowup(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, apply Rollbackable) (rollback func(), err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("FuturesMatcher check blow up is panic with %v, callstack is \n%v", rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if err != nil && rollback != nil {
			rollback()
			rollback = nil
		}
	}()
	depth := f.bookVal.Depth(1)
	if len(depth.Asks) < 1 || len(depth.Bids) < 1 {
		//current depth is too little, skip blowup
		rollback, err = apply()
		return
	}

	//try apply and check new depth will having blowup holding
	rollback, err = apply()
	if err != nil {
		return
	}
	depth = f.bookVal.Depth(1)
	if len(depth.Asks) < 1 || len(depth.Bids) < 1 {
		return
	}
	holdings, err := gexdb.ListHoldingForBlowupOverCall(tx, ctx, f.Symbol, depth.Asks[0][0], depth.Bids[0][0])
	if err != nil {
		err = NewErrMatcher(err, "[checkBlowup] list blowup holding by %v,%v,%v", f.Symbol, depth.Asks[0][0], depth.Bids[0][0])
		return
	}
	if len(holdings) < 1 {
		//skip blowup
		return
	}

	//rollback first and blowup all holding
	rollback()
	rollback = nil
	var rb func()
	var rollbackAll RollbackQueue
	for _, holding := range holdings {
		rb, err = f.blowupHolding(tx, ctx, changed, holding, depth.Asks[0][0], depth.Bids[0][0])
		if err != nil {
			err = NewErrMatcher(err, "[checkBlowup] blowup holding by %v,%v,%v", converter.JSON(holding), depth.Asks[0][0], depth.Bids[0][0])
			break
		}
		rollbackAll = append(rollbackAll, rb)
		rollback = rollbackAll.Call
	}
	if err != nil {
		return
	}
	changed.AddBlowup(holdings...)

	//apply again
	rb, err = apply()
	if err == nil {
		rollbackAll = append(rollbackAll, rb)
		rollback = rollbackAll.Call
	}
	return
}

func (f *FuturesMatcher) freeBlowup(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, startDepth *orderbook.Depth) (err error) {
	depth := f.bookVal.Depth(1)
	if len(depth.Asks) < 1 || len(depth.Bids) < 1 || //not depth
		(len(startDepth.Asks) > 0 && len(startDepth.Bids) > 0 && depth.Asks[0][0] == startDepth.Asks[0][0] && depth.Bids[0][0] == startDepth.Bids[0][0]) { //depth not changed
		return
	}
	ask, bid := depth.Asks[0][0], depth.Bids[0][0]
	holdings, err := gexdb.ListHoldingForBlowupFreeCall(tx, ctx, f.Symbol, ask, bid)
	if err != nil {
		err = NewErrMatcher(err, "[freeBlowup] list blowup holding by %v,%v,%v", f.Symbol, ask, bid)
		return
	}
	for _, holding := range holdings {
		var marginPrice decimal.Decimal
		if holding.Amount.IsPositive() {
			marginPrice = bid
		} else {
			marginPrice = ask
		}
		marginFree := marginPrice.Sub(holding.Blowup).Mul(holding.Amount).Round(f.PrecisionPrice)
		if marginFree.Div(holding.MarginUsed).LessThanOrEqual(f.MarginAdd) {
			continue
		}
		if marginFree.GreaterThanOrEqual(holding.MarginAdded) {
			marginFree = holding.MarginAdded
		}
		balance := &gexdb.Balance{
			UserID: holding.UserID,
			Area:   f.Area,
			Asset:  f.Quote,
			Free:   marginFree,
			Locked: decimal.Zero.Sub(marginFree),
			Margin: decimal.Zero.Sub(marginFree),
		}
		err = gexdb.IncreaseBalanceCall(tx, ctx, balance)
		if err != nil {
			err = NewErrMatcher(err, "[freeBlowup] free margin by %v fail", converter.JSON(balance))
			return
		}
		holding.MarginAdded = holding.MarginAdded.Sub(marginFree)
		holding.Blowup = holding.CalcBlowup(f.PrecisionPrice, f.MarginMax)
		err = holding.UpdateFilter(tx, ctx, "margin_added,blowup#all")
		if err != nil {
			err = NewErrMatcher(err, "[blowupHolding] update holding by %v fail", converter.JSON(holding))
			return
		}
		changed.AddBalance(balance)
		changed.AddHolding(holding)
	}
	return
}

func (f *FuturesMatcher) blowupHolding(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, holding *gexdb.Holding, ask, bid decimal.Decimal) (rollback func(), err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("FuturesMatcher blowup holding by %v,%v,%v is panic with %v, callstack is \n%v", converter.JSON(holding), ask, bid, rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if err != nil && rollback != nil {
			rollback()
			rollback = nil
		}
	}()
	balance, err := gexdb.FindBalanceByAssetCall(tx, ctx, holding.UserID, f.Area, f.Quote)
	if err != nil {
		err = NewErrMatcher(err, "[blowupHolding] find balance by %v,%v fail", holding.UserID, f.Quote)
		return
	}
	marginAdd := holding.CalcMargin(f.PrecisionPrice).Mul(f.MarginAdd)
	var marginPrice decimal.Decimal
	if holding.Amount.IsPositive() {
		marginPrice = bid
	} else {
		marginPrice = ask
	}
	marginAdd = marginAdd.Add(holding.Blowup.Sub(marginPrice).Mul(holding.Amount).Round(f.PrecisionPrice))
	if balance.Free.LessThan(marginAdd) {
		marginAdd = balance.Free
	}
	newHolding := holding.Copy()
	newHolding.MarginAdded = newHolding.MarginAdded.Add(marginAdd)
	newBlowup := newHolding.CalcBlowup(f.PrecisionPrice, f.MarginMax)
	if (holding.Amount.IsPositive() && newBlowup.LessThan(bid)) || (holding.Amount.IsNegative() && newBlowup.GreaterThan(ask)) {
		//should add margin
		balance := &gexdb.Balance{
			UserID: holding.UserID,
			Area:   f.Area,
			Asset:  f.Quote,
			Free:   decimal.Zero.Sub(marginAdd),
			Locked: marginAdd,
			Margin: marginAdd,
		}
		err = gexdb.IncreaseBalanceCall(tx, ctx, balance)
		if err != nil {
			err = NewErrMatcher(err, "[blowupHolding] add margin by %v fail", converter.JSON(balance))
			return
		}
		holding.MarginAdded = marginAdd
		holding.Blowup = holding.CalcBlowup(f.PrecisionPrice, f.MarginMax)
		err = holding.UpdateFilter(tx, ctx, "margin_added,blowup")
		if err != nil {
			err = NewErrMatcher(err, "[blowupHolding] update holding by %v fail", converter.JSON(holding))
			return
		}
		changed.AddBalance(balance)
		changed.AddHolding(holding)
		return
	}
	order := &gexdb.Order{
		OrderID:  f.NewOrderID(),
		Type:     gexdb.OrderTypeBlowup,
		UserID:   holding.UserID,
		Creator:  0,
		Symbol:   f.Symbol,
		Quantity: holding.Amount.Abs(),
	}
	var bookSide orderbook.Side
	if holding.Amount.IsNegative() {
		bookSide = orderbook.Buy
		order.Side = gexdb.OrderSideBuy
	} else {
		bookSide = orderbook.Sell
		order.Side = gexdb.OrderSideSell
	}
	doneOrder, partOrder, partFilled, _, rollback, _ := f.bookVal.ProcessMarketQuantityOrder(bookSide, holding.Amount.Abs())

	totalQuantity := decimal.Zero
	totalPrice := decimal.Zero
	for _, order := range doneOrder {
		totalQuantity = totalQuantity.Add(order.Quantity())
		totalPrice = totalPrice.Add(order.Price().Mul(order.Quantity()))
	}

	if partFilled.Sign() > 0 {
		totalQuantity = totalQuantity.Add(partFilled)
		totalPrice = totalPrice.Add(partOrder.Price().Mul(partFilled))
	}
	order.Filled = totalQuantity
	order.TotalPrice = totalPrice
	order.Owned = order.Quantity.Sub(order.Filled)
	order.Unhedged = order.Owned
	if totalPrice.IsPositive() && totalQuantity.IsPositive() {
		order.AvgPrice = totalPrice.DivRound(totalQuantity, f.PrecisionPrice)
	}
	order.FeeBalance = f.Quote
	order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
	if order.Side == gexdb.OrderSideBuy {
		order.Holding = order.Filled
	} else {
		order.Holding = decimal.Zero.Sub(order.Filled)
	}
	order.Transaction.Trans = f.allTrans(order, order.Price, doneOrder, partOrder, partFilled)
	if order.Quantity.Equal(order.Filled) {
		order.Status = gexdb.OrderStatusDone
	} else {
		order.Status = gexdb.OrderStatusPartCanceled
	}

	if len(doneOrder) > 0 {
		err = f.doneBookOrder(tx, ctx, changed, order.OrderID, doneOrder...)
	}
	if err == nil && partOrder != nil {
		err = f.partBookOrder(tx, ctx, changed, order.OrderID, partOrder, partFilled)
	}
	if err != nil {
		err = NewErrMatcher(err, "[blowupHolding] sync order by %v fail", converter.JSON(order))
		return
	}

	marginClear := holding.MarginUsed.Add(holding.MarginAdded)
	balance = &gexdb.Balance{
		UserID: holding.UserID,
		Area:   f.Area,
		Asset:  f.Quote,
		Free:   decimal.Zero.Sub(balance.Free),
		Locked: decimal.Zero.Sub(marginClear),
		Margin: decimal.Zero.Sub(marginClear),
	}
	err = gexdb.IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		err = NewErrMatcher(err, "[blowupHolding] blowup balance by %v fail", converter.JSON(balance))
		return
	}

	holding.Amount = decimal.Zero
	holding.Open = decimal.Zero
	holding.Blowup = decimal.Zero
	holding.MarginUsed = decimal.Zero
	holding.MarginAdded = decimal.Zero
	err = holding.UpdateFilter(tx, ctx, "amount,open,blowup,margin_used,margin_added#all")
	if err != nil {
		err = NewErrMatcher(err, "[blowupHolding] blowup holding by %v fail", converter.JSON(holding))
		return
	}
	changed.AddBalance(balance)
	changed.AddHolding(holding)

	//save order
	err = gexdb.AddOrderCall(tx, ctx, order)
	if err != nil {
		err = NewErrMatcher(err, "[blowupHolding] add order by %v fail", converter.JSON(order))
		return
	}
	changed.AddOrder(order)
	changed.AddMatched(doneOrder, partOrder, nil)
	return
}

func (f *FuturesMatcher) allTrans(base *gexdb.Order, price decimal.Decimal, doneOrders []*orderbook.Order, partOrder *orderbook.Order, partFilled decimal.Decimal) (trans []*gexdb.OrderTransactionItem) {
	for _, doneOrder := range doneOrders {
		tran := &gexdb.OrderTransactionItem{
			OrderID:    base.OrderID,
			Filled:     doneOrder.Quantity(),
			Price:      price,
			TotalPrice: price.Mul(doneOrder.Quantity()),
			FeeBalance: f.Quote,
			FeeFilled:  price.Mul(doneOrder.Quantity()).Mul(base.FeeRate),
			CreateTime: xsql.TimeNow(),
		}
		trans = append(trans, tran)
	}
	if partOrder != nil {
		tran := &gexdb.OrderTransactionItem{
			OrderID:    partOrder.ID(),
			Filled:     partFilled,
			Price:      price,
			TotalPrice: price.Mul(partFilled),
			FeeBalance: f.Quote,
			FeeFilled:  price.Mul(partFilled).Mul(base.FeeRate),
			CreateTime: xsql.TimeNow(),
		}
		trans = append(trans, tran)
	}
	return
}

func (f *FuturesMatcher) doneBookOrder(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, baseOrderID string, bookOrders ...*orderbook.Order) (err error) {
	var order *gexdb.Order
	for _, bookOrder := range bookOrders {
		order, err = gexdb.FindOrderByOrderIDCall(tx, ctx, bookOrder.ID(), false)
		if err != nil {
			err = NewErrMatcher(err, "[doneBookOrder] find order by %v fail", bookOrder.ID())
			break
		}
		tran := &gexdb.OrderTransactionItem{
			OrderID:    baseOrderID,
			Filled:     bookOrder.Quantity(),
			Price:      order.Price,
			TotalPrice: order.Price.Mul(bookOrder.Quantity()),
			FeeBalance: f.Quote,
			FeeFilled:  order.Price.Mul(bookOrder.Quantity()).Mul(order.FeeRate),
			CreateTime: xsql.TimeNow(),
		}
		order.Transaction.Trans = append(order.Transaction.Trans, tran)
		order.Filled = order.Filled.Add(tran.Filled)
		order.TotalPrice = order.Filled.Mul(order.Price)
		order.FeeBalance = f.Quote
		order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
		order.Status = gexdb.OrderStatusDone
		profit, xerr := f.syncHoldingByPartDone(tx, ctx, changed, order, tran.Filled)
		if xerr != nil {
			err = NewErrMatcher(xerr, "[doneBookOrder] sync holding by %v,%v fail", converter.JSON(order), tran.Filled)
			break
		}
		order.Profit = order.Profit.Add(profit)
		err = f.updateOrder(tx, ctx, order, gexdb.OrderStatusPending, gexdb.OrderStatusPartialled)
		if err != nil {
			err = NewErrMatcher(err, "[doneBookOrder] update order by %v fail", converter.JSON(order))
			break
		}
		changed.DoneOrderIDs[order.UserID] = append(changed.DoneOrderIDs[order.UserID], order.TID)
	}
	return
}

func (f *FuturesMatcher) partBookOrder(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, baseOrderID string, partOrder *orderbook.Order, partDone decimal.Decimal) (err error) {
	order, err := gexdb.FindOrderByOrderIDCall(tx, ctx, partOrder.ID(), false)
	if err != nil {
		err = NewErrMatcher(err, "[partOrder] find order by %v fail", partOrder.ID())
		return
	}
	tran := &gexdb.OrderTransactionItem{
		OrderID:    baseOrderID,
		Filled:     partDone,
		Price:      order.Price,
		TotalPrice: order.Price.Mul(partDone),
		FeeBalance: f.Quote,
		FeeFilled:  order.Price.Mul(partDone).Mul(order.FeeRate),
		CreateTime: xsql.TimeNow(),
	}
	order.Transaction.Trans = append(order.Transaction.Trans, tran)
	order.Filled = order.Filled.Add(partDone)
	order.TotalPrice = order.Filled.Mul(order.Price)
	order.FeeBalance = f.Quote
	order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
	order.Status = gexdb.OrderStatusPartialled
	_, err = f.syncHoldingByPartDone(tx, ctx, changed, order, tran.Filled)
	if err != nil {
		err = NewErrMatcher(err, "[partBookOrder] sync holding by %v,%v fail", converter.JSON(order), tran.Filled)
		return
	}
	err = f.updateOrder(tx, ctx, order, gexdb.OrderStatusPending, gexdb.OrderStatusPartialled)
	if err != nil {
		err = NewErrMatcher(err, "[partBookOrder] update order by %v fail", converter.JSON(order))
		return
	}
	changed.AddOrder(order)
	return
}

func (f *FuturesMatcher) updateOrder(tx *pgx.Tx, ctx context.Context, order *gexdb.Order, status ...gexdb.OrderStatus) (err error) {
	err = gexdb.UpdateOrderFilterWherefCall(tx, ctx, order, "filled,total_price,in_filled,out_filled,fee_filled,fee_rate,transaction,status", "order_id=$%v,status=any($%v)", order.OrderID, status)
	if err != nil {
		err = NewErrMatcher(err, "[updateOrder] update order by %v,%v", converter.JSON(order), converter.JSON(status))
	}
	return
}

func (f *FuturesMatcher) syncBalanceByOrderAdd(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, order *gexdb.Order) (err error) {
	holding, err := gexdb.FindHoldlingBySymbolCall(tx, ctx, order.UserID, order.Symbol, true)
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderAdd] find holding by %v,%v fail", order.UserID, order.Symbol)
		return
	}
	orders, err := f.listUserOrder(tx, ctx, order.UserID) //only limit order
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderCancel] list user order by %v fail", order.UserID)
		return
	}
	oldLocked := f.calcHoldingLocked(holding, orders, nil)
	newLocked := f.calcHoldingLocked(holding, orders, order)
	if newLocked.LessThanOrEqual(oldLocked) {
		//not changed: close only will change after sync holding
		return
	}
	//haivng more open
	balance := &gexdb.Balance{
		UserID: order.UserID,
		Area:   gexdb.BalanceAreaFutures,
		Asset:  f.Quote,
		Locked: newLocked.Sub(oldLocked),
		Free:   oldLocked.Sub(newLocked),
	}
	err = gexdb.IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderChanged] change balance %v fail", converter.JSON(balance))
		return
	}
	changed.AddBalance(balance)
	return
}

func (f *FuturesMatcher) syncBalanceByOrderCancel(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, order *gexdb.Order) (err error) {
	holding, err := gexdb.FindHoldlingBySymbolCall(tx, ctx, order.UserID, order.Symbol, true)
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderCancel] find holding by %v,%v fail", order.UserID, order.Symbol)
		return
	}
	oldOrders, err := f.listUserOrder(tx, ctx, order.UserID) //only limit order
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderCancel] list user order by %v fail", order.UserID)
		return
	}
	oldLocked := f.calcHoldingLocked(holding, oldOrders, nil)
	newOrders := []*gexdb.Order{}
	for _, oldOrder := range oldOrders {
		if oldOrder.OrderID == order.OrderID {
			continue
		}
		newOrders = append(newOrders, oldOrder)
	}
	newLocked := f.calcHoldingLocked(holding, newOrders, nil)
	//haivng cancel
	balance := &gexdb.Balance{
		UserID: order.UserID,
		Area:   gexdb.BalanceAreaFutures,
		Asset:  f.Quote,
		Locked: newLocked.Sub(oldLocked),
		Free:   oldLocked.Sub(newLocked),
	}
	err = gexdb.IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderCancel] change balance %v fail", converter.JSON(balance))
		return
	}
	changed.AddBalance(balance)
	return
}

func (f *FuturesMatcher) calcHoldingLocked(holding *gexdb.Holding, orders []*gexdb.Order, newOrder *gexdb.Order) (total decimal.Decimal) {
	holdingAmount := holding.Amount
	totalPrice := decimal.Zero
	fee := decimal.Zero
	closeOnly := true
	for _, order := range orders {
		remain := order.Quantity.Sub(order.Filled)
		fee = fee.Add(remain.Mul(order.Price).Mul(order.FeeRate))
		if order.Side == gexdb.OrderSideSell {
			remain = decimal.Zero.Sub(remain)
		}
		if closeOnly && holdingAmount.Sign() != 0 && holdingAmount.Sign() != remain.Sign() { //only first close order
			if remain.Abs().LessThanOrEqual(holdingAmount.Abs()) {
				holdingAmount = holdingAmount.Add(remain)
				continue
			}
			remain = remain.Add(holdingAmount)
			holdingAmount = decimal.Zero
		}
		totalPrice = totalPrice.Add(remain.Abs().Mul(order.Price))
		closeOnly = false
	}
	if newOrder != nil {
		var remain decimal.Decimal
		if newOrder.Price.IsPositive() {
			fee = fee.Add(newOrder.Quantity.Mul(newOrder.Price).Mul(newOrder.FeeRate))
			remain = newOrder.Quantity.Sub(newOrder.Filled)
		} else {
			fee = fee.Add(newOrder.TotalPrice.Mul(newOrder.FeeRate))
			remain = newOrder.Filled
		}
		if newOrder.Side == gexdb.OrderSideSell {
			remain = decimal.Zero.Sub(remain)
		}
		if closeOnly && holdingAmount.Sign() != 0 && holdingAmount.Sign() != remain.Sign() {
			if remain.Abs().LessThanOrEqual(holdingAmount.Abs()) {
				remain = decimal.Zero
			} else {
				remain = remain.Add(holdingAmount)
			}
		}
		if newOrder.Price.IsPositive() {
			totalPrice = totalPrice.Add(remain.Abs().Mul(newOrder.Price))
		} else {
			totalPrice = totalPrice.Add(remain.Abs().Mul(newOrder.AvgPrice))
		}
	}
	totalPrice = totalPrice.Add(holding.Amount.Abs().Mul(holding.Open))
	margin := totalPrice.DivRound(decimal.NewFromInt(int64(holding.Lever)), f.PrecisionPrice)
	total = margin.Add(fee)
	return
}

func (f *FuturesMatcher) syncHoldingByPartDone(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, order *gexdb.Order, partDone decimal.Decimal) (profit decimal.Decimal, err error) {
	if partDone.IsZero() {
		return
	}
	holding, err := gexdb.FindHoldlingBySymbolCall(tx, ctx, order.UserID, order.Symbol, true)
	if err != nil {
		err = NewErrMatcher(err, "[syncHoldingByPartDone] find holding by %v,%v fail", order.UserID, order.Symbol)
		return
	}
	partHolding := partDone
	if order.Side == gexdb.OrderSideSell {
		partHolding = decimal.Zero.Sub(partHolding)
	}
	balance := &gexdb.Balance{
		UserID: order.UserID,
		Area:   f.Area,
		Asset:  f.Quote,
	}
	remain := holding.Amount.Add(partHolding)
	if holding.Amount.Sign() == 0 || holding.Amount.Sign() == partHolding.Sign() { //open new or open more
		margin := holding.CalcMargin(f.PrecisionPrice)
		holding.Open = holding.Amount.Mul(holding.Open).Add(partHolding.Mul(order.AvgPrice)).Div(holding.Amount.Add(partHolding)).Abs().Round(f.PrecisionPrice)
		holding.Amount = holding.Amount.Add(partHolding)
		marginChange := holding.CalcMargin(f.PrecisionPrice).Sub(margin)
		balance.Margin = marginChange
	} else if remain.Sign() == 0 || remain.Sign() == holding.Amount.Sign() { //close only
		profit = partHolding.Mul(decimal.NewFromInt(-1)).Mul(order.AvgPrice.Sub(holding.Open))
		marginUsed := holding.CalcMargin(f.PrecisionPrice)
		marginAdded := partHolding.Div(holding.Amount).Mul(holding.MarginAdded)
		holding.Amount = remain
		holding.MarginAdded = holding.MarginAdded.Sub(marginAdded)
		marginChange := holding.CalcMargin(f.PrecisionPrice).Sub(marginUsed).Sub(marginAdded)
		balance.Margin = marginChange
		balance.Locked = marginChange
		balance.Free = decimal.Zero.Sub(marginChange)
		balance.Free = balance.Free.Add(profit)
	} else { //close first and open new
		//close first
		profit = holding.Amount.Mul(order.AvgPrice.Sub(holding.Open))
		marginUsed := holding.CalcMargin(f.PrecisionPrice)
		marginAdded := holding.MarginAdded
		marginChange := decimal.Zero.Sub(marginUsed).Sub(marginAdded)
		balance.Margin = marginChange
		balance.Locked = marginChange
		balance.Free = decimal.Zero.Sub(marginChange)
		balance.Free = balance.Free.Add(profit)
		//open new
		holding.Open = order.AvgPrice
		holding.Amount = remain
		holding.MarginAdded = decimal.Zero
		marginOpen := holding.CalcMargin(f.PrecisionPrice)
		balance.Margin = balance.Margin.Add(marginOpen)
	}
	fee := partHolding.Abs().Mul(order.AvgPrice).Mul(order.FeeRate)
	balance.Locked = balance.Locked.Sub(fee)
	holding.MarginUsed = holding.CalcMargin(f.PrecisionPrice)
	holding.Blowup = holding.CalcBlowup(f.PrecisionPrice, f.MarginMax)
	if holding.Amount.Sign() == 0 {
		holding.Open = decimal.Zero
	}
	err = gexdb.IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		err = NewErrMatcher(err, "[syncHoldingByPartDone] change balance %v fail", converter.JSON(balance))
		return
	}
	err = holding.UpdateFilter(tx, ctx, "amount,open,margin_used,blowup#all")
	if err != nil {
		err = NewErrMatcher(err, "[syncHoldingByPartDone] change holding %v fail", converter.JSON(holding))
		return
	}
	changed.AddBalance(balance)
	changed.AddHolding(holding)
	return
}

func (f *FuturesMatcher) syncUserOrder(changed *MatcherEvent) {
	for _, order := range changed.Orders {
		switch order.Status {
		case gexdb.OrderStatusCanceled, gexdb.OrderStatusPartCanceled, gexdb.OrderStatusDone:
			if f.bookUser[order.UserID] != nil {
				delete(f.bookUser[order.UserID], order.TID)
			}
		case gexdb.OrderStatusPending, gexdb.OrderStatusPartialled:
			if f.bookUser[order.UserID] == nil {
				f.bookUser[order.UserID] = map[int64]int{}
			}
			f.bookUser[order.UserID][order.TID] = 1
		}
		if len(f.bookUser[order.UserID]) < 1 {
			delete(f.bookUser, order.UserID)
		}
	}
	for userID, orderIDs := range changed.DoneOrderIDs {
		for _, orderID := range orderIDs {
			if f.bookUser[userID] != nil {
				delete(f.bookUser[userID], orderID)
			}
		}
		if len(f.bookUser[userID]) < 1 {
			delete(f.bookUser, userID)
		}
	}
}

func (f *FuturesMatcher) listUserOrder(caller crud.Queryer, ctx context.Context, userID int64) (orders []*gexdb.Order, err error) {
	orderIDs := []int64{}
	for orderID := range f.bookUser[userID] {
		orderIDs = append(orderIDs, orderID)
	}
	orders, _, err = gexdb.ListOrderByIDCall(caller, ctx, orderIDs...)
	if err == nil {
		xsort.SortFunc(orders, func(x, y int) bool {
			return orders[x].CreateTime.AsTime().Before(orders[y].CreateTime.AsTime())
		})
	}
	return
}

func (f *FuturesMatcher) Depth(max int) (depth *orderbook.Depth) {
	f.bookLock.RLock()
	defer f.bookLock.RUnlock()
	depth = f.bookVal.Depth(max)
	return
}
