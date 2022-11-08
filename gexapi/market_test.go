package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
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

func TestFavoritesSymbol(t *testing.T) {
	user := testAddUser(gexdb.UserRoleNormal, "TestFavoritesSymbol")
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *user.Account, "123")
	//symbol
	ts.Should(t, "code", define.Success, "symbols", xmap.ShouldIsNil, "days", xmap.ShouldIsNil).GetMap("/usr/listFavoritesSymbol")
	//add
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/addFavoritesSymbol?symbol=%v", "")
	ts.Should(t, "code", define.Success).GetMap("/usr/addFavoritesSymbol?symbol=%v", "spot.YWEUSDT")
	ts.Should(t, "code", define.Success).GetMap("/usr/addFavoritesSymbol?symbol=%v", "spot.YWEUSDT")
	ts.Should(t, "code", define.Success).GetMap("/usr/addFavoritesSymbol?symbol=%v", "futures.YWEUSDT")
	//list
	listFavoritesSymbol, _ := ts.Should(t, "code", define.Success, "symbols/0/symbol", "spot.YWEUSDT").GetMap("/usr/listFavoritesSymbol")
	fmt.Printf("listFavoritesSymbol--->%v\n", converter.JSON(listFavoritesSymbol))
	//switch
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/switchFavoritesSymbol?symbol=%v", "")
	ts.Should(t, "code", define.Success).GetMap("/usr/switchFavoritesSymbol?symbol=%v", "futures.YWEUSDT")
	ts.Should(t, "code", define.Success, "symbols/0/symbol", "futures.YWEUSDT").GetMap("/usr/listFavoritesSymbol")
	ts.Should(t, "code", define.Success).GetMap("/usr/switchFavoritesSymbol?symbol=%v&to=%v", "futures.YWEUSDT", "spot.YWEUSDT")
	ts.Should(t, "code", define.Success, "symbols/0/symbol", "spot.YWEUSDT").GetMap("/usr/listFavoritesSymbol")
	//remove
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/removeFavoritesSymbol?symbol=%v", "")
	ts.Should(t, "code", define.Success).GetMap("/usr/removeFavoritesSymbol?symbol=%v", "futures.YWEUSDT")
	ts.Should(t, "code", define.Success, "symbols/0/symbol", "spot.YWEUSDT").GetMap("/usr/listFavoritesSymbol")
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()

	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/listFavoritesSymbol")
	pgx.MockerSetCall("Pool.Begin", 1).Should(t, "code", define.ServerError).GetMap("/usr/addFavoritesSymbol?symbol=%v", "spot.YWEUSDT")
	pgx.MockerSetCall("Pool.Begin", 1).Should(t, "code", define.ServerError).GetMap("/usr/switchFavoritesSymbol?symbol=%v", "spot.YWEUSDT")
	pgx.MockerSetCall("Pool.Begin", 1).Should(t, "code", define.ServerError).GetMap("/usr/removeFavoritesSymbol?symbol=%v", "spot.YWEUSDT")
}
