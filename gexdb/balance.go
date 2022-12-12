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
	added, err = TouchMultiBalanceCall(caller, ctx, BalanceAreaArray{area}, assets, userIDs...)
	return
}

func TouchMultiBalanceCall(caller crud.Queryer, ctx context.Context, areaAll BalanceAreaArray, assets []string, userIDs ...int64) (added int64, err error) {
	upsertArg := []interface{}{0, 0, time.Now(), time.Now(), BalanceStatusNormal}
	values := []string{}

	for _, userID := range userIDs {
		for _, area := range areaAll {
			for _, asset := range assets {
				upsertArg = append(upsertArg, userID, area, asset)
				values = append(values, fmt.Sprintf("($1,$2,$3,$4,$5,$%d,$%d,$%d)", len(upsertArg)-2, len(upsertArg)-1, len(upsertArg)))
			}
		}
	}
	upsertSQL := fmt.Sprintf(`
		insert into gex_balance(free,locked,update_time,create_time,status,user_id,area,asset)
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
		`select tid,free,locked,margin,create_time,status from gex_balance where user_id=$1 and area=$2 and asset=$3 for update`,
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

func ListAreaBalance(ctx context.Context, userID int64, area BalanceAreaArray, asset string, status BalanceStatusArray) (balanceList []*Balance, balanceMap map[BalanceArea]*Balance, err error) {
	err = ScanBalanceFilterWheref(
		ctx, "#all",
		"user_id=$%v,area=any($%v),asset=$%v,status=any($%v)",
		[]interface{}{userID, area, asset, status},
		"", &balanceList, &balanceMap, "area",
	)
	return
}

func TransferChange(ctx context.Context, creator, userID int64, from, to BalanceArea, asset string, value decimal.Decimal) (err error) {
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
	err = TransferChangeCall(tx, ctx, creator, userID, from, to, asset, value)
	return
}

func TransferChangeCall(caller crud.Queryer, ctx context.Context, creator, userID int64, from, to BalanceArea, asset string, value decimal.Decimal) (err error) {
	_, err = TouchMultiBalanceCall(caller, ctx, BalanceAreaArray{from, to}, []string{asset}, userID)
	if err != nil {
		return
	}
	fromBalance := &Balance{
		UserID: userID,
		Area:   from,
		Asset:  asset,
		Free:   decimal.Zero.Sub(value),
	}
	err = IncreaseBalanceCall(caller, ctx, fromBalance)
	if err != nil {
		return
	}
	toBalance := &Balance{
		UserID: userID,
		Area:   to,
		Asset:  asset,
		Free:   value,
	}
	err = IncreaseBalanceCall(caller, ctx, toBalance)
	if err != nil {
		return
	}
	_, err = AddBalanceRecordCall(
		caller, ctx,
		&BalanceRecord{
			Creator:   creator,
			BalanceID: fromBalance.TID,
			Type:      BalanceRecordTypeChange,
			Target:    int(to),
			Changed:   decimal.Zero.Sub(value),
		},
		&BalanceRecord{
			Creator:   creator,
			BalanceID: toBalance.TID,
			Type:      BalanceRecordTypeChange,
			Target:    int(from),
			Changed:   value,
		},
	)
	return
}

func ChangeBalance(ctx context.Context, creator, userID int64, area BalanceArea, asset string, changed decimal.Decimal) (balance *Balance, err error) {
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
	balance, err = ChangeBalanceCall(tx, ctx, creator, userID, area, asset, changed)
	return
}

func ChangeBalanceCall(caller crud.Queryer, ctx context.Context, creator, userID int64, area BalanceArea, asset string, changed decimal.Decimal) (balance *Balance, err error) {
	_, err = TouchBalanceCall(caller, ctx, area, []string{asset}, userID)
	if err != nil {
		return
	}
	balance = &Balance{
		UserID: userID,
		Area:   area,
		Asset:  asset,
		Free:   changed,
	}
	err = IncreaseBalanceCall(caller, ctx, balance)
	if err != nil {
		return
	}
	_, err = AddBalanceRecordCall(caller, ctx, &BalanceRecord{
		Creator:   creator,
		BalanceID: balance.TID,
		Type:      BalanceRecordTypeChange,
		Changed:   changed,
	})
	return
}

func ListBalanceAsset(ctx context.Context, area BalanceAreaArray) (assets []string, err error) {
	sql := `select distinct asset from gex_balance`
	sql, args := crud.JoinWheref(sql, nil, "area=any($%v)", area)
	err = crud.Query(Pool, ctx, MetaWithBalance(string("")), "asset#all", sql, args, &assets, "asset")
	return
}

/**
 * @apiDefine BalanceUnifySearcher
 * @apiParam  {String} [area] the balance area filter, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiParam  {Number} [asset] the balance asset filter, multi with comma
 * @apiParam  {Number} [status] the balance status filter, multi with comma, all type supported is <a href="#metadata-Balance">BalanceStatusAll</a>
 * @apiParam  {String} [key] the search key
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type BalanceUnifySearcher struct {
	Model Balance `json:"model" from:"gex_balance b join gex_user u on b.user_id=u.tid"`
	Where struct {
		UserID int64              `json:"user_id" cmp:"b.user_id=$%v" valid:"user_id,o|i,r:0;"`
		Area   BalanceAreaArray   `json:"area" cmp:"b.area=any($%v)" valid:"area,o|i,e:0;"`
		Asset  []string           `json:"asset" cmp:"b.asset=any($%v)" valid:"asset,o|s,l:0;"`
		Status BalanceStatusArray `json:"status" cmp:"b.status=any($%v)" valid:"status,o|i,e:;"`
		Key    string             `json:"key" cmp:"(u.tid::text like $%v or u.name like $%v or u.phone like $%v or u.account like $%v)" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by b.update_time desc" valid:"order,o|s,l:0;"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Balances []*Balance `json:"balances"`
		UserIDs  []int64    `json:"user_ids" scan:"user_id"`
	} `json:"query" filter:"b.#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"b.count(tid)#all"`
}

