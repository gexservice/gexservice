package gexdb

import (
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/shopspring/decimal"
)

type BalanceRecordSyncer struct {
	LatestPrice func(asset string) decimal.Decimal
	TriggerHour int
	Last        time.Time
	Assets      []string
}

func NewBalanceRecordSyncer(triggerHour int, latestPrice func(asset string) decimal.Decimal) (syncer *BalanceRecordSyncer) {
	syncer = &BalanceRecordSyncer{
		TriggerHour: triggerHour,
		LatestPrice: latestPrice,
	}
	return
}

func (b *BalanceRecordSyncer) Proc() (err error) {
	// currentHour := time.Now().Hour()
	// if currentHour != b.TriggerHour || time.Since(b.Last) < 6*time.Hour {
	// 	err = pgx.ErrNoRows
	// 	return
	// }
	// price := b.LatestPrice("xx")
	// if price.Sign() <= 0 {
	// 	xlog.Warnf("BalanceRecordSyncer sync fail with not latest price")
	// 	err = pgx.ErrNoRows
	// 	return
	// }
	// syncTime := xsql.TimeStartOfToday()
	// syncSQL := `
	// 	insert into exs_balance_history(user_id,asset,current,update_time,create_time,status)
	// 	select m.user_id,m.asset,m.free+m.locked+(y.free+y.locked)*$1,$2,$2,$3 from exs_balance m join exs_balance y on m.asset=$4 and y.asset=$5 and m.user_id=y.user_id
	// 	on conflict(user_id,asset,create_time) do update set current=excluded.current
	// `
	// syncArg := []interface{}{price, syncTime, BalanceHistoryStatusNormal, BalanceAssetMMK, BalanceAssetYWE}
	// _, affected, err := Pool().Exec(context.Background(), syncSQL, syncArg...)
	// if err != nil {
	// 	xlog.Errorf("BalanceRecordSyncer sync fail with %v", err)
	// 	return
	// }
	// xlog.Infof("BalanceRecordSyncer sync %v balance record to history", affected)
	err = pgx.ErrNoRows
	return
}
