package gexdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

func TestBalance(t *testing.T) {
	asset := "TEST"
	user := testAddUser("TestBalance")
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
		Locked: decimal.NewFromFloat(100),
		Margin: decimal.NewFromFloat(100),
	}
	err = IncreaseBalance(ctx, balance)
	if err != nil || !balance.Free.Equal(decimal.NewFromFloat(100)) || !balance.Locked.Equal(decimal.NewFromFloat(100)) || !balance.Margin.Equal(decimal.NewFromFloat(100)) {
		t.Error(err)
		return
	}
	findBalance, err := FindBalanceByAsset(ctx, balance.UserID, balance.Area, balance.Asset)
	if err != nil || findBalance.TID != balance.TID {
		t.Error(err)
		return
	}
	balances, _, err := ListUserBalance(ctx, balance.UserID, BalanceAreaSpot, []string{asset}, BalanceStatusAll)
	if err != nil || len(balances) != 1 {
		t.Error(err)
		return
	}
	searcher := BalanceUnifySearcher{}
	searcher.Where.Area = BalanceAreaAll
	searcher.Where.Asset = []string{balance.Asset}
	searcher.Where.Key = "Test"
	err = searcher.Apply(ctx)
	if err != nil || len(searcher.Query.Balances) < 1 || searcher.Count.Total < 1 {
		t.Error(err)
		return
	}
	areaBalances, _, err := ListAreaBalance(ctx, balance.UserID, BalanceAreaArray{BalanceAreaSpot}, asset, BalanceStatusAll)
	if err != nil || len(areaBalances) != 1 {
		t.Error(err)
		return
	}
	if !IsErrBalanceNotEnought(IncreaseBalance(ctx,
		&Balance{
			UserID: user.TID,
			Area:   BalanceAreaSpot,
			Asset:  asset,
			Free:   decimal.NewFromFloat(-200),
		})) {
		t.Error("err")
		return
	}
	if !IsErrBalanceNotEnought(IncreaseBalance(ctx,
		&Balance{
			UserID: user.TID,
			Area:   BalanceAreaSpot,
			Asset:  asset,
			Locked: decimal.NewFromFloat(-200),
		})) {
		t.Error("err")
		return
	}
	if !IsErrBalanceNotEnought(IncreaseBalance(ctx,
		&Balance{
			UserID: user.TID,
			Area:   BalanceAreaSpot,
			Asset:  asset,
			Margin: decimal.NewFromFloat(-200),
		})) {
		t.Error("err")
		return
	}
	// balance = &Balance{UserID: user.TID, Asset: BalanceAssetYWE, Status: BalanceStatusNormal, Locked: decimal.NewFromFloat(-11100)}
	// err = IncreaseBalanceCall(Pool(), ctx, balance)
	// if err != ErrBalanceNotEnought {
	// 	t.Error(err)
	// 	return
	// }

	err = IncreaseBalance(ctx, &Balance{UserID: 10})
	if !IsErrBalanceNotFound(err) {
		t.Error(err)
		return
	}
	balance.Status = BalanceStatusLocked
	err = UpdateBalanceFilter(ctx, balance, "status")
	if err != nil {
		t.Error(err)
		return
	}
	err = IncreaseBalance(ctx, &Balance{
		UserID: user.TID,
		Area:   BalanceAreaSpot,
		Asset:  asset,
		Free:   decimal.NewFromFloat(100),
		Locked: decimal.NewFromFloat(100),
		Margin: decimal.NewFromFloat(100),
	})
	if err == nil {
		t.Error(err)
		return
	}

	balanceAll, err := CountAllBalance(ctx, asset)
	if err != nil || len(balanceAll) < 1 {
		t.Error(err)
		return
	}

	balanceArea, err := CountAreaBalance(ctx, BalanceAreaSpot, time.Time{}, time.Now())
	if err != nil || len(balanceArea) < 1 {
		t.Error(err)
		return
	}

	balanceUser, err := CountUserBalance(ctx, asset, user.TID)
	if err != nil || len(balanceUser) < 1 {
		t.Error(err)
		return
	}

}

