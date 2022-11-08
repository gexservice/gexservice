package gexdb

import (
	"fmt"
	"testing"
	"time"

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

	balanceAll, err := CountBalance(ctx, BalanceAreaSpot, time.Now().Add(-time.Hour), time.Now())
	if err != nil || len(balanceAll) < 1 {
		t.Error(err)
		return
	}

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

	// //
	// //test error
	// pgx.MockerStart()
	// defer pgx.MockerStop()
	// pgx.MockerClear()

	// //list balance error
	// pgx.MockerSet("Pool.Query", 1)
	// _, err = ListUserBalance(ctx, balance.UserID, BalanceAssetAll, BalanceStatusAll)
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()
	// pgx.MockerSet("Rows.Scan", 1)
	// _, err = ListUserBalance(ctx, balance.UserID, BalanceAssetAll, BalanceStatusAll)
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// //increase balance error
	// balance = &Balance{UserID: userNone.TID, Asset: BalanceAssetYWE, Status: BalanceStatusNormal, Free: decimal.NewFromFloat(100)}
	// err = IncreaseBalanceCall(Pool(), ctx, balance)
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// Pool().Exec(ctx, `update gex_balance set status=$1 where user_id=$2 and asset=$3`, BalanceStatusLocked, user.TID, BalanceAssetMMK)
	// balance = &Balance{UserID: user.TID, Asset: BalanceAssetMMK, Status: BalanceStatusNormal, Free: decimal.NewFromFloat(100)}
	// err = IncreaseBalanceCall(Pool(), ctx, balance)
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()

	// //count balance error
	// pgx.MockerSet("Pool.Query", 1)
	// _, err = CountBalance(ctx, time.Now().Add(-time.Hour), time.Now())
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()
	// pgx.MockerSet("Rows.Scan", 1)
	// _, err = CountBalance(ctx, time.Now().Add(-time.Hour), time.Now())
	// if err == nil {
	// 	t.Error(err)
	// 	return
	// }
	// pgx.MockerClear()
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
	err = searcher.Apply(ctx)
	if err != nil || len(searcher.Query.Records) < 1 || searcher.Count.Total < 1 {
		fmt.Println("-->", searcher.Query.Records)
		fmt.Println("-->", searcher.Count.Total)
		t.Error(err)
		return
	}
}

// func TestChangeBalance(t *testing.T) {
// 	user := testAddUser("TestChangeBalance")
// 	TouchBalance(ctx, BalanceAssetAll, user.TID)
// 	//dec error
// 	_, _, err := ChangeBalance(ctx, 100, user.TID, BalanceAssetYWE, decimal.NewFromFloat(-100))
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	//inc
// 	balance, order, err := ChangeBalance(ctx, 100, user.TID, BalanceAssetYWE, decimal.NewFromFloat(100))
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Printf("balance->%v,order->%v\n", balance.TID, order.TID)
// 	//dec
// 	balance, order, err = ChangeBalance(ctx, 100, user.TID, BalanceAssetYWE, decimal.NewFromFloat(-100))
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Printf("balance->%v,order->%v\n", balance.TID, order.TID)
// 	//dec error
// 	_, _, err = ChangeBalance(ctx, 100, user.TID, BalanceAssetYWE, decimal.NewFromFloat(-100))
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}

// 	//
// 	//test error
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	pgx.MockerClear()

// 	_, _, err = ChangeBalance(ctx, 100, user.TID, BalanceAssetMMK, decimal.NewFromFloat(100))
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	_, _, err = ChangeBalance(ctx, 100, user.TID, BalanceAssetYWE, decimal.NewFromFloat(100))
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Row.Scan", 1)
// 	_, _, err = ChangeBalance(ctx, 100, user.TID, BalanceAssetYWE, decimal.NewFromFloat(-100))
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// }
