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

// func TestChangeBalance(t *testing.T) {
// 	user := testAddUser(gexdb.UserRoleNormal, "TestChangeBalance")
// 	gexdb.TouchBalance(ctx, gexdb.BalanceAssetAll, user.TID)
// 	clearCookie()
// 	login, err := ts.GetMap("/pub/login?username=%v&password=%v", "admin", "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	changeUserBalance, err := ts.GetMap("/usr/changeUserBalance?user_id=%v&asset=%v&changed=100", user.TID, gexdb.BalanceAssetYWE)
// 	if err != nil || changeUserBalance.Int64("code") != define.Success {
// 		t.Errorf("err:%v,changeUserBalance:%v", err, changeUserBalance)
// 		return
// 	}
// 	fmt.Printf("changeUserBalance->%v\n", converter.JSON(changeUserBalance))

// 	//
// 	//test error
// 	var res xmap.M
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	pgx.MockerClear()

// 	res, err = ts.GetMap("/usr/changeUserBalance?user_id=%v&asset=%v&changed=xx", user.TID, gexdb.BalanceAssetYWE)
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Row.Scan", 1)
// 	res, err = ts.GetMap("/usr/changeUserBalance?user_id=%v&asset=%v&changed=100", user.TID, gexdb.BalanceAssetYWE)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	res, err = ts.GetMap("/usr/changeUserBalance?user_id=%v&asset=%v&changed=100", user.TID, gexdb.BalanceAssetYWE)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//not access
// 	clearCookie()
// 	login, err = ts.GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	res, err = ts.GetMap("/usr/changeUserBalance?user_id=%v&asset=%v&changed=100", user.TID, gexdb.BalanceAssetYWE)
// 	if err != nil || res.Int64("code") != define.NotAccess {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// }

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
