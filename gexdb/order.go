package gexdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xsql"
	"github.com/shopspring/decimal"
)

var MarchineID = 1
var seqOrderID uint16
var lckOrderID = sync.RWMutex{}

func NewOrderID() (orderID string) {
	lckOrderID.Lock()
	defer lckOrderID.Unlock()
	seqOrderID++
	timeStr := time.Now().Format("20060102150405")
	return fmt.Sprintf("%v%02d%05d", timeStr, MarchineID, seqOrderID)
}

func FindOrderByOrderID(ctx context.Context, userID int64, orderID string) (order *Order, err error) {
	order, err = FindOrderByOrderIDCall(Pool(), ctx, userID, orderID, false)
	return
}

func FindOrderByOrderIDCall(caller crud.Queryer, ctx context.Context, userID int64, orderID string, lock bool) (order *Order, err error) {
	orderIDInt, _ := strconv.ParseInt(orderID, 10, 64)
	querySQL := crud.QuerySQL(&Order{}, "#all")
	and, args := crud.AppendWheref(nil, nil, "user_id=$%v", userID)
	or, args := crud.AppendWheref(nil, args, "tid=$%v,order_id=$%v", orderIDInt, orderID)
	and = append(and, "("+strings.Join(or, " or ")+")")
	querySQL = crud.JoinWhere(querySQL, and, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Order{}, "#all", querySQL, args, &order)
	return
}

func ClearCanceledOrder(ctx context.Context, userID int64, symbol string, before time.Time) (removed int64, err error) {
	sql := `delete from gex_order`
	sql, args := crud.JoinWheref(sql, nil, "user_id=$%v,symbol=$%v,create_time<$%v,status=$%v", userID, symbol, before, OrderStatusCanceled)
	_, removed, err = Pool().Exec(ctx, sql, args...)
	return
}

func CountOrderFee(ctx context.Context, start, end time.Time) (fee map[string]decimal.Decimal, err error) {
	//not using sql sum for percision loss
	fee = map[string]decimal.Decimal{}
	err = crud.QueryWheref(
		Pool, ctx, MetaWithOrder(string(""), decimal.Zero), "fee_balance,fee_filled#all",
		"update_time>=$%v,update_time<$%v,status=any($%v)", []interface{}{start, end, OrderStatusArray{OrderStatusPartCanceled, OrderStatusDone}},
		"", 0, 0,
		func(v []interface{}) {
			balance := *(v[0].(*string))
			filled := *(v[1].(*decimal.Decimal))
			having, ok := fee[balance]
			if !ok {
				having = decimal.Zero
				fee[balance] = having
			}
			fee[balance] = having.Add(filled)
		},
	)
	return
}

func CancelTriggerOrder(ctx context.Context, userID int64, symbol string, orderID int64) (updated int64, err error) {
	updated, err = crud.UpdateWheref(Pool, ctx, &Order{Status: OrderStatusCanceled}, "status", "user_id=$%v,symbol=$%v,tid=$%v,status=$%v", userID, symbol, orderID, OrderStatusWaiting)
	return
}

func ListOrderForTrigger(ctx context.Context, symbol string, ask, bid decimal.Decimal) (orders []*Order, err error) {
	orders, err = ListOrderForTriggerCall(Pool(), ctx, symbol, ask, bid)
	return
}

func ListOrderForTriggerCall(caller crud.Queryer, ctx context.Context, symbol string, ask, bid decimal.Decimal) (orders []*Order, err error) {
	querySQL := crud.QuerySQL(&Order{}, "#all")
	var args []interface{}
	var and, or []string
	and, args = crud.AppendWheref(and, args, "type=$%v,symbol=$%v,status=$%v", OrderTypeTrigger, symbol, OrderStatusWaiting)
	if bid.IsPositive() {
		args = append(args, OrderSideSell, OrderTriggerTypeStopProfit, bid)
		or = append(or, fmt.Sprintf("(side=$%v and trigger_type=$%v and trigger_price<=$%v)", len(args)-2, len(args)-1, len(args)))
	}
	if bid.IsPositive() {
		args = append(args, OrderSideSell, OrderTriggerTypeStopLoss, bid)
		or = append(or, fmt.Sprintf("(side=$%v and trigger_type=$%v and trigger_price>=$%v)", len(args)-2, len(args)-1, len(args)))
	}
	if ask.IsPositive() {
		args = append(args, OrderSideBuy, OrderTriggerTypeStopProfit, ask)
		or = append(or, fmt.Sprintf("(side=$%v and trigger_type=$%v and trigger_price>=$%v)", len(args)-2, len(args)-1, len(args)))
	}
	if ask.IsPositive() {
		args = append(args, OrderSideBuy, OrderTriggerTypeStopLoss, ask)
		or = append(or, fmt.Sprintf("(side=$%v and trigger_type=$%v and trigger_price<=$%v)", len(args)-2, len(args)-1, len(args)))
	}
	if len(or) < 1 {
		err = fmt.Errorf("ask/bid is zero")
		return
	}
	and = append(and, "("+strings.Join(or, " or ")+")")
	querySQL = crud.JoinWhere(querySQL, and, " and ", "order by update_time asc")
	err = crud.Query(caller, ctx, &Order{}, "#all", querySQL, args, &orders)
	return
}

