package gexdb

import (
	"context"

	"github.com/codingeasygo/crud"
	"github.com/jackc/pgx/v4"
)

var ApplyWithdraw = func(withdraw *Withdraw) (err error) {
	panic("not init")
}

func ProcWithdraw() (err error) {
	ctx := context.Background()
	var withdraw *Withdraw
	err = crud.QueryWheref(Pool, ctx, &Withdraw{}, "#all", "type=$%v,processed<$%v,status=$%v", crud.Args(WithdrawTypeWithdraw, 1, WithdrawStatusConfirmed), "order by update_time asc", 0, 1, &withdraw)
	if err != nil {
		return
	}
	if withdraw == nil {
		err = pgx.ErrNoRows
		return
	}
	err = ApplyWithdraw(withdraw)
	if err != nil {
		return
	}
	withdraw.Processed = 1
	err = withdraw.UpdateFilter(Pool, ctx, "processed")
	return
}

var ApplyWallet = func(method WalletMethod) (address string, err error) {
	panic("not init")
}

func LoadWalletByMethod(ctx context.Context, userID int64, method WalletMethod) (wallet *Wallet, err error) {
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
	_, err = FindUserCall(tx, ctx, userID, true) //lock
	if err != nil {
		return
	}
	wallet, err = FindWalletWherefCall(tx, ctx, false, "user_id=$%v,method=$%v", userID, method)
	if err != nil && err != pgx.ErrNoRows {
		return
	}
	if err == nil {
		return
	}
	address, err := ApplyWallet(method)
	if err != nil {
		return
	}
	wallet = &Wallet{
		UserID:  userID,
		Method:  method,
		Address: address,
		Status:  WalletStatusNormal,
	}
	err = wallet.Insert(tx, ctx)
	return
}
