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

func TestBalance(t *testing.T) {
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	//
	loadBalanceOverview, _ := ts.Should(t, "code", define.Success, "total_value", xmap.ShouldIsNoZero).GetMap("/usr/loadBalanceOverview")
	fmt.Printf("loadBalanceOverview--->%v\n", converter.JSON(loadBalanceOverview))
	listBalance, _ := ts.Should(t, "code", define.Success, "total_value", xmap.ShouldIsNoZero).GetMap("/usr/listBalance?area=%d", gexdb.BalanceAreaSpot)
	fmt.Printf("listBalance--->%v\n", converter.JSON(listBalance))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/searchBalance?area=xxx")
	searchBalance, _ := ts.Should(t, "code", define.Success).GetMap("/usr/searchBalance")
	fmt.Printf("searchBalance--->%v\n", converter.JSON(searchBalance))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/searchBalanceRecord?type=xxx")
	searchBalanceRecord, _ := ts.Should(t, "code", define.Success).GetMap("/usr/searchBalanceRecord")
	fmt.Printf("searchBalanceRecord--->%v\n", converter.JSON(searchBalanceRecord))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/loadBalanceOverview")
	pgx.Should(t, "code", define.ArgsInvalid).GetMap("/usr/listBalance?area=%d", 1)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/listBalance?area=%d", gexdb.BalanceAreaSpot)
	pgx.MockerSetCall("Pool.Query", 1, 2, 3, 4).Should(t, "code", define.ServerError).GetMap("/usr/searchBalance")
	pgx.MockerSetCall("Pool.Query", 1, 2).Should(t, "code", define.ServerError).GetMap("/usr/searchBalanceRecord")
}

func TestTransferBalance(t *testing.T) {
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/transferBalance?from=%d&to=%d&asset=%v&value=10", 1, 1, spotBalanceQuote)
	transferBalance, _ := ts.Should(t, "code", define.Success).GetMap("/usr/transferBalance?from=%d&to=%d&asset=%v&value=10", gexdb.BalanceAreaFunds, gexdb.BalanceAreaSpot, spotBalanceQuote)
	fmt.Printf("transferBalance->%v\n", converter.JSON(transferBalance))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Tx.Exec", 1).Should(t, "code", define.ServerError).GetMap("/usr/transferBalance?from=%d&to=%d&asset=%v&value=10", gexdb.BalanceAreaFunds, gexdb.BalanceAreaSpot, spotBalanceQuote)
}

func TestChangeUserBalance(t *testing.T) {
	user := testAddUser(gexdb.UserRoleNormal, "TestChangeBalance")
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/admin/changeUserBalance?user_id=%v&area=%d&asset=%v&changed=100", user.TID, 1, spotBalanceQuote)
	changeUserBalance, _ := ts.Should(t, "code", define.Success).GetMap("/admin/changeUserBalance?user_id=%v&area=%d&asset=%v&changed=100", user.TID, gexdb.BalanceAreaFunds, spotBalanceQuote)
	fmt.Printf("changeUserBalance->%v\n", converter.JSON(changeUserBalance))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Tx.Exec", 1).Should(t, "code", define.ServerError).GetMap("/admin/changeUserBalance?user_id=%v&area=%d&asset=%v&changed=100", user.TID, gexdb.BalanceAreaFunds, spotBalanceQuote)
}
