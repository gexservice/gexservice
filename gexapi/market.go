package gexapi

import (
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xtime"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/market"
)

//Market is struct to market impl
/**
 *
 * @api {WS} /ws/market Ws Market
 * @apiName WsMarket
 * @apiGroup Market
 *
 * @apiParam  {String} action subscribe action, supported is "sub.kline"/"sub.depth"/"sub.ticker"
 * @apiParam  {Arrasy} [intervals] the kline interval, only for "sub.kline", supported is "5min"/"30min"/"1hour"/"4hour"/"day"/"week"/"mon"
 * @apiParam  {Number} [max] the depth max size
 *
 * @apiSuccess (Success) {Number} code the response code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Number} tid the int order id
 * @apiSuccess (Success) {String} action the received action, supported is "sub.kline"/"sub.depth"/"notify.kline"/"notify.depth"
 * @apiSuccess (Success) {Object} depth the received depth data, only for "notify.depth"
 * @apiSuccess (Success) {String} depth.symbol the received depth symbol
 * @apiSuccess (Success) {Array} depth.bids the received depth bids data, the inner data is ["price","quantity"]
 * @apiSuccess (Success) {Array} depth.asks the received depth asks data, the inner data is ["price","quantity"]
 * @apiSuccess (Success) {Object} kline the received kline data, only for "notify.kline"
 * @apiSuccess (Success) {String} kline.symbol the received kline symbol
 * @apiSuccess (Success) {String} kline.start_time the received kline id, the timeline
 * @apiSuccess (Success) {String} kline.volume the received kline total traded price
 * @apiSuccess (Success) {String} kline.amount the received kline total traded quantity
 * @apiSuccess (Success) {String} kline.count the received kline total traded count
 * @apiSuccess (Success) {String} kline.open the received kline open price
 * @apiSuccess (Success) {String} kline.close the received kline close price
 * @apiSuccess (Success) {String} kline.high the received kline high price
 * @apiSuccess (Success) {String} kline.low the received kline low price
 *
 * @apiParamExample  {JSON} Subscribe-KLine:
 * {
 *     "action": "sub.kline",
 *     "symbols": [
 *         {
 *             "symbol": "spot.YWEUSDT",
 *             "interval": "5min"
 *         }
 *     ]
 * }
 * @apiParamExample  {JSON} Subscribe-Depth:
 * {
 *     "action": "sub.depth",
 *     "symbols": [
 *         {
 *             "symbol": "spot.YWEUSDT",
 *             "max": 30
 *         }
 *     ]
 * }
 * @apiParamExample  {JSON} Subscribe-Ticker:
 * {
 *     "action": "sub.ticker",
 *     "symbols": ["spot.YWEUSDT"]
 * }
 *
 * @apiSuccessExample {JSON} Reponse-Depth:
 * {
 *     "action": "sub.kline",
 *     "code": 0
 * }
 *
 * @apiSuccessExample {JSON} Reponse-Depth:
 * {
 *     "action": "sub.depth",
 *     "code": 0
 * }
 * @apiSuccessExample {JSON} Reponse-Ticker:
 * {
 *     "action": "sub.ticker",
 *     "code": 0
 * }
 *
 * @apiSuccessExample {JSON} Notify-Depth:
 * {
 *     "code": 0,
 *     "action": "notify.depth",
 *     "depth": {
 *         "symbol": "spot.YWEUSDT",
 *         "bids": [
 *             [
 *                 "90",
 *                 "2"
 *             ]
 *         ],
 *         "asks": [
 *             [
 *                 "100",
 *                 "2"
 *             ]
 *         ]
 *     }
 * }
 * @apiSuccessExample {JSON} Notify-Ticker:
 * {
 *     "code": 0,
 *     "action": "notify.ticker",
 *     "ticker": {
 *         "symbol": "spot.YWEUSDT",
 *         "bid": [
 *             "90",
 *             "2"
 *         ],
 *         "ask": [
 *             "100",
 *             "2"
 *         ],
 *         "close": "100"
 *     }
 * }
 *
 * @apiSuccessExample {JSON} Notify-KLine:
 * {
 *     "action": "notify.kline",
 *     "code": 0,
 *     "kline": {
 *         "tid": 0,
 *         "symbol": "YWEMMK",
 *         "interval": "5min",
 *         "amount": "1",
 *         "count": 1,
 *         "open": "100",
 *         "close": "100",
 *         "low": "100",
 *         "high": "100",
 *         "volume": "100",
 *         "start_time": 1632577200000,
 *         "update_time": 1632577495577
 *     }
 * }
 *
 */
var MarketOnline *OnlineHander

//ListKLineH is http handler
/**
 *
 * @api {GET} /pub/listKLine List KLine
 * @apiName ListKLine
 * @apiGroup Market
 *
 * @apiParam  {String} interval the kline interval, supported is "5min"/"30min"/"1hour"/"4hour"/"1day"/"1week"/"1mon"
 * @apiParam  {Number} start_time filter kline kline.start_time>=start_time
 * @apiParam  {Number} end_time filter kline kline.start_time<end_time
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Array} lines the all kline array
 * @apiSuccess (Success) {String} lines.start_time the received kline id, the timeline
 * @apiSuccess (Success) {String} lines.volume the received kline total traded price
 * @apiSuccess (Success) {String} lines.amount the received kline total traded quantity
 * @apiSuccess (Success) {String} lines.count the received kline total traded count
 *
 * @apiParamExample  {Query} QueryOrder:
 * interval=100&start_time=100&end_time=1632578100000
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "lines": [
 *         {
 *             "tid": 0,
 *             "symbol": "YWEMMK",
 *             "interval": "5min",
 *             "amount": "1",
 *             "count": 1,
 *             "open": "100",
 *             "close": "100",
 *             "low": "100",
 *             "high": "100",
 *             "volume": "100",
 *             "start_time": 1632578100000,
 *             "update_time": 1632578330897
 *         }
 *     ]
 * }
 *
 */
func ListKLineH(s *web.Session) web.Result {
	var symbol, interval string
	var startTime, endTime int64
	var err = s.ValidFormat(`
		symbol,R|S,L:0;
		interval,R|S,L:0;
		start_time,O|I,R:0;
		end_time,O|I,R:0;
	`, &symbol, &interval, &startTime, &endTime)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	lines, err := market.ListKLine(s.R.Context(), symbol, interval, xtime.TimeUnix(startTime), xtime.TimeUnix(endTime))
	if err != nil {
		xlog.Warnf("ListKLineH list kline fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":  0,
		"lines": lines,
	})
}

//LoadDepthH is http handler
/**
 *
 * @api {GET} /pub/loadDepth Load Depth
 * @apiName LoadDepth
 * @apiGroup Market
 *
 * @apiParam  {Number} [max] max depth
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} depth the depth info
 * @apiSuccess (Success) {Array} depth.bids the depth bid array
 * @apiSuccess (Success) {Array} depth.asks the depth ask array
 *
 * @apiParamExample  {Query} QueryOrder:
 * max=10
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "depth": {
 *         "bids": [
 *             [
 *                 "86.5",
 *                 "0.5"
 *             ]
 *         ],
 *         "asks": [
 *             [
 *                 "125",
 *                 "1"
 *             ]
 *         ]
 *     }
 * }
 *
 */
func LoadDepthH(s *web.Session) web.Result {
	var symbol string
	var max int = 8
	var err = s.ValidFormat(`
		symbol,R|S,L:0;
		max,O|I,R:0;
	`, &symbol, &max)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	depth := market.LoadDepth(symbol, max)
	return s.SendJSON(xmap.M{
		"code":  0,
		"depth": depth,
	})
}
