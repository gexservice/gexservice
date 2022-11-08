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

func TestBalance(t *testing.T) {
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
	//
	loadBalanceOverview, _ := ts.Should(t, "code", define.Success, "total_value", xmap.ShouldIsNoZero).GetMap("/usr/loadBalanceOverview")
	fmt.Printf("loadBalanceOverview--->%v\n", converter.JSON(loadBalanceOverview))
	listBalance, _ := ts.Should(t, "code", define.Success, "total_value", xmap.ShouldIsNoZero).GetMap("/usr/listBalance?area=%d", gexdb.BalanceAreaSpot)
	fmt.Printf("listBalance--->%v\n", converter.JSON(listBalance))
	//
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/listBalanceRecord?type=xxx")
	listBalanceRecord, _ := ts.Should(t, "code", define.Success).GetMap("/usr/listBalanceRecord")
	fmt.Printf("listBalanceRecord--->%v\n", converter.JSON(listBalanceRecord))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/loadBalanceOverview")
	pgx.Should(t, "code", define.ArgsInvalid).GetMap("/usr/listBalance?area=%d", 1)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/listBalance?area=%d", gexdb.BalanceAreaSpot)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/listBalanceRecord")
}
