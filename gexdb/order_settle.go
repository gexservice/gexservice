package gexdb

// import (
// 	"time"

// 	"github.com/codingeasygo/crud/pgx"
// 	"github.com/shopspring/decimal"
// 	"github.com/gexservice/gexservice/base/basedb"
// 	"github.com/gexservice/gexservice/base/xlog"
// )

// func ProcSettleOrderFee() (err error) {
// 	var brokerCommRate float64
// 	err = basedb.LoadConf(ConfigBrokerCommRate, &brokerCommRate)
// 	if err != nil {
// 		xlog.Errorf("ProcSettleOrderFee get broker common rate config fail with %v", err)
// 		err = pgx.ErrNoRows
// 		return
// 	}
// 	var orderID int64
// 	updateSQL := `
// 		update gex_order set fee_settled_next=$1
// 		from (select tid from gex_order where type=$2 and status=any($3) and fee_settled_status=$4 and fee_settled_next<$5 order by update_time asc limit 1) o
// 		where gex_order.tid=o.tid
// 		returning o.tid
// 	`
// 	err = Pool().QueryRow(updateSQL, time.Now().Add(time.Minute), OrderTypeTrade, []int{OrderStatusPartcanceled, OrderStatusDone}, 0, time.Now()).Scan(&orderID)
// 	if err != nil {
// 		if err != pgx.ErrNoRows {
// 			xlog.Errorf("ProcSettleOrderFee query settle order fail with %v", err)
// 		}
// 		return
// 	}

// 	//
// 	var userRole int
// 	var brokerID int64
// 	var feeBalance string
// 	var feeFilled, feeBroker decimal.Decimal
// 	tx, err := Pool().Begin()
// 	if err != nil {
// 		xlog.Errorf("ProcSettleOrderFee begin tx fail with %v", err)
// 		return
// 	}
// 	defer func() {
// 		if err == nil {
// 			err = tx.Commit()
// 		} else {
// 			tx.Rollback()
// 		}
// 		if err != nil && err != pgx.ErrNoRows {
// 			xlog.Errorf("ProcSettleOrderFee proc settle order %v fail with %v", orderID, err)
// 		} else if brokerID > 0 && feeBroker.IsPositive() {
// 			xlog.Infof(
// 				"ProcSettleOrderFee proc settle order %v with %v %v fee success, settle %v %v to broker %v",
// 				orderID, feeFilled, feeBalance, feeBroker, feeBalance, brokerID,
// 			)
// 		}
// 	}()
// 	querySQL := `
// 		select u.role,u.broker_id,o.fee_balance,o.fee_filled from gex_order o join gex_user u on o.user_id=u.tid
// 		where o.tid=$1 and o.type=$2 and o.status=any($3) and o.fee_settled_status=$4  for update
// 	`
// 	queryArg := []interface{}{orderID, OrderTypeTrade, []int{OrderStatusPartcanceled, OrderStatusDone}, 0}
// 	err = tx.QueryRow(querySQL, queryArg...).Scan(&userRole, &brokerID, &feeBalance, &feeFilled)
// 	if err != nil {
// 		return
// 	}
// 	if brokerID > 0 && userRole == UserRoleNormal {
// 		feeBroker = feeFilled.Mul(decimal.NewFromFloat(brokerCommRate))
// 		comms := []*OrderComm{
// 			{
// 				OrderID:   orderID,
// 				UserID:    0,
// 				Type:      OrderCommTypeSys,
// 				InBalance: feeBalance,
// 				InFee:     feeFilled.Sub(feeBroker),
// 				Status:    OrderCommStatusNormal,
// 			},
// 			{
// 				OrderID:   orderID,
// 				UserID:    brokerID,
// 				Type:      OrderCommTypeBroker,
// 				InBalance: feeBalance,
// 				InFee:     feeBroker,
// 				Status:    OrderCommStatusNormal,
// 			},
// 		}
// 		err = AddOrderCommCall(tx, comms...)
// 		if err != nil {
// 			return
// 		}
// 		balance := &Balance{
// 			UserID: brokerID,
// 			Asset:  feeBalance,
// 			Free:   feeBroker,
// 		}
// 		err = IncreaseBalanceCall(tx, balance)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	err = tx.ExecRow(`update gex_order set fee_settled_status=$1 where tid=$2`, 1, orderID)
// 	return
// }
