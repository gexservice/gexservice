package gexapi

import (
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
)

var CreatePrepayOrder = func(merchOrderID, title string, totalAmount float64, transCurrency string) (result xmap.M, err error) {
	result = xmap.New()
	return
}

//CreateTopupOrderH is http handler
/**
 *
 * @api {GET} /usr/createTopupOrder Create Topup Order
 * @apiName CreateTopupOrder
 * @apiGroup Order
 *
 * @apiParam  {Number} amount the top up amount
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} order the created order info
 * @apiSuccess (Success) {Number} order.tid the int order id
 * @apiSuccess (Success) {String} order.type the order type, all type supported is <a href="#metadata-Order">OrderTypeAll</a>
 * @apiSuccess (Success) {String} order.order_id the string order id
 * @apiSuccess (Success) {String} order.in_balance the in balance key
 * @apiSuccess (Success) {String} order.in_filled the in balance filled amount
 * @apiSuccess (Success) {Number} order.create_time the order create time
 * @apiSuccess (Success) {Number} order.update_time the order update time
 * @apiSuccess (Success) {Number} order.status the order status, all status supported is <a href="#metadata-Order">OrderStatusAll</a>
 * @apiSuccess (Success) {Object} order.prepay_result the order prepay result argument, it will be used to start pay on minapp
 *
 * @apiParamExample  {Query} CreateTopupOrder:
 * amount=100
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "order": {
 *         "avg_price": "0",
 *         "create_time": 1632572929348,
 *         "creator": 100002,
 *         "fee_balance": "",
 *         "fee_filled": "0",
 *         "filled": "100",
 *         "in_balance": "MMK",
 *         "in_filled": "100",
 *         "notify_result": null,
 *         "order_id": "202109252028490100001",
 *         "out_balance": "",
 *         "out_filled": "0",
 *         "prepay_result": {
 *             "order_info": "abc=1",
 *             "prepay_id": "KBZ00b9a7127afe8193e124e59460c7c4bf37102005062",
 *             "resp": {
 *                 "code": "0",
 *                 "merch_order_id": "613c276c285c666229000001",
 *                 "msg": "success",
 *                 "nonce_str": "MB7TQ7VMWBY0MBNBBEA7WGRNG3YSLQUJ",
 *                 "prepay_id": "KBZ00b9a7127afe8193e124e59460c7c4bf37102005062",
 *                 "result": "SUCCESS",
 *                 "sign": "BA4173F40D3789D3EE6929CBB663FA3400DDB35E4446634EAB4FFCF71121A892",
 *                 "sign_type": "SHA256"
 *             },
 *             "sign": "e8258d2fd2b54b8a9950e676707358d7b5efdf46da992738961c80165ddb62a6",
 *             "sign_type": "SHA256",
 *             "trade_type": "APPH5"
 *         },
 *         "price": "0",
 *         "quantity": "100",
 *         "status": 300,
 *         "tid": 1000,
 *         "total_price": "0",
 *         "transaction": {},
 *         "type": 100,
 *         "update_time": 1632572929348,
 *         "user_id": 100002
 *     }
 * }
 *
 */
func CreateTopupOrderH(s *web.Session) web.Result {
	// var amount float64
	// err := s.ValidFormat(`
	// 	amount,R|F,R:0;
	// `, &amount)
	// if err != nil {
	// 	return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	// }
	// totalAmount := decimal.NewFromFloat(amount)
	// userID := s.Int64("user_id")
	// order := &gexdb.Order{
	// 	Type:      gexdb.OrderTypeTopup,
	// 	UserID:    userID,
	// 	Creator:   userID,
	// 	Quantity:  totalAmount,
	// 	InBalance: gexdb.BalanceAssetMMK,
	// 	Status:    gexdb.OrderStatusPending,
	// }
	// err = gexdb.CreateOrder(order)
	// if err != nil {
	// 	xlog.Errorf("CreateTopupOrderH create db order fail with %v by %v", err, converter.JSON(order))
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// result, err := CreatePrepayOrder(order.OrderID, "Top UP", amount, order.InBalance)
	// if err != nil {
	// 	xlog.Errorf("CreateTopupOrderH create prepay order fail with %v by %v", err, converter.JSON(order))
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// order.PrepayResult = xsql.M(result)
	// err = gexdb.UpdateOrderPrepay(order.TID, order.PrepayResult)
	// if err != nil {
	// 	xlog.Errorf("CreateTopupOrderH update prepay to db order fail with %v by %v", err, converter.JSON(order))
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// xlog.Infof("CreateTopupOrderH create order from %v success with %v", s.R.RemoteAddr, converter.JSON(order))
	// return s.SendJSON(xmap.M{
	// 	"code":  0,
	// 	"order": order,
	// })
	return web.Return
}

//MockPayTopupOrderH is http handler
/**
 *
 * @api {GET} /mock/payTopupOrder Mock Pay Topup Order
 * @apiName MockPayTopupOrder
 * @apiGroup Mock
 *
 * @apiParam  {String} order_id the order id to mock pay
 * @apiParam  {Number} amount the pay amount to mock pay
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see common define
 *
 * @apiParamExample  {Query} CreateTopupOrder:
 * order_id=100&amount=100
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0
 * }
 *
 */
func MockPayTopupOrderH(s *web.Session) web.Result {
	// var orderID string
	// var amount float64
	// err := s.ValidFormat(`
	// 	order_id,R|S,L:0;
	// 	amount,R|F,R:0;
	// `, &orderID, &amount)
	// if err != nil {
	// 	return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	// }
	// err = gexdb.PayTopupOrder(orderID, decimal.NewFromFloat(amount), xmap.M{"mock": 1, "amount": amount, "remote": s.R.RemoteAddr})
	// if err != nil {
	// 	xlog.Errorf("MockPayTopupOrderH mock pay topup order fail with %v by orderID:%v,amount:%v", err, orderID, amount)
	// 	return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	// }
	// xlog.Warnf("MockPayTopupOrderH mock pay topup order from %v success with orderID:%v,amount:%v", s.R.RemoteAddr, orderID, amount)
	// return s.SendJSON(xmap.M{
	// 	"code": 0,
	// })
	return web.Return
}
