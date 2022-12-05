package gexdb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/define"
	"github.com/jackc/pgx/v4"
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

func CreateWithdraw(ctx context.Context, withdraw *Withdraw) (err error) {
	if len(withdraw.OrderID) < 1 {
		withdraw.OrderID = fmt.Sprintf("withdraw_%v", NewOrderID())
	}
	withdraw.Type = WithdrawTypeWithdraw
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
		UserID: withdraw.UserID,
		Area:   BalanceAreaSpot,
		Asset:  withdraw.Asset,
		Free:   decimal.Zero.Sub(withdraw.Quantity),
		Locked: withdraw.Quantity,
	}
	err = IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		return
	}
	review, err := LoadWithdrawReviewCall(tx, ctx)
	if err != nil {
		return
	}
	reviewMin := review.Float64Def(-1, withdraw.Asset)
	if reviewMin < 0 {
		reviewMin = review.Float64Def(-1, "*")
	}
	if reviewMin < 0 || withdraw.Quantity.LessThan(decimal.NewFromFloat(reviewMin)) {
		withdraw.Status = WithdrawStatusConfirmed
	} else {
		withdraw.Status = WithdrawStatusPending
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
		err = fmt.Errorf("order is not withdraw")
		return
	}
	if withdraw.Status != WithdrawStatusPending {
		err = fmt.Errorf("withdraw is not pending")
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

func ConfirmWithdraw(ctx context.Context, orderID string) (err error) {
	err = crud.UpdateRowWheref(Pool, ctx, &Withdraw{Status: WithdrawStatusConfirmed}, "status", "order_id=$%v,status=$%v", orderID, WithdrawStatusPending)
	return
}

func DoneWithdraw(ctx context.Context, orderID string, success bool, result xmap.M) (withdraw *Withdraw, err error) {
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
	if withdraw.Status == WithdrawStatusDone {
		if success {
			return
		}
		err = fmt.Errorf("withdraw %v status is %v", orderID, withdraw.Status)
		return
	} else if withdraw.Status == WithdrawStatusCanceled {
		if !success {
			return
		}
		err = fmt.Errorf("withdraw %v status is %v", orderID, withdraw.Status)
		return
	}
	for k, v := range result {
		withdraw.Result[k] = v
	}
	if success {
		withdraw.Status = WithdrawStatusDone
	} else {
		withdraw.Status = WithdrawStatusCanceled
	}
	err = withdraw.UpdateFilter(tx, ctx, "result,status")
	if err != nil {
		return
	}
	balance := &Balance{
		UserID: withdraw.UserID,
		Area:   BalanceAreaSpot,
		Asset:  withdraw.Asset,
		Locked: decimal.Zero.Sub(withdraw.Quantity),
	}
	if !success {
		balance.Free = withdraw.Quantity
	}
	err = IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		return
	}
	if success {
		_, err = AddBalanceRecordCall(tx, ctx, &BalanceRecord{
			Creator:   withdraw.UserID,
			BalanceID: balance.TID,
			Type:      BalanceRecordTypeWithdraw,
			Changed:   withdraw.Quantity,
			Transaction: xsql.M{
				"txid":     withdraw.OrderID,
				"withdraw": withdraw.TID,
			},
		})
	}
	return
}

func ReceiveTopup(ctx context.Context, method WalletMethod, address string, txid, asset string, amount decimal.Decimal, result xmap.M) (topup *Withdraw, skip bool, err error) {
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
	wallet, err := FindWalletWherefCall(tx, ctx, true, "method=$%v,address=$%v", method, address)
	if err != nil {
		if err == pgx.ErrNoRows {
			skip = true
			err = nil
		}
		return
	}
	orderID := fmt.Sprintf("topup_%v", txid)
	topup, err = FindWithdrawByOrderIDCall(tx, ctx, orderID, false)
	if err != nil && err != pgx.ErrNoRows {
		return
	}
	if err == nil { //received
		return
	}
	topup = &Withdraw{
		OrderID:   orderID,
		Type:      WithdrawTypeTopup,
		UserID:    wallet.UserID,
		Method:    WithdrawMethod(method),
		Asset:     asset,
		Quantity:  amount,
		Processed: 1,
		Result:    xsql.M(result),
		Status:    WithdrawStatusDone,
	}
	err = topup.Insert(tx, ctx)
	if err != nil {
		return
	}
	balance := &Balance{
		UserID: topup.UserID,
		Area:   BalanceAreaSpot,
		Asset:  topup.Asset,
		Free:   topup.Quantity,
	}
	_, err = TouchBalanceCall(tx, ctx, balance.Area, []string{balance.Asset}, balance.UserID)
	if err != nil {
		return
	}
	err = IncreaseBalanceCall(tx, ctx, balance)
	if err != nil {
		return
	}
	_, err = AddBalanceRecordCall(tx, ctx, &BalanceRecord{
		Creator:   topup.UserID,
		BalanceID: balance.TID,
		Type:      BalanceRecordTypeTopup,
		Changed:   topup.Quantity,
		Transaction: xsql.M{
			"txid":  topup.OrderID,
			"topup": topup.TID,
		},
	})
	return
}

/**
 * @apiDefine WithdrawUnifySearcher
 * @apiParam  {Number} [type] the withdraw type filter, multi with comma, all type supported is <a href="#metadata-Withdraw">WithdrawTypeAll</a>
 * @apiParam  {Number} [asset] the balance asset filter, multi with comma
 * @apiParam  {Number} [start_time] the time filter
 * @apiParam  {Number} [end_time] the time filter
 * @apiParam  {Number} [status] the withdraw status filter, multi with comma, all type supported is <a href="#metadata-Withdraw">WithdrawStatusAll</a>
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type WithdrawUnifySearcher struct {
	Model Withdraw `json:"model"`
	Where struct {
		UserID    int64               `json:"user_id" cmp:"user_id=$%v" valid:"user_id,o|i,r:0;"`
		Type      WithdrawTypeArray   `json:"type" cmp:"type=any($%v)" valid:"type,o|i,e:;"`
		Asset     []string            `json:"asset" cmp:"asset=any($%v)" valid:"asset,o|s,l:0;"`
		StartTime xsql.Time           `json:"start_time" cmp:"update_time>=$%v" valid:"start_time,o|i,r:-1;"`
		EndTime   xsql.Time           `json:"end_time" cmp:"update_time<$%v" valid:"end_time,o|i,r:-1;"`
		Status    WithdrawStatusArray `json:"status" cmp:"status=any($%v)" valid:"status,o|i,e:;"`
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
