package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/uuid"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xtime"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
)

func TestWithdraw(t *testing.T) {
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	//
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap(xmap.M{}, "/usr/createWithdraw")
	createWithdraw, _ := ts.Should(t, "code", define.Success).GetMap("/usr/createWithdraw?method=tron&asset=%v&quantity=%v&receiver=test&trade_pass=123", spotBalanceQuote, "1")
	fmt.Printf("createWithdraw--->%v\n", converter.JSON(createWithdraw))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/searchWithdraw?asset=%v&type=xx", spotBalanceQuote)
	searchWithdraw, _ := ts.Should(t, "code", define.Success).GetMap("/usr/searchWithdraw?asset=%v", spotBalanceQuote)
	fmt.Printf("searchWithdraw--->%v\n", converter.JSON(searchWithdraw))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/cancelWithdraw?order_id=%v", "")
	cancelWithdraw, _ := ts.Should(t, "code", define.Success).GetMap("/usr/cancelWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	fmt.Printf("cancelWithdraw--->%v\n", converter.JSON(cancelWithdraw))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/confirmWithdraw?order_id=%v", "")
	createWithdraw, _ = ts.Should(t, "code", define.Success).GetMap("/usr/createWithdraw?method=tron&asset=%v&quantity=%v&receiver=test&trade_pass=123", spotBalanceQuote, "1")
	ts.Should(t, "code", define.NotAccess).GetMap("/usr/confirmWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	ts.Should(t, "code", define.Success).GetMap("/usr/confirmWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Pool.Query", 1, 2).Should(t, "code", define.ServerError).GetMap("/usr/searchWithdraw?asset=%v", spotBalanceQuote)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/cancelWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	pgx.MockerSetCall("Pool.Exec", 1).Should(t, "code", define.ServerError).GetMap("/usr/confirmWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))

	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", gexdb.CodeTradePasswordInvalid).GetMap("/usr/createWithdraw?method=tron&asset=%v&quantity=%v&receiver=test&trade_pass=123", spotBalanceQuote, "1")
	pgx.MockerSetCall("Rows.Scan", 2).Should(t, "code", define.ServerError).GetMap("/usr/createWithdraw?method=tron&asset=%v&quantity=%v&receiver=test&trade_pass=123", spotBalanceQuote, "1")
}

func TestGoldbar(t *testing.T) {
	basedb.StoreConf(ctx, gexdb.ConfigGoldbarFee, "0.001")
	basedb.StoreConf(ctx, gexdb.ConfigGoldbarRate, "10")
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	//
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap(xmap.M{}, "/usr/createGoldbar")
	createGoldbar, _ := ts.Should(t, "code", define.Success).GetMap("/usr/createGoldbar?pickup_amount=%v&pickup_time=%v&trade_pass=123&pickup_name=name&pickup_phone=phone&pickup_address=addr", "1", xtime.Now())
	fmt.Printf("createGoldbar--->%v\n", converter.JSON(createGoldbar))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/searchGoldbar?status=%v", "xx")
	searchGoldbar, _ := ts.Should(t, "code", define.Success).GetMap("/usr/searchGoldbar")
	fmt.Printf("searchGoldbar--->%v\n", converter.JSON(searchGoldbar))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/cancelGoldbar?order_id=%v", "")
	cancelGoldbar, _ := ts.Should(t, "code", define.Success).GetMap("/usr/cancelGoldbar?order_id=%v", createGoldbar.StrDef("", "/goldbar/order_id"))
	fmt.Printf("cancelGoldbar--->%v\n", converter.JSON(cancelGoldbar))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/confirmGoldbar?order_id=%v", "")
	createGoldbar, _ = ts.Should(t, "code", define.Success).GetMap("/usr/createGoldbar?pickup_amount=%v&pickup_time=%v&trade_pass=123&pickup_name=name&pickup_phone=phone&pickup_address=addr", "1", xtime.Now())
	ts.Should(t, "code", define.NotAccess).GetMap("/usr/confirmGoldbar?order_id=%v", createGoldbar.StrDef("", "/goldbar/order_id"))
	ts.Should(t, "code", define.NotAccess).PostJSONMap(xmap.M{
		"code":     createGoldbar.StrDef("", "/goldbar/receiver"),
		"order_id": createGoldbar.StrDef("", "/goldbar/order_id"),
	}, "/usr/doneGoldbar")
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	ts.Should(t, "code", define.Success).GetMap("/usr/confirmGoldbar?order_id=%v", createGoldbar.StrDef("", "/goldbar/order_id"))
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap("xx", "/usr/doneGoldbar")
	ts.Should(t, "code", define.Success).PostJSONMap(xmap.M{
		"code":     createGoldbar.StrDef("", "/goldbar/receiver"),
		"order_id": createGoldbar.StrDef("", "/goldbar/order_id"),
	}, "/usr/doneGoldbar")
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Pool.Query", 1, 2).Should(t, "code", define.ServerError).GetMap("/usr/searchGoldbar")
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/cancelGoldbar?order_id=%v", createGoldbar.StrDef("", "/goldbar/order_id"))
	pgx.MockerSetCall("Pool.Exec", 1).Should(t, "code", define.ServerError).GetMap("/usr/confirmGoldbar?order_id=%v", createGoldbar.StrDef("", "/goldbar/order_id"))
	pgx.MockerSetCall("Pool.Begin", 1).Should(t, "code", define.ServerError).PostJSONMap(xmap.M{
		"code":     createGoldbar.StrDef("", "/goldbar/receiver"),
		"order_id": createGoldbar.StrDef("", "/goldbar/order_id"),
	}, "/usr/doneGoldbar")

	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", gexdb.CodeTradePasswordInvalid).GetMap("/usr/createGoldbar?pickup_amount=%v&pickup_time=%v&trade_pass=123&pickup_name=name&pickup_phone=phone&pickup_address=addr", "1", xtime.Now())
	pgx.MockerSetCall("Rows.Scan", 2).Should(t, "code", define.ServerError).GetMap("/usr/createGoldbar?pickup_amount=%v&pickup_time=%v&trade_pass=123&pickup_name=name&pickup_phone=phone&pickup_address=addr", "1", xtime.Now())
}

func TestLoadTopupAddress(t *testing.T) {
	gexdb.AssignWallet = func(method gexdb.WalletMethod) (address string, err error) {
		address = uuid.New()
		return
	}
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/loadTopupAddress?method=%v", "xx")
	loadTopupAddress, _ := ts.Should(t, "code", define.Success).GetMap("/usr/loadTopupAddress?method=%v", gexdb.WalletMethodTron)
	fmt.Printf("loadTopupAddress--->%v\n", converter.JSON(loadTopupAddress))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/loadTopupAddress?method=%v", gexdb.WalletMethodTron)
}
