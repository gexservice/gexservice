package gexdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xsql"
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

func ListUserHolding(ctx context.Context, userID int64, symbolOnly []string) (holdings []*Holding, symbols []string, err error) {
	holdings, symbols, err = ListUserHoldingCall(Pool(), ctx, userID, symbolOnly)
	return
}

func ListUserHoldingCall(caller crud.Queryer, ctx context.Context, userID int64, symbolOnly []string) (holdings []*Holding, symbols []string, err error) {
	sql := crud.QuerySQL(&Holding{}, "#all")
	where, args := crud.AppendWheref(nil, nil, "user_id=$%v,amount!=$%v#all", userID, 0)
	where, args = crud.AppendWheref(where, args, "symbol=any($%v)", xsql.StringArray(symbolOnly))
	sql = crud.JoinWhere(sql, where, "and")
	err = crud.Query(caller, ctx, &Holding{}, "#all", sql, args, &holdings, &symbols, "symbol")
	return
}

func ListHoldingForBlowupOverCall(caller crud.Queryer, ctx context.Context, symbol string, ask, bid decimal.Decimal, lock bool) (holdings []*Holding, err error) {
	querySQL := crud.QuerySQL(&Holding{}, "#all")
	var args []interface{}
	var and, or []string
	and, args = crud.AppendWhere(and, args, true, "symbol=$%v", symbol)
	or, args = crud.AppendWhere(or, args, bid.IsPositive(), "(amount>0 and blowup>=$%v)", bid)
	or, args = crud.AppendWhere(or, args, ask.IsPositive(), "(amount<0 and blowup<=$%v)", ask)
	and = append(and, "("+strings.Join(or, " or ")+")")
	and, args = crud.AppendWhere(and, args, true, "status=$%v", HoldingStatusNormal)
	querySQL = crud.JoinWhere(querySQL, and, " and ", "order by update_time asc")
	if lock {
		querySQL += " for update "
	}
	err = crud.Query(caller, ctx, &Holding{}, "#all", querySQL, args, &holdings)
	return
}

func ListHoldingForBlowupFreeCall(caller crud.Queryer, ctx context.Context, symbol string, ask, bid decimal.Decimal, lock bool) (holdings []*Holding, err error) {
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
	if lock {
		querySQL += " for update "
	}
	err = crud.Query(caller, ctx, &Holding{}, "#all", querySQL, args, &holdings)
	return
}

func ListHoldingByUser(ctx context.Context, userIDs []int64, symbolOnly []string) (holdings map[int64][]*Holding, symbols []string, err error) {
	holdings, symbols, err = ListHoldingByUserCall(Pool(), ctx, userIDs, symbolOnly)
	return
}

func ListHoldingByUserCall(caller crud.Queryer, ctx context.Context, userIDs []int64, symbolOnly []string) (holdings map[int64][]*Holding, symbols []string, err error) {
	sql := crud.QuerySQL(&Holding{}, "#all")
	where, args := crud.AppendWheref(nil, nil, "user_id=any($%v),amount!=$%v#all", xsql.Int64Array(userIDs), 0)
	where, args = crud.AppendWheref(where, args, "symbol=any($%v)", xsql.StringArray(symbolOnly))
	sql = crud.JoinWhere(sql, where, "and")
	symbolAll := map[string]bool{}
	holdings = map[int64][]*Holding{}
	err = crud.Query(caller, ctx, &Holding{}, "#all", sql, args, func(holding *Holding) {
		if !symbolAll[holding.Symbol] {
			symbols = append(symbols, holding.Symbol)
			symbolAll[holding.Symbol] = true
		}
		holdings[holding.UserID] = append(holdings[holding.UserID], holding)
	})
	return
}

/**
 * @apiDefine HoldingUnifySearcher
 * @apiParam  {Number} [symbol] the holding symbol filter, multi with comma
 * @apiParam  {Number} [status] the holding status filter, multi with comma, all type supported is <a href="#metadata-Balance">HoldingStatusAll</a>
 * @apiParam  {String} [key] the search key
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type HoldingUnifySearcher struct {
	Model Holding `json:"model" from:"gex_holding h join gex_user u on h.user_id=u.tid"`
	Where struct {
		UserID int64              `json:"user_id" cmp:"h.user_id=$%v" valid:"user_id,o|i,r:0;"`
		Symbol []string           `json:"symbol" cmp:"h.symbol=any($%v)" valid:"symbol,o|s,l:0;"`
		Status HoldingStatusArray `json:"status" cmp:"h.status=any($%v)" valid:"status,o|i,e:;"`
		Key    string             `json:"key" cmp:"(u.tid::text like $%v or u.name like $%v or u.phone like $%v or u.account like $%v)" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by h.update_time desc" valid:"order,o|s,l:0;"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Holdings []*Holding `json:"holdings"`
		UserIDs  []int64    `json:"user_ids" scan:"user_id"`
	} `json:"query" filter:"h.#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"h.count(tid)#all"`
}

func (b *HoldingUnifySearcher) Apply(ctx context.Context) (err error) {
	b.Page.Order = ""
	if len(b.Where.Key) > 0 {
		b.Where.Key = "%" + b.Where.Key + "%"
	}
	err = crud.ApplyUnify(Pool(), ctx, b)
	return
}
