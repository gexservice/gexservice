package gexapi

import (
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
	"github.com/shopspring/decimal"
)

//CreateWithdrawH is http handler
/**
 *
 * @api {GET} /usr/createWithdraw Create Withdraw
 * @apiName CreateWithdraw
 * @apiGroup Withdraw
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
 *     "code": 0,
 *     "withdraw": {
 *         "asset": "USDT",
 *         "create_time": 1667896705684,
 *         "creator": 100004,
 *         "order_id": "202211081638250100013",
 *         "quantity": "1",
 *         "status": 100,
 *         "tid": 1006,
 *         "type": 100,
 *         "update_time": 1667896705684,
 *         "user_id": 100004
 *     }
 * }
 */
func CreateWithdrawH(s *web.Session) web.Result {
	var asset string
	var quantity decimal.Decimal
	err := s.ValidFormat(`
		asset,r|s,l:0;
		quantity,r|f,r:0;
	`, &asset, &quantity)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	withdraw, err := gexdb.CreateWithdraw(s.R.Context(), userID, asset, quantity)
	if err != nil {
		xlog.Errorf("CreateWithdrawH create withdraw by %v,%v,%v fail with %v", userID, asset, quantity, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     define.Success,
		"withdraw": withdraw,
	})
}

//CancelWithdrawH is http handler
/**
 *
 * @api {GET} /usr/cancelWithdraw Cancel Withdraw
 * @apiName CancelWithdraw
 * @apiGroup Withdraw
 *
 * @apiParam  {String} order_id the withdraw order id
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0
 * }
 */
func CancelWithdrawH(s *web.Session) web.Result {
	var orderID string
	err := s.ValidFormat(`
		order_id,r|s,l:0;
	`, &orderID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	withdraw, err := gexdb.CancelWithdraw(s.R.Context(), userID, orderID)
	if err != nil {
		xlog.Errorf("CreateWithdrawH create withdraw by %v,%v fail with %v", userID, orderID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     define.Success,
		"withdraw": withdraw,
	})
}

//ListWithdrawH is http handler
/**
 *
 * @api {GET} /usr/listWithdraw List Withdraw
 * @apiName ListWithdraw
 * @apiGroup Withdraw
 *
 *
 * @apiUse WithdrawUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Withdraw) {Array} withdraws the withdraw records
 * @apiUse WithdrawObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "total": 1,
 *     "withdraws": [
 *         {
 *             "asset": "USDT",
 *             "create_time": 1667896770794,
 *             "creator": 100004,
 *             "order_id": "202211081639300100013",
 *             "quantity": "1",
 *             "status": 100,
 *             "tid": 1006,
 *             "type": 100,
 *             "update_time": 1667896770794,
 *             "user_id": 100004
 *         }
 *     ]
 * }
 */
func ListWithdrawH(s *web.Session) web.Result {
	searcher := &gexdb.WithdrawUnifySearcher{}
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
		"code":      define.Success,
		"withdraws": searcher.Query.Withdraws,
		"total":     searcher.Count.Total,
	})
}
