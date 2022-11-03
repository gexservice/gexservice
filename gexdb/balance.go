package gexdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

func TouchBalance(ctx context.Context, area BalanceArea, assets []string, userIDs ...int64) (added int64, err error) {
	added, err = TouchBalanceCall(Pool(), ctx, area, assets, userIDs...)
	return
}

func TouchBalanceCall(caller crud.Queryer, ctx context.Context, area BalanceArea, assets []string, userIDs ...int64) (added int64, err error) {
	upsertArg := []interface{}{0, 0, time.Now(), time.Now(), BalanceStatusNormal, area}
	values := []string{}

	for _, userID := range userIDs {
		for _, asset := range assets {
			upsertArg = append(upsertArg, userID, asset)
			values = append(values, fmt.Sprintf("($1,$2,$3,$4,$5,$%d,$6,$%d)", len(upsertArg)-1, len(upsertArg)))
		}
	}
	upsertSQL := fmt.Sprintf(`
		insert into exs_balance(free,locked,update_time,create_time,status,user_id,area,asset)
		values %v
		on conflict(user_id,area,asset) do nothing
	`, strings.Join(values, ","))

	_, added, err = caller.Exec(ctx, upsertSQL, upsertArg...)
	return
}

func IncreaseBalance(ctx context.Context, balance *Balance) (err error) {
	err = IncreaseBalanceCall(Pool(), ctx, balance)
	return
}

func IncreaseBalanceCall(caller crud.Queryer, ctx context.Context, balance *Balance) (err error) {
	var free, locked, margin decimal.Decimal
	err = caller.QueryRow(
		ctx,
		`select tid,free,locked,margin,create_time,status from exs_balance where user_id=$1 and area=$2 and asset=$3 for update`,
		balance.UserID, balance.Area, balance.Asset,
	).Scan(
		&balance.TID, &free, &locked, &margin, &balance.CreateTime, &balance.Status,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = ErrBalanceNotFound(fmt.Sprintf("balance %v in %v not found", balance.Asset, balance.Area))
		}
		return
	}
	if balance.Status != BalanceStatusNormal {
		err = fmt.Errorf("balance status is %v", balance.Status)
		return
	}
	if free.Add(balance.Free).IsNegative() {
		err = ErrBalanceNotEnought(fmt.Sprintf("%v balance %v.free %v in %v not enought to %v", balance.UserID, balance.Asset, balance.Area, free, balance.Free))
		return
	}
	if locked.Add(balance.Locked).IsNegative() {
		err = ErrBalanceNotEnought(fmt.Sprintf("%v balance %v.locked %v in %v not enought to %v", balance.UserID, balance.Asset, balance.Area, locked, balance.Locked))
		return
	}
	if margin.Add(balance.Margin).IsNegative() {
		err = ErrBalanceNotEnought(fmt.Sprintf("%v balance %v.margin %v in %v not enought to %v", balance.UserID, balance.Asset, balance.Area, margin, balance.Margin))
		return
	}
	balance.Free = free.Add(balance.Free)
	balance.Locked = locked.Add(balance.Locked)
	balance.Margin = margin.Add(balance.Margin)
	balance.UpdateTime = xsql.TimeNow()
	err = balance.UpdateFilter(caller, ctx, "free,locked,margin,update_time#all")
	return
}

func FindBalanceByAsset(ctx context.Context, userID int64, area BalanceArea, asset string) (balance *Balance, err error) {
	balance, err = FindBalanceByAssetCall(Pool(), ctx, userID, area, asset)
	return
}

func FindBalanceByAssetCall(caller crud.Queryer, ctx context.Context, userID int64, area BalanceArea, asset string) (balance *Balance, err error) {
	balance, err = FindBalanceWherefCall(caller, ctx, false, "user_id=$%v,area=$%v,asset=$%v#all", userID, area, asset)
	return
}