func ListPendingOrder(ctx context.Context, userID int64, area, symbol string) (orders []*Order, err error) {
	orders, err = ListPendingOrderCall(Pool(), ctx, userID, area, symbol)
	return
}

func ListPendingOrderCall(caller crud.Queryer, ctx context.Context, userID int64, area, symbol string) (orders []*Order, err error) {
	if len(area) > 0 {
		area += "%"
	}
	err = crud.QueryWheref(caller, ctx, &Order{}, "#all", "user_id=$%v,symbol like $%v,symbol=$%v,status=any($%v)", []interface{}{userID, area, symbol, OrderStatusArray{OrderStatusWaiting, OrderStatusPending, OrderStatusPartialled}}, "", 0, 0, &orders)
	return
}

func CountPendingOrderCall(caller crud.Queryer, ctx context.Context, userID int64, symbol string) (having int64, err error) {
	err = crud.CountWheref(caller, ctx, MetaWithOrder(having), "count(tid)", "user_id=$%v,symbol=$%v,status=any($%v)", []interface{}{userID, symbol, OrderStatusArray{OrderStatusWaiting, OrderStatusPending, OrderStatusPartialled}}, "", &having, "tid")
	return
}

/**
 * @apiDefine OrderUnifySearcher
 * @apiParam  {String} [side] the side filter, multi with comma, all type supported is <a href="#metadata-Order">OrderSideAll</a>
 * @apiParam  {Number} [type] the type filter, multi with comma, all type supported is <a href="#metadata-Order">OrderTypeAll</a>
 * @apiParam  {String} [area] the symbol area filter
 * @apiParam  {String} [symbol] the symbol filter
 * @apiParam  {Number} [start_time] the time filter
 * @apiParam  {Number} [end_time] the time filter
 * @apiParam  {Number} [status] the status filter, multi with comma, all status supported is <a href="#metadata-Order">OrderStatusAll</a>
 * @apiParam  {String} [key] search key
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type OrderUnifySearcher struct {
	Model Order `json:"model"`
	Where struct {
		UserID    xsql.Int64Array  `json:"user_id" cmp:"user_id=any($%v)" valid:"user_id,o|i,r:0;"`
		Creator   xsql.Int64Array  `json:"creator" cmp:"creator=any($%v)" valid:"creator,o|i,r:0;"`
		Area      string           `json:"area" cmp:"symbol like $%v" valid:"area,o|s,l:0;"`
		Symbol    string           `json:"symbol" cmp:"symbol=$%v"  valid:"symbol,o|s,l:0;"`
		Side      OrderSideArray   `json:"side" cmp:"side=any($%v)" valid:"side,o|s,e:0;"`
		Type      OrderTypeArray   `json:"type" cmp:"type=any($%v)" valid:"type,o|i,e:;"`
		StartTime xsql.Time        `json:"start_time" cmp:"update_time>=$%v" valid:"start_time,o|i,r:-1;"`
		EndTime   xsql.Time        `json:"end_time" cmp:"update_time<$%v" valid:"end_time,o|i,r:-1;"`
		Status    OrderStatusArray `json:"status" cmp:"status=any($%v)" valid:"status,o|i,e:;"`
		Key       string           `json:"key" cmp:"(tid::text ilike $%v or order_id ilike $%v)" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by update_time desc" valid:"order,o|s,l:0;"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Orders   []*Order `json:"orders"`
		OrderIDs []int64  `json:"order_ids" scan:"tid"`
		UserIDs  []int64  `json:"user_ids" scan:"user_id"`
	} `json:"query" filter:"^transaction#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"count(tid)#all"`
}

func (o *OrderUnifySearcher) Apply(ctx context.Context) (err error) {
	if len(o.Where.Key) > 0 {
		o.Where.Key = "%" + o.Where.Key + "%"
	}
	if len(o.Where.Area) > 0 {
		o.Where.Area = o.Where.Area + "%"
	}
	o.Page.Order = crud.BuildOrderby(OrderOrderbyAll, o.Page.Order)
	err = crud.ApplyUnify(Pool(), ctx, o)
	return
}
