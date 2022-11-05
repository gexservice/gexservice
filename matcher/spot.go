package matcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/centny/orderbook"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/debug"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

type SpotMatcher struct {
	Timeout           time.Duration
	PrecisionQuantity int32
	PrecisionPrice    int32
	Area              gexdb.BalanceArea
	Symbol            string
	Base              string
	Quote             string
	Fee               FeeCache
	NewOrderID        func() string
	PrepareProcess    func(ctx context.Context, matcher *SpotMatcher, userID int64) error
	Monitor           MatcherMonitor
	BestAsk           []decimal.Decimal
	BestBid           []decimal.Decimal
	bookVal           *orderbook.OrderBook
	bookLock          sync.RWMutex
}

func NewSpotMatcher(symbol, base, quote string, monitor MatcherMonitor) (matcher *SpotMatcher) {
	matcher = &SpotMatcher{
		Timeout:           5 * time.Second,
		PrecisionQuantity: 2,
		PrecisionPrice:    2,
		Area:              gexdb.BalanceAreaSpot,
		Symbol:            symbol,
		Base:              base,
		Quote:             quote,
		Fee:               ConstFee(0.002),
		NewOrderID:        gexdb.NewOrderID,
		PrepareProcess:    func(ctx context.Context, matcher *SpotMatcher, userID int64) error { return nil },
		Monitor:           monitor,
		bookVal:           orderbook.NewOrderBook(),
		bookLock:          sync.RWMutex{},
	}
	return
}

func (s *SpotMatcher) Bootstrap(ctx context.Context) (changed *MatcherEvent, err error) {
	changed = NewMatcherEvent(s.Symbol)
	var tx *pgx.Tx
	s.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("SpotMatcher bootstrap is panic with %v,\n%v", rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if tx != nil {
			if err == nil {
				err = tx.Commit(ctx)
			} else {
				tx.Rollback(ctx)
			}
		}
		s.bookLock.Unlock()
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] begin tx fail")
		return
	}
	var orders []*gexdb.Order
	err = gexdb.ScanOrderFilterWheref(ctx, "#all", "symbol=$%v,status=any($%v)", []interface{}{s.Symbol, gexdb.OrderStatusArray{gexdb.OrderStatusPending, gexdb.OrderStatusPartialled}}, "", &orders)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] query pending order by %v fail", converter.JSON([]interface{}{s.Symbol, gexdb.OrderStatusArray{gexdb.OrderStatusPending, gexdb.OrderStatusPartialled}}))
		return
	}
	for _, order := range orders {
		xlog.Infof("SpotMatcher bootstrap start cancel pending order %v", converter.JSON(order))
		if order.Filled.IsPositive() {
			order.Status = gexdb.OrderStatusPartCanceled
		} else {
			order.Status = gexdb.OrderStatusCanceled
		}
		//free balance
		err = s.syncBalanceByOrderDone(tx, ctx, changed, order)
		if err != nil {
			err = NewErrMatcher(err, "[ProcessCancel] sync balance by order %v fail", converter.JSON(order))
			return
		}

		//change status
		err = order.UpdateFilter(tx, ctx, "status")
		if err != nil {
			err = NewErrMatcher(err, "[ProcessCancel] change order status by order %v fail", converter.JSON(order))
			return
		}
		xlog.Infof("SpotMatcher bootstrap cancel pending order %v is success", order.OrderID)
		changed.AddOrder(order)
	}
	return
}

func (s *SpotMatcher) ProcessCancel(ctx context.Context, userID int64, orderID string) (order *gexdb.Order, err error) {
	args := &gexdb.Order{
		OrderID: orderID,
		UserID:  userID,
		Creator: userID,
		Status:  gexdb.OrderStatusCanceled,
	}
	order, err = s.ProcessOrder(ctx, args)
	return
}

func (s *SpotMatcher) ProcessMarket(ctx context.Context, userID int64, side gexdb.OrderSide, total, quantity decimal.Decimal) (order *gexdb.Order, err error) {
	args := &gexdb.Order{
		OrderID:    s.NewOrderID(),
		Type:       gexdb.OrderTypeTrade,
		UserID:     userID,
		Creator:    userID,
		Symbol:     s.Symbol,
		Side:       side,
		Quantity:   quantity,
		TotalPrice: total,
	}
	order, err = s.ProcessOrder(ctx, args)
	return
}

