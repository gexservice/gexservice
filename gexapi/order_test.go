package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

func TestOrder(t *testing.T) {
	symbol := "spot.YWEUSDT"
	{ //balance not enought
		clearCookie()
		ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc3.Account, "123")
		ts.Should(t, "code", gexdb.CodeBalanceNotEnought).GetMap("/usr/placeOrder?type=%v&symbol=%v&side=%v&quantity=1&price=10", gexdb.OrderTypeTrade, symbol, gexdb.OrderSideBuy)
	}
	{ //buy cancel
		clearCookie()
		ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc0.Account, "123")
		ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/placeOrder?type=%v&symbol=%v&side=%v&quantity=1&price=100", gexdb.OrderTypeTrade, symbol, 1)
		buyOrder, _ := ts.Should(t, "code", define.Success, "/order/tid", xmap.ShouldIsNoZero).GetMap("/usr/placeOrder?type=%v&symbol=%v&side=%v&quantity=1&price=10", gexdb.OrderTypeTrade, symbol, gexdb.OrderSideBuy)
		orderID := buyOrder.StrDef("", "/order/order_id")
		ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/cancelOrder?symbol=%v&order_id=%v", "", orderID)
		cancelOrder, _ := ts.Should(t, "code", define.Success, "/order/status", gexdb.OrderStatusCanceled).GetMap("/usr/cancelOrder?symbol=%v&order_id=%v", symbol, orderID)
		fmt.Printf("cancelOrder--->%v\n", converter.JSON(cancelOrder))
		ts.Should(t, "code", gexdb.CodeOrderNotCancelable).GetMap("/usr/cancelOrder?symbol=%v&order_id=%v", symbol, orderID)
		ts.Should(t, "code", define.Success, "/order/status", gexdb.OrderStatusCanceled).GetMap("/usr/queryOrder?order_id=%v", orderID)
	}
	{ //buy cancel(post)
		clearCookie()
		ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc0.Account, "123")
		ts.Should(t, "code", define.ArgsInvalid).PostJSONMap(&gexdb.Order{}, "/usr/placeOrder")
		buyArgs := &gexdb.Order{Type: gexdb.OrderTypeTrade, Symbol: symbol, Side: gexdb.OrderSideBuy, Quantity: decimal.NewFromFloat(1), Price: decimal.NewFromFloat(10)}
		buyOrder, _ := ts.Should(t, "code", define.Success, "/order/tid", xmap.ShouldIsNoZero).PostJSONMap(buyArgs, "/usr/placeOrder")
		orderID := buyOrder.StrDef("", "/order/order_id")
		ts.Should(t, "code", define.Success).GetMap("/usr/cancelOrder?symbol=%v&order_id=%v", symbol, orderID)
	}
	{ //buy sell
		clearCookie()
		ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc0.Account, "123")
		buyOrder, _ := ts.Should(t, "code", define.Success, "/order/tid", xmap.ShouldIsNoZero).GetMap("/usr/placeOrder?type=%v&symbol=%v&side=%v&quantity=1&price=95", gexdb.OrderTypeTrade, symbol, gexdb.OrderSideBuy)
		fmt.Printf("buyOrder--->%v\n", converter.JSON(buyOrder))
		//
		clearCookie()
		ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc2.Account, "123")
		sellOrder, _ := ts.Should(t, "code", define.Success, "/order/tid", xmap.ShouldIsNoZero).GetMap("/usr/placeOrder?type=%v&symbol=%v&side=%v&quantity=1&price=95", gexdb.OrderTypeTrade, symbol, gexdb.OrderSideSell)
		fmt.Printf("sellOrder--->%v\n", converter.JSON(sellOrder))

		ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/queryOrder?order_id=%v", "")
		ts.Should(t, "code", define.NotAccess).GetMap("/usr/queryOrder?order_id=%v", buyOrder.StrDef("", "/order/order_id"))
		queryOrder, _ := ts.Should(t, "code", define.Success, "/order/status", gexdb.OrderStatusDone).GetMap("/usr/queryOrder?order_id=%v", sellOrder.StrDef("", "/order/order_id"))
		fmt.Printf("queryOrder--->%v\n", converter.JSON(queryOrder))
	}
	//search
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc0.Account, "123")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/searchOrder?type=10")
	searchOrder, _ := ts.Should(t, "code", define.Success, "/orders", xmap.ShouldIsNoEmpty).GetMap("/usr/searchOrder")
	fmt.Printf("searchOrder--->%v\n", converter.JSON(searchOrder))
	orderID := searchOrder.StrDef("", "/orders/0/order_id")
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()

	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc0.Account, "123")
	pgx.MockerClear()

	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/placeOrder?type=%v&symbol=%v&side=%v&quantity=1&price=10", gexdb.OrderTypeTrade, symbol, gexdb.OrderSideBuy)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/cancelOrder?symbol=%v&order_id=%v", symbol, orderID)
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/searchOrder")
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/queryOrder?order_id=%v", orderID)

	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", *userabc2.Account, "123")
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 2).Should(t, "code", define.ServerError).GetMap("/usr/queryOrder?order_id=%v", orderID)
}
