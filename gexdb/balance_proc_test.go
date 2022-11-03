package gexdb

// import (
// 	"testing"
// 	"time"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/codingeasygo/util/xtime"
// 	"github.com/shopspring/decimal"
// )

// func TestBalanceProc(t *testing.T) {
// 	user := testAddUser("TestBalanceProc")
// 	TouchBalance(ctx, BalanceAssetAll, user.TID)
// 	err := IncreaseBalanceCall(Pool(), ctx, &Balance{UserID: user.TID, Asset: BalanceAssetYWE, Free: decimal.NewFromFloat(1)})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	err = IncreaseBalanceCall(Pool(), ctx, &Balance{UserID: user.TID, Asset: BalanceAssetMMK, Free: decimal.NewFromFloat(50)})
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	//
// 	syncer := NewBalanceRecordSyncer(time.Now().Hour(), func(asset string) decimal.Decimal { return decimal.NewFromFloat(100) })
// 	err = syncer.Proc()
// 	if err != pgx.ErrNoRows {
// 		t.Error(err)
// 		return
// 	}
// 	histories, err := ListUserBalanceHistory(ctx, user.TID, BalanceAssetMMK, xtime.TimeStartOfToday(), time.Now())
// 	if err != nil || len(histories) != 1 {
// 		t.Error(err)
// 		return
// 	}
// 	current, _ := histories[0].Current.Float64()
// 	if current != 150 {
// 		t.Error(err)
// 		return
// 	}

// 	//
// 	//test error
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	pgx.MockerClear()

// 	//proc error
// 	pgx.MockerSet("Pool.Exec", 1)
// 	err = syncer.Proc()
// 	if err == pgx.ErrNoRows {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//list history error
// 	pgx.MockerSet("Pool.Query", 1)
// 	_, err = ListUserBalanceHistory(ctx, user.TID, BalanceAssetMMK, xtime.TimeStartOfToday(), time.Now())
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Rows.Scan", 1)
// 	_, err = ListUserBalanceHistory(ctx, user.TID, BalanceAssetMMK, xtime.TimeStartOfToday(), time.Now())
// 	if err == nil {
// 		t.Error(err)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//not hour
// 	syncer = NewBalanceRecordSyncer(1, func(asset string) decimal.Decimal { return decimal.NewFromFloat(100) })
// 	syncer.Proc()

// 	//not price
// 	syncer = NewBalanceRecordSyncer(time.Now().Hour(), func(asset string) decimal.Decimal { return decimal.NewFromFloat(0) })
// 	syncer.Proc()
// }
