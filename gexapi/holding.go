package gexapi

import (
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/market"
	"github.com/gexservice/gexservice/matcher"
)

//ListHoldingH is http handler
/**
 *
 * @api {GET} /usr/listHolding List Holding
 * @apiName ListHolding
 * @apiGroup Balance
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Balance) {Object} balance the user balance info
 * @apiUse BalanceObject
 * @apiSuccess (Holding) {Array} holdings the user holding info
 * @apiUse HoldingObject
 * @apiSuccess (Ticker) {Array} tickers the symbol ticker info
 * @apiSuccess (Ticker) {Array} tickers.ask the ticker ask info by [price,qty]
 * @apiSuccess (Ticker) {Array} tickers.bid the ticker ask info by [price,qty]
 * @apiSuccess (Unprofit) {Object} unprofits the symbol unprofit info, mapping by symbol as key
 * @apiSuccess (Unprofit) {Object} unprofits.total the total unprofit
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "balance": {
 *         "area": 300,
 *         "asset": "USDT",
 *         "create_time": 1667736566986,
 *         "free": "899.8",
 *         "locked": "100.1",
 *         "margin": "50",
 *         "status": 100,
 *         "tid": 1014,
 *         "update_time": 1667736567069,
 *         "user_id": 100002
 *     },
 *     "code": 0,
 *     "holdings": [
 *         {
 *             "amount": "-0.5",
 *             "blowup": "199",
 *             "create_time": 1667736567041,
 *             "lever": 1,
 *             "margin_added": "0",
 *             "margin_used": "50",
 *             "open": "100",
 *             "status": 100,
 *             "symbol": "futures.YWEUSDT",
 *             "tid": 1000,
 *             "update_time": 1667736567041,
 *             "user_id": 100002
 *         }
 *     ],
 *     "tickers": {
 *         "futures.YWEUSDT": {
 *             "ask": [
 *                 "100",
 *                 "0.5"
 *             ],
 *             "bid": [
 *                 "90",
 *                 "0.5"
 *             ],
 *             "symbol": "futures.YWEUSDT"
 *         }
 *     },
 *     "unprofits": {
 *         "futures.YWEUSDT": "0",
 *         "total": "0"
 *     }
 * }
 */
func ListHoldingH(s *web.Session) web.Result {
	userID := s.Value("user_id").(int64)
	balance, err := gexdb.FindBalanceByAsset(s.R.Context(), userID, gexdb.BalanceAreaFutures, matcher.Quote)
	if err != nil {
		xlog.Errorf("ListHoldingH find balance by %v,%v fail with %v", userID, matcher.Quote, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	holdings, err := gexdb.ListUserHolding(s.R.Context(), userID)
	if err != nil {
		xlog.Errorf("ListHoldingH list holding by %v fail with %v", userID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	unprofits, tickers := market.CalcHoldingUnprofit(s.R.Context(), holdings...)
	return s.SendJSON(xmap.M{
		"code":      0,
		"balance":   balance,
		"holdings":  holdings,
		"unprofits": unprofits,
		"tickers":   tickers,
	})
}
