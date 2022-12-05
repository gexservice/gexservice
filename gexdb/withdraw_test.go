package gexdb

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/uuid"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/shopspring/decimal"
)

func TestWithdraw(t *testing.T) {
	asset := "TEST"
	user := testAddUser("TestWithdraw")
	basedb.StoreConf(ctx, ConfigWithdrawReview, converter.JSON(xmap.M{asset: 100}))
	added, err := TouchBalance(ctx, BalanceAreaSpot, []string{asset, "x1"}, user.TID)
	if err != nil || added != 2 {
		t.Error(err)
		return
	}
	balance := &Balance{
		UserID: user.TID,
		Area:   BalanceAreaSpot,
		Asset:  asset,
		Free:   decimal.NewFromFloat(10000),
	}
	err = IncreaseBalance(ctx, balance)
	if err != nil || !balance.Free.Equal(decimal.NewFromFloat(10000)) {
		t.Error(err)
		return
	}
	balance = &Balance{
		UserID: user.TID,
		Area:   BalanceAreaSpot,
		Asset:  "x1",
		Free:   decimal.NewFromFloat(100),
	}
	err = IncreaseBalance(ctx, balance)
	if err != nil || !balance.Free.Equal(decimal.NewFromFloat(100)) {
		t.Error(err)
		return
	}
	withdraw := &Withdraw{
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(100),
	}
	err = CreateWithdraw(ctx, withdraw)
	if err != nil || withdraw.Status != WithdrawStatusPending {
		t.Error(err)
		return
	}
	withdraw, err = CancelWithdraw(ctx, user.TID, withdraw.OrderID)
	if err != nil || withdraw.Status != WithdrawStatusCanceled {
		t.Error(err)
		return
	}
	withdraw00 := &Withdraw{
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(100),
	}
	err = CreateWithdraw(ctx, withdraw00)
	if err != nil || withdraw00.Status != WithdrawStatusPending {
		t.Error(err)
		return
	}
	err = ConfirmWithdraw(ctx, withdraw00.OrderID)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = DoneWithdraw(ctx, withdraw00.OrderID, false, xmap.M{"A": 123})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = DoneWithdraw(ctx, withdraw00.OrderID, false, xmap.M{"A": 123})
	if err != nil {
		t.Error(err)
		return
	}
	withdraw01 := &Withdraw{
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(100),
	}
	err = CreateWithdraw(ctx, withdraw01)
	if err != nil || withdraw01.Status != WithdrawStatusPending {
		t.Error(err)
		return
	}
	err = ConfirmWithdraw(ctx, withdraw01.OrderID)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = DoneWithdraw(ctx, withdraw01.OrderID, true, xmap.M{"A": 123})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = DoneWithdraw(ctx, withdraw01.OrderID, true, xmap.M{"A": 123})
	if err != nil {
		t.Error(err)
		return
	}

	_, err = DoneWithdraw(ctx, withdraw00.OrderID, true, xmap.M{"A": 123})
	if err == nil {
		t.Error(err)
		return
	}
	_, err = DoneWithdraw(ctx, withdraw01.OrderID, false, xmap.M{"A": 123})
	if err == nil {
		t.Error(err)
		return
	}
	//
	searcher := WithdrawUnifySearcher{}
	searcher.Where.Type = WithdrawTypeArray{WithdrawTypeWithdraw}
	searcher.Where.Asset = []string{asset}
	searcher.Where.Status = WithdrawStatusAll
	err = searcher.Apply(ctx)
	if err != nil || len(searcher.Query.Withdraws) < 1 {
		t.Error(err)
		return
	}
	withdraw = &Withdraw{
		UserID:   user.TID,
		Asset:    "x1",
		Quantity: decimal.NewFromFloat(1),
	}
	err = CreateWithdraw(ctx, withdraw)
	if err != nil || withdraw.Status != WithdrawStatusConfirmed {
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
	withdraw2 := &Withdraw{
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(1000),
	}
	err = CreateWithdraw(ctx, withdraw2)
	if err != nil {
		t.Error(err)
		return
	}
	withdraw3 := &Withdraw{
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(1),
	}
	err = CreateWithdraw(ctx, withdraw3)
	if err != nil {
		t.Error(err)
		return
	}
	//cancel type is not correct
	withdraw4 := &Withdraw{
		OrderID:  NewOrderID(),
		Type:     WithdrawTypeTopup,
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(1),
	}
	err = AddWithdraw(ctx, withdraw4)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = CancelWithdraw(ctx, user.TID, withdraw4.OrderID)
	if err == nil {
		t.Error(err)
		return
	}
	//
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1, 2).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		withdraw = &Withdraw{
			UserID:   user.TID,
			Asset:    asset,
			Quantity: decimal.NewFromFloat(10),
		}
		err = CreateWithdraw(ctx, withdraw)
		return
	})
	pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = CancelWithdraw(ctx, user.TID, withdraw2.OrderID)
		return
	})
	pgx.MockerSetCall("Rows.Scan", 2, "Tx.Exec", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = CancelWithdraw(ctx, user.TID, withdraw2.OrderID)
		return
	})
	pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1, 2, 3, "Tx.Exec", 1, 2, 3).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = DoneWithdraw(ctx, withdraw3.OrderID, true, xmap.M{"A": 123})
		return
	})
}