func (s *SpotMatcher) ProcessLimit(ctx context.Context, userID int64, side gexdb.OrderSide, quantity, price decimal.Decimal) (order *gexdb.Order, err error) {
	if price.Sign() <= 0 {
		err = fmt.Errorf("process limit userID/quantity/price is required or too small")
		err = NewErrMatcher(err, "[ProcessLimit] args invalid")
		return
	}
	args := &gexdb.Order{
		OrderID:  s.NewOrderID(),
		Type:     gexdb.OrderTypeTrade,
		UserID:   userID,
		Creator:  userID,
		Symbol:   s.Symbol,
		Side:     side,
		Quantity: quantity,
		Price:    price,
	}
	order, err = s.ProcessOrder(ctx, args)
	return
}

func (s *SpotMatcher) ProcessOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	if args.Status == gexdb.OrderStatusCanceled {
		if len(args.OrderID) < 1 {
			err = fmt.Errorf("process cancel orderID is required")
			err = NewErrMatcher(err, "[ProcessCancel] args invalid")
			return
		}
		if args.UserID <= 0 {
			err = fmt.Errorf("process cancel userID is required")
			err = NewErrMatcher(err, "[ProcessCancel] args invalid")
			return
		}
		order, err = s.processCancelOrder(ctx, args)
		return
	}
	err = s.PrepareProcess(ctx, s, args.UserID)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessOrder] prepare process fail")
		return
	}
	if args.Price.IsPositive() {
		args.Quantity = args.Quantity.Round(s.PrecisionQuantity)
		args.Price = args.Price.Round(s.PrecisionPrice)
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
		args.FeeRate, err = s.Fee.LoadFee(ctx, args.UserID, s.Symbol)
		if err == nil {
			order, err = s.processLimitOrder(ctx, args)
		}
	} else {
		args.Quantity = args.Quantity.Round(s.PrecisionQuantity)
		args.TotalPrice = args.TotalPrice.Round(s.PrecisionPrice)
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
			err = fmt.Errorf("process sell order.Quantity quantity is required or too small")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		if args.Side == gexdb.OrderSideBuy && args.TotalPrice.IsPositive() && len(s.BestAsk) > 0 && args.TotalPrice.DivRound(s.BestAsk[0], s.PrecisionQuantity).Sign() == 0 {
			err = fmt.Errorf("process buy market invest is too small")
			err = NewErrMatcher(err, "[ProcessOrder] args invalid")
			return
		}
		args.FeeRate, err = s.Fee.LoadFee(ctx, args.UserID, s.Symbol)
		if err == nil {
			order, err = s.processMarketOrder(ctx, args)
		}
	}
	return
}

func (s *SpotMatcher) processCancelOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	changed := NewMatcherEvent(s.Symbol)
	var tx *pgx.Tx
	var cancelOrder *orderbook.Order
	var rollback func()
	s.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("SpotMatcher process cancel by %v,%v is panic with %v,\n%v", args.UserID, args.OrderID, rerr, debug.CallStatck())
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
		cancel()
		changed.Depth = s.bookVal.Depth(30)
		s.BestAsk, s.BestBid = bestPrice(changed.Depth)
		s.bookLock.Unlock()

		//monitor
		if err == nil && s.Monitor != nil {
			changed.AddOrder(order)
			changed.AddMatched(nil, nil, cancelOrder)
			s.Monitor.OnMatched(ctx, changed)
		}
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] begin tx fail")
		return
	}

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

	//free balance
	err = s.syncBalanceByOrderDone(tx, ctx, changed, order)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] sync balance by order %v fail", converter.JSON(order))
		return
	}

	//change status
	err = order.UpdateFilter(tx, ctx, "status")
	if err != nil {
		err = NewErrMatcher(err, "[ProcessCancel] change order status by order %v fail", converter.JSON(order))
		return
	}

	//cancel order
	cancelOrder, rollback = s.bookVal.CancelOrder(order.OrderID)
	return
}

