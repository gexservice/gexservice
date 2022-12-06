package gexapi

import (
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/market"
	"github.com/gexservice/gexservice/matcher"
)

//PlaceOrderH is http handler
/**
 *
 * @api {POST} /usr/placeOrder Place Order
 * @apiName PlaceOrder
 * @apiGroup Order
 *
 * @apiParam  {String} symbol the symbol to trade
 * @apiParam  {String} type the type, all type supported is <a href="#metadata-Order">OrderTypeAll</a>
 * @apiParam  {String} side the trade side, all type supported is <a href="#metadata-Order">OrderSideAll</a>
 * @apiParam  {Number} [price] the limit price to buy/sell, price>0 is limit order, price=0 is market order
 * @apiParam  {Number} [total_price] the total price to buy, only supported when side=OrderSideBuy and price=0
 * @apiParam  {Number} [quantity] the total quantity to trade, required when price>0
 * @apiParam  {Number} [trigger_type] the trigger type, required when type=OrderTypeTrigger, all type supported is <a href="#metadata-Order">OrderTriggerTypeAll</a>
 * @apiParam  {Number} [trigger_price] the trigger price, required when type=OrderTypeTrigger, it is ms when trigger_type=OrderTriggerTypeAfterOpen
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a> or <a href="#metadata-ExReturnCode">ExReturnCode</a>
 * @apiSuccess (Order) {Object} order the created order info
 * @apiUse OrderObject
 *
 * @apiParamExample  {Query} Market Buy By Total:
 * type=OrderTypeTrade&symbol=YWKUSDT&side=OrderSideBuy&total_price=10
 *
 * @apiParamExample  {Query} Market Buy By Quantity:
 * type=OrderTypeTrade&symbol=YWKUSDT&side=OrderSideBuy&quantity=1
 *
 * @apiParamExample  {Query} Limit Buy:
 * type=OrderTypeTrade&symbol=YWKUSDT&side=OrderSideBuy&quantity=1&price=100
 *
 * @apiParamExample  {Query} Market Sell:
 * type=OrderTypeTrade&symbol=YWKUSDT&side=OrderSideSell&quantity=1
 *
 * @apiParamExample  {Query} Limit Sell:
 * type=OrderTypeTrade&symbol=YWKUSDT&side=OrderSideSell&quantity=1&price=100
 *
 * @apiParamExample  {Query} Trigger Limit Sell:
 * type=OrderTypeTrigger&symbol=YWKUSDT&side=OrderSideSell&quantity=1&price=100&trigger_type=OrderTriggerTypeStopProfit&trigger_price=100
 *
 * @apiParamExample  {Query} After Trigger:
 * type=OrderTypeTrigger&symbol=YWKUSDT&side=OrderSideSell&quantity=1&trigger_type=OrderTriggerTypeAfterOpen&trigger_price=60000
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "order": {
 *         "avg_price": "95",
 *         "create_time": 1667475452051,
 *         "creator": 100002,
 *         "fee_balance": "YWE",
 *         "fee_filled": "0",
 *         "fee_settled_next": 0,
 *         "filled": "0",
 *         "holding": "0",
 *         "in_balance": "YWE",
 *         "in_filled": "0",
 *         "order_id": "202211031937320100009",
 *         "out_balance": "USDT",
 *         "out_filled": "0",
 *         "owned": "0",
 *         "price": "95",
 *         "profit": "0",
 *         "quantity": "1",
 *         "side": "buy",
 *         "status": 200,
 *         "symbol": "spot.YWEUSDT",
 *         "tid": 1005,
 *         "total_price": "0",
 *         "transaction": {},
 *         "trigger_price": "0",
 *         "type": 100,
 *         "unhedged": "0",
 *         "update_time": 1667475452051,
 *         "user_id": 100002
 *     }
 * }
 *
 */
func PlaceOrderH(s *web.Session) web.Result {
	var err error
	var args = &gexdb.Order{}
	filter := "tid,type,symbol,side,quantity,price,total_price,trigger_type,trigger_price,status#all"
	if s.R.Method == "GET" {
		err = s.Valid(args, filter, "")
	} else {
		_, err = s.RecvValidJSON(args, filter, "")
	}
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	args.UserID = userID
	args.Creator = userID
	order, err := matcher.ProcessOrder(s.R.Context(), args)
	if err != nil {
		xlog.Errorf("PlaceOrderH process order by %v, err is \n%v", converter.JSON(args), matcher.ErrStack(err))
		code := define.ServerError
		if matcher.IsErrBalanceNotEnought(err) {
			code = gexdb.CodeBalanceNotEnought
		}
		return util.ReturnCodeLocalErr(s, code, "srv-err", err)
	}
	xlog.Infof("PlaceOrderH user %v process order success with %v", order.UserID, order.Info())
	return s.SendJSON(xmap.M{
		"code":  0,
		"order": order,
	})
}

