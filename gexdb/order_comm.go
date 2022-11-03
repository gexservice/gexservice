package gexdb

// import (
// 	"context"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/codingeasygo/crud"
// 	"github.com/codingeasygo/util/xsql"
// 	"github.com/codingeasygo/util/xtime"
// 	"github.com/shopspring/decimal"
// )

// func AddMultiOrderCommCall(caller crud.Queryer, ctx context.Context, comms ...*OrderComm) (err error) {
// 	for _, comm := range comms {
// 		comm.CreateTime = xsql.TimeNow()
// 		err = comm.Insert(caller, ctx)
// 		if err != nil {
// 			break
// 		}
// 	}
// 	return
// }

// func ListUserOrderComm(userID int64, orderIDs ...int64) (comms map[int64]*OrderComm, total map[string]decimal.Decimal, err error) {
// 	querySQL := `
// 		select
// 			tid,order_id,user_id,type,in_balance,in_fee,create_time,status
// 		from exs_order_comm where user_id=$1 and order_id=any($2)
// 	`
// 	args := []interface{}{userID, xsql.Int64Array(orderIDs).DbArray()}
// 	rows, err := Pool().Query(querySQL, args...)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	comms = map[int64]*OrderComm{}
// 	total = map[string]decimal.Decimal{}
// 	for rows.Next() {
// 		comm := &OrderComm{}
// 		err = rows.Scan(
// 			&comm.TID, &comm.OrderID, &comm.UserID, &comm.Type, &comm.InBalance, &comm.InFee, &comm.CreateTime, &comm.Status,
// 		)
// 		if err != nil {
// 			return
// 		}
// 		total[comm.InBalance] = total[comm.InBalance].Add(comm.InFee)
// 		comms[comm.OrderID] = comm
// 	}
// 	countSQL := `
// 		select in_balance,sum(in_fee)
// 		from exs_order_comm where user_id=$1 and order_id=any($2)
// 		group by in_balance
// 	`
// 	total, err = countOrderComm(countSQL, args...)
// 	return
// }

// func CountUserOrderComm(userID int64, startTime, endTime time.Time) (comms map[string]decimal.Decimal, err error) {
// 	where := []string{"user_id=$1"}
// 	args := []interface{}{userID}

// 	if xtime.Timestamp(startTime) > 0 {
// 		args = append(args, startTime)
// 		where = append(where, fmt.Sprintf(" create_time>=$%d ", len(args)))
// 	}
// 	if xtime.Timestamp(endTime) > 0 {
// 		args = append(args, endTime)
// 		where = append(where, fmt.Sprintf(" create_time<=$%d ", len(args)))
// 	}
// 	querySQL := `
// 		select in_balance,sum(in_fee)
// 		from exs_order_comm where
// 	`
// 	querySQL += strings.Join(where, " and ")
// 	querySQL += ` group by in_balance `
// 	comms, err = countOrderComm(querySQL, args...)
// 	return
// }

// func countOrderComm(query string, args ...interface{}) (comms map[string]decimal.Decimal, err error) {
// 	rows, err := Pool().Query(query, args...)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	comms = map[string]decimal.Decimal{}
// 	for rows.Next() {
// 		var balance string
// 		var fee decimal.Decimal
// 		err = rows.Scan(&balance, &fee)
// 		if err != nil {
// 			return
// 		}
// 		comms[balance] = fee
// 	}
// 	return
// }

// func CountMyUserOrderComm(userID int64, myUserIDs ...int64) (comms map[int64]map[string]decimal.Decimal, err error) {
// 	where := []string{"c.user_id=$1"}
// 	args := []interface{}{userID}
// 	if len(myUserIDs) > 0 {
// 		args = append(args, xsql.Int64Array(myUserIDs).DbArray())
// 		where = append(where, fmt.Sprintf("o.user_id=any($%v)", len(args)))
// 	}
// 	querySQL := `
// 		select o.user_id,c.in_balance,sum(c.in_fee)
// 		from exs_order_comm c join exs_order o on c.order_id=o.tid where
// 	`
// 	querySQL += strings.Join(where, " and ")
// 	querySQL += ` group by o.user_id,c.in_balance `
// 	rows, err := Pool().Query(querySQL, args...)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	comms = map[int64]map[string]decimal.Decimal{}
// 	for rows.Next() {
// 		var myUserID int64
// 		var balance string
// 		var fee decimal.Decimal
// 		err = rows.Scan(&myUserID, &balance, &fee)
// 		if err != nil {
// 			return
// 		}
// 		if comms[myUserID] == nil {
// 			comms[myUserID] = map[string]decimal.Decimal{}
// 		}
// 		comms[myUserID][balance] = fee
// 	}
// 	return
// }
