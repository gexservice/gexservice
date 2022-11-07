package matcher

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/centny/orderbook"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xprop"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/baseupgrade"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/gexupgrade"
	"github.com/shopspring/decimal"
)

var ctx = context.Background()

func init() {
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/gexservice")
	if err != nil {
		panic(err)
	}
	gexdb.Pool = pgx.Pool
	basedb.SYS = "gex"
	basedb.Pool = pgx.Pool
	_, _, err = gexdb.Pool().Exec(ctx, gexupgrade.DROP)
	if err != nil {
		panic(err)
	}
	_, _, err = gexdb.Pool().Exec(ctx, strings.ReplaceAll(baseupgrade.DROP, "_sys_", "gex_"))
	if err != nil {
		panic(err)
	}
	_, err = basedb.CheckDb()
	if err != nil {
		panic(err)
	}
	_, err = gexdb.CheckDb(ctx)
	if err != nil {
		panic(err)
	}
	go http.ListenAndServe(":6062", nil)
}

func clear() {
	_, _, err := gexdb.Pool().Exec(ctx, gexupgrade.CLEAR)
	if err != nil {
		panic(err)
	}
}

func testAddUser(prefix string) (user *gexdb.User) {
	account, phone, password := prefix+"_acc", prefix+"_123", "123"
	image := prefix + "_image"
	user = &gexdb.User{
		Type:      gexdb.UserTypeNormal,
		Role:      gexdb.UserRoleNormal,
		Name:      &prefix,
		Account:   &account,
		Phone:     &phone,
		Image:     &image,
		Password:  &password,
		TradePass: &password,
		External:  xsql.M{"abc": 1},
		Status:    gexdb.UserStatusNormal,
	}
	err := gexdb.AddUser(ctx, user)
	if err != nil {
		panic(err)
	}
	return
}

func assetOrderStatus(orderID string, status gexdb.OrderStatus) {
	order, err := gexdb.FindOrderByOrderID(ctx, orderID)
	if err != nil {
		panic(err)
	}
	if order.Status != status {
		panic(fmt.Sprintf("status is %v", order.Status))
	}
}

func assetBalanceMargin(userID int64, area gexdb.BalanceArea, asset string, margin decimal.Decimal) {
	_, balances, err := gexdb.ListUserBalance(ctx, userID, area, []string{asset}, nil)
	if err != nil {
		panic(err)
	}
	if !balances[asset].Margin.Equal(margin) {
		panic(fmt.Sprintf("margin is %v", balances[asset].Margin))
	}
}

func assetBalanceLocked(userID int64, area gexdb.BalanceArea, asset string, locked decimal.Decimal) {
	_, balances, err := gexdb.ListUserBalance(ctx, userID, area, []string{asset}, nil)
	if err != nil {
		panic(err)
	}
	if !balances[asset].Locked.Equal(locked) {
		panic(fmt.Sprintf("locked is %v", balances[asset].Locked))
	}
}

func assetBalanceFree(userID int64, area gexdb.BalanceArea, asset string, free decimal.Decimal) {
	_, balances, err := gexdb.ListUserBalance(ctx, userID, area, []string{asset}, nil)
	if err != nil {
		panic(err)
	}
	if !balances[asset].Free.Equal(free) {
		panic(fmt.Sprintf("user %v %v free is %v", balances[asset].UserID, asset, balances[asset].Free))
	}
}

func assetHoldingAmount(userID int64, symbol string, amount decimal.Decimal) {
	holding, err := gexdb.FindHoldlingBySymbol(ctx, userID, symbol)
	if err != nil {
		panic(err)
	}
	if !holding.Amount.Equal(amount) {
		panic(fmt.Sprintf("amount is %v", holding.Amount))
	}
}

func assetDepthMust(depth *orderbook.Depth, bids, asks int) {
	if len(depth.Bids) != bids {
		panic(fmt.Sprintf("bids is not %v", asks))
	}
	if len(depth.Asks) != asks {
		panic(fmt.Sprintf("asks is not %v", asks))
	}
}

func assetDepthEmpty(depth *orderbook.Depth) {
	if len(depth.Asks) > 0 || len(depth.Bids) > 0 {
		panic("not empty")
	}
}

func TestErrMatcher(t *testing.T) {
	err := fmt.Errorf("error")
	notEnought := NewErrMatcher(gexdb.ErrBalanceNotEnought("Not Enought"), "abc")
	notFound := NewErrMatcher(gexdb.ErrBalanceNotFound("Not Found"), "abc")
	notCancelable := ErrNotCancelable("Not Cancelable")
	if !IsErrBalanceNotEnought(notEnought) {
		t.Error(notEnought)
		return
	}
	if IsErrBalanceNotEnought(err) {
		t.Error(err)
		return
	}
	if !IsErrBalanceNotFound(notFound) {
		t.Error(notFound)
		return
	}
	if IsErrBalanceNotFound(err) {
		t.Error(err)
		return
	}
	if !IsErrNotCancelable(notCancelable) {
		t.Error(notFound)
		return
	}
	if IsErrNotCancelable(err) {
		t.Error(err)
		return
	}
	fmt.Printf("err->%v\n", notEnought.Error())
	fmt.Printf("string->%v\n", notEnought.String())
	fmt.Printf("print->%v\n", notEnought)
	fmt.Printf("stack->\n%v\n", ErrStack(notEnought))
	fmt.Printf("stack->\n%v\n", ErrStack(notFound))
	fmt.Printf("stack->\n%v\n", ErrStack(fmt.Errorf("error")))
	fmt.Printf("stack->\n%v\n", ErrStack(nil))
}

func ParallelTest(total, max int64, call func(i int64)) (elapsed time.Duration, avg float64) {
	waiter := sync.WaitGroup{}
	queue := make(chan int64, total)
	for i := int64(0); i < max; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for {
				v := <-queue
				if v < 0 {
					break
				}
				call(v)
			}
		}()
	}
	begin := time.Now()
	for i := int64(0); i < total; i++ {
		queue <- i
	}
	for i := int64(0); i < max; i++ {
		queue <- -1
	}
	waiter.Wait()
	elapsed = time.Since(begin)
	avg = float64(total) / elapsed.Seconds()
	return
}

func TestShared(t *testing.T) {
	clear()
	config := xprop.NewConfig()
	config.LoadPropString(matcherConfig)
	err := Bootstrap(config)
	if err != nil {
		t.Error(err)
		return
	}
	ProcessCancel(ctx, 0, "", "")
	ProcessLimit(ctx, 0, "", gexdb.OrderSideBuy, decimal.Zero, decimal.Zero)
	ProcessMarket(ctx, 0, "", gexdb.OrderSideBuy, decimal.Zero, decimal.Zero)
	ProcessOrder(ctx, &gexdb.Order{})
	for _, symobl := range Shared.Symbols {
		fmt.Printf("--->%v\n", symobl)
	}
}

func TestSome(t *testing.T) {
	bestPrice(nil)
}
