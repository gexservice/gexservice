package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/uuid"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
)

func TestWithdraw(t *testing.T) {
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	//
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap(xmap.M{}, "/usr/createWithdraw")
	createWithdraw, _ := ts.Should(t, "code", define.Success).GetMap("/usr/createWithdraw?method=tron&asset=%v&quantity=%v&receiver=test", spotBalanceQuote, "1")
	fmt.Printf("createWithdraw--->%v\n", converter.JSON(createWithdraw))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/listWithdraw?asset=%v&type=xx", spotBalanceQuote)
	listWithdraw, _ := ts.Should(t, "code", define.Success).GetMap("/usr/listWithdraw?asset=%v", spotBalanceQuote)
	fmt.Printf("listWithdraw--->%v\n", converter.JSON(listWithdraw))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/cancelWithdraw?order_id=%v", "")
	cancelWithdraw, _ := ts.Should(t, "code", define.Success).GetMap("/usr/cancelWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	fmt.Printf("cancelWithdraw--->%v\n", converter.JSON(cancelWithdraw))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/confirmWithdraw?order_id=%v", "")
	createWithdraw, _ = ts.Should(t, "code", define.Success).GetMap("/usr/createWithdraw?method=tron&asset=%v&quantity=%v&receiver=test", spotBalanceQuote, "1")
	ts.Should(t, "code", define.NotAccess).GetMap("/usr/confirmWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	ts.Should(t, "code", define.Success).GetMap("/usr/confirmWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 2).Should(t, "code", define.ServerError).GetMap("/usr/listWithdraw?asset=%v", spotBalanceQuote)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/createWithdraw?method=tron&asset=%v&quantity=%v&receiver=test", spotBalanceQuote, "1")
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/cancelWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
	pgx.MockerSetCall("Pool.Exec", 1).Should(t, "code", define.ServerError).GetMap("/usr/confirmWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
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