func (s *SpotMatcher) processMarketOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	//begin tx
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	changed := NewMatcherEvent(s.Symbol)
	var tx *pgx.Tx
	var doneOrder []*orderbook.Order
	var partOrder *orderbook.Order
	var partFilled decimal.Decimal
	var rollback func()
	s.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("SpotMatcher process market by %v,%v,%v,%v is panic with %v,\n%v", args.UserID, args.Side, args.TotalPrice, args.Quantity, rerr, debug.CallStatck())
			err = fmt.Errorf("%v", rerr)
		}
		if tx != nil {
			if err == nil {
				err = tx.Commit(ctx)
			} else {
				tx.Rollback(ctx)
			}
			// if err != nil && order != nil && order.Status == 0 {
			// 	order.Status = gexdb.OrderStatusCanceled
			// }
		}
		if err != nil && rollback != nil {
			rollback()
		}
		cancel()
		changed.Depth = s.bookVal.Depth(30)
		s.BestAsk, s.BestBid = bestPrice(changed.Depth)
		s.bookLock.Unlock()

		//monitor
		if err == nil && s.Monitor != nil {
			changed.AddOrder(order)
			changed.AddMatched(doneOrder, partOrder, nil)
			s.Monitor.OnMatched(ctx, changed)
		}
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] begin tx fail")
		return
	}

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
			OrderID: s.NewOrderID(),
			Type:    gexdb.OrderTypeTrade,
			UserID:  args.UserID,
			Creator: args.UserID,
			Symbol:  s.Symbol,
			Side:    args.Side,
			FeeRate: args.FeeRate,
		}
	}

	if order.Side == gexdb.OrderSideBuy {
		order.FeeBalance = s.Base
		if args.TotalPrice.IsPositive() {
			doneOrder, partOrder, partFilled, _, rollback, _ = s.bookVal.ProcessMarketPriceBuy(args.TotalPrice, s.PrecisionPrice)
		} else {
			doneOrder, partOrder, partFilled, _, rollback, _ = s.bookVal.ProcessMarketQuantityOrder(orderbook.Buy, args.Quantity)
		}
	} else {
		order.FeeBalance = s.Quote
		doneOrder, partOrder, partFilled, _, rollback, _ = s.bookVal.ProcessMarketQuantityOrder(orderbook.Sell, args.Quantity)
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
		order.AvgPrice = totalPrice.DivRound(totalQuantity, s.PrecisionPrice)
		if args.TotalPrice.IsPositive() {
			order.Quantity = args.TotalPrice.DivRound(order.AvgPrice, s.PrecisionPrice)
		} else {
			order.Quantity = args.Quantity
		}
	}
	order.Filled = totalQuantity
	order.TotalPrice = totalPrice

	if order.Side == gexdb.OrderSideBuy {
		order.InBalance = s.Base
		order.InFilled = order.Filled.Sub(order.Filled.Mul(order.FeeRate))
		order.OutBalance = s.Quote
		order.OutFilled = order.TotalPrice
		order.FeeBalance = s.Base
		order.FeeFilled = order.Filled.Mul(order.FeeRate)
	} else {
		order.InBalance = s.Quote
		order.InFilled = order.TotalPrice.Sub(order.TotalPrice.Mul(order.FeeRate))
		order.OutBalance = s.Base
		order.OutFilled = order.Filled
		order.FeeBalance = s.Quote
		order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
	}
	if totalPrice.IsPositive() && totalQuantity.IsPositive() {
		order.Transaction.Trans = s.allTrans(order, order.Price, doneOrder, partOrder, partFilled)
		if order.Quantity.Equal(order.Filled) {
			order.Status = gexdb.OrderStatusDone
		} else {
			order.Status = gexdb.OrderStatusPartCanceled
		}
	} else {
		order.Status = gexdb.OrderStatusCanceled
	}

	//reduce balance
	err = s.syncBalanceByOrderDone(tx, ctx, changed, order)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] sync balance by order %v", converter.JSON(order))
		return
	}

	//sync book order
	if totalPrice.IsPositive() && totalQuantity.IsPositive() {
		if len(doneOrder) > 0 {
			err = s.doneBookOrder(tx, ctx, changed, order, doneOrder...)
		}
		if err == nil && partOrder != nil {
			err = s.partBookOrder(tx, ctx, changed, order, partOrder, partFilled)
		}
		if err != nil {
			err = NewErrMatcher(err, "[ProcessMarket] sync book order fail")
			return
		}
	}

	//save order
	if order.TID > 0 {
		err = order.UpdateFilter(tx, ctx, "")
	} else {
		err = gexdb.AddOrderCall(tx, ctx, order)
	}
	if err != nil {
		err = NewErrMatcher(err, "[ProcessMarket] create order fail")
		return
	}
	return
}

