package gexdb

// import (
// 	"github.com/codingeasygo/util/converter"
// 	"github.com/codingeasygo/util/xmap"
// 	"github.com/shopspring/decimal"
// 	"github.com/gexservice/gexservice/base/xlog"
// )

// func PayTopupOrder(orderID string, totalAmount decimal.Decimal, result xmap.M) (err error) {
// 	tx, err := Pool().Begin()
// 	if err != nil {
// 		return
// 	}
// 	defer func() {
// 		if err == nil {
// 			err = tx.Commit()
// 		} else {
// 			tx.Rollback()
// 		}
// 	}()
// 	order := Order{OrderID: orderID, Filled: totalAmount, InFilled: totalAmount}
// 	err = tx.QueryRow(`select tid,type,user_id,in_balance,status from exs_order where order_id=$1 for update`, orderID).Scan(&order.TID, &order.Type, &order.UserID, &order.InBalance, &order.Status)
// 	if err != nil {
// 		xlog.Errorf("PayOrder receive pay notify fail with lock order %v, notify result is %v", err, converter.JSON(result))
// 		return
// 	}
// 	if order.Status != OrderStatusPending {
// 		xlog.Warnf("PayOrder receive pay notify fail with order status is %v, skipped notify result is %v", order.Status, converter.JSON(result))
// 		return
// 	}
// 	if order.Type != OrderTypeTopup {
// 		xlog.Warnf("PayOrder receive pay notify fail with order type is %v, skipped notify result is %v", order.Type, converter.JSON(result))
// 		return
// 	}
// 	err = tx.ExecRow(`update exs_order set filled=$1,in_filled=$1,status=$2,notify_result=$3 where tid=$4`, order.InFilled, OrderStatusDone, converter.JSON(result), order.TID)
// 	if err != nil {
// 		xlog.Errorf("PayOrder receive pay notify fail with update order %v, notify result is %v", err, converter.JSON(result))
// 		return
// 	}
// 	balance := &Balance{
// 		UserID: order.UserID,
// 		Asset:  order.InBalance,
// 		Free:   order.Filled,
// 		Status: BalanceStatusNormal,
// 	}
// 	err = IncreaseBalanceCall(tx, balance)
// 	if err != nil {
// 		xlog.Errorf("PayOrder receive pay notify fail with increase blance %v, notify result is %v", err, converter.JSON(result))
// 		return
// 	}
// 	xlog.Infof("PayOrder receive pay notify on order %v success by increase %v %v to user %v", order.OrderID, order.InFilled, order.InBalance, order.UserID)
// 	return
// }
