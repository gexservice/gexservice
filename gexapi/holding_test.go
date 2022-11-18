package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
)

func TestHolding(t *testing.T) {
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/loadHolding?symbol=")
	loadHolding, _ := ts.Should(t, "code", define.Success, "holding/lever", 5).GetMap("/usr/loadHolding?symbol=XX")
	fmt.Printf("loadHolding--->%v\n", converter.JSON(loadHolding))
	//
	listHolding, _ := ts.Should(t, "code", define.Success, "balance/free", xmap.ShouldIsNoZero).GetMap("/usr/listHolding")
	fmt.Printf("listHolding--->%v\n", converter.JSON(listHolding))
	//change lever
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/changeHoldingLever?symbol=%v&lever=0", "futures.YWEUSDT")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/changeHoldingLever?symbol=%v&lever=100", "futures.YWEUSDT")
	ts.Should(t, "code", gexdb.CodeOrderPending).GetMap("/usr/changeHoldingLever?symbol=%v&lever=1", "futures.YWEUSDT")
	//
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc3", "123")
	ts.Should(t, "code", define.Success).GetMap("/usr/changeHoldingLever?symbol=%v&lever=1", "futures.YWEUSDT")
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	pgx.MockerClear()

	pgx.MockerSetCall("Pool.Exec", 1).Should(t, "code", define.ServerError).GetMap("/usr/loadHolding?symbol=XX")
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/loadHolding?symbol=XX")
	pgx.MockerSetCall("Rows.Scan", 1, 2, "Pool.Exec", 1).Should(t, "code", define.ServerError).GetMap("/usr/listHolding")
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/changeHoldingLever?symbol=%v&lever=1", "futures.YWEUSDT")
}