func (s *SpotMatcher) processLimitOrder(ctx context.Context, args *gexdb.Order) (order *gexdb.Order, err error) {
	//begin tx
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	changed := NewMatcherEvent(s.Symbol)
	var tx *pgx.Tx
	var doneOrder []*orderbook.Order
	var partOrder *orderbook.Order
	var partFilled decimal.Decimal
	var rollback func()
	s.bookLock.Lock()
	defer func() {
		if rerr := recover(); rerr != nil {
			xlog.Errorf("SpotMatcher process limit by %v,%v,%v,%v is panic with %v,\n%v", args.UserID, args.Side, args.Quantity, args.Price, rerr, debug.CallStatck())
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
		cancel()
		changed.Depth = s.bookVal.Depth(30)
		s.BestAsk, s.BestBid = bestPrice(changed.Depth)
		s.bookLock.Unlock()

		//montiro
		if err == nil && s.Monitor != nil {
			changed.AddOrder(order)
			changed.AddMatched(doneOrder, partOrder, nil)
			s.Monitor.OnMatched(ctx, changed)
		}
	}()

	tx, err = gexdb.Pool().Begin(ctx)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] begin tx fail")
		return
	}
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
			OrderID:  s.NewOrderID(),
			Type:     gexdb.OrderTypeTrade,
			UserID:   args.UserID,
			Creator:  args.UserID,
			Symbol:   s.Symbol,
			Side:     args.Side,
			Quantity: args.Quantity,
			Price:    args.Price,
			FeeRate:  args.FeeRate,
		}
	}

	//lock balance
	lockedBalance := &gexdb.Balance{
		Area:   s.Area,
		UserID: order.UserID,
	}
	if order.Side == gexdb.OrderSideBuy {
		lockedBalance.Asset = s.Quote
		lockedBalance.Free = decimal.Zero.Sub(order.Quantity.Mul(order.Price))
		lockedBalance.Locked = order.Quantity.Mul(order.Price)
	} else {
		lockedBalance.Asset = s.Base
		lockedBalance.Free = decimal.Zero.Sub(order.Quantity)
		lockedBalance.Locked = order.Quantity
	}
	err = gexdb.IncreaseBalanceCall(tx, ctx, lockedBalance)
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] lock balance fail by %v", converter.JSON(lockedBalance))
		return
	}

	//process order

	var bookSide orderbook.Side
	if order.Side == gexdb.OrderSideBuy {
		bookSide = orderbook.Buy
		order.FeeBalance = s.Base
	} else {
		bookSide = orderbook.Sell
		order.FeeBalance = s.Quote
	}
	doneOrder, partOrder, partFilled, rollback, err = s.bookVal.ProcessLimitOrder(bookSide, order.OrderID, order.Quantity, order.Price)
	if err != nil {
		err = fmt.Errorf("process limit order fail with %v", err)
		err = NewErrMatcher(err, "[ProcessLimit] process limit order by %v fail", converter.JSON(order))
		return
	}

	//sync done partial ordr
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
	if len(refDoneOrder) > 0 {
		err = s.doneBookOrder(tx, ctx, changed, order, refDoneOrder...)
	}
	if err == nil && partOrder != nil && partOrder.ID() != order.OrderID {
		err = s.partBookOrder(tx, ctx, changed, order, partOrder, partFilled)
	}
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] sync order fail")
		return
	}

	//save order
	if order.Side == gexdb.OrderSideBuy {
		order.InBalance = s.Base
		order.InFilled = order.Filled.Sub(order.Filled.Mul(order.FeeRate))
		order.OutBalance = s.Quote
		order.OutFilled = order.TotalPrice
		order.FeeBalance = s.Base
		order.FeeFilled = order.Filled.Mul(order.FeeRate)
	} else {
		order.InBalance = s.Quote
		order.InFilled = order.TotalPrice.Sub(order.TotalPrice.Mul(order.FeeRate))
		order.OutBalance = s.Base
		order.OutFilled = order.Filled
		order.FeeBalance = s.Quote
		order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
	}
	order.Transaction.Trans = s.allTrans(order, order.Price, refDoneOrder, partOrder, partFilled)

	//unlock balance
	if order.Status == gexdb.OrderStatusDone {
		err = s.syncBalanceByOrderDone(tx, ctx, changed, order)
		if err != nil {
			err = NewErrMatcher(err, "[ProcessLimit] sync balance by order %v", converter.JSON(order))
			return
		}
	} else {
		changed.AddBalance(lockedBalance)
	}

	if order.TID > 0 {
		err = order.UpdateFilter(tx, ctx, "")
	} else {
		err = gexdb.AddOrderCall(tx, ctx, order)
	}
	if err != nil {
		err = NewErrMatcher(err, "[ProcessLimit] create order %v", converter.JSON(order))
		return
	}
	return
}

