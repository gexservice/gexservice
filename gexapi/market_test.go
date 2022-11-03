package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/define"
)

func TestMarket(t *testing.T) {
	symbol := "spot.YWEUSDT"

	//symbol
	listSymbol, _ := ts.Should(t, "code", define.Success, "symbols", xmap.ShouldIsNoEmpty, "days", xmap.ShouldIsNoEmpty).GetMap("/pub/listSymbol")
	fmt.Printf("listSymbol--->%v\n", converter.JSON(listSymbol))
	loadSymbol, _ := ts.Should(t, "code", define.Success, "symbol", xmap.ShouldIsNoNil, "day", xmap.ShouldIsNoNil).GetMap("/pub/loadSymbol?symbol=%v", symbol)
	fmt.Printf("loadSymbol--->%v\n", converter.JSON(loadSymbol))

	//kline
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/pub/listKLine?symbol=%v&interval=5min&start_time=xx&end_time=%v", symbol, xsql.TimeNow().Timestamp())
	listKLine, _ := ts.Should(t, "code", define.Success, "lines", xmap.ShouldIsNoEmpty).GetMap("/pub/listKLine?symbol=%v&interval=5min&start_time=100&end_time=%v", symbol, xsql.TimeNow().Timestamp())
	fmt.Printf("listKLine--->%v\n", converter.JSON(listKLine))

	ts.Should(t, "code", define.ArgsInvalid).GetMap("/pub/loadDepth?max=%v", -10)
	loadDepthRes, _ := ts.Should(t, "code", define.Success, "/depth/bids", xmap.ShouldIsNoEmpty).GetMap("/pub/loadDepth?symbol=%v&max=%v", symbol, 10)
	fmt.Printf("loadDepthRes--->%v\n", converter.JSON(loadDepthRes))

	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()

	pgx.MockerSetCall("Pool.Query", 1).Should(t, "code", define.ServerError).GetMap("/pub/listKLine?symbol=%v&interval=5min&start_time=100&end_time=%v", symbol, xsql.TimeNow().Timestamp())

}
