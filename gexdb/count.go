package gexdb

import (
	"context"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

func CountUser(ctx context.Context, start, end time.Time) (total int64, err error) {
	err = crud.CountWheref(Pool, ctx, MetaWithUser(int64(0)), "count(tid)#all", "create_time>=$%v,create_time<$%v", crud.Args(start, end), "", &total, "tid")
	return
}

func CountAreaBalance(ctx context.Context, area BalanceAreaArray, asset string, start, end time.Time) (balanceAll []*Balance, balanceMap map[BalanceArea]map[string]*Balance, err error) {
	sql := `select area,asset,sum(free),sum(locked) as total from gex_balance`
	sql, args := crud.JoinWheref(sql, nil, "area=any($%v),asset=$%v,update_time>=$%v,update_time<$%v", area, asset, start, end)
	sql += " group by area,asset order by sum(free+locked) desc"
	balanceMap = map[BalanceArea]map[string]*Balance{}
	err = crud.Query(Pool, ctx, &Balance{}, "area,asset,free,locked#all", sql, args, func(balance *Balance) {
		if balanceMap[balance.Area] == nil {
			balanceMap[balance.Area] = map[string]*Balance{}
		}
		balanceMap[balance.Area][balance.Asset] = balance
		balanceAll = append(balanceAll, balance)
	})
	return
}

func CountUserBalance(ctx context.Context, asset string, userIDs ...int64) (balances map[int64]*Balance, err error) {
	sql := `select user_id,asset,sum(free),sum(locked) from gex_balance`
	sql, args := crud.JoinWheref(sql, nil, "asset=$%v,user_id=any($%v)", asset, xsql.Int64Array(userIDs))
	sql += " group by user_id,asset"
	err = crud.Query(Pool, ctx, &Balance{}, "user_id,asset,free,locked#all", sql, args, &balances, "user_id")
	return
}

func CountOrderFee(ctx context.Context, area OrderArea, start, end time.Time) (feeAll []xmap.M, feeMap map[OrderArea]map[string]decimal.Decimal, err error) {
	feeMap = map[OrderArea]map[string]decimal.Decimal{}
	err = crud.CountWheref(
		Pool, ctx,
		MetaWithOrder(OrderAreaNone, string(""), decimal.Zero), "area,fee_balance,sum(fee_filled)#all",
		"area=$%v,update_time>=$%v,update_time<$%v,status=any($%v)", crud.Args(area, start, end, OrderStatusArray{OrderStatusDone, OrderStatusPartCanceled}),
		"group by area,fee_balance order by area asc,fee_balance asc",
		func(v []interface{}) {
			area := *(v[0].(*OrderArea))
			balance := *(v[1].(*string))
			filled := *(v[2].(*decimal.Decimal))
			if feeMap[area] == nil {
				feeMap[area] = map[string]decimal.Decimal{}
			}
			feeMap[area][balance] = filled
			feeAll = append(feeAll, xmap.M{
				"area":        area,
				"fee_balance": balance,
				"fee_filled":  filled,
			})
		},
	)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return
}

func CountOrderVolume(ctx context.Context, area OrderArea, start, end time.Time) (orders []*Order, err error) {
	err = crud.CountWheref(
		Pool, ctx,
		&Order{}, "area,symbol,sum(quantity),sum(filled),sum(total_price)#all",
		"area=$%v,update_time>=$%v,update_time<$%v,status=any($%v)", crud.Args(area, start, end, OrderStatusArray{OrderStatusDone, OrderStatusPartCanceled}),
		"group by area,symbol",
		&orders,
	)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return
}

func CountHolding(ctx context.Context, side int, start, end time.Time) (holdingAll []*Holding, holdingMap map[string]*Holding, err error) {
	sql := `select symbol,sum(amount) as total from gex_holding`
	where, args := crud.AppendWheref(nil, nil, "update_time>=$%v,update_time<$%v", start, end)
	if side > 0 {
		where = append(where, "amount>0")
		sql = crud.JoinWhere(sql, where, "and")
		sql += " group by symbol order by sum(amount) desc"
	} else {
		where = append(where, "amount<0")
		sql = crud.JoinWhere(sql, where, "and")
		sql += " group by symbol order by sum(amount) asc"
	}
	err = crud.Query(Pool, ctx, &Holding{}, "symbol,amount#all", sql, args, &holdingAll, &holdingMap, "symbol")
	return
}