func (s *SpotMatcher) allTrans(base *gexdb.Order, price decimal.Decimal, doneOrders []*orderbook.Order, partOrder *orderbook.Order, partFilled decimal.Decimal) (trans []*gexdb.OrderTransactionItem) {
	for _, doneOrder := range doneOrders {
		tran := &gexdb.OrderTransactionItem{
			OrderID:    doneOrder.ID(),
			Filled:     doneOrder.Quantity(),
			Price:      price,
			TotalPrice: price.Mul(doneOrder.Quantity()),
			CreateTime: xsql.TimeNow(),
		}
		if base.Side == gexdb.OrderSideBuy {
			tran.FeeBalance = s.Base
			tran.FeeFilled = tran.Filled.Mul(base.FeeRate)
		} else {
			tran.FeeBalance = s.Quote
			tran.FeeFilled = tran.TotalPrice.Mul(base.FeeRate)
		}
		trans = append(trans, tran)
	}
	if partOrder != nil {
		tran := &gexdb.OrderTransactionItem{
			OrderID:    partOrder.ID(),
			Filled:     partFilled,
			Price:      price,
			TotalPrice: price.Mul(partFilled),
			CreateTime: xsql.TimeNow(),
		}
		if base.Side == gexdb.OrderSideBuy {
			tran.FeeBalance = s.Base
			tran.FeeFilled = tran.Filled.Mul(base.FeeRate)
		} else {
			tran.FeeBalance = s.Quote
			tran.FeeFilled = tran.TotalPrice.Mul(base.FeeRate)
		}
		trans = append(trans, tran)
	}
	return
}

func (s *SpotMatcher) doneBookOrder(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, base *gexdb.Order, bookOrders ...*orderbook.Order) (err error) {
	for _, bookOrder := range bookOrders {
		var order *gexdb.Order
		order, err = gexdb.FindOrderFilterWherefCall(tx, ctx, false, "order_id,type,user_id,side,quantity,filled,price,fee_rate,transaction#all", "order_id=$%v", bookOrder.ID())
		if err != nil {
			err = NewErrMatcher(err, "[doneBookOrder] find order by %v fail", bookOrder.ID())
			break
		}
		tran := &gexdb.OrderTransactionItem{
			OrderID:    base.OrderID,
			Filled:     bookOrder.Quantity(),
			Price:      order.Price,
			TotalPrice: order.Price.Mul(bookOrder.Quantity()),
			CreateTime: xsql.TimeNow(),
		}
		if order.Side == gexdb.OrderSideBuy {
			tran.FeeBalance = s.Base
			tran.FeeFilled = tran.Filled.Mul(order.FeeRate)
		} else {
			tran.FeeBalance = s.Quote
			tran.FeeFilled = tran.TotalPrice.Mul(order.FeeRate)
		}
		order.Transaction.Trans = append(order.Transaction.Trans, tran)
		order.Filled = order.Filled.Add(tran.Filled)
		order.TotalPrice = order.Filled.Mul(order.Price)
		if bookOrder.Side() == orderbook.Buy {
			order.InFilled = order.Filled.Sub(order.Filled.Mul(order.FeeRate))
			order.OutFilled = order.TotalPrice
			order.FeeFilled = order.Filled.Mul(order.FeeRate)
		} else {
			order.InFilled = order.TotalPrice.Sub(order.TotalPrice.Mul(order.FeeRate))
			order.OutFilled = order.Filled
			order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
		}
		order.Status = gexdb.OrderStatusDone
		err = s.updateOrder(tx, ctx, order, gexdb.OrderStatusPending, gexdb.OrderStatusPartialled)
		if err != nil {
			err = NewErrMatcher(err, "[doneBookOrder] update order by %v fail", converter.JSON(order))
			break
		}

		err = s.syncBalanceByOrderDone(tx, ctx, changed, order)
		if err != nil {
			err = NewErrMatcher(err, "[doneBookOrder] sync balance by order %v fail", converter.JSON(order))
			break
		}
	}
	return
}

