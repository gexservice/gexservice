package gexapi

import (
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

//CreateWithdrawH is http handler
/**
 *
 * @api {GET} /usr/createWithdraw Create Withdraw
 * @apiName CreateWithdraw
 * @apiGroup Withdraw
 *
 * @apiUse WithdrawUpdate
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Withdraw) {Object} withdraw the withdraw info
 * @apiUse WithdrawObject
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
	var withdraw gexdb.Withdraw
	err := s.Valid(&withdraw, "method,asset,quantity,receiver#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	withdraw.UserID = userID
	withdraw.Creator = userID
	err = gexdb.CreateWithdraw(s.R.Context(), &withdraw)
	if err != nil {
		xlog.Errorf("CreateWithdrawH create withdraw by %v fail with %v", converter.JSON(withdraw), err)
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

//ConfirmWithdrawH is http handler
/**
 *
 * @api {GET} /usr/confirmWithdraw Confirm Withdraw
 * @apiName ConfirmWithdraw
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
func ConfirmWithdrawH(s *web.Session) web.Result {
	var orderID string
	err := s.ValidFormat(`
		order_id,r|s,l:0;
	`, &orderID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	userID := s.Int64("user_id")
	err = gexdb.ConfirmWithdraw(s.R.Context(), orderID)
	if err != nil {
		xlog.Errorf("ConfirmWithdrawH confirm withdraw by %v,%v fail with %v", userID, orderID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
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
	if !AdminAccess(s) {
		userID := s.Int64("user_id")
		searcher.Where.UserID = userID
	}
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

//CreateGoldbarH is http handler
/**
 *
 * @api {GET} /usr/createGoldbar Create Goldbar
 * @apiName CreateGoldbar
 * @apiGroup Withdraw
 *
 * @apiUse WithdrawUpdate
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Withdraw) {Object} goldbar the withdraw info
 * @apiUse WithdrawObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "goldbar": {
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
func CreateGoldbarH(s *web.Session) web.Result {
	var pickupAmount int64
	var pickupTime int64
	var pickupAddress string
	err := s.ValidFormat(`
		pickup_amount,R|I,R:0;
		pickup_time,R|I,R:0;
		pickup_address,R|S,L:0;
	`, &pickupAmount, &pickupTime, &pickupAddress)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	goldbar, err := gexdb.CreateGoldbar(s.R.Context(), userID, pickupAmount, pickupTime, pickupAddress)
	if err != nil {
		xlog.Errorf("CreateGoldbarH create goldbar by %v fail with %v", converter.JSON(goldbar), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":    define.Success,
		"goldbar": goldbar,
	})
}

//CancelGoldbarH is http handler
/**
 *
 * @api {GET} /usr/cancelGoldbar Cancel Goldbar
 * @apiName CancelGoldbar
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
func CancelGoldbarH(s *web.Session) web.Result {
	var orderID string
	err := s.ValidFormat(`
		order_id,r|s,l:0;
	`, &orderID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	withdraw, err := gexdb.CancelGoldbar(s.R.Context(), userID, orderID)
	if err != nil {
		xlog.Errorf("CreateGoldbarH create withdraw by %v,%v fail with %v", userID, orderID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     define.Success,
		"withdraw": withdraw,
	})
}

//ConfirmGoldbarH is http handler
/**
 *
 * @api {GET} /usr/confirmGoldbar Confirm Goldbar
 * @apiName ConfirmGoldbar
 * @apiGroup Goldbar
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
func ConfirmGoldbarH(s *web.Session) web.Result {
	var orderID string
	err := s.ValidFormat(`
		order_id,r|s,l:0;
	`, &orderID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	userID := s.Int64("user_id")
	err = gexdb.ConfirmGoldbar(s.R.Context(), orderID)
	if err != nil {
		xlog.Errorf("ConfirmGoldbarH confirm withdraw by %v,%v fail with %v", userID, orderID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//DoneGoldbarH is http handler
/**
 *
 * @api {GET} /usr/doneGoldbar Done Goldbar
 * @apiName DoneGoldbar
 * @apiGroup Goldbar
 *
 * @apiParam  {String} order_id the goldbar order id
 * @apiParam  {String} code the goldbar code
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0
 * }
 */
func DoneGoldbarH(s *web.Session) web.Result {
	args := struct {
		OrderID string `json:"order_id"`
		Code    string `json:"code"`
		Result  xmap.M `json:"result"`
	}{}
	_, err := s.RecvJSON(&args)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	goldbar, err := gexdb.DoneGoldbar(s.R.Context(), args.OrderID, args.Code, args.Result)
	if err != nil {
		xlog.Errorf("DoneGoldbarH done goldbar by %%v fail with %v", converter.JSON(goldbar), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":    define.Success,
		"goldbar": goldbar,
	})
}

//ListGoldbarH is http handler
/**
 *
 * @api {GET} /usr/listGoldbar List Goldbar
 * @apiName ListGoldbar
 * @apiGroup Withdraw
 *
 *
 * @apiUse WithdrawUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Withdraw) {Array} goldbars the withdraw records
 * @apiUse WithdrawObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "total": 1,
 *     "goldbars": [
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
func ListGoldbarH(s *web.Session) web.Result {
	searcher := &gexdb.WithdrawUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	searcher.Where.Type = gexdb.WithdrawTypeArray{gexdb.WithdrawTypeGoldbar}
	if !AdminAccess(s) {
		userID := s.Int64("user_id")
		searcher.Where.UserID = userID
	}
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Errorf("SearchOrderH searcher order fail with %v by %v", err, converter.JSON(searcher))
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     define.Success,
		"goldbars": searcher.Query.Withdraws,
		"total":    searcher.Count.Total,
	})
}

//LoadTopupAddressH is http handler
/**
 *
 * @api {GET} /usr/loadTopupAddress Load Topup Address
 * @apiName LoadTopupAddress
 * @apiGroup Withdraw
 *
 * @apiParam {WalletMethod} method the wallet method, all suported is <a href="#metadata-Wallet">WalletMethodAll</a>
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Wallet) {Object} wallet the wallet info
 * @apiUse WalletObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "wallet": {
 *         "address": "638cb19b285c660c48000002",
 *         "create_time": 1670164891118,
 *         "method": "tron",
 *         "status": 100,
 *         "tid": 1000,
 *         "update_time": 1670164891118,
 *         "user_id": 100004
 *     }
 * }
 */
func LoadTopupAddressH(s *web.Session) web.Result {
	var method gexdb.WalletMethod
	err := s.ValidFormat(`
		method,R|S,E:0;
	`, &method)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	wallet, err := gexdb.LoadWalletByMethod(s.R.Context(), userID, method)
	if err != nil {
		xlog.Errorf("LoadTopupAddressH load wallet by %v,%v fail with %v", userID, method, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":   define.Success,
		"wallet": wallet,
	})
}
