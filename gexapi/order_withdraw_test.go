package gexapi

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/codingeasygo/util/converter"
// 	"github.com/codingeasygo/util/uuid"
// 	"github.com/codingeasygo/util/xmap"
// 	"github.com/shopspring/decimal"
// 	"github.com/gexservice/gexservice/base/define"
// 	"github.com/gexservice/gexservice/gexdb"
// )

// func TestWithdrawOrder(t *testing.T) {
// 	userMMK := testAddUser(gexdb.UserRoleNormal, "TestWithdrawOrder-MMK")
// 	_, err := gexdb.TouchBalance(ctx, gexdb.BalanceAssetAll, userMMK.TID)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
// 		UserID: userMMK.TID,
// 		Asset:  gexdb.BalanceAssetMMK,
// 		Free:   decimal.NewFromFloat(10000),
// 		Status: gexdb.BalanceStatusNormal,
// 	})
// 	{ //user invalid
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userMMK.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createWithdrawOrder, err := ts.GetMap("/usr/createWithdrawOrder?quantity=1&password=123")
// 		if err != nil || createWithdrawOrder.Int64("code") == 0 {
// 			t.Errorf("err:%v,createWithdrawOrder:%v", err, createWithdrawOrder)
// 			return
// 		}
// 	}
// 	gexdb.Pool().ExecRow(ctx, `update exs_user set kbz_openid=$1 where tid=$2`, uuid.New(), userMMK.TID)
// 	{ //create user cancel
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userMMK.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createWithdrawOrder, err := ts.GetMap("/usr/createWithdrawOrder?quantity=1&password=123")
// 		if err != nil || createWithdrawOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,createWithdrawOrder:%v", err, createWithdrawOrder)
// 			return
// 		}
// 		fmt.Printf("createWithdrawOrder--->%v\n", converter.JSON(createWithdrawOrder))
// 		orderID := createWithdrawOrder.StrDef("", "/order/order_id")

// 		cancelWithdrawOrder, err := ts.GetMap("/usr/cancelWithdrawOrder?order_id=%v", orderID)
// 		if err != nil || cancelWithdrawOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,cancelWithdrawOrder:%v", err, cancelWithdrawOrder)
// 			return
// 		}
// 		fmt.Printf("cancelWithdrawOrder--->%v\n", converter.JSON(cancelWithdrawOrder))
// 	}
// 	{ //create admin cancel
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userMMK.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createWithdrawOrder, err := ts.GetMap("/usr/createWithdrawOrder?quantity=1&password=123")
// 		if err != nil || createWithdrawOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,createWithdrawOrder:%v", err, createWithdrawOrder)
// 			return
// 		}
// 		fmt.Printf("createWithdrawOrder--->%v\n", converter.JSON(createWithdrawOrder))
// 		orderID := createWithdrawOrder.StrDef("", "/order/order_id")

// 		clearCookie()
// 		login, err = ts.GetMap("/pub/login?username=%v&password=%v", "admin", "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		cancelWithdrawOrder, err := ts.GetMap("/usr/cancelWithdrawOrder?order_id=%v", orderID)
// 		if err != nil || cancelWithdrawOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,cancelWithdrawOrder:%v", err, cancelWithdrawOrder)
// 			return
// 		}
// 		fmt.Printf("cancelWithdrawOrder--->%v\n", converter.JSON(cancelWithdrawOrder))
// 	}
// 	{ //normal error
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userMMK.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createWithdrawOrder, err := ts.GetMap("/usr/createWithdrawOrder?quantity=1&password=1234")
// 		if err != nil || createWithdrawOrder.Int64("code") != 8000 {
// 			t.Errorf("err:%v,createWithdrawOrder:%v", err, createWithdrawOrder)
// 			return
// 		}
// 		createWithdrawOrder, err = ts.GetMap("/usr/createWithdrawOrder?quantity=100000&password=123")
// 		if err != nil || createWithdrawOrder.Int64("code") != 9000 {
// 			t.Errorf("err:%v,createWithdrawOrder:%v", err, createWithdrawOrder)
// 			return
// 		}
// 	}
// 	//
// 	//test error
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	pgx.MockerClear()
// 	var res xmap.M

// 	//to user
// 	clearCookie()
// 	login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userMMK.Account, "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//create error
// 	res, err = ts.GetMap("/usr/createWithdrawOrder?quantity=x")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Row.Scan", 2)
// 	res, err = ts.GetMap("/usr/createWithdrawOrder?quantity=1&password=123")
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	res, err = ts.GetMap("/usr/createWithdrawOrder?quantity=1&password=123")
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//create test data
// 	createWithdrawOrder, err := ts.GetMap("/usr/createWithdrawOrder?quantity=1&password=123")
// 	if err != nil || createWithdrawOrder.Int64("code") != 0 {
// 		t.Errorf("err:%v,createWithdrawOrder:%v", err, createWithdrawOrder)
// 		return
// 	}
// 	orderID := createWithdrawOrder.StrDef("", "/order/order_id")
// 	pgx.MockerClear()

// 	//cancel error
// 	res, err = ts.GetMap("/usr/cancelWithdrawOrder?order_id=%v", "")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Row.Scan", 1)
// 	res, err = ts.GetMap("/usr/cancelWithdrawOrder?order_id=%v", orderID)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	res, err = ts.GetMap("/usr/cancelWithdrawOrder?order_id=%v", orderID)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// }
