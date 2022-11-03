package gexapi

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/codingeasygo/util/converter"
// 	"github.com/codingeasygo/util/xmap"
// 	"github.com/gexservice/gexservice/base/define"
// )

// var CreatePrepayOrderErr error

// func init() {
// 	CreatePrepayOrder("", "", 100, "")
// 	CreatePrepayOrder = func(merchOrderID, title string, totalAmount float64, transCurrency string) (result xmap.M, err error) {
// 		result = xmap.New()
// 		err = CreatePrepayOrderErr
// 		return
// 	}
// }

// func TestTopupOrder(t *testing.T) {
// 	login, err := ts.GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	//
// 	//test create popup order
// 	createTopupOrderRes, err := ts.GetMap("/usr/createTopupOrder?amount=%v", 100)
// 	if err != nil || createTopupOrderRes.Int64("code") != 0 {
// 		t.Errorf("err:%v,createTopupOrderRes:%v", err, createTopupOrderRes)
// 		return
// 	}
// 	fmt.Printf("createTopupOrderRes--->%v\n", converter.JSON(createTopupOrderRes))
// 	topupOrderID := createTopupOrderRes.StrDef("", "/order/order_id")

// 	mockPayTopupOrderRes, err := ts.GetMap("/usr/mockPayTopupOrder?order_id=%v&amount=%v", topupOrderID, 100)
// 	if err != nil || mockPayTopupOrderRes.Int64("code") != 0 {
// 		t.Errorf("err:%v,mockPayTopupOrderRes:%v", err, mockPayTopupOrderRes)
// 		return
// 	}
// 	fmt.Printf("mockPayTopupOrderRes--->%v\n", converter.JSON(mockPayTopupOrderRes))

// 	//
// 	//test search order
// 	searchOrderRes, err := ts.GetMap("/usr/searchOrder")
// 	if err != nil || searchOrderRes.Int64("code") != 0 {
// 		t.Errorf("err:%v,searchOrderRes:%v", err, searchOrderRes)
// 		return
// 	}
// 	fmt.Printf("searchOrderRes--->%v\n", converter.JSON(searchOrderRes))

// 	//
// 	//test query order
// 	queryOrderRes, err := ts.GetMap("/usr/queryOrder?order_id=%v", topupOrderID)
// 	if err != nil || queryOrderRes.Int64("code") != 0 {
// 		t.Errorf("err:%v,queryOrderRes:%v", err, queryOrderRes)
// 		return
// 	}
// 	fmt.Printf("queryOrderRes--->%v\n", converter.JSON(queryOrderRes))

// 	//login to other
// 	clearCookie()
// 	login, err = ts.GetMap("/pub/login?username=%v&password=%v", "abc2", "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	queryOrderRes, err = ts.GetMap("/usr/queryOrder?order_id=%v", topupOrderID)
// 	if err != nil || queryOrderRes.Int64("code") != define.NotAccess {
// 		t.Errorf("err:%v,queryOrderRes:%v", err, queryOrderRes)
// 		return
// 	}

// 	//login to admin
// 	clearCookie()
// 	login, err = ts.GetMap("/pub/login?username=%v&password=%v", "admin", "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	queryOrderRes, err = ts.GetMap("/usr/queryOrder?order_id=%v", topupOrderID)
// 	if err != nil || queryOrderRes.Int64("code") != define.Success {
// 		t.Errorf("err:%v,queryOrderRes:%v", err, queryOrderRes)
// 		return
// 	}
// 	searchOrderRes, err = ts.GetMap("/usr/searchOrder")
// 	if err != nil || searchOrderRes.Int64("code") != 0 {
// 		t.Errorf("err:%v,searchOrderRes:%v", err, searchOrderRes)
// 		return
// 	}
// 	fmt.Printf("searchOrderRes--->%v\n", converter.JSON(searchOrderRes))

// 	//
// 	//test error
// 	pgx.MockerStart()
// 	defer pgx.MockerStop()
// 	pgx.MockerClear()
// 	var res xmap.M

// 	//
// 	//test create topup error
// 	res, err = ts.GetMap("/usr/createTopupOrder?amount=%v", "xxxx")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Row.Scan", 1)
// 	res, err = ts.GetMap("/usr/createTopupOrder?amount=%v", 100)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	CreatePrepayOrderErr = fmt.Errorf("test error")
// 	res, err = ts.GetMap("/usr/createTopupOrder?amount=%v", 100)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	CreatePrepayOrderErr = nil
// 	pgx.MockerClear()

// 	pgx.MockerSet("Pool.Exec", 1)
// 	res, err = ts.GetMap("/usr/createTopupOrder?amount=%v", 100)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//
// 	//test mock pay topup order error
// 	res, err = ts.GetMap("/usr/mockPayTopupOrder?order_id=%v", "")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Pool.Begin", 1)
// 	res, err = ts.GetMap("/usr/mockPayTopupOrder?order_id=%v&amount=%v", topupOrderID, 100)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//
// 	//test query order error
// 	res, err = ts.GetMap("/usr/queryOrder?order_id=%v", "")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Row.Scan", 1)
// 	res, err = ts.GetMap("/usr/queryOrder?order_id=%v", topupOrderID)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	pgx.MockerSet("Row.Scan", 2)
// 	res, err = ts.GetMap("/usr/queryOrder?order_id=%v", topupOrderID)
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()

// 	//
// 	//test search order error
// 	clearCookie() //to admin
// 	login, err = ts.GetMap("/pub/login?username=%v&password=%v", "admin", "123")
// 	if err != nil || login.Int64("code") != 0 {
// 		t.Errorf("err:%v,code:%v", err, login)
// 		return
// 	}
// 	res, err = ts.GetMap("/usr/searchOrder?user_id=xx")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	res, err = ts.GetMap("/usr/searchOrder?type=xx")
// 	if err != nil || res.Int64("code") != define.ArgsInvalid {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Row.Scan", 1)
// 	res, err = ts.GetMap("/usr/searchOrder")
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Query", 1)
// 	res, err = ts.GetMap("/usr/searchOrder")
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// 	pgx.MockerSet("Pool.Query", 2)
// 	res, err = ts.GetMap("/usr/searchOrder")
// 	if err != nil || res.Int64("code") != define.ServerError {
// 		t.Errorf("err:%v,res:%v", err, res)
// 		return
// 	}
// 	pgx.MockerClear()
// }