func (b *BalanceUnifySearcher) Apply(ctx context.Context) (err error) {
	b.Page.Order = ""
	if len(b.Where.Key) > 0 {
		b.Where.Key = "%" + b.Where.Key + "%"
	}
	err = crud.ApplyUnify(Pool(), ctx, b)
	return
}

func AddBalanceRecordCall(caller crud.Queryer, ctx context.Context, records ...*BalanceRecord) (added int64, err error) {
	if len(records) < 1 {
		return
	}
	talbe, fileds, _, _ := crud.InsertArgs(records[0], "^tid#all", nil)
	insertVal := []string{}
	insertArg := []interface{}{}
	now := xsql.TimeNow()
	for _, record := range records {
		record.UpdateTime = now
		record.CreateTime = now
		record.Status = BalanceRecordStatusNormal
		var param []string
		_, _, param, insertArg = crud.InsertArgs(record, "^tid#all", insertArg)
		insertVal = append(insertVal, "("+strings.Join(param, ",")+")")
	}
	insertSQL := fmt.Sprintf(`insert into %v(%v) values %v`, talbe, strings.Join(fileds, ","), strings.Join(insertVal, ","))
	_, added, err = caller.Exec(ctx, insertSQL, insertArg...)
	return
}

/**
 * @apiDefine BalanceRecordUnifySearcher
 * @apiParam  {String} [area] the balance area filter, all type supported is <a href="#metadata-Balance">BalanceAreaAll</a>
 * @apiParam  {Number} [asset] the balance asset filter, multi with comma
 * @apiParam  {Number} [type] the balance record type filter, multi with comma, all type supported is <a href="#metadata-BalanceRecord">BalanceRecordTypeAll</a>
 * @apiParam  {Number} [start_time] the time filter
 * @apiParam  {Number} [end_time] the time filter
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type BalanceRecordUnifySearcher struct {
	Model BalanceRecordItem `json:"model" from:"gex_balance_record r join gex_balance b on b.tid=r.balance_id join gex_user u on b.user_id=u.tid"`
	Where struct {
		UserID    int64                  `json:"user_id" cmp:"b.user_id=$%v" valid:"user_id,o|i,r:0;"`
		Area      BalanceArea            `json:"area" cmp:"b.area=$%v" valid:"area,o|i,e:0;"`
		Asset     []string               `json:"asset" cmp:"b.asset=any($%v)" valid:"asset,o|s,l:0;"`
		Type      BalanceRecordTypeArray `json:"type" cmp:"r.type=any($%v)" valid:"type,o|i,e:;"`
		StartTime xsql.Time              `json:"start_time" cmp:"r.update_time>=$%v" valid:"start_time,o|i,r:-1;"`
		EndTime   xsql.Time              `json:"end_time" cmp:"r.update_time<$%v" valid:"end_time,o|i,r:-1;"`
		Key       string                 `json:"key" cmp:"(u.tid::text like $%v or u.name like $%v or u.phone like $%v or u.account like $%v)" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by r.update_time desc" valid:"order,o|s,l:0;"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Records []*BalanceRecordItem `json:"records"`
		UserIDs []int64              `json:"user_id" scan:"user_id"`
	} `json:"query" filter:"b.user_id,asset#all|r.type,target,changed,update_time#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"r.count(tid)#all"`
}

func (b *BalanceRecordUnifySearcher) Apply(ctx context.Context) (err error) {
	b.Page.Order = ""
	if len(b.Where.Key) > 0 {
		b.Where.Key = "%" + b.Where.Key + "%"
	}
	err = crud.ApplyUnify(Pool(), ctx, b)
	return
}
