package gexdb

import (
	"testing"
	"time"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xsql"
)

func TestKLine(t *testing.T) {
	IntervString(5 * time.Minute)
	IntervString(30 * time.Minute)
	IntervString(1 * time.Hour)
	IntervString(4 * time.Hour)
	IntervString(24 * time.Hour)
	IntervString(7 * 24 * time.Hour)
	IntervString(30 * 24 * time.Hour)
	IntervString(130 * 24 * time.Hour)
	StringInterv("5min")
	StringInterv("30min")
	StringInterv("1hour")
	StringInterv("4hour")
	StringInterv("1day")
	StringInterv("1week")
	StringInterv("1mon")
	StringInterv("xxxx")
	//
	line := &KLine{StartTime: xsql.TimeNow()}
	_, err := AddMultiKLine(ctx, line)
	if err != nil {
		t.Error(err)
		return
	}
	lines, err := ListKLine(ctx, line.Symbol, line.Interv, time.Now().Add(-time.Hour), time.Now())
	if err != nil || len(lines) == 0 {
		t.Error(err)
		return
	}

	_, err = AddMultiKLine(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()

	//list kline error
	pgx.MockerSet("Pool.Query", 1)
	_, err = ListKLine(ctx, line.Symbol, line.Interv, time.Now().Add(-time.Hour), time.Now())
	if err == nil {
		t.Error(err)
		return
	}
	pgx.MockerClear()
	pgx.MockerSet("Rows.Scan", 1)
	_, err = ListKLine(ctx, line.Symbol, line.Interv, time.Now().Add(-time.Hour), time.Now())
	if err == nil {
		t.Error(err)
		return
	}
	pgx.MockerClear()
}
