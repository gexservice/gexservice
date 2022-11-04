package gexdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/shopspring/decimal"
)

func TouchHolding(ctx context.Context, symbols []string, userIDs ...int64) (added int64, err error) {
	added, err = TouchHoldingCall(Pool(), ctx, symbols, userIDs...)
	return
}

func TouchHoldingCall(caller crud.Queryer, ctx context.Context, symbols []string, userIDs ...int64) (added int64, err error) {
	upsertArg := []interface{}{time.Now(), time.Now(), HoldingStatusNormal}
	values := []string{}

	for _, userID := range userIDs {
		for _, symbol := range symbols {
			upsertArg = append(upsertArg, userID, symbol)
			values = append(values, fmt.Sprintf("($1,$2,$3,$%d,$%d)", len(upsertArg)-1, len(upsertArg)))
		}
	}
	upsertSQL := fmt.Sprintf(`
		insert into gex_holding(update_time,create_time,status,user_id,symbol)
		values %v
		on conflict(user_id,symbol) do nothing
	`, strings.Join(values, ","))

	_, added, err = caller.Exec(ctx, upsertSQL, upsertArg...)
	return
}

func FindHoldlingBySymbol(ctx context.Context, userID int64, symbol string) (holding *Holding, err error) {
	holding, err = FindHoldlingBySymbolCall(Pool(), ctx, userID, symbol, false)
	return
}

func FindHoldlingBySymbolCall(caller crud.Queryer, ctx context.Context, userID int64, symbol string, lock bool) (holding *Holding, err error) {
	holding, err = FindHoldingFilterWherefCall(caller, ctx, lock, "#all", "user_id=$%v,symbol=$%v#all", userID, symbol)
	return
}

func ListHoldingForBlowupOverCall(caller crud.Queryer, ctx context.Context, symbol string, ask, bid decimal.Decimal) (holdings []*Holding, err error) {
	querySQL := crud.QuerySQL(&Holding{}, "#all")
	var args []interface{}
	var and, or []string
	and, args = crud.AppendWhere(and, args, true, "symbol=$%v", symbol)
	or, args = crud.AppendWhere(or, args, bid.IsPositive(), "(amount>0 and blowup>=$%v)", bid)
	or, args = crud.AppendWhere(or, args, ask.IsPositive(), "(amount<0 and blowup<=$%v)", ask)
	and = append(and, "("+strings.Join(or, " or ")+")")
	and, args = crud.AppendWhere(and, args, true, "status=$%v", HoldingStatusNormal)
	querySQL = crud.JoinWhere(querySQL, and, " and ", "order by update_time asc")
	err = crud.Query(caller, ctx, &Holding{}, "#all", querySQL, args, &holdings)
	return
}

func ListHoldingForBlowupFreeCall(caller crud.Queryer, ctx context.Context, symbol string, ask, bid decimal.Decimal) (holdings []*Holding, err error) {
	querySQL := crud.QuerySQL(&Holding{}, "#all")
	var args []interface{}
	var and, or []string
	and, args = crud.AppendWhere(and, args, true, "symbol=$%v", symbol)
	and, args = crud.AppendWhere(and, args, true, "margin_added>$%v", 0)
	or, args = crud.AppendWhere(or, args, bid.IsPositive(), "(amount>0 and blowup<$%v)", bid)
	or, args = crud.AppendWhere(or, args, ask.IsPositive(), "(amount<0 and blowup>$%v)", ask)
	and = append(and, "("+strings.Join(or, " or ")+")")
	and, args = crud.AppendWhere(and, args, true, "status=$%v", HoldingStatusNormal)
	querySQL = crud.JoinWhere(querySQL, and, " and ", "order by update_time asc")
	err = crud.Query(caller, ctx, &Holding{}, "#all", querySQL, args, &holdings)
	return
}
