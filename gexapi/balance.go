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
	"github.com/shopspring/decimal"
)

//LoadBalanceOverviewH is http handler
/**
 *
 * @api {GET} /usr/loadBalanceOverview Load Balance Overview
 * @apiName LoadBalanceOverview
 * @apiGroup Balance
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} total_value the user total estimated value by quote
 * @apiSuccess (Success) {Array} area_values the user all area estimated value
 * @apiSuccess (Success) {BalanceArea} area_values.area the blalance area, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiSuccess (Success) {Decimal} area_values.value the blalance area area estimated value
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "area_values": [
 *         {
 *             "area": 100,
 *             "value": "0"
 *         },
 *         {
 *             "area": 200,
 *             "value": "2000"
 *         },
 *         {
 *             "area": 300,
 *             "value": "0"
 *         }
 *     ],
 *     "code": 0,
 *     "total_value": "2000"
 * }
 */
func LoadBalanceOverviewH(s *web.Session) web.Result {
	userID := s.Value("user_id").(int64)
	totalValue, areaValues, err := market.CalcBalanceOverview(s.R.Context(), userID)
	if err != nil {
		xlog.Errorf("LoadBalanceOverviewH calc user %v balance overview fail with %v", userID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":        0,
		"total_value": totalValue,
		"area_values": areaValues,
	})
}

//ListBalanceH is http handler
/**
 *
 * @api {GET} /usr/listBalance List Balance
 * @apiName ListBalance
 * @apiGroup Balance
 *
 *
 * @apiParam  {Number} area the balance area to list, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} total_value the user total estimated value by quote
 * @apiSuccess (Balance) {Array} balances the user balance info
 * @apiUse BalanceObject
 * @apiSuccess (Success) {Object} values balance estimated value, mapping by key is balances.asset to value is estimated value
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "balances": [
 *         {
 *             "area": 200,
 *             "asset": "USDT",
 *             "create_time": 1667547245486,
 *             "free": "1000",
 *             "locked": "0",
 *             "margin": "0",
 *             "status": 100,
 *             "tid": 1005,
 *             "update_time": 1667547245498,
 *             "user_id": 100002
 *         },
 *         {
 *             "area": 200,
 *             "asset": "YWE",
 *             "create_time": 1667547245486,
 *             "free": "999",
 *             "locked": "1",
 *             "margin": "0",
 *             "status": 100,
 *             "tid": 1004,
 *             "update_time": 1667547245506,
 *             "user_id": 100002
 *         }
 *     ],
 *     "code": 0,
 *     "total_value": "2000",
 *     "values": {
 *         "USDT": "1000",
 *         "YWE": "0"
 *     }
 * }
 */
func ListBalanceH(s *web.Session) web.Result {
	var area gexdb.BalanceArea
	err := s.ValidFormat(`
		area,R|I,e:;
	`, &area)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Value("user_id").(int64)
	totalValue, todayWinned, balances, values, err := market.CalcBalanceTotalValue(s.R.Context(), userID, area)
	if err != nil {
		xlog.Errorf("ListBalanceH calc user %v %v balance overview fail with %v", userID, area, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":         0,
		"total_value":  totalValue,
		"today_winned": todayWinned,
		"balances":     balances,
		"values":       values,
	})
}

//TransferBalanceH is http handler
/**
 *
 * @api {GET} /usr/transferBalance Transfer Balance
 * @apiName TransferBalance
 * @apiGroup Balance
 *
 * @apiParam  {Number} from the balance transfer from area, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiParam  {Number} to the balance transfer to area, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiParam  {String} asset the balance asset to change
 * @apiParam  {Number} value the transfer value
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0
 * }
 */
func TransferBalanceH(s *web.Session) web.Result {
	var from, to gexdb.BalanceArea
	var asset string
	var value decimal.Decimal
	err := s.ValidFormat(`
		from,r|i,e:;
		to,r|i,e:;
		asset,r|s,l:0;
		value,r|f,n:;
	`, &from, &to, &asset, &value)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	err = gexdb.TransferChange(s.R.Context(), userID, userID, from, to, asset, value)
	if err != nil {
		xlog.Errorf("ChangeBalanceH change balance by %v,%v,%v,%v fail with %v", from, to, asset, value, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//ChangeUserBalanceH is http handler
/**
 *
 * @api {GET} /admin/changeUserBalance Change User Balance
 * @apiName ChangeUserBalance
 * @apiGroup Balance
 *
 * @apiParam  {Number} user_id the user id to change balance
 * @apiParam  {Number} area the balance area to change, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiParam  {String} asset the balance asset to change
 * @apiParam  {Number} changed the value to increase or descrease
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0
 * }
 */
func ChangeUserBalanceH(s *web.Session) web.Result {
	var userID int64
	var area gexdb.BalanceArea
	var asset string
	var changed decimal.Decimal
	err := s.ValidFormat(`
		user_id,r|i,r:0;
		area,r|i,e:;
		asset,r|s,l:0;
		changed,r|f,n:;
	`, &userID, &area, &asset, &changed)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	creator := s.Int64("user_id")
	_, err = gexdb.ChangeBalance(s.R.Context(), creator, userID, area, asset, changed)
	if err != nil {
		xlog.Errorf("ChangeBalanceH change balance by %v,%v,%v fail with %v", area, asset, changed, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//ListBalanceRecordH is http handler
/**
 *
 * @api {GET} /usr/listBalanceRecord List Balance Record
 * @apiName ListBalanceRecord
 * @apiGroup Balance
 *
 *
 * @apiUse BalanceRecordUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (BalanceRecordItem) {Array} records the balance records
 * @apiUse BalanceRecordItemObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "records": [
 *         {
 *             "asset": "USDT",
 *             "changed": "0.1",
 *             "tid": 0,
 *             "type": 100,
 *             "update_time": 1667873432495
 *         }
 *     ],
 *     "total": 1
 * }
 */
func ListBalanceRecordH(s *web.Session) web.Result {
	searcher := &gexdb.BalanceRecordUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	searcher.Where.UserID = userID
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Errorf("SearchOrderH searcher order fail with %v by %v", err, converter.JSON(searcher))
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":    define.Success,
		"records": searcher.Query.Records,
		"total":   searcher.Count.Total,
	})
}