func TestTransferChange(t *testing.T) {
	user := testAddUser("TestTransferChange")
	//inc
	_, err := ChangeBalance(ctx, 100, user.TID, BalanceAreaFunds, "test", decimal.NewFromFloat(100))
	if err != nil {
		t.Error(err)
		return
	}
	//from->to
	err = TransferChange(ctx, 100, user.TID, BalanceAreaFunds, BalanceAreaSpot, "test", decimal.NewFromFloat(100))
	if err != nil {
		t.Error(err)
		return
	}
	//transfer error
	err = TransferChange(ctx, 100, user.TID, BalanceAreaFunds, BalanceAreaSpot, "test", decimal.NewFromFloat(100))
	if !IsErrBalanceNotEnought(err) {
		t.Error(err)
		return
	}
	//to->from
	err = TransferChange(ctx, 100, user.TID, BalanceAreaFunds, BalanceAreaSpot, "test", decimal.NewFromFloat(-100))
	if err != nil {
		t.Error(err)
		return
	}
	//transfer error
	err = TransferChange(ctx, 100, user.TID, BalanceAreaFunds, BalanceAreaSpot, "test", decimal.NewFromFloat(-100))
	if !IsErrBalanceNotEnought(err) {
		t.Error(err)
		return
	}
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerSetCall("Pool.Begin", 1, "Tx.Exec", 1, 2).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		err = TransferChange(ctx, 100, user.TID, BalanceAreaFunds, BalanceAreaSpot, "test", decimal.NewFromFloat(100))
		return
	})
}

func TestChangeBalance(t *testing.T) {
	user := testAddUser("TestChangeBalance")
	//inc
	_, err := ChangeBalance(ctx, 100, user.TID, BalanceAreaSpot, "test", decimal.NewFromFloat(100))
	if err != nil {
		t.Error(err)
		return
	}
	//dec
	_, err = ChangeBalance(ctx, 100, user.TID, BalanceAreaSpot, "test", decimal.NewFromFloat(-100))
	if err != nil {
		t.Error(err)
		return
	}
	//dec error
	_, err = ChangeBalance(ctx, 100, user.TID, BalanceAreaSpot, "test", decimal.NewFromFloat(-100))
	if !IsErrBalanceNotEnought(err) {
		t.Error(err)
		return
	}

	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerSetCall("Pool.Begin", 1, "Tx.Exec", 1, 2).ShouldError(t).Call(func(trigger int) (res xmap.M, err error) {
		_, err = ChangeBalance(ctx, 100, user.TID, BalanceAreaSpot, "test", decimal.NewFromFloat(100))
		return
	})
}

func TestBalanceRecord(t *testing.T) {
	asset := "TEST"
	user := testAddUser("TestBalanceRecord")
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
		Locked: decimal.NewFromFloat(100),
		Margin: decimal.NewFromFloat(100),
	}
	err = IncreaseBalance(ctx, balance)
	if err != nil || !balance.Free.Equal(decimal.NewFromFloat(100)) || !balance.Locked.Equal(decimal.NewFromFloat(100)) || !balance.Margin.Equal(decimal.NewFromFloat(100)) {
		t.Error(err)
		return
	}

	_, err = AddBalanceRecordCall(Pool(), ctx, &BalanceRecord{
		Creator:   user.TID,
		BalanceID: balance.TID,
		Type:      BalanceRecordTypeTradeFee,
		Changed:   decimal.NewFromFloat(1),
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = AddBalanceRecordCall(Pool(), ctx)
	if err != nil {
		t.Error(err)
		return
	}
	searcher := BalanceRecordUnifySearcher{}
	searcher.Where.Area = BalanceAreaSpot
	searcher.Where.Asset = []string{balance.Asset}
	searcher.Where.StartTime = xsql.Time(time.Now().Add(-time.Hour))
	searcher.Where.EndTime = xsql.Time(time.Now())
	searcher.Where.Key = "Test"
	err = searcher.Apply(ctx)
	if err != nil || len(searcher.Query.Records) < 1 || searcher.Count.Total < 1 {
		fmt.Println("-->", searcher.Query.Records)
		fmt.Println("-->", searcher.Count.Total)
		t.Error(err)
		return
	}
}