//CancelOrderH is http handler
/**
 *
 * @api {GET} /usr/cancelOrder Cancel Order
 * @apiName CancelOrder
 * @apiGroup Order
 *
 * @apiParam  {String} order_id the order id
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a> or <a href="#metadata-ExReturnCode">ExReturnCode</a>
 * @apiSuccess (Order) {Object} order the cancel order info
 * @apiUse OrderObject
 *
 * @apiParamExample  {Query} Cancel Order:
 * order_id=100
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "order": {
 *         "avg_price": "10",
 *         "create_time": 1667475452026,
 *         "creator": 100002,
 *         "fee_balance": "YWE",
 *         "fee_filled": "0",
 *         "fee_settled_next": 0,
 *         "filled": "0",
 *         "holding": "0",
 *         "in_balance": "YWE",
 *         "in_filled": "0",
 *         "order_id": "202211031937320100007",
 *         "out_balance": "USDT",
 *         "out_filled": "0",
 *         "owned": "0",
 *         "price": "10",
 *         "profit": "0",
 *         "quantity": "1",
 *         "side": "buy",
 *         "status": 420,
 *         "symbol": "spot.YWEUSDT",
 *         "tid": 1003,
 *         "total_price": "0",
 *         "transaction": {},
 *         "trigger_price": "0",
 *         "type": 100,
 *         "unhedged": "0",
 *         "update_time": 1667475452032,
 *         "user_id": 100002
 *     }
 * }
 */
