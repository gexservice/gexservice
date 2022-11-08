package gexdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xsql"
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
	err = crud.QueryRow(caller, ctx, &Withdraw{}, "#all", querySQL, args, &withdraw)
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
		Area:   BalanceAreaSpot,
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
		Area:   BalanceAreaSpot,
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

/**
 * @apiDefine WithdrawUnifySearcher
 * @apiParam  {Number} [type] the withdraw type filter, multi with comma, all type supported is <a href="#metadata-Withdraw">WithdrawTypeAll</a>
 * @apiParam  {Number} [asset] the balance asset filter, multi with comma
 * @apiParam  {Number} [start_time] the time filter
 * @apiParam  {Number} [end_time] the time filter
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type WithdrawUnifySearcher struct {
	Model Withdraw `json:"model"`
	Where struct {
		UserID    int64             `json:"user_id" cmp:"user_id=$%v" valid:"user_id,o|i,r:0;"`
		Type      WithdrawTypeArray `json:"type" cmp:"type=any($%v)" valid:"type,o|i,e:;"`
		Asset     []string          `json:"asset" cmp:"asset=any($%v)" valid:"asset,o|s,l:0;"`
		StartTime xsql.Time         `json:"start_time" cmp:"update_time>=$%v" valid:"start_time,o|i,r:-1;"`
		EndTime   xsql.Time         `json:"end_time" cmp:"update_time<$%v" valid:"end_time,o|i,r:-1;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by update_time desc" valid:"order,o|s,l:0;"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Withdraws []*Withdraw `json:"withdraws"`
	} `json:"query" filter:"#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"r.count(tid)#all"`
}

func (w *WithdrawUnifySearcher) Apply(ctx context.Context) (err error) {
	w.Page.Order = ""
	err = crud.ApplyUnify(Pool(), ctx, w)
	return
}
