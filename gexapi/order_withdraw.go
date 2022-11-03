package gexapi

import (
	"github.com/codingeasygo/web"
)

//CreateWithdrawOrderH is http handler
/**
 *
 * @api {GET} /usr/createWithdrawOrder Create Withdraw Order
 * @apiName CreateWithdrawOrder
 * @apiGroup Order
 *
 * @apiParam  {Number} quantity the withdraw quantity
 * @apiParam  {String} password the trade password
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>, 8000 is trade password error, 9000 is balance not enought
 * @apiSuccess (Success) {Object} order the created order info
 * @apiSuccess (Success) {Number} order.tid the int order id
 * @apiSuccess (Success) {String} order.type the order type, all type supported is <a href="#metadata-Order">OrderTypeAll</a>
 * @apiSuccess (Success) {String} order.order_id the string order id
 * @apiSuccess (Success) {String} order.transaction.code the order verify code
 * @apiSuccess (Success) {Number} order.create_time the order create time
 * @apiSuccess (Success) {Number} order.update_time the order update time
 * @apiSuccess (Success) {Number} order.status the order status, all status supported is <a href="#metadata-Order">OrderStatusAll</a>
 *
 * @apiParamExample  {Query} CreateTopupOrder:
 * quantity=1&password=11
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "order": {
 *         "avg_price": "0",
 *         "create_time": 1632668243553,
 *         "creator": 100005,
 *         "fee_balance": "",
 *         "fee_filled": "0",
 *         "filled": "0",
 *         "in_balance": "",
 *         "in_filled": "0",
 *         "notify_result": {},
 *         "order_id": "202109262257230100001",
 *         "out_balance": "YWE",
 *         "out_filled": "0",
 *         "prepay_result": {},
 *         "price": "0",
 *         "quantity": "1",
 *         "side": "",
 *         "status": 100,
 *         "tid": 1000,
 *         "total_price": "0",
 *         "transaction": {
 *             "code": "5E48A8"
 *         },
 *         "type": 400,
 *         "update_time": 1632668243553,
 *         "user_id": 100005
 *     }
 * }
 *
 */
func CreateWithdrawOrderH(s *web.Session) web.Result {
	// var quantity float64
	// var password string
	// err := s.ValidFormat(`
	// 	quantity,R|I,R:0;
	// 	password,R|S,L:0;
	// `, &quantity, &password)
	// if err != nil {
	// 	return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	// }
	// userID := s.Int64("user_id")
	// err = gexdb.UserVerifyTradePassword(userID, xhash.SHA1([]byte(password)))
	// if err != nil {
	// 	xlog.Errorf("CreateWithdrawOrderH verify trad pass fail with %v", err)
	// 	code := define.ServerError
	// 	if err == pgx.ErrNoRows {
	// 		code = 8000
	// 	}
	// 	return util.ReturnCodeLocalErr(s, code, "srv-err", err)
	// }
	// user, err := gexdb.FindUser(s.R.Context(), userID)
	// if err != nil {
	// 	xlog.Errorf("CreateWithdrawOrderH find user fail with %v", err)
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// if user.KbzOpenid == nil || len(*user.KbzOpenid) < 1 {
	// 	err = fmt.Errorf("not kbzpay user")
	// 	xlog.Errorf("CreateWithdrawOrderH check user fail with %v", err)
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// order, err := gexdb.CreateWithdrawOrder(userID, decimal.NewFromFloat(quantity).Truncate(2))
	// if err != nil {
	// 	xlog.Errorf("CreateWithdrawOrderH create withdraw order fail with %v", err)
	// 	code := define.ServerError
	// 	if err == gexdb.ErrBalanceNotEnought {
	// 		code = 9000
	// 	}
	// 	return util.ReturnCodeLocalErr(s, code, "srv-err", err)
	// }
	// xlog.Infof("CreateWithdrawOrderH create order from %v success with %v", s.R.RemoteAddr, converter.JSON(order))
	// return s.SendJSON(xmap.M{
	// 	"code":  0,
	// 	"order": order,
	// })
	return web.Return
}

//CancelWithdrawOrderH is http handler
/**
 *
 * @api {GET} /usr/cancelWithdrawOrder Cancel Withdraw Order
 * @apiName CancelWithdrawOrder
 * @apiGroup Order
 *
 * @apiParam  {String} order_id the withdraw order id
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} order the created order info
 * @apiSuccess (Success) {Number} order.tid the int order id
 * @apiSuccess (Success) {String} order.type the order type, all type supported is <a href="#metadata-Order">OrderTypeAll</a>
 * @apiSuccess (Success) {String} order.order_id the string order id
 * @apiSuccess (Success) {String} order.transaction.code the order verify code
 * @apiSuccess (Success) {Number} order.create_time the order create time
 * @apiSuccess (Success) {Number} order.update_time the order update time
 * @apiSuccess (Success) {Number} order.status the order status, all status supported is <a href="#metadata-Order">OrderStatusAll</a>
 *
 * @apiParamExample  {Query} CreateTopupOrder:
 * order_id=100
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "order": {
 *         "avg_price": "0",
 *         "create_time": 1632668243553,
 *         "creator": 100005,
 *         "fee_balance": "",
 *         "fee_filled": "0",
 *         "filled": "0",
 *         "in_balance": "",
 *         "in_filled": "0",
 *         "notify_result": {},
 *         "order_id": "202109262257230100001",
 *         "out_balance": "YWE",
 *         "out_filled": "0",
 *         "prepay_result": {},
 *         "price": "0",
 *         "quantity": "1",
 *         "side": "",
 *         "status": 100,
 *         "tid": 1000,
 *         "total_price": "0",
 *         "transaction": {
 *             "code": "5E48A8"
 *         },
 *         "type": 400,
 *         "update_time": 1632668243553,
 *         "user_id": 100005
 *     }
 * }
 *
 */
func CancelWithdrawOrderH(s *web.Session) web.Result {
	// var orderID string
	// err := s.ValidFormat(`
	// 	order_id,R|S,L:0;
	// `, &orderID)
	// if err != nil {
	// 	return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	// }
	// userID := s.Int64("user_id")
	// targetUserID := userID
	// user, err := gexdb.FindUser(userID)
	// if err != nil {
	// 	xlog.Errorf("CancelWithdrawOrderH find target user(%v) err: %v", userID, err)
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// if user.Type == gexdb.UserTypeAdmin {
	// 	targetUserID = 0
	// }
	// order, err := gexdb.CancelWithdrawOrder(targetUserID, orderID)
	// if err != nil {
	// 	xlog.Errorf("CancelWithdrawOrderH cancel withdraw order by user:%v fail with %v", targetUserID, err)
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// xlog.Infof("CancelWithdrawOrderH cancel order by user %v from %v success with %v", targetUserID, s.R.RemoteAddr, converter.JSON(order))
	// return s.SendJSON(xmap.M{
	// 	"code":  0,
	// 	"order": order,
	// })
	return web.Return
}
