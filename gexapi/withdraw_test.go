package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/gexservice/gexservice/base/define"
)

func TestWithdraw(t *testing.T) {
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/createWithdraw?asset=%v&quantity=%v", spotBalanceQuote, "")
	createWithdraw, _ := ts.Should(t, "code", define.Success).GetMap("/usr/createWithdraw?asset=%v&quantity=%v", spotBalanceQuote, "1")
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
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/listWithdraw?asset=%v", spotBalanceQuote)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/createWithdraw?asset=%v&quantity=%v", spotBalanceQuote, "1")
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/cancelWithdraw?order_id=%v", createWithdraw.StrDef("", "/withdraw/order_id"))
}