func CancelOrderH(s *web.Session) web.Result {
	var orderID string
	err := s.ValidFormat(`
		order_id,R|S,L:0;
	`, &orderID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	order, err := gexdb.FindOrderByOrderID(s.R.Context(), userID, orderID)
	if err != nil {
		xlog.Errorf("CancelOrderH find order fail with %v by %v,%v", err, userID, orderID)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	if order.Status == gexdb.OrderStatusWaiting {
		_, err = gexdb.CancelTriggerOrder(s.R.Context(), userID, order.Symbol, order.TID)
		if err == nil {
			order.Status = gexdb.OrderStatusCanceled
		}
	} else {
		order, err = matcher.ProcessCancel(s.R.Context(), userID, order.Symbol, orderID)
	}
	if err != nil {
		code := define.ServerError
		if matcher.IsErrNotCancelable(err) {
			code = gexdb.CodeOrderNotCancelable
		} else {
			xlog.Errorf("CancelOrderH cancel order  by user:%v,order_id:%v, err is \n%v", userID, orderID, matcher.ErrStack(err))
		}
		return util.ReturnCodeLocalErr(s, code, "srv-err", err)
	}
	xlog.Infof("PlaceOrderH user %v cancel order success with %v", order.UserID, order.Info())
	return s.SendJSON(xmap.M{
		"code":  0,
		"order": order,
	})
}

//CancelAllOrderH is http handler
/**
 *
 * @api {GET} /usr/cancelAllOrder Cancel All Order
 * @apiName CancelAllOrder
 * @apiGroup Order
 *
 * @apiParam  {String} [area] the order symbol area
 * @apiParam  {String} [symbol] the order symbol
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a> or <a href="#metadata-ExReturnCode">ExReturnCode</a>
 *
 * @apiParamExample  {Query} Cancel Order:
 * symbol=100
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0
 * }
 */
func CancelAllOrderH(s *web.Session) web.Result {
	var area string
	var symbol string
	err := s.ValidFormat(`
		area,O|S,L:0;
		symbol,O|S,L:0~32;
	`, &area, &symbol)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	orders, err := gexdb.ListPendingOrder(s.R.Context(), userID, area, symbol)
	if err != nil {
		xlog.Errorf("CancelAllOrderH list pending order fail with %v by %v,%v", err, userID, symbol)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	for _, order := range orders {
		if order.Status == gexdb.OrderStatusWaiting {
			_, err = gexdb.CancelTriggerOrder(s.R.Context(), userID, order.Symbol, order.TID)
		} else {
			_, err = matcher.ProcessCancel(s.R.Context(), userID, order.Symbol, order.OrderID)
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	xlog.Infof("PlaceOrderH user %v cancel %v order success with %v", userID, symbol)
	return s.SendJSON(xmap.M{
		"code": 0,
	})
}

//SearchOrderH is http handler
/**
 *
 * @api {GET} /usr/searchOrder Search Order
 * @apiName SearchOrder
 * @apiGroup Order
 *
 * @apiUse OrderUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Order) {Array} orders the order array
 * @apiUse OrderObject
 * @apiSuccess (SymbolInfo) {Array} symbols the symbol info, mapping by symbol key
 * @apiUse SymbolInfoObject
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "orders": [
 *         {
 *             "avg_price": "95",
 *             "create_time": 1667475452051,
 *             "creator": 100002,
 *             "fee_balance": "YWE",
 *             "fee_filled": "0.002",
 *             "fee_settled_next": 0,
 *             "filled": "1",
 *             "holding": "0",
 *             "in_balance": "YWE",
 *             "in_filled": "0.998",
 *             "order_id": "202211031937320100009",
 *             "out_balance": "USDT",
 *             "out_filled": "95",
 *             "owned": "0",
 *             "price": "95",
 *             "profit": "0",
 *             "quantity": "1",
 *             "side": "buy",
 *             "status": 400,
 *             "symbol": "spot.YWEUSDT",
 *             "tid": 1005,
 *             "total_price": "95",
 *             "trigger_price": "0",
 *             "type": 100,
 *             "unhedged": "0",
 *             "update_time": 1667475452051,
 *             "user_id": 100002
 *         }
 *     ],
 *     "total": 4
 * }
 */
func SearchOrderH(s *web.Session) web.Result {
	searcher := &gexdb.OrderUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	if !AdminAccess(s) {
		searcher.Where.UserID = xsql.Int64Array{userID}
	}
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Errorf("SearchOrderH searcher order fail with %v by %v", err, converter.JSON(searcher))
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	_, symbols, _ := market.ListSymbol("", searcher.Query.Symbols, "")
	return s.SendJSON(xmap.M{
		"code":    define.Success,
		"orders":  searcher.Query.Orders,
		"symbols": symbols,
		"total":   searcher.Count.Total,
	})
}

//QueryOrderH is http handler
/**
 *
 * @api {GET} /usr/queryOrder Query Order
 * @apiName QueryOrder
 * @apiGroup Order
 *
 * @apiParam  {String} order_id the order id, it can be order.tid or order.order_id
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Order) {Object} order the order info
 * @apiUse OrderObject
 * @apiSuccess (SymbolInfo) {Object} symbol the symbol info
 * @apiUse SymbolInfoObject
 *
 * @apiParamExample  {Query} QueryOrder:
 * order_id=100
 *
 *
 * @apiSuccessExample {JSON} Order-Response:
 * {
 *     "code": 0,
 *     "order": {
 *         "avg_price": "95",
 *         "create_time": 1667475452063,
 *         "creator": 100004,
 *         "fee_balance": "USDT",
 *         "fee_filled": "0.19",
 *         "fee_settled_next": 0,
 *         "filled": "1",
 *         "holding": "0",
 *         "in_balance": "USDT",
 *         "in_filled": "94.81",
 *         "order_id": "202211031937320100010",
 *         "out_balance": "YWE",
 *         "out_filled": "1",
 *         "owned": "0",
 *         "price": "95",
 *         "profit": "0",
 *         "quantity": "1",
 *         "side": "sell",
 *         "status": 400,
 *         "symbol": "spot.YWEUSDT",
 *         "tid": 1006,
 *         "total_price": "95",
 *         "transaction": {
 *             "trans": [
 *                 {
 *                     "order_id": "202211031937320100009",
 *                     "create_time": 1667475452061,
 *                     "fee_balance": "USDT",
 *                     "fee_filled": "0.19",
 *                     "filled": "1",
 *                     "price": "95",
 *                     "total_price": "95"
 *                 }
 *             ]
 *         },
 *         "trigger_price": "0",
 *         "type": 100,
 *         "unhedged": "0",
 *         "update_time": 1667475452063,
 *         "user_id": 100004
 *     }
 * }
 */
func QueryOrderH(s *web.Session) web.Result {
	var orderID string
	err := s.ValidFormat(`
		order_id,R|S,L:0;
	`, &orderID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	order, err := gexdb.FindOrderByOrderID(s.R.Context(), 0, orderID)
	if err != nil {
		xlog.Errorf("QueryOrderH find order fail with %v by %v", err, orderID)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	if order.Creator != userID {
		user, err := gexdb.FindUser(s.R.Context(), userID)
		if err != nil {
			xlog.Errorf("QueryOrderH find current user(%v) err: %v", userID, err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		if user.Type != gexdb.UserTypeAdmin {
			return util.ReturnCodeLocalErr(s, define.NotAccess, "srv-err", define.ErrNotAccess)
		}
	}
	symbol, _ := market.LoadSymbol(order.Symbol)
	return s.SendJSON(xmap.M{
		"code":   0,
		"order":  order,
		"symbol": symbol,
	})
}