func (s *SpotMatcher) partBookOrder(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, base *gexdb.Order, partOrder *orderbook.Order, partDone decimal.Decimal) (err error) {
	order, err := gexdb.FindOrderFilterWherefCall(tx, ctx, false, "order_id,type,user_id,side,quantity,filled,price,fee_rate,transaction#all", "order_id=$%v", partOrder.ID())
	if err != nil {
		err = NewErrMatcher(err, "[partBookOrder] find order by %v fail", partOrder.ID())
		return
	}

	tran := &gexdb.OrderTransactionItem{
		OrderID:    partOrder.ID(),
		Filled:     partDone,
		Price:      order.Price,
		TotalPrice: order.Price.Mul(partDone),
		CreateTime: xsql.TimeNow(),
	}
	if base.Side == gexdb.OrderSideBuy {
		tran.FeeBalance = s.Base
		tran.FeeFilled = tran.Filled.Mul(order.FeeRate)
	} else {
		tran.FeeBalance = s.Quote
		tran.FeeFilled = tran.TotalPrice.Mul(order.FeeRate)
	}
	order.Transaction.Trans = append(order.Transaction.Trans, tran)

	order.Filled = order.Filled.Add(partDone)
	order.TotalPrice = order.Filled.Mul(order.Price)
	if partOrder.Side() == orderbook.Buy {
		order.InFilled = order.Filled.Sub(order.Filled.Mul(order.FeeRate))
		order.OutFilled = order.TotalPrice
		order.FeeFilled = order.Filled.Mul(order.FeeRate)
	} else {
		order.InFilled = order.TotalPrice.Sub(order.TotalPrice.Mul(order.FeeRate))
		order.OutFilled = order.Filled
		order.FeeFilled = order.TotalPrice.Mul(order.FeeRate)
	}
	order.Status = gexdb.OrderStatusPartialled
	err = s.updateOrder(tx, ctx, order, gexdb.OrderStatusPending, gexdb.OrderStatusPartialled)
	if err != nil {
		err = NewErrMatcher(err, "[partBookOrder] update order by %v fail", converter.JSON(order))
		return
	}
	changed.AddOrder(order)
	return
}

func (s *SpotMatcher) updateOrder(tx *pgx.Tx, ctx context.Context, order *gexdb.Order, status ...gexdb.OrderStatus) (err error) {
	err = gexdb.UpdateOrderFilterWherefCall(tx, ctx, order, "filled,total_price,in_filled,out_filled,fee_filled,fee_rate,transaction,status", "order_id=$%v,status=any($%v)", order.OrderID, status)
	if err != nil {
		err = NewErrMatcher(err, "[updateOrder] upda order by %v fail", converter.JSON(order))
	}
	return
}

func (s *SpotMatcher) syncBalanceByOrderDone(tx *pgx.Tx, ctx context.Context, changed *MatcherEvent, order *gexdb.Order) (err error) {
	var in, out *gexdb.Balance
	if order.Side == gexdb.OrderSideBuy {
		in = &gexdb.Balance{
			UserID: order.UserID,
			Area:   s.Area,
			Asset:  s.Base,
			Free:   order.InFilled,
		}
		if order.Price.IsPositive() { //limit buy
			out = &gexdb.Balance{
				UserID: order.UserID,
				Area:   s.Area,
				Asset:  s.Quote,
				Free:   order.Quantity.Sub(order.Filled).Mul(order.Price),
				Locked: decimal.Zero.Sub(order.Quantity.Mul(order.Price)),
			}
		} else { //market buy
			out = &gexdb.Balance{
				UserID: order.UserID,
				Area:   s.Area,
				Asset:  s.Quote,
				Free:   decimal.Zero.Sub(order.TotalPrice),
			}
		}
	} else {
		in = &gexdb.Balance{
			UserID: order.UserID,
			Area:   s.Area,
			Asset:  s.Quote,
			Free:   order.InFilled,
		}
		if order.Price.IsPositive() { //limit sell
			out = &gexdb.Balance{
				UserID: order.UserID,
				Area:   s.Area,
				Asset:  s.Base,
				Free:   order.Quantity.Sub(order.Filled),
				Locked: decimal.Zero.Sub(order.Quantity),
			}
		} else { //market sell
			out = &gexdb.Balance{
				UserID: order.UserID,
				Area:   s.Area,
				Asset:  s.Base,
				Free:   decimal.Zero.Sub(order.Filled),
			}
		}
	}
	err = gexdb.IncreaseBalanceCall(tx, ctx, in)
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderDone] change in balance %v fail", converter.JSON(in))
		return
	}
	err = gexdb.IncreaseBalanceCall(tx, ctx, out)
	if err != nil {
		err = NewErrMatcher(err, "[syncBalanceByOrderDone] change out balance %v fail", converter.JSON(out))
		return
	}
	changed.AddBalance(in, out)
	return
}

func (s *SpotMatcher) Depth(max int) (depth *orderbook.Depth) {
	s.bookLock.RLock()
	defer s.bookLock.RUnlock()
	depth = s.bookVal.Depth(max)
	return
}
