package gexdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codingeasygo/crud"
	"github.com/gexservice/gexservice/base/define"
	"github.com/shopspring/decimal"
)

func FindWithdrawByOrderIDCall(caller crud.Queryer, ctx context.Context, orderID string, lock bool) (withdraw *Withdraw, err error) {
	orderIDInt, _ := strconv.ParseInt(orderID, 10, 64)
	querySQL := crud.QuerySQL(&Withdraw{}, "#all")
	querySQL, args := crud.JoinWheref(querySQL, nil, "tid=$%v,order_id=$%v#+or", orderIDInt, orderID)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Order{}, "#all", querySQL, args, &withdraw)
	return
}

func CreateWithdraw(ctx context.Context, userID int64, asset string, quantity decimal.Decimal) (withdraw *Withdraw, err error) {
	tx, err := Pool().Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit(ctx)
		} else {
			tx.Rollback(ctx)
		}
	}()
	balance := &Balance{
		UserID: userID,
		Asset:  asset,
		Free:   decimal.Zero.Sub(quantity),
		Locked: quantity,
	}
	err = IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		return
	}
	withdraw = &Withdraw{
		OrderID:  NewOrderID(),
		Type:     WithdrawTypeWithdraw,
		UserID:   userID,
		Creator:  userID,
		Asset:    asset,
		Quantity: quantity,
		Status:   WithdrawStatusPending,
	}
	err = AddWithdrawCall(tx, ctx, withdraw)
	return
}

func CancelWithdraw(ctx context.Context, userID int64, orderID string) (withdraw *Withdraw, err error) {
	tx, err := Pool().Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit(ctx)
		} else {
			tx.Rollback(ctx)
		}
	}()
	withdraw, err = FindWithdrawByOrderIDCall(tx, ctx, orderID, true)
	if err != nil {
		return
	}
	if userID > 0 && withdraw.UserID != userID {
		err = define.ErrNotAccess
		return
	}
	if withdraw.Type != WithdrawTypeWithdraw {
		err = fmt.Errorf("order is withdraw")
		return
	}
	if withdraw.Status != WithdrawStatusPending {
		err = fmt.Errorf("order is not pending")
		return
	}
	withdraw.Status = WithdrawStatusCanceled
	free := withdraw.Quantity
	balance := &Balance{
		UserID: withdraw.UserID,
		Asset:  withdraw.Asset,
		Free:   free,
		Locked: decimal.Zero.Sub(free),
	}
	err = IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		return
	}
	err = withdraw.UpdateFilter(tx, ctx, "status")
	return
}

// const (
// 	WithdrawVerifyPending  = 100
// 	WithdrawVerifyFail     = 200
// 	WithdrawVerifyNotFound = 210
// 	WithdrawVerifyDone     = 300
// )

// var WithdrawVerifyOrder = func(orderID string) (result int, info xmap.M, err error) {
// 	panic("not impl")
// }

// var WithdrawApplyOrder = func(user *User, order *Order) (result int, info xmap.M, err error) {
// 	panic("not impl")
// }

// var ProcWithdrawApplyDelay = 5 * time.Minute

// func ProcWithdrawApply() (err error) {
// 	defer func() {
// 		if perr := recover(); perr != nil {
// 			xlog.Errorf("ProcWithdrawApply proc withdraw apply is panic with %v, callstatck is \n%v", perr, debug.CallStatck())
// 		}
// 		if err != nil && err != pgx.ErrNoRows {
// 			xlog.Errorf("ProcWithdrawApply proc withdraw apply fail with %v", err)
// 		}
// 	}()
// 	var orderID string
// 	var oldStatus int
// 	updateSQL := `
// 		update exs_order set withdraw_status=$5,withdraw_next=$6,status=$7
// 		from (select tid,order_id,status from exs_order where type=$1 and status=any($2) and withdraw_status=any($3) and withdraw_next<$4 order by update_time asc limit 1) o
// 		where exs_order.tid=o.tid
// 		returning o.order_id,o.status
// 	`
// 	err = Pool().QueryRow(
// 		updateSQL,
// 		OrderTypeWithdraw, xsql.IntArray([]int{OrderStatusPending, OrderStatusApply}).DbArray(), xsql.IntArray([]int{OrderWithdrawStatusPending, OrderWithdrawStatusApply}).DbArray(), time.Now(),
// 		OrderWithdrawStatusApply, time.Now().Add(ProcWithdrawApplyDelay), OrderStatusApply,
// 	).Scan(&orderID, &oldStatus)
// 	if err != nil {
// 		if err != pgx.ErrNoRows {
// 			xlog.Errorf("ProcWithdrawApply query withdraw order fail with %v", err)
// 		}
// 		return
// 	}
// 	order, err := FindOrderByOrderID(orderID)
// 	if err != nil {
// 		err = fmt.Errorf("find order fail with %v", err)
// 		return
// 	}
// 	user, err := FindUser(order.UserID)
// 	if err != nil {
// 		err = fmt.Errorf("find user fail with %v", err)
// 		return
// 	}
// 	var result int
// 	var info xmap.M
// 	if oldStatus == OrderStatusApply {
// 		xlog.Infof("ProcWithdrawApply start verify withdraw order by %v", orderID)
// 		result, info, err = WithdrawVerifyOrder(orderID)
// 	} else {
// 		xlog.Infof("ProcWithdrawApply start apply withdraw order by %v", converter.JSON(order))
// 		result, info, err = WithdrawApplyOrder(user, order)
// 	}
// 	if err != nil {
// 		xlog.Errorf("ProcWithdrawApply call server fail with err:%v,result:%v,info:%v", err, result, converter.JSON(info))
// 		return
// 	}
// 	switch result {
// 	case WithdrawVerifyPending:
// 		xlog.Infof("ProcWithdrawApply apply order(%v) is pending by result:%v,info:%v", orderID, result, converter.JSON(info))
// 		err = UpdateOrderPrepay(order.TID, xsql.M(info))
// 	case WithdrawVerifyFail, WithdrawVerifyNotFound:
// 		xlog.Warnf("ProcWithdrawApply apply order(%v) is fail by result:%v,info:%v, will retry next time", orderID, result, converter.JSON(info))
// 		err = Pool().ExecRow(`update exs_order set withdraw_status=$2,withdraw_next=$3,status=$4 where tid=$1`, order.TID, OrderWithdrawStatusPending, time.Now().Add(time.Minute), OrderStatusPending)
// 	case WithdrawVerifyDone:
// 		xlog.Infof("ProcWithdrawApply apply order(%v) is done by result:%v,info:%v", orderID, result, converter.JSON(info))
// 		err = Pool().ExecRow(`update exs_order set filled=quantity,out_filled=quantity,withdraw_status=$2,notify_result=$3,status=$4 where tid=$1`, order.TID, OrderWithdrawStatusDone, xsql.M(info), OrderStatusDone)
// 	default:
// 		err = fmt.Errorf("unknow result code(%v)", result)
// 	}
// 	return
// }
