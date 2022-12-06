package gexdb

import (
	"context"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

func CountUser(ctx context.Context, start, end time.Time) (total int64, err error) {
	err = crud.CountWheref(Pool, ctx, MetaWithUser(int64(0)), "count(tid)#all", "create_time>=$%v,create_time<$%v", crud.Args(start, end), "", &total, "tid")
	return
}

func CountAreaBalance(ctx context.Context, area BalanceArea, start, end time.Time) (balances []*Balance, err error) {
	sql := `select asset,sum(free+locked) as total from gex_balance`
	sql, args := crud.JoinWheref(sql, nil, "area=$%v,update_time>=$%v,update_time<$%v", area, start, end)
	sql += " group by asset order by sum(free+locked) desc"
	err = crud.Query(Pool, ctx, &Balance{}, "asset,free#all", sql, args, &balances)
	return
}

func CountAllBalance(ctx context.Context, asset string) (balances []*Balance, err error) {
	sql := `select area,asset,sum(free),sum(locked) as total from gex_balance`
	sql, args := crud.JoinWheref(sql, nil, "asset=$%v", asset)
	sql += " group by area,asset order by area,asset"
	err = crud.Query(Pool, ctx, &Balance{}, "area,asset,free,locked#all", sql, args, &balances)
	return
}

func CountUserBalance(ctx context.Context, asset string, userIDs ...int64) (balances map[int64]*Balance, err error) {
	sql := `select user_id,asset,sum(free),sum(locked) from gex_balance`
	sql, args := crud.JoinWheref(sql, nil, "asset=$%v,user_id=any($%v)", asset, xsql.Int64Array(userIDs))
	sql += " group by user_id,asset"
	err = crud.Query(Pool, ctx, &Balance{}, "user_id,asset,free,locked#all", sql, args, &balances, "user_id")
	return
}

func CountOrderFee(ctx context.Context, area OrderArea, start, end time.Time) (fee map[OrderArea]map[string]decimal.Decimal, err error) {
	fee = map[OrderArea]map[string]decimal.Decimal{}
	err = crud.CountWheref(
		Pool, ctx,
		MetaWithOrder(OrderAreaNone, string(""), decimal.Zero), "area,fee_balance,sum(fee_filled)#all",
		"area=$%v,update_time>=$%v,update_time<$%v,status=any($%v)", crud.Args(area, start, end, OrderStatusArray{OrderStatusDone, OrderStatusPartCanceled}),
		"group by area,fee_balance",
		func(v []interface{}) {
			area := *(v[0].(*OrderArea))
			balance := *(v[1].(*string))
			filled := *(v[2].(*decimal.Decimal))
			if fee[area] == nil {
				fee[area] = map[string]decimal.Decimal{}
			}
			fee[area][balance] = filled
		},
	)
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
	return
}

func CountHolding(ctx context.Context, side int, start, end time.Time) (holdings []*Holding, err error) {
	sql := `select symbol,sum(amount) as total from gex_holding`
	where, args := crud.AppendWheref(nil, nil, "update_time>=$%v,update_time<$%v", start, end)
	if side > 0 {
		where = append(where, "amount>0")
	} else {
		where = append(where, "amount<0")
	}
	sql = crud.JoinWhere(sql, where, "and")
	sql += " group by symbol order by symbol asc"
	err = crud.Query(Pool, ctx, &Holding{}, "symbol,amount#all", sql, args, &holdings)
	return
}