func ListUserBalance(ctx context.Context, userID int64, area BalanceArea, assets []string, status BalanceStatusArray) (balanceList []*Balance, balanceMap map[string]*Balance, err error) {
	err = ScanBalanceFilterWheref(
		ctx, "#all",
		"user_id=$%v,area=$%v,asset=any($%v),status=any($%v)",
		[]interface{}{userID, area, xsql.StringArray(assets), status},
		"", &balanceList, &balanceMap, "asset",
	)
	return
}

func CountBalance(ctx context.Context, area BalanceArea, start, end time.Time) (balances map[string]decimal.Decimal, err error) {
	//not using sql sum for percision loss
	balances = map[string]decimal.Decimal{}
	err = crud.QueryWheref(
		Pool, ctx, &Balance{}, "asset,free,locked#all",
		"area=$%v,update_time>=$%v,update_time<$%v",
		[]interface{}{area, start, end},
		"", 0, 0,
		func(balance *Balance) {
			having, ok := balances[balance.Asset]
			if !ok {
				having = decimal.Zero
				balances[balance.Asset] = having
			}
			balances[balance.Asset] = having.Add(balance.Free).Add(balance.Locked)
		},
	)
	return
}

// func ListUserBalanceHistory(ctx context.Context, userID int64, asset string, startTime, endTime time.Time) (histories []*BalanceHistory, err error) {
// 	err = crud.QueryWheref(
// 		Pool, ctx, &BalanceHistory{}, "#all",
// 		"user_id=$%v,asset=$%v,create_time>=$%v,create_time<=$%v",
// 		[]interface{}{userID, asset, startTime, endTime},
// 		"order by create_time asc", 0, 0,
// 		&histories,
// 	)
// 	return
// }

// func ChangeBalance(ctx context.Context, creator, userID int64, asset string, changed decimal.Decimal) (balance *Balance, order *Order, err error) {
// 	tx, err := Pool().Begin(ctx)
// 	if err != nil {
// 		return
// 	}
// 	defer func() {
// 		if err == nil {
// 			err = tx.Commit(ctx)
// 		} else {
// 			tx.Rollback(ctx)
// 		}
// 	}()
// 	balance, order, err = ChangeBalanceCall(tx, ctx, creator, userID, asset, changed)
// 	return
// }

// func ChangeBalanceCall(caller crud.Queryer, ctx context.Context, creator, userID int64, asset string, changed decimal.Decimal) (balance *Balance, order *Order, err error) {
// 	balance = &Balance{
// 		UserID: userID,
// 		Asset:  asset,
// 		Free:   changed,
// 		Status: BalanceStatusNormal,
// 	}
// 	order = &Order{
// 		OrderID:  NewOrderID(),
// 		UserID:   balance.UserID,
// 		Creator:  creator,
// 		Quantity: changed,
// 		Filled:   changed,
// 		Status:   OrderStatusDone,
// 	}
// 	// switch balance.Asset {
// 	// // case BalanceAssetYWE:
// 	// // 	order.Type = OrderTypeChangeYWE
// 	// // case BalanceAssetMMK:
// 	// // 	order.Type = OrderTypeMMK
// 	// default:
// 	// 	err = fmt.Errorf("balance asset %v is not supported", balance.Asset)
// 	// 	return
// 	// }
// 	if balance.Free.LessThan(decimal.Zero) {
// 		order.OutBalance = balance.Asset
// 		order.OutFilled = balance.Free.Abs()
// 		var having decimal.Decimal
// 		err = caller.QueryRow(ctx, `select free from exs_balance where user_id=$1 and asset=$2 for update`, balance.UserID, balance.Asset).Scan(&having)
// 		if err != nil {
// 			return
// 		}
// 		if having.LessThan(balance.Free.Abs()) {
// 			err = ErrBalanceNotEnought(fmt.Errorf("not enought"))
// 			return
// 		}
// 	} else {
// 		order.InBalance = balance.Asset
// 		order.InFilled = balance.Free
// 	}
// 	err = IncreaseBalanceCall(caller, ctx, balance)
// 	if err == nil {
// 		err = AddOrderCall(caller, ctx, order)
// 	}
// 	return
// }
