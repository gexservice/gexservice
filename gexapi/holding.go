package gexapi

import (
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/market"
	"github.com/gexservice/gexservice/matcher"
	"github.com/shopspring/decimal"
)

//ListHoldingH is http handler
/**
 *
 * @api {GET} /usr/listHolding List Holding
 * @apiName ListHolding
 * @apiGroup Balance
 *
 *
 * @apiParam  {String} [symbol] the symbol to list
 * @apiParam  {Number} [target_user_id] the target user id to list holding
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
	var symbolOnly []string
	var targetUserID int64
	err := s.ValidFormat(`
		symbol,o|s,l:0;
		target_user_id,o|i,r:0;
	`, &symbolOnly, &targetUserID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Value("user_id").(int64)
	if targetUserID > 0 && AdminAccess(s) {
		userID = targetUserID
	}
	_, err = gexdb.TouchBalance(s.R.Context(), gexdb.BalanceAreaFutures, []string{matcher.Quote}, userID)
	if err != nil {
		xlog.Errorf("ListHoldingH touch balance by %v,%v fail with %v", matcher.Quote, userID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	balance, err := gexdb.FindBalanceByAsset(s.R.Context(), userID, gexdb.BalanceAreaFutures, matcher.Quote)
	if err != nil {
		xlog.Errorf("ListHoldingH find balance by %v,%v fail with %v", userID, matcher.Quote, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	holdings, symbols, err := gexdb.ListUserHolding(s.R.Context(), userID, symbolOnly)
	if err != nil {
		xlog.Errorf("ListHoldingH list holding by %v fail with %v", userID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	var symbolInfoes map[string]*matcher.SymbolInfo
	if len(symbols) > 0 {
		_, symbolInfoes, _ = market.ListSymbol("", symbols, "")
	}
	var unprofits map[string]decimal.Decimal
	var tickers map[string]*gexdb.Ticker
	if len(holdings) > 0 {
		unprofits, tickers = market.CalcHoldingUnprofit(s.R.Context(), holdings...)
	}
	return s.SendJSON(xmap.M{
		"code":      0,
		"balance":   balance,
		"holdings":  holdings,
		"unprofits": unprofits,
		"symbols":   symbolInfoes,
		"tickers":   tickers,
	})
}

//LoadHoldingH is http handler
/**
 *
 * @api {GET} /usr/loadHolding Load Holding
 * @apiName LoadHolding
 * @apiGroup Balance
 *
 * @apiParam  {String} symbol the symbol to add
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Holding) {Object} holding the user holding info
 * @apiUse HoldingObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "holding": {
 *         "amount": "0",
 *         "blowup": "0",
 *         "create_time": 1668413078999,
 *         "lever": 5,
 *         "margin_added": "0",
 *         "margin_used": "0",
 *         "open": "0",
 *         "status": 100,
 *         "symbol": "XX",
 *         "tid": 1002,
 *         "update_time": 1668413078999,
 *         "user_id": 100002
 *     }
 * }
 */
func LoadHoldingH(s *web.Session) web.Result {
	var symbol string
	err := s.ValidFormat(`
		symbol,r|s,l:0;
	`, &symbol)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Value("user_id").(int64)
	_, err = gexdb.TouchHolding(s.R.Context(), []string{symbol}, userID)
	if err != nil {
		xlog.Errorf("LoadHoldingH touch holding by %v,%v fail with %v", symbol, userID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	holding, err := gexdb.FindHoldlingBySymbol(s.R.Context(), userID, symbol)
	if err != nil {
		xlog.Errorf("LoadHoldingH find holding by %v,%v fail with %v", symbol, userID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	symbolInfo, day := market.LoadSymbol(symbol)
	return s.SendJSON(xmap.M{
		"code":    0,
		"holding": holding,
		"symbol":  symbolInfo,
		"day":     day,
	})
}

//ChangeHoldingLeverH is http handler
/**
 *
 * @api {GET} /usr/changeHoldingLever Change Holding Lever
 * @apiName ChangeHoldingLever
 * @apiGroup Balance
 *
 * @apiParam  {String} symbol the symbol to change lever
 * @apiParam  {Number} lever the new lever to change, must be 0<lever<100
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a> or <a href="#metadata-ExReturnCode">ExReturnCode</a>
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0
 * }
 */
func ChangeHoldingLeverH(s *web.Session) web.Result {
	var symbol string
	var lever int
	err := s.ValidFormat(`
		symbol,r|s,l:0;
		lever,r|i,r:0~100;
	`, &symbol, &lever)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Value("user_id").(int64)
	err = matcher.ChangeLever(s.R.Context(), userID, symbol, lever)
	if err != nil {
		xlog.Errorf("ChangeHoldingLeverH change holding lever by %v,%v,%v fail with %v", userID, symbol, lever, err)
		code, ok := matcher.IsErrCode(err)
		if !ok {
			code = define.ServerError
		}
		return util.ReturnCodeLocalErr(s, code, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//SearchHoldingH is http handler
/**
 *
 * @api {GET} /usr/searchHolding Search Holding
 * @apiName SearchHolding
 * @apiGroup Balance
 *
 *
 * @apiUse HoldingUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} unprofits the balance unprofit, mapping by user id, then mapping by symbol
 * @apiSuccess (Holding) {Array} holdings the holding records
 * @apiUse HoldingObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "holdings": [
 *         {
 *             "amount": "0",
 *             "blowup": "0",
 *             "create_time": 1670327985918,
 *             "lever": 1,
 *             "margin_added": "0",
 *             "margin_used": "0",
 *             "open": "0",
 *             "status": 100,
 *             "symbol": "futures.YWEUSDT",
 *             "tid": 1004,
 *             "update_time": 1670327985918,
 *             "user_id": 100005
 *         }
 *     ],
 *     "total": 1,
 *     "unprofits": {},
 *     "users": {
 *         "100005": {
 *             "account": "abc3",
 *             "create_time": 1670327985411,
 *             "favorites": {},
 *             "image": "abc3_image",
 *             "name": "abc3_name",
 *             "phone": "abc3_123",
 *             "role": 100,
 *             "status": 100,
 *             "tid": 100005,
 *             "type": 100,
 *             "update_time": 1670327985411
 *         }
 *     }
 * }
 */
func SearchHoldingH(s *web.Session) web.Result {
	searcher := &gexdb.HoldingUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	if !AdminAccess(s) {
		searcher.Where.UserID = userID
	}
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Errorf("SearchHoldingH search holding fail with %v by %v", err, converter.JSON(searcher))
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	var users map[int64]*gexdb.User
	var unprofits map[int64]map[string]decimal.Decimal
	if len(searcher.Query.UserIDs) > 0 {
		unprofits, err = market.ListHoldingUnprofit(s.R.Context(), searcher.Query.UserIDs...)
		if err != nil {
			xlog.Errorf("SearchHoldingH list holding profits fail with %v by %v", err, converter.JSON(searcher))
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
		_, users, err = gexdb.ListUserByID(s.R.Context(), searcher.Query.UserIDs...)
		if err != nil {
			xlog.Errorf("SearchHoldingH list holding user fail with %v by %v", err, converter.JSON(searcher))
			return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
		}
	}
	return s.SendJSON(xmap.M{
		"code":      define.Success,
		"holdings":  searcher.Query.Holdings,
		"unprofits": unprofits,
		"users":     users,
		"total":     searcher.Count.Total,
	})
}
