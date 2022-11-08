package gexdb

import (
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/shopspring/decimal"
)

func TestWithdraw(t *testing.T) {
	asset := "TEST"
	user := testAddUser("TestWithdraw")
	added, err := TouchBalance(ctx, BalanceAreaSpot, []string{asset}, user.TID)
	if err != nil || added != 1 {
		t.Error(err)
		return
	}
	balance := &Balance{
		UserID: user.TID,
		Area:   BalanceAreaSpot,
		Asset:  asset,
		Free:   decimal.NewFromFloat(100),
	}
	err = IncreaseBalance(ctx, balance)
	if err != nil || !balance.Free.Equal(decimal.NewFromFloat(100)) {
		t.Error(err)
		return
	}
	withdraw, err := CreateWithdraw(ctx, user.TID, asset, decimal.NewFromFloat(100))
	if err != nil || withdraw.Status != WithdrawStatusPending {
		t.Error(err)
		return
	}
	withdraw, err = CancelWithdraw(ctx, user.TID, withdraw.OrderID)
	if err != nil || withdraw.Status != WithdrawStatusCanceled {
		t.Error(err)
		return
	}
	//
	searcher := WithdrawUnifySearcher{}
	searcher.Where.Type = WithdrawTypeArray{WithdrawTypeWithdraw}
	searcher.Where.Asset = []string{asset}
	err = searcher.Apply(ctx)
	if err != nil || len(searcher.Query.Withdraws) < 1 {
		t.Error(err)
		return
	}
	//
	//test error
	_, err = CancelWithdraw(ctx, 11, withdraw.OrderID)
	if err == nil {
		t.Error(err)
		return
	}
	_, err = CancelWithdraw(ctx, user.TID, withdraw.OrderID)
	if err == nil {
		t.Error(err)
		return
	}
	withdraw2, err := CreateWithdraw(ctx, user.TID, asset, decimal.NewFromFloat(10))
	if err != nil {
		t.Error(err)
		return
	}
	withdraw3 := &Withdraw{
		OrderID:  NewOrderID(),
		Type:     WithdrawTypeTopup,
		UserID:   user.TID,
		Creator:  user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(1),
		Status:   WithdrawStatusPending,
	}
	err = AddWithdrawCall(Pool, ctx, withdraw3)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = CancelWithdraw(ctx, user.TID, withdraw3.OrderID)
	if err == nil {
		t.Error(err)
		return
	}
	//
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = CreateWithdraw(ctx, user.TID, asset, decimal.NewFromFloat(10))
		return
	})
	pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = CancelWithdraw(ctx, user.TID, withdraw.OrderID)
		return
	})
	pgx.MockerSetCall("Rows.Scan", 2).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = CancelWithdraw(ctx, user.TID, withdraw2.OrderID)
		return
	})
}
