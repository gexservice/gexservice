package gexapi

import (
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/gexservice/gexservice/market"
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
 * @apiSuccess (Success) {Array} area_values.area the blalance area, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiSuccess (Success) {Array} area_values.value the blalance area area estimated value
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "balances": {
 *         "MMK": {
 *             "asset": "MMK",
 *             "create_time": 1632661193389,
 *             "free": "100",
 *             "locked": "0",
 *             "status": 100,
 *             "tid": 1005,
 *             "update_time": 1632661193391,
 *             "user_id": 100002
 *         },
 *         "YWE": {
 *             "asset": "YWE",
 *             "create_time": 1632661193389,
 *             "free": "100",
 *             "locked": "0",
 *             "status": 100,
 *             "tid": 1004,
 *             "update_time": 1632661193390,
 *             "user_id": 100002
 *         }
 *     },
 *     "code": 0,
 *     "estimated": {
 *         "all_balance": "MMK",
 *         "all_daily": "0",
 *         "all_free": "10100"
 *     },
 *     "user": {
 *         "account": "abc0",
 *         "broker_id": 0,
 *         "create_time": 1632661193387,
 *         "external": {},
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1632661193387
 *     }
 * }
 *
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
 *     "balances": {
 *         "MMK": {
 *             "asset": "MMK",
 *             "create_time": 1632661193389,
 *             "free": "100",
 *             "locked": "0",
 *             "status": 100,
 *             "tid": 1005,
 *             "update_time": 1632661193391,
 *             "user_id": 100002
 *         },
 *         "YWE": {
 *             "asset": "YWE",
 *             "create_time": 1632661193389,
 *             "free": "100",
 *             "locked": "0",
 *             "status": 100,
 *             "tid": 1004,
 *             "update_time": 1632661193390,
 *             "user_id": 100002
 *         }
 *     },
 *     "code": 0,
 *     "estimated": {
 *         "all_balance": "MMK",
 *         "all_daily": "0",
 *         "all_free": "10100"
 *     },
 *     "user": {
 *         "account": "abc0",
 *         "broker_id": 0,
 *         "create_time": 1632661193387,
 *         "external": {},
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1632661193387
 *     }
 * }
 *
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
	totalValue, balances, values, err := market.CalcBalanceTotalValue(s.R.Context(), userID, area)
	if err != nil {
		xlog.Errorf("ListBalanceH calc user %v %v balance overview fail with %v", userID, area, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":        0,
		"total_value": totalValue,
		"balances":    balances,
		"values":      values,
	})
}
