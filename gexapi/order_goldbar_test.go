package gexapi

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/codingeasygo/util/converter"
// 	"github.com/codingeasygo/util/xmap"
// 	"github.com/shopspring/decimal"
// 	"github.com/gexservice/gexservice/base/define"
// 	"github.com/gexservice/gexservice/gexdb"
// )

// func TestGlodbarOrder(t *testing.T) {
// 	userYWE := testAddUser(gexdb.UserRoleNormal, "TestGlodbarOrder-YWE")
// 	_, err := gexdb.TouchBalance(ctx, gexdb.BalanceAssetAll, userYWE.TID)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	gexdb.IncreaseBalanceCall(gexdb.Pool(), ctx, &gexdb.Balance{
// 		UserID: userYWE.TID,
// 		Asset:  gexdb.BalanceAssetYWE,
// 		Free:   decimal.NewFromFloat(10000),
// 		Status: gexdb.BalanceStatusNormal,
// 	})
// 	{ //create user cancel
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userYWE.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createGoldbarOrder, err := ts.GetMap("/usr/createGoldbarOrder?quantity=1&password=123&city=city&address=address")
// 		if err != nil || createGoldbarOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,createGoldbarOrder:%v", err, createGoldbarOrder)
// 			return
// 		}
// 		fmt.Printf("createGoldbarOrder--->%v\n", converter.JSON(createGoldbarOrder))
// 		orderID := createGoldbarOrder.StrDef("", "/order/order_id")

// 		cancelGoldbarOrder, err := ts.GetMap("/usr/cancelGoldbarOrder?order_id=%v", orderID)
// 		if err != nil || cancelGoldbarOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,cancelGoldbarOrder:%v", err, cancelGoldbarOrder)
// 			return
// 		}
// 		fmt.Printf("cancelGoldbarOrder--->%v\n", converter.JSON(cancelGoldbarOrder))
// 	}
// 	{ //create admin cancel
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userYWE.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createGoldbarOrder, err := ts.GetMap("/usr/createGoldbarOrder?quantity=1&password=123&city=city&address=address")
// 		if err != nil || createGoldbarOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,createGoldbarOrder:%v", err, createGoldbarOrder)
// 			return
// 		}
// 		fmt.Printf("createGoldbarOrder--->%v\n", converter.JSON(createGoldbarOrder))
// 		orderID := createGoldbarOrder.StrDef("", "/order/order_id")

// 		clearCookie()
// 		login, err = ts.GetMap("/pub/login?username=%v&password=%v", "admin", "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		cancelGoldbarOrder, err := ts.GetMap("/usr/cancelGoldbarOrder?order_id=%v", orderID)
// 		if err != nil || cancelGoldbarOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,cancelGoldbarOrder:%v", err, cancelGoldbarOrder)
// 			return
// 		}
// 		fmt.Printf("cancelGoldbarOrder--->%v\n", converter.JSON(cancelGoldbarOrder))
// 	}
// 	{ //create admin verify
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userYWE.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createGoldbarOrder, err := ts.GetMap("/usr/createGoldbarOrder?quantity=1&password=123&city=city&address=address")
// 		if err != nil || createGoldbarOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,createGoldbarOrder:%v", err, createGoldbarOrder)
// 			return
// 		}
// 		fmt.Printf("createGoldbarOrder--->%v\n", converter.JSON(createGoldbarOrder))
// 		orderID := createGoldbarOrder.StrDef("", "/order/order_id")
// 		code := createGoldbarOrder.StrDef("", "/order/transaction/code")

// 		verifyGoldbarOrder, err := ts.GetMap("/usr/verifyGoldbarOrder?order_id=%v&code=%v", orderID, code)
// 		if err != nil || verifyGoldbarOrder.Int64("code") != define.NotAccess {
// 			t.Errorf("err:%v,verifyGoldbarOrder:%v", err, verifyGoldbarOrder)
// 			return
// 		}

// 		clearCookie()
// 		login, err = ts.GetMap("/pub/login?username=%v&password=%v", "admin", "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		verifyGoldbarOrder, err = ts.GetMap("/usr/verifyGoldbarOrder?order_id=%v&code=%v", orderID, code)
// 		if err != nil || verifyGoldbarOrder.Int64("code") != 0 {
// 			t.Errorf("err:%v,verifyGoldbarOrder:%v", err, verifyGoldbarOrder)
// 			return
// 		}
// 		fmt.Printf("verifyGoldbarOrder--->%v\n", converter.JSON(verifyGoldbarOrder))
// 		searchOrderRes, err := ts.GetMap("/usr/searchOrder?type=%v", gexdb.OrderTypeGoldbar)
// 		if err != nil || searchOrderRes.Int64("code") != 0 {
// 			t.Errorf("err:%v,searchOrderRes:%v", err, searchOrderRes)
// 			return
// 		}
// 	}
// 	{ //normal error
// 		clearCookie()
// 		login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userYWE.Account, "123")
// 		if err != nil || login.Int64("code") != 0 {
// 			t.Errorf("err:%v,code:%v", err, login)
// 			return
// 		}
// 		createGoldbarOrder, err := ts.GetMap("/usr/createGoldbarOrder?quantity=1&password=1234&city=city&address=address")
// 		if err != nil || createGoldbarOrder.Int64("code") != 8000 {
// 			t.Errorf("err:%v,createGoldbarOrder:%v", err, createGoldbarOrder)
// 			return
// 		}
// 		createGoldbarOrder, err = ts.GetMap("/usr/createGoldbarOrder?quantity=100000&password=123&city=city&address=address")
// 		if err != nil || createGoldbarOrder.Int64("code") != 9000 {
// 			t.Errorf("err:%v,createGoldbarOrder:%v", err, createGoldbarOrder)
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
// 	login, err := ts.GetMap("/pub/login?username=%v&password=%v", *userYWE.Account, "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//create error
// 	res, err = ts.GetMap("/usr/createGoldbarOrder?quantity=x")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	res, err = ts.GetMap("/usr/createGoldbarOrder?quantity=1&password=123&city=city&address=address")
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//create test data
// 	createGoldbarOrder, err := ts.GetMap("/usr/createGoldbarOrder?quantity=1&password=123&city=city&address=address")
// 	if err != nil || createGoldbarOrder.Int64("code") != 0 {
// 		t.Errorf("err:%v,createGoldbarOrder:%v", err, createGoldbarOrder)
// 		return
// 	}
// 	orderID := createGoldbarOrder.StrDef("", "/order/order_id")
// 	code := createGoldbarOrder.StrDef("", "/order/transaction/code")
// 	pgx.MockerClear()

// 	//cancel error
// 	res, err = ts.GetMap("/usr/cancelGoldbarOrder?order_id=%v", "")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Row.Scan", 1)
// 	res, err = ts.GetMap("/usr/cancelGoldbarOrder?order_id=%v", orderID)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	res, err = ts.GetMap("/usr/cancelGoldbarOrder?order_id=%v", orderID)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//to admin
// 	clearCookie()
// 	login, err = ts.GetMap("/pub/login?username=%v&password=%v", "admin", "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//verify error
// 	res, err = ts.GetMap("/usr/verifyGoldbarOrder?order_id=%v", "")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Row.Scan", 1)
// 	res, err = ts.GetMap("/usr/verifyGoldbarOrder?order_id=%v&code=%v", orderID, code)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Begin", 1)
// 	res, err = ts.GetMap("/usr/verifyGoldbarOrder?order_id=%v&code=%v", orderID, code)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// }