func TestProcWithdraw(t *testing.T) {
	clear()
	func() {
		defer func() {
			recover()
		}()
		ApplyWithdraw(nil)
	}()
	ApplyWithdraw = func(withdraw *Withdraw) (err error) { return nil }
	asset := "TEST"
	user := testAddUser("TestProcWithdraw")
	basedb.StoreConf(ctx, ConfigWithdrawReview, converter.JSON(xmap.M{asset: 100}))
	added, err := TouchBalance(ctx, BalanceAreaSpot, []string{asset, "x1"}, user.TID)
	if err != nil || added != 2 {
		t.Error(err)
		return
	}
	balance := &Balance{
		UserID: user.TID,
		Area:   BalanceAreaSpot,
		Asset:  asset,
		Free:   decimal.NewFromFloat(10000),
	}
	err = IncreaseBalance(ctx, balance)
	if err != nil || !balance.Free.Equal(decimal.NewFromFloat(10000)) {
		t.Error(err)
		return
	}
	withdraw := &Withdraw{
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(1),
	}
	err = CreateWithdraw(ctx, withdraw)
	if err != nil || withdraw.Status != WithdrawStatusConfirmed {
		t.Error(err)
		return
	}
	err = ProcWithdraw()
	if err != nil {
		t.Error(err)
		return
	}
	err = ProcWithdraw()
	if err != pgx.ErrNoRows {
		t.Error(err)
		return
	}

	ApplyWithdraw = func(withdraw *Withdraw) (err error) { return fmt.Errorf("test") }
	withdraw = &Withdraw{
		UserID:   user.TID,
		Asset:    asset,
		Quantity: decimal.NewFromFloat(1),
	}
	err = CreateWithdraw(ctx, withdraw)
	if err != nil || withdraw.Status != WithdrawStatusConfirmed {
		t.Error(err)
		return
	}
	err = ProcWithdraw()
	if err == nil {
		t.Error(err)
		return
	}

	//
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Pool.Query", 1).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = ProcWithdraw()
		return
	})
}

func TestLoadWallet(t *testing.T) {
	clear()
	func() {
		defer func() {
			recover()
		}()
		AssignWallet(WalletMethodTron)
	}()
	AssignWallet = func(method WalletMethod) (address string, err error) { return "xxx", nil }
	user := testAddUser("TestLoadWallet")
	wallet1, err := LoadWalletByMethod(ctx, user.TID, WalletMethodTron)
	if err != nil || wallet1.TID < 1 {
		t.Error(err)
		return
	}
	wallet2, err := LoadWalletByMethod(ctx, user.TID, WalletMethodTron)
	if err != nil || wallet2.TID != wallet1.TID {
		t.Error(err)
		return
	}

	//
	pgx.MockerStart()
	defer pgx.MockerStop()
	user = testAddUser("TestLoadWallet-1")
	pgx.MockerClear()

	pgx.MockerSetCall("Pool.Begin", 1, "Rows.Scan", 1, 2).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = LoadWalletByMethod(ctx, user.TID, WalletMethodTron)
		return
	})

	AssignWallet = func(method WalletMethod) (address string, err error) { return "xxx", fmt.Errorf("error") }
	_, err = LoadWalletByMethod(ctx, user.TID, WalletMethodTron)
	if err == nil {
		t.Error(err)
		return
	}
}

func TestTopup(t *testing.T) {
	user := testAddUser("TestTopup")
	AssignWallet = func(method WalletMethod) (address string, err error) { return uuid.New(), nil }
	wallet, err := LoadWalletByMethod(ctx, user.TID, WalletMethodTron)
	if err != nil || wallet.TID < 1 {
		t.Error(err)
		return
	}
	txid := uuid.New()
	_, _, err = ReceiveTopup(ctx, wallet.Method, wallet.Address, txid, "TEST", decimal.NewFromFloat(100), xmap.M{})
	if err != nil {
		t.Error(err)
		return
	}
	_, _, err = ReceiveTopup(ctx, wallet.Method, wallet.Address, txid, "TEST", decimal.NewFromFloat(100), xmap.M{})
	if err != nil {
		t.Error(err)
		return
	}
	_, skip, err := ReceiveTopup(ctx, wallet.Method, uuid.New(), txid, "TEST", decimal.NewFromFloat(100), xmap.M{})
	if err != nil || !skip {
		t.Error(err)
		return
	}

	//
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	txid = uuid.New()

	pgx.MockerSetCall("Pool.Begin", 1, "Tx.Exec", 1, 2, 3, "Rows.Scan", 1, 2, 3, 4).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, _, err = ReceiveTopup(ctx, wallet.Method, wallet.Address, txid, "TEST", decimal.NewFromFloat(100), xmap.M{})
		return
	})
}
