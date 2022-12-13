package gexapi

import (
	"time"

	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xtime"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

var overviewCache = xmap.M{}
var overviewLast = time.Time{}

//LoadOverviewH is http handler
/**
 *
 * @api {GET} /usr/loadOverview Load Overview
 * @apiName LoadOverview
 * @apiGroup Count
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} user the user count info
 * @apiSuccess (Success) {Number} user.total the total user count info
 * @apiSuccess (Success) {Number} user.today the today user register count info
 * @apiSuccess (Success) {Number} user.online the online user count info
 * @apiSuccess (Success) {Object} fee the fee count info, fee.today is today fee info, fee.total is total fee info
 * @apiSuccess (Success) {Object} fee.xxx the fee count info by area, all suported is <a href="#metadata-Order">OrderAreaAll</a>
 * @apiSuccess (Success) {Object} fee.xxx.xxx the fee count info by balance asset
 * @apiSuccess (Success) {Array} assets the all asset key
 * @apiSuccess (Success) {Object} spot the spot balance info
 * @apiSuccess (Success) {Array} spot.balances the total spot balance info
 * @apiSuccess (Success) {Object} futures the futures holding info
 * @apiSuccess (Success) {Array} futures.balances the total futures balance info
 * @apiSuccess (Success) {Array} futures.symbols the all symbol key
 * @apiSuccess (Success) {Array} futures.buy the futures buy holding info
 * @apiSuccess (Success) {Array} futures.sell the futures sell holding info
 * @apiSuccess (Success) {Object} trade the trade count info
 * @apiSuccess (Success) {Array} trade.xxx the area trade count info, trade.spot is spot count info, trade.futures is futures count info
 * @apiSuccess (Success) {Decimal} trade.xxx.filled the traded amount
 * @apiSuccess (Success) {Decimal} trade.xxx.total_price the traded total volume
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "assets": [
 *         "USDT",
 *         "YWE"
 *     ],
 *     "code": 0,
 *     "fee": {
 *         "today": {
 *             "200": {
 *                 "USDT": "0.29000000000000004"
 *             }
 *         },
 *         "total": {
 *             "200": {
 *                 "USDT": "0.29000000000000004"
 *             }
 *         }
 *     },
 *     "futures": {
 *         "balances": [
 *             {
 *                 "area": 100,
 *                 "asset": "USDT",
 *                 "create_time": 0,
 *                 "free": "1000",
 *                 "locked": "0",
 *                 "margin": "0",
 *                 "update_time": 0
 *             },
 *             {
 *                 "area": 100,
 *                 "asset": "YWE",
 *                 "create_time": 0,
 *                 "free": "0",
 *                 "locked": "0",
 *                 "margin": "0",
 *                 "update_time": 0
 *             }
 *         ],
 *         "buy": {
 *             "futures.YWEUSDT": {
 *                 "amount": "0.5",
 *                 "blowup": "0",
 *                 "create_time": 0,
 *                 "margin_added": "0",
 *                 "margin_used": "0",
 *                 "open": "0",
 *                 "symbol": "futures.YWEUSDT",
 *                 "update_time": 0
 *             }
 *         },
 *         "sell": {
 *             "futures.YWEUSDT": {
 *                 "amount": "-0.5",
 *                 "blowup": "0",
 *                 "create_time": 0,
 *                 "margin_added": "0",
 *                 "margin_used": "0",
 *                 "open": "0",
 *                 "symbol": "futures.YWEUSDT",
 *                 "update_time": 0
 *             }
 *         },
 *         "symbols": [
 *             "futures.YWEUSDT"
 *         ]
 *     },
 *     "spot": {
 *         "balances": [
 *             {
 *                 "area": 200,
 *                 "asset": "USDT",
 *                 "create_time": 0,
 *                 "free": "2954.71",
 *                 "locked": "45",
 *                 "margin": "0",
 *                 "update_time": 0
 *             },
 *             {
 *                 "area": 200,
 *                 "asset": "YWE",
 *                 "create_time": 0,
 *                 "free": "1999.997",
 *                 "locked": "0",
 *                 "margin": "0",
 *                 "update_time": 0
 *             }
 *         ]
 *     },
 *     "trade": {
 *         "futures": [
 *             {
 *                 "area": 300,
 *                 "avg_price": "0",
 *                 "create_time": 0,
 *                 "fee_filled": "0",
 *                 "fee_rate": "0",
 *                 "fee_settled_next": 0,
 *                 "filled": "0.5",
 *                 "holding": "0",
 *                 "in_filled": "0",
 *                 "out_filled": "0",
 *                 "owned": "0",
 *                 "price": "0",
 *                 "profit": "0",
 *                 "quantity": "0.5",
 *                 "symbol": "futures.YWEUSDT",
 *                 "total_price": "50",
 *                 "transaction": {},
 *                 "trigger_price": "0",
 *                 "trigger_time": 0,
 *                 "unhedged": "0",
 *                 "update_time": 0
 *             }
 *         ],
 *         "spot": [
 *             {
 *                 "area": 200,
 *                 "avg_price": "0",
 *                 "create_time": 0,
 *                 "fee_filled": "0",
 *                 "fee_rate": "0",
 *                 "fee_settled_next": 0,
 *                 "filled": "3",
 *                 "holding": "0",
 *                 "in_filled": "0",
 *                 "out_filled": "0",
 *                 "owned": "0",
 *                 "price": "0",
 *                 "profit": "0",
 *                 "quantity": "3.5",
 *                 "symbol": "spot.YWEUSDT",
 *                 "total_price": "290",
 *                 "transaction": {},
 *                 "trigger_price": "0",
 *                 "trigger_time": 0,
 *                 "unhedged": "0",
 *                 "update_time": 0
 *             }
 *         ]
 *     },
 *     "user": {
 *         "online": 0,
 *         "today": 6,
 *         "total": 7
 *     }
 * }
 */
func LoadOverviewH(s *web.Session) web.Result {
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	if time.Since(overviewLast) < time.Minute {
		return s.SendJSON(overviewCache)
	}
	var err error
	result := xmap.M{}
	{
		result["assets"], err = gexdb.ListBalanceAsset(s.R.Context(), nil)
		if err != nil {
			xlog.Errorf("LoadOverviewH list balance asset fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
	}
	{ //user
		user := xmap.M{}
		user["total"], err = gexdb.CountUser(s.R.Context(), time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count user total fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		user["today"], err = gexdb.CountUser(s.R.Context(), xtime.TimeStartOfToday(), time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count user today fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		user["online"] = MarketOnline.Size()
		result["user"] = user
	}
	{ //fee
		fee := xmap.M{}
		_, fee["total"], err = gexdb.CountOrderFee(s.R.Context(), 0, time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count user total fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		_, fee["today"], err = gexdb.CountOrderFee(s.R.Context(), 0, xtime.TimeStartOfToday(), time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count user today fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		result["fee"] = fee
	}
	{ //spot
		spot := xmap.M{}
		spot["balances"], _, err = gexdb.CountAreaBalance(s.R.Context(), gexdb.BalanceAreaArray{gexdb.BalanceAreaSpot}, "", time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count user total fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		result["spot"] = spot
	}
	{ //futures
		futures := xmap.M{}
		futures["balances"], _, err = gexdb.CountAreaBalance(s.R.Context(), gexdb.BalanceAreaArray{gexdb.BalanceAreaFunds}, "", time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count user total fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		futures["symbols"], err = gexdb.ListHoldingSymbol(s.R.Context())
		if err != nil {
			xlog.Errorf("LoadOverviewH list holding symbol fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		_, futures["buy"], err = gexdb.CountHolding(s.R.Context(), 1, time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count buy holding fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		_, futures["sell"], err = gexdb.CountHolding(s.R.Context(), -1, time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count sell holding fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		result["futures"] = futures
	}
	{ //trade
		trade := xmap.M{}
		trade["spot"], err = gexdb.CountOrderVolume(s.R.Context(), gexdb.OrderAreaSpot, time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count trade spot total fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		trade["futures"], err = gexdb.CountOrderVolume(s.R.Context(), gexdb.OrderAreaFutures, time.Time{}, time.Time{})
		if err != nil {
			xlog.Errorf("LoadOverviewH count trade futures total fail with %v", err)
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		result["trade"] = trade
	}
	result["code"] = define.Success
	overviewCache = result
	overviewLast = time.Now()
	return s.SendJSON(overviewCache)
}

//ListBalanceCountH is http handler
/**
 *
 * @api {GET} /usr/listBalanceCount List Balance Count
 * @apiName ListBalanceCount
 * @apiGroup Count
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Balance) {Array} balances the user balance info
 * @apiUse BalanceObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "balances": [
 *         {
 *             "area": 100,
 *             "asset": "USDT",
 *             "create_time": 0,
 *             "free": "1000",
 *             "locked": "0",
 *             "margin": "0",
 *             "update_time": 0
 *         },
 *         {
 *             "area": 100,
 *             "asset": "YWE",
 *             "create_time": 0,
 *             "free": "0",
 *             "locked": "0",
 *             "margin": "0",
 *             "update_time": 0
 *         },
 *         {
 *             "area": 200,
 *             "asset": "USDT",
 *             "create_time": 0,
 *             "free": "2954.71",
 *             "locked": "45",
 *             "margin": "0",
 *             "update_time": 0
 *         },
 *         {
 *             "area": 200,
 *             "asset": "YWE",
 *             "create_time": 0,
 *             "free": "1999.997",
 *             "locked": "0",
 *             "margin": "0",
 *             "update_time": 0
 *         },
 *         {
 *             "area": 300,
 *             "asset": "USDT",
 *             "create_time": 0,
 *             "free": "1960.61",
 *             "locked": "39.19",
 *             "margin": "0",
 *             "update_time": 0
 *         }
 *     ],
 *     "code": 0
 * }
 */
func ListBalanceCountH(s *web.Session) web.Result {
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	asset := s.Argument("asset")
	balances, _, err := gexdb.CountAreaBalance(s.R.Context(), nil, asset, time.Time{}, time.Time{})
	if err != nil {
		xlog.Errorf("ListBalanceCountH count all balance fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     define.Success,
		"balances": balances,
	})
}
