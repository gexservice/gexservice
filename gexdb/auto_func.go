//auto gen func by autogen
package gexdb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/attrvalid"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xsql"
)

var GetQueryer interface{} = func() crud.Queryer { return Pool() }

//Validable is interface to valid
type Validable interface {
	Valid() error
}

//BalanceFilterOptional is crud filter
const BalanceFilterOptional = ""

//BalanceFilterRequired is crud filter
const BalanceFilterRequired = ""

//BalanceFilterInsert is crud filter
const BalanceFilterInsert = ""

//BalanceFilterUpdate is crud filter
const BalanceFilterUpdate = "update_time"

//BalanceFilterFind is crud filter
const BalanceFilterFind = "#all"

//BalanceFilterScan is crud filter
const BalanceFilterScan = "#all"

//EnumValid will valid value by BalanceArea
func (o *BalanceArea) EnumValid(v interface{}) (err error) {
	var target BalanceArea
	targetType := reflect.TypeOf(BalanceArea(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(BalanceArea)
	}
	for _, value := range BalanceAreaAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", BalanceAreaAll)
}

//EnumValid will valid value by BalanceAreaArray
func (o *BalanceAreaArray) EnumValid(v interface{}) (err error) {
	var target BalanceArea
	targetType := reflect.TypeOf(BalanceArea(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(BalanceArea)
	}
	for _, value := range BalanceAreaAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", BalanceAreaAll)
}

//DbArray will join value to database array
func (o BalanceAreaArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o BalanceAreaArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by BalanceStatus
func (o *BalanceStatus) EnumValid(v interface{}) (err error) {
	var target BalanceStatus
	targetType := reflect.TypeOf(BalanceStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(BalanceStatus)
	}
	for _, value := range BalanceStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", BalanceStatusAll)
}

//EnumValid will valid value by BalanceStatusArray
func (o *BalanceStatusArray) EnumValid(v interface{}) (err error) {
	var target BalanceStatus
	targetType := reflect.TypeOf(BalanceStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(BalanceStatus)
	}
	for _, value := range BalanceStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", BalanceStatusAll)
}

//DbArray will join value to database array
func (o BalanceStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o BalanceStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithBalance will return gex_balance meta data
func MetaWithBalance(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_balance"), fields...)
	return
}

//MetaWith will return gex_balance meta data
func (balance *Balance) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_balance"), fields...)
	return
}

//Meta will return gex_balance meta data
func (balance *Balance) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(balance, "#all")
	return
}

//Valid will valid by filter
func (balance *Balance) Valid() (err error) {
	if reflect.ValueOf(balance.TID).IsZero() {
		err = attrvalid.Valid(balance, BalanceFilterInsert+"#all", BalanceFilterOptional)
	} else {
		err = attrvalid.Valid(balance, BalanceFilterUpdate, "")
	}
	return
}

//Insert will add gex_balance to database
func (balance *Balance) Insert(caller interface{}, ctx context.Context) (err error) {

	if balance.UpdateTime.Timestamp() < 1 {
		balance.UpdateTime = xsql.TimeNow()
	}

	if balance.CreateTime.Timestamp() < 1 {
		balance.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, balance, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_balance to database
func (balance *Balance) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = balance.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_balance to database
func (balance *Balance) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = balance.UpdateFilterWheref(caller, ctx, BalanceFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_balance to database
func (balance *Balance) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	balance.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(balance, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", balance.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, balance, sql, where, "and", args)
	return
}

//AddBalance will add gex_balance to database
func AddBalance(ctx context.Context, balance *Balance) (err error) {
	err = AddBalanceCall(GetQueryer, ctx, balance)
	return
}

//AddBalance will add gex_balance to database
func AddBalanceCall(caller interface{}, ctx context.Context, balance *Balance) (err error) {
	err = balance.Insert(caller, ctx)
	return
}

//UpdateBalanceFilter will update gex_balance to database
func UpdateBalanceFilter(ctx context.Context, balance *Balance, filter string) (err error) {
	err = UpdateBalanceFilterCall(GetQueryer, ctx, balance, filter)
	return
}

//UpdateBalanceFilterCall will update gex_balance to database
func UpdateBalanceFilterCall(caller interface{}, ctx context.Context, balance *Balance, filter string) (err error) {
	err = balance.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateBalanceWheref will update gex_balance to database
func UpdateBalanceWheref(ctx context.Context, balance *Balance, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateBalanceWherefCall(GetQueryer, ctx, balance, formats, formatArgs...)
	return
}

//UpdateBalanceWherefCall will update gex_balance to database
func UpdateBalanceWherefCall(caller interface{}, ctx context.Context, balance *Balance, formats string, formatArgs ...interface{}) (err error) {
	err = balance.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateBalanceFilterWheref will update gex_balance to database
func UpdateBalanceFilterWheref(ctx context.Context, balance *Balance, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateBalanceFilterWherefCall(GetQueryer, ctx, balance, filter, formats, formatArgs...)
	return
}

//UpdateBalanceFilterWherefCall will update gex_balance to database
func UpdateBalanceFilterWherefCall(caller interface{}, ctx context.Context, balance *Balance, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = balance.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindBalanceCall will find gex_balance by id from database
func FindBalance(ctx context.Context, balanceID int64) (balance *Balance, err error) {
	balance, err = FindBalanceCall(GetQueryer, ctx, balanceID, false)
	return
}

//FindBalanceCall will find gex_balance by id from database
func FindBalanceCall(caller interface{}, ctx context.Context, balanceID int64, lock bool) (balance *Balance, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", balanceID)
	balance, err = FindBalanceWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindBalanceWhereCall will find gex_balance by where from database
func FindBalanceWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (balance *Balance, err error) {
	querySQL := crud.QuerySQL(&Balance{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Balance{}, "#all", querySQL, args, &balance)
	return
}

//FindBalanceWheref will find gex_balance by where from database
func FindBalanceWheref(ctx context.Context, format string, args ...interface{}) (balance *Balance, err error) {
	balance, err = FindBalanceWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindBalanceWherefCall will find gex_balance by where from database
func FindBalanceWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (balance *Balance, err error) {
	balance, err = FindBalanceFilterWherefCall(GetQueryer, ctx, lock, "#all", format, args...)
	return
}

//FindBalanceFilterWheref will find gex_balance by where from database
func FindBalanceFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (balance *Balance, err error) {
	balance, err = FindBalanceFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindBalanceFilterWherefCall will find gex_balance by where from database
func FindBalanceFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (balance *Balance, err error) {
	querySQL := crud.QuerySQL(&Balance{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Balance{}, filter, querySQL, queryArgs, &balance)
	return
}

//ListBalanceByID will list gex_balance by id from database
func ListBalanceByID(ctx context.Context, balanceIDs ...int64) (balanceList []*Balance, balanceMap map[int64]*Balance, err error) {
	balanceList, balanceMap, err = ListBalanceByIDCall(GetQueryer, ctx, balanceIDs...)
	return
}

//ListBalanceByIDCall will list gex_balance by id from database
func ListBalanceByIDCall(caller interface{}, ctx context.Context, balanceIDs ...int64) (balanceList []*Balance, balanceMap map[int64]*Balance, err error) {
	if len(balanceIDs) < 1 {
		balanceMap = map[int64]*Balance{}
		return
	}
	err = ScanBalanceByIDCall(caller, ctx, balanceIDs, &balanceList, &balanceMap, "tid")
	return
}

//ListBalanceFilterByID will list gex_balance by id from database
func ListBalanceFilterByID(ctx context.Context, filter string, balanceIDs ...int64) (balanceList []*Balance, balanceMap map[int64]*Balance, err error) {
	balanceList, balanceMap, err = ListBalanceFilterByIDCall(GetQueryer, ctx, filter, balanceIDs...)
	return
}

//ListBalanceFilterByIDCall will list gex_balance by id from database
func ListBalanceFilterByIDCall(caller interface{}, ctx context.Context, filter string, balanceIDs ...int64) (balanceList []*Balance, balanceMap map[int64]*Balance, err error) {
	if len(balanceIDs) < 1 {
		balanceMap = map[int64]*Balance{}
		return
	}
	err = ScanBalanceFilterByIDCall(caller, ctx, filter, balanceIDs, &balanceList, &balanceMap, "tid")
	return
}

//ScanBalanceByID will list gex_balance by id from database
func ScanBalanceByID(ctx context.Context, balanceIDs []int64, dest ...interface{}) (err error) {
	err = ScanBalanceByIDCall(GetQueryer, ctx, balanceIDs, dest...)
	return
}

//ScanBalanceByIDCall will list gex_balance by id from database
func ScanBalanceByIDCall(caller interface{}, ctx context.Context, balanceIDs []int64, dest ...interface{}) (err error) {
	err = ScanBalanceFilterByIDCall(caller, ctx, "#all", balanceIDs, dest...)
	return
}

//ScanBalanceFilterByID will list gex_balance by id from database
func ScanBalanceFilterByID(ctx context.Context, filter string, balanceIDs []int64, dest ...interface{}) (err error) {
	err = ScanBalanceFilterByIDCall(GetQueryer, ctx, filter, balanceIDs, dest...)
	return
}

//ScanBalanceFilterByIDCall will list gex_balance by id from database
func ScanBalanceFilterByIDCall(caller interface{}, ctx context.Context, filter string, balanceIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Balance{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(balanceIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Balance{}, filter, querySQL, nil, dest...)
	return
}

//ScanBalanceWherefCall will list gex_balance by format from database
func ScanBalanceWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanBalanceWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanBalanceWherefCall will list gex_balance by format from database
func ScanBalanceWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanBalanceFilterWherefCall(caller, ctx, "#all", format, args, suffix, dest...)
	return
}

//ScanBalanceFilterWheref will list gex_balance by format from database
func ScanBalanceFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanBalanceFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanBalanceFilterWherefCall will list gex_balance by format from database
func ScanBalanceFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Balance{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &Balance{}, filter, querySQL, args, dest...)
	return
}

//BalanceHistoryFilterOptional is crud filter
const BalanceHistoryFilterOptional = ""

//BalanceHistoryFilterRequired is crud filter
const BalanceHistoryFilterRequired = ""

//BalanceHistoryFilterInsert is crud filter
const BalanceHistoryFilterInsert = ""

//BalanceHistoryFilterUpdate is crud filter
const BalanceHistoryFilterUpdate = "update_time"

//BalanceHistoryFilterFind is crud filter
const BalanceHistoryFilterFind = "#all"

//BalanceHistoryFilterScan is crud filter
const BalanceHistoryFilterScan = "#all"

//EnumValid will valid value by BalanceHistoryStatus
func (o *BalanceHistoryStatus) EnumValid(v interface{}) (err error) {
	var target BalanceHistoryStatus
	targetType := reflect.TypeOf(BalanceHistoryStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(BalanceHistoryStatus)
	}
	for _, value := range BalanceHistoryStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", BalanceHistoryStatusAll)
}

//EnumValid will valid value by BalanceHistoryStatusArray
func (o *BalanceHistoryStatusArray) EnumValid(v interface{}) (err error) {
	var target BalanceHistoryStatus
	targetType := reflect.TypeOf(BalanceHistoryStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(BalanceHistoryStatus)
	}
	for _, value := range BalanceHistoryStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", BalanceHistoryStatusAll)
}

//DbArray will join value to database array
func (o BalanceHistoryStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o BalanceHistoryStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithBalanceHistory will return gex_balance_history meta data
func MetaWithBalanceHistory(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_balance_history"), fields...)
	return
}

//MetaWith will return gex_balance_history meta data
func (balanceHistory *BalanceHistory) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_balance_history"), fields...)
	return
}

//Meta will return gex_balance_history meta data
func (balanceHistory *BalanceHistory) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(balanceHistory, "#all")
	return
}

//Valid will valid by filter
func (balanceHistory *BalanceHistory) Valid() (err error) {
	if reflect.ValueOf(balanceHistory.TID).IsZero() {
		err = attrvalid.Valid(balanceHistory, BalanceHistoryFilterInsert+"#all", BalanceHistoryFilterOptional)
	} else {
		err = attrvalid.Valid(balanceHistory, BalanceHistoryFilterUpdate, "")
	}
	return
}

//Insert will add gex_balance_history to database
func (balanceHistory *BalanceHistory) Insert(caller interface{}, ctx context.Context) (err error) {

	if balanceHistory.UpdateTime.Timestamp() < 1 {
		balanceHistory.UpdateTime = xsql.TimeNow()
	}

	if balanceHistory.CreateTime.Timestamp() < 1 {
		balanceHistory.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, balanceHistory, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_balance_history to database
func (balanceHistory *BalanceHistory) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = balanceHistory.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_balance_history to database
func (balanceHistory *BalanceHistory) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = balanceHistory.UpdateFilterWheref(caller, ctx, BalanceHistoryFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_balance_history to database
func (balanceHistory *BalanceHistory) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	balanceHistory.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(balanceHistory, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", balanceHistory.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, balanceHistory, sql, where, "and", args)
	return
}

//AddBalanceHistory will add gex_balance_history to database
func AddBalanceHistory(ctx context.Context, balanceHistory *BalanceHistory) (err error) {
	err = AddBalanceHistoryCall(GetQueryer, ctx, balanceHistory)
	return
}

//AddBalanceHistory will add gex_balance_history to database
func AddBalanceHistoryCall(caller interface{}, ctx context.Context, balanceHistory *BalanceHistory) (err error) {
	err = balanceHistory.Insert(caller, ctx)
	return
}

//UpdateBalanceHistoryFilter will update gex_balance_history to database
func UpdateBalanceHistoryFilter(ctx context.Context, balanceHistory *BalanceHistory, filter string) (err error) {
	err = UpdateBalanceHistoryFilterCall(GetQueryer, ctx, balanceHistory, filter)
	return
}

//UpdateBalanceHistoryFilterCall will update gex_balance_history to database
func UpdateBalanceHistoryFilterCall(caller interface{}, ctx context.Context, balanceHistory *BalanceHistory, filter string) (err error) {
	err = balanceHistory.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateBalanceHistoryWheref will update gex_balance_history to database
func UpdateBalanceHistoryWheref(ctx context.Context, balanceHistory *BalanceHistory, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateBalanceHistoryWherefCall(GetQueryer, ctx, balanceHistory, formats, formatArgs...)
	return
}

//UpdateBalanceHistoryWherefCall will update gex_balance_history to database
func UpdateBalanceHistoryWherefCall(caller interface{}, ctx context.Context, balanceHistory *BalanceHistory, formats string, formatArgs ...interface{}) (err error) {
	err = balanceHistory.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateBalanceHistoryFilterWheref will update gex_balance_history to database
func UpdateBalanceHistoryFilterWheref(ctx context.Context, balanceHistory *BalanceHistory, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateBalanceHistoryFilterWherefCall(GetQueryer, ctx, balanceHistory, filter, formats, formatArgs...)
	return
}

//UpdateBalanceHistoryFilterWherefCall will update gex_balance_history to database
func UpdateBalanceHistoryFilterWherefCall(caller interface{}, ctx context.Context, balanceHistory *BalanceHistory, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = balanceHistory.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindBalanceHistoryCall will find gex_balance_history by id from database
func FindBalanceHistory(ctx context.Context, balanceHistoryID int64) (balanceHistory *BalanceHistory, err error) {
	balanceHistory, err = FindBalanceHistoryCall(GetQueryer, ctx, balanceHistoryID, false)
	return
}

//FindBalanceHistoryCall will find gex_balance_history by id from database
func FindBalanceHistoryCall(caller interface{}, ctx context.Context, balanceHistoryID int64, lock bool) (balanceHistory *BalanceHistory, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", balanceHistoryID)
	balanceHistory, err = FindBalanceHistoryWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindBalanceHistoryWhereCall will find gex_balance_history by where from database
func FindBalanceHistoryWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (balanceHistory *BalanceHistory, err error) {
	querySQL := crud.QuerySQL(&BalanceHistory{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &BalanceHistory{}, "#all", querySQL, args, &balanceHistory)
	return
}

//FindBalanceHistoryWheref will find gex_balance_history by where from database
func FindBalanceHistoryWheref(ctx context.Context, format string, args ...interface{}) (balanceHistory *BalanceHistory, err error) {
	balanceHistory, err = FindBalanceHistoryWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindBalanceHistoryWherefCall will find gex_balance_history by where from database
func FindBalanceHistoryWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (balanceHistory *BalanceHistory, err error) {
	balanceHistory, err = FindBalanceHistoryFilterWherefCall(GetQueryer, ctx, lock, "#all", format, args...)
	return
}

//FindBalanceHistoryFilterWheref will find gex_balance_history by where from database
func FindBalanceHistoryFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (balanceHistory *BalanceHistory, err error) {
	balanceHistory, err = FindBalanceHistoryFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindBalanceHistoryFilterWherefCall will find gex_balance_history by where from database
func FindBalanceHistoryFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (balanceHistory *BalanceHistory, err error) {
	querySQL := crud.QuerySQL(&BalanceHistory{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &BalanceHistory{}, filter, querySQL, queryArgs, &balanceHistory)
	return
}

//ListBalanceHistoryByID will list gex_balance_history by id from database
func ListBalanceHistoryByID(ctx context.Context, balanceHistoryIDs ...int64) (balanceHistoryList []*BalanceHistory, balanceHistoryMap map[int64]*BalanceHistory, err error) {
	balanceHistoryList, balanceHistoryMap, err = ListBalanceHistoryByIDCall(GetQueryer, ctx, balanceHistoryIDs...)
	return
}

//ListBalanceHistoryByIDCall will list gex_balance_history by id from database
func ListBalanceHistoryByIDCall(caller interface{}, ctx context.Context, balanceHistoryIDs ...int64) (balanceHistoryList []*BalanceHistory, balanceHistoryMap map[int64]*BalanceHistory, err error) {
	if len(balanceHistoryIDs) < 1 {
		balanceHistoryMap = map[int64]*BalanceHistory{}
		return
	}
	err = ScanBalanceHistoryByIDCall(caller, ctx, balanceHistoryIDs, &balanceHistoryList, &balanceHistoryMap, "tid")
	return
}

//ListBalanceHistoryFilterByID will list gex_balance_history by id from database
func ListBalanceHistoryFilterByID(ctx context.Context, filter string, balanceHistoryIDs ...int64) (balanceHistoryList []*BalanceHistory, balanceHistoryMap map[int64]*BalanceHistory, err error) {
	balanceHistoryList, balanceHistoryMap, err = ListBalanceHistoryFilterByIDCall(GetQueryer, ctx, filter, balanceHistoryIDs...)
	return
}

//ListBalanceHistoryFilterByIDCall will list gex_balance_history by id from database
func ListBalanceHistoryFilterByIDCall(caller interface{}, ctx context.Context, filter string, balanceHistoryIDs ...int64) (balanceHistoryList []*BalanceHistory, balanceHistoryMap map[int64]*BalanceHistory, err error) {
	if len(balanceHistoryIDs) < 1 {
		balanceHistoryMap = map[int64]*BalanceHistory{}
		return
	}
	err = ScanBalanceHistoryFilterByIDCall(caller, ctx, filter, balanceHistoryIDs, &balanceHistoryList, &balanceHistoryMap, "tid")
	return
}

//ScanBalanceHistoryByID will list gex_balance_history by id from database
func ScanBalanceHistoryByID(ctx context.Context, balanceHistoryIDs []int64, dest ...interface{}) (err error) {
	err = ScanBalanceHistoryByIDCall(GetQueryer, ctx, balanceHistoryIDs, dest...)
	return
}

//ScanBalanceHistoryByIDCall will list gex_balance_history by id from database
func ScanBalanceHistoryByIDCall(caller interface{}, ctx context.Context, balanceHistoryIDs []int64, dest ...interface{}) (err error) {
	err = ScanBalanceHistoryFilterByIDCall(caller, ctx, "#all", balanceHistoryIDs, dest...)
	return
}

//ScanBalanceHistoryFilterByID will list gex_balance_history by id from database
func ScanBalanceHistoryFilterByID(ctx context.Context, filter string, balanceHistoryIDs []int64, dest ...interface{}) (err error) {
	err = ScanBalanceHistoryFilterByIDCall(GetQueryer, ctx, filter, balanceHistoryIDs, dest...)
	return
}

//ScanBalanceHistoryFilterByIDCall will list gex_balance_history by id from database
func ScanBalanceHistoryFilterByIDCall(caller interface{}, ctx context.Context, filter string, balanceHistoryIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&BalanceHistory{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(balanceHistoryIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &BalanceHistory{}, filter, querySQL, nil, dest...)
	return
}

//ScanBalanceHistoryWherefCall will list gex_balance_history by format from database
func ScanBalanceHistoryWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanBalanceHistoryWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanBalanceHistoryWherefCall will list gex_balance_history by format from database
func ScanBalanceHistoryWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanBalanceHistoryFilterWherefCall(caller, ctx, "#all", format, args, suffix, dest...)
	return
}

//ScanBalanceHistoryFilterWheref will list gex_balance_history by format from database
func ScanBalanceHistoryFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanBalanceHistoryFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanBalanceHistoryFilterWherefCall will list gex_balance_history by format from database
func ScanBalanceHistoryFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&BalanceHistory{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &BalanceHistory{}, filter, querySQL, args, dest...)
	return
}

//HoldingFilterOptional is crud filter
const HoldingFilterOptional = ""

//HoldingFilterRequired is crud filter
const HoldingFilterRequired = ""

//HoldingFilterInsert is crud filter
const HoldingFilterInsert = ""

//HoldingFilterUpdate is crud filter
const HoldingFilterUpdate = "update_time"

//HoldingFilterFind is crud filter
const HoldingFilterFind = "#all"

//HoldingFilterScan is crud filter
const HoldingFilterScan = "#all"

//EnumValid will valid value by HoldingStatus
func (o *HoldingStatus) EnumValid(v interface{}) (err error) {
	var target HoldingStatus
	targetType := reflect.TypeOf(HoldingStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(HoldingStatus)
	}
	for _, value := range HoldingStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", HoldingStatusAll)
}

//EnumValid will valid value by HoldingStatusArray
func (o *HoldingStatusArray) EnumValid(v interface{}) (err error) {
	var target HoldingStatus
	targetType := reflect.TypeOf(HoldingStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(HoldingStatus)
	}
	for _, value := range HoldingStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", HoldingStatusAll)
}

//DbArray will join value to database array
func (o HoldingStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o HoldingStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithHolding will return gex_holding meta data
func MetaWithHolding(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_holding"), fields...)
	return
}

//MetaWith will return gex_holding meta data
func (holding *Holding) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_holding"), fields...)
	return
}

//Meta will return gex_holding meta data
func (holding *Holding) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(holding, "#all")
	return
}

//Valid will valid by filter
func (holding *Holding) Valid() (err error) {
	if reflect.ValueOf(holding.TID).IsZero() {
		err = attrvalid.Valid(holding, HoldingFilterInsert+"#all", HoldingFilterOptional)
	} else {
		err = attrvalid.Valid(holding, HoldingFilterUpdate, "")
	}
	return
}

//Insert will add gex_holding to database
func (holding *Holding) Insert(caller interface{}, ctx context.Context) (err error) {

	if holding.UpdateTime.Timestamp() < 1 {
		holding.UpdateTime = xsql.TimeNow()
	}

	if holding.CreateTime.Timestamp() < 1 {
		holding.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, holding, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_holding to database
func (holding *Holding) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = holding.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_holding to database
func (holding *Holding) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = holding.UpdateFilterWheref(caller, ctx, HoldingFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_holding to database
func (holding *Holding) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	holding.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(holding, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", holding.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, holding, sql, where, "and", args)
	return
}

//UpdateHoldingFilter will update gex_holding to database
func UpdateHoldingFilter(ctx context.Context, holding *Holding, filter string) (err error) {
	err = UpdateHoldingFilterCall(GetQueryer, ctx, holding, filter)
	return
}

//UpdateHoldingFilterCall will update gex_holding to database
func UpdateHoldingFilterCall(caller interface{}, ctx context.Context, holding *Holding, filter string) (err error) {
	err = holding.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateHoldingWheref will update gex_holding to database
func UpdateHoldingWheref(ctx context.Context, holding *Holding, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateHoldingWherefCall(GetQueryer, ctx, holding, formats, formatArgs...)
	return
}

//UpdateHoldingWherefCall will update gex_holding to database
func UpdateHoldingWherefCall(caller interface{}, ctx context.Context, holding *Holding, formats string, formatArgs ...interface{}) (err error) {
	err = holding.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateHoldingFilterWheref will update gex_holding to database
func UpdateHoldingFilterWheref(ctx context.Context, holding *Holding, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateHoldingFilterWherefCall(GetQueryer, ctx, holding, filter, formats, formatArgs...)
	return
}

//UpdateHoldingFilterWherefCall will update gex_holding to database
func UpdateHoldingFilterWherefCall(caller interface{}, ctx context.Context, holding *Holding, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = holding.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindHoldingCall will find gex_holding by id from database
func FindHolding(ctx context.Context, holdingID int64) (holding *Holding, err error) {
	holding, err = FindHoldingCall(GetQueryer, ctx, holdingID, false)
	return
}

//FindHoldingCall will find gex_holding by id from database
func FindHoldingCall(caller interface{}, ctx context.Context, holdingID int64, lock bool) (holding *Holding, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", holdingID)
	holding, err = FindHoldingWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindHoldingWhereCall will find gex_holding by where from database
func FindHoldingWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (holding *Holding, err error) {
	querySQL := crud.QuerySQL(&Holding{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Holding{}, "#all", querySQL, args, &holding)
	return
}

//FindHoldingWheref will find gex_holding by where from database
func FindHoldingWheref(ctx context.Context, format string, args ...interface{}) (holding *Holding, err error) {
	holding, err = FindHoldingWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindHoldingWherefCall will find gex_holding by where from database
func FindHoldingWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (holding *Holding, err error) {
	holding, err = FindHoldingFilterWherefCall(GetQueryer, ctx, lock, "#all", format, args...)
	return
}

//FindHoldingFilterWheref will find gex_holding by where from database
func FindHoldingFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (holding *Holding, err error) {
	holding, err = FindHoldingFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindHoldingFilterWherefCall will find gex_holding by where from database
func FindHoldingFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (holding *Holding, err error) {
	querySQL := crud.QuerySQL(&Holding{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Holding{}, filter, querySQL, queryArgs, &holding)
	return
}

//ListHoldingByID will list gex_holding by id from database
func ListHoldingByID(ctx context.Context, holdingIDs ...int64) (holdingList []*Holding, holdingMap map[int64]*Holding, err error) {
	holdingList, holdingMap, err = ListHoldingByIDCall(GetQueryer, ctx, holdingIDs...)
	return
}

//ListHoldingByIDCall will list gex_holding by id from database
func ListHoldingByIDCall(caller interface{}, ctx context.Context, holdingIDs ...int64) (holdingList []*Holding, holdingMap map[int64]*Holding, err error) {
	if len(holdingIDs) < 1 {
		holdingMap = map[int64]*Holding{}
		return
	}
	err = ScanHoldingByIDCall(caller, ctx, holdingIDs, &holdingList, &holdingMap, "tid")
	return
}

//ListHoldingFilterByID will list gex_holding by id from database
func ListHoldingFilterByID(ctx context.Context, filter string, holdingIDs ...int64) (holdingList []*Holding, holdingMap map[int64]*Holding, err error) {
	holdingList, holdingMap, err = ListHoldingFilterByIDCall(GetQueryer, ctx, filter, holdingIDs...)
	return
}

//ListHoldingFilterByIDCall will list gex_holding by id from database
func ListHoldingFilterByIDCall(caller interface{}, ctx context.Context, filter string, holdingIDs ...int64) (holdingList []*Holding, holdingMap map[int64]*Holding, err error) {
	if len(holdingIDs) < 1 {
		holdingMap = map[int64]*Holding{}
		return
	}
	err = ScanHoldingFilterByIDCall(caller, ctx, filter, holdingIDs, &holdingList, &holdingMap, "tid")
	return
}

//ScanHoldingByID will list gex_holding by id from database
func ScanHoldingByID(ctx context.Context, holdingIDs []int64, dest ...interface{}) (err error) {
	err = ScanHoldingByIDCall(GetQueryer, ctx, holdingIDs, dest...)
	return
}

//ScanHoldingByIDCall will list gex_holding by id from database
func ScanHoldingByIDCall(caller interface{}, ctx context.Context, holdingIDs []int64, dest ...interface{}) (err error) {
	err = ScanHoldingFilterByIDCall(caller, ctx, "#all", holdingIDs, dest...)
	return
}

//ScanHoldingFilterByID will list gex_holding by id from database
func ScanHoldingFilterByID(ctx context.Context, filter string, holdingIDs []int64, dest ...interface{}) (err error) {
	err = ScanHoldingFilterByIDCall(GetQueryer, ctx, filter, holdingIDs, dest...)
	return
}

//ScanHoldingFilterByIDCall will list gex_holding by id from database
func ScanHoldingFilterByIDCall(caller interface{}, ctx context.Context, filter string, holdingIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Holding{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(holdingIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Holding{}, filter, querySQL, nil, dest...)
	return
}

//ScanHoldingWherefCall will list gex_holding by format from database
func ScanHoldingWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanHoldingWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanHoldingWherefCall will list gex_holding by format from database
func ScanHoldingWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanHoldingFilterWherefCall(caller, ctx, "#all", format, args, suffix, dest...)
	return
}

//ScanHoldingFilterWheref will list gex_holding by format from database
func ScanHoldingFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanHoldingFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanHoldingFilterWherefCall will list gex_holding by format from database
func ScanHoldingFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Holding{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &Holding{}, filter, querySQL, args, dest...)
	return
}

//KLineFilterOptional is crud filter
const KLineFilterOptional = ""

//KLineFilterRequired is crud filter
const KLineFilterRequired = ""

//KLineFilterInsert is crud filter
const KLineFilterInsert = ""

//KLineFilterUpdate is crud filter
const KLineFilterUpdate = "update_time"

//KLineFilterFind is crud filter
const KLineFilterFind = "#all"

//KLineFilterScan is crud filter
const KLineFilterScan = "#all"

//MetaWithKLine will return gex_kline meta data
func MetaWithKLine(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_kline"), fields...)
	return
}

//MetaWith will return gex_kline meta data
func (kLine *KLine) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_kline"), fields...)
	return
}

//Meta will return gex_kline meta data
func (kLine *KLine) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(kLine, "#all")
	return
}

//Valid will valid by filter
func (kLine *KLine) Valid() (err error) {
	if reflect.ValueOf(kLine.TID).IsZero() {
		err = attrvalid.Valid(kLine, KLineFilterInsert+"#all", KLineFilterOptional)
	} else {
		err = attrvalid.Valid(kLine, KLineFilterUpdate, "")
	}
	return
}

//Insert will add gex_kline to database
func (kLine *KLine) Insert(caller interface{}, ctx context.Context) (err error) {

	if kLine.UpdateTime.Timestamp() < 1 {
		kLine.UpdateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, kLine, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_kline to database
func (kLine *KLine) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = kLine.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_kline to database
func (kLine *KLine) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = kLine.UpdateFilterWheref(caller, ctx, KLineFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_kline to database
func (kLine *KLine) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	kLine.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(kLine, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", kLine.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, kLine, sql, where, "and", args)
	return
}

//AddKLine will add gex_kline to database
func AddKLine(ctx context.Context, kLine *KLine) (err error) {
	err = AddKLineCall(GetQueryer, ctx, kLine)
	return
}

//AddKLine will add gex_kline to database
func AddKLineCall(caller interface{}, ctx context.Context, kLine *KLine) (err error) {
	err = kLine.Insert(caller, ctx)
	return
}

//UpdateKLineFilter will update gex_kline to database
func UpdateKLineFilter(ctx context.Context, kLine *KLine, filter string) (err error) {
	err = UpdateKLineFilterCall(GetQueryer, ctx, kLine, filter)
	return
}

//UpdateKLineFilterCall will update gex_kline to database
func UpdateKLineFilterCall(caller interface{}, ctx context.Context, kLine *KLine, filter string) (err error) {
	err = kLine.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateKLineWheref will update gex_kline to database
func UpdateKLineWheref(ctx context.Context, kLine *KLine, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateKLineWherefCall(GetQueryer, ctx, kLine, formats, formatArgs...)
	return
}

//UpdateKLineWherefCall will update gex_kline to database
func UpdateKLineWherefCall(caller interface{}, ctx context.Context, kLine *KLine, formats string, formatArgs ...interface{}) (err error) {
	err = kLine.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateKLineFilterWheref will update gex_kline to database
func UpdateKLineFilterWheref(ctx context.Context, kLine *KLine, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateKLineFilterWherefCall(GetQueryer, ctx, kLine, filter, formats, formatArgs...)
	return
}

//UpdateKLineFilterWherefCall will update gex_kline to database
func UpdateKLineFilterWherefCall(caller interface{}, ctx context.Context, kLine *KLine, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = kLine.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindKLineCall will find gex_kline by id from database
func FindKLine(ctx context.Context, kLineID int64) (kLine *KLine, err error) {
	kLine, err = FindKLineCall(GetQueryer, ctx, kLineID, false)
	return
}

//FindKLineCall will find gex_kline by id from database
func FindKLineCall(caller interface{}, ctx context.Context, kLineID int64, lock bool) (kLine *KLine, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", kLineID)
	kLine, err = FindKLineWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindKLineWhereCall will find gex_kline by where from database
func FindKLineWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (kLine *KLine, err error) {
	querySQL := crud.QuerySQL(&KLine{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &KLine{}, "#all", querySQL, args, &kLine)
	return
}

//FindKLineWheref will find gex_kline by where from database
func FindKLineWheref(ctx context.Context, format string, args ...interface{}) (kLine *KLine, err error) {
	kLine, err = FindKLineWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindKLineWherefCall will find gex_kline by where from database
func FindKLineWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (kLine *KLine, err error) {
	kLine, err = FindKLineFilterWherefCall(GetQueryer, ctx, lock, "#all", format, args...)
	return
}

//FindKLineFilterWheref will find gex_kline by where from database
func FindKLineFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (kLine *KLine, err error) {
	kLine, err = FindKLineFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindKLineFilterWherefCall will find gex_kline by where from database
func FindKLineFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (kLine *KLine, err error) {
	querySQL := crud.QuerySQL(&KLine{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &KLine{}, filter, querySQL, queryArgs, &kLine)
	return
}

//ListKLineByID will list gex_kline by id from database
func ListKLineByID(ctx context.Context, kLineIDs ...int64) (kLineList []*KLine, kLineMap map[int64]*KLine, err error) {
	kLineList, kLineMap, err = ListKLineByIDCall(GetQueryer, ctx, kLineIDs...)
	return
}

//ListKLineByIDCall will list gex_kline by id from database
func ListKLineByIDCall(caller interface{}, ctx context.Context, kLineIDs ...int64) (kLineList []*KLine, kLineMap map[int64]*KLine, err error) {
	if len(kLineIDs) < 1 {
		kLineMap = map[int64]*KLine{}
		return
	}
	err = ScanKLineByIDCall(caller, ctx, kLineIDs, &kLineList, &kLineMap, "tid")
	return
}

//ListKLineFilterByID will list gex_kline by id from database
func ListKLineFilterByID(ctx context.Context, filter string, kLineIDs ...int64) (kLineList []*KLine, kLineMap map[int64]*KLine, err error) {
	kLineList, kLineMap, err = ListKLineFilterByIDCall(GetQueryer, ctx, filter, kLineIDs...)
	return
}

//ListKLineFilterByIDCall will list gex_kline by id from database
func ListKLineFilterByIDCall(caller interface{}, ctx context.Context, filter string, kLineIDs ...int64) (kLineList []*KLine, kLineMap map[int64]*KLine, err error) {
	if len(kLineIDs) < 1 {
		kLineMap = map[int64]*KLine{}
		return
	}
	err = ScanKLineFilterByIDCall(caller, ctx, filter, kLineIDs, &kLineList, &kLineMap, "tid")
	return
}

//ScanKLineByID will list gex_kline by id from database
func ScanKLineByID(ctx context.Context, kLineIDs []int64, dest ...interface{}) (err error) {
	err = ScanKLineByIDCall(GetQueryer, ctx, kLineIDs, dest...)
	return
}

//ScanKLineByIDCall will list gex_kline by id from database
func ScanKLineByIDCall(caller interface{}, ctx context.Context, kLineIDs []int64, dest ...interface{}) (err error) {
	err = ScanKLineFilterByIDCall(caller, ctx, "#all", kLineIDs, dest...)
	return
}

//ScanKLineFilterByID will list gex_kline by id from database
func ScanKLineFilterByID(ctx context.Context, filter string, kLineIDs []int64, dest ...interface{}) (err error) {
	err = ScanKLineFilterByIDCall(GetQueryer, ctx, filter, kLineIDs, dest...)
	return
}

//ScanKLineFilterByIDCall will list gex_kline by id from database
func ScanKLineFilterByIDCall(caller interface{}, ctx context.Context, filter string, kLineIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&KLine{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(kLineIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &KLine{}, filter, querySQL, nil, dest...)
	return
}

//ScanKLineWherefCall will list gex_kline by format from database
func ScanKLineWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanKLineWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanKLineWherefCall will list gex_kline by format from database
func ScanKLineWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanKLineFilterWherefCall(caller, ctx, "#all", format, args, suffix, dest...)
	return
}

//ScanKLineFilterWheref will list gex_kline by format from database
func ScanKLineFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanKLineFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanKLineFilterWherefCall will list gex_kline by format from database
func ScanKLineFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&KLine{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &KLine{}, filter, querySQL, args, dest...)
	return
}

//OrderFilterOptional is crud filter
const OrderFilterOptional = "tid,quantity,price,total_price,trigger_type,trigger_price,status"

//OrderFilterRequired is crud filter
const OrderFilterRequired = ""

//OrderFilterInsert is crud filter
const OrderFilterInsert = "tid,quantity,price,total_price,trigger_type,trigger_price,status"

//OrderFilterUpdate is crud filter
const OrderFilterUpdate = "update_time,tid,quantity,price,total_price,trigger_type,trigger_price,status"

//OrderFilterFind is crud filter
const OrderFilterFind = "#all"

//OrderFilterScan is crud filter
const OrderFilterScan = "^transaction#all"

//EnumValid will valid value by OrderType
func (o *OrderType) EnumValid(v interface{}) (err error) {
	var target OrderType
	targetType := reflect.TypeOf(OrderType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderType)
	}
	for _, value := range OrderTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderTypeAll)
}

//EnumValid will valid value by OrderTypeArray
func (o *OrderTypeArray) EnumValid(v interface{}) (err error) {
	var target OrderType
	targetType := reflect.TypeOf(OrderType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderType)
	}
	for _, value := range OrderTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderTypeAll)
}

//DbArray will join value to database array
func (o OrderTypeArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o OrderTypeArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by OrderSide
func (o *OrderSide) EnumValid(v interface{}) (err error) {
	var target OrderSide
	targetType := reflect.TypeOf(OrderSide(""))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderSide)
	}
	for _, value := range OrderSideAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderSideAll)
}

//EnumValid will valid value by OrderSideArray
func (o *OrderSideArray) EnumValid(v interface{}) (err error) {
	var target OrderSide
	targetType := reflect.TypeOf(OrderSide(""))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderSide)
	}
	for _, value := range OrderSideAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderSideAll)
}

//DbArray will join value to database array
func (o OrderSideArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o OrderSideArray) InArray() (res string) {
	res = "'" + converter.JoinSafe(o, "','", converter.JoinPolicyDefault) + "'"
	return
}

//EnumValid will valid value by OrderTriggerType
func (o *OrderTriggerType) EnumValid(v interface{}) (err error) {
	var target OrderTriggerType
	targetType := reflect.TypeOf(OrderTriggerType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderTriggerType)
	}
	for _, value := range OrderTriggerTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderTriggerTypeAll)
}

//EnumValid will valid value by OrderTriggerTypeArray
func (o *OrderTriggerTypeArray) EnumValid(v interface{}) (err error) {
	var target OrderTriggerType
	targetType := reflect.TypeOf(OrderTriggerType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderTriggerType)
	}
	for _, value := range OrderTriggerTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderTriggerTypeAll)
}

//DbArray will join value to database array
func (o OrderTriggerTypeArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o OrderTriggerTypeArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by OrderStatus
func (o *OrderStatus) EnumValid(v interface{}) (err error) {
	var target OrderStatus
	targetType := reflect.TypeOf(OrderStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderStatus)
	}
	for _, value := range OrderStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderStatusAll)
}

//EnumValid will valid value by OrderStatusArray
func (o *OrderStatusArray) EnumValid(v interface{}) (err error) {
	var target OrderStatus
	targetType := reflect.TypeOf(OrderStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderStatus)
	}
	for _, value := range OrderStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderStatusAll)
}

//DbArray will join value to database array
func (o OrderStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o OrderStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithOrder will return gex_order meta data
func MetaWithOrder(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_order"), fields...)
	return
}

//MetaWith will return gex_order meta data
func (order *Order) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_order"), fields...)
	return
}

//Meta will return gex_order meta data
func (order *Order) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(order, "#all")
	return
}

//Valid will valid by filter
func (order *Order) Valid() (err error) {
	if reflect.ValueOf(order.TID).IsZero() {
		err = attrvalid.Valid(order, OrderFilterInsert+"#all", OrderFilterOptional)
	} else {
		err = attrvalid.Valid(order, OrderFilterUpdate, "")
	}
	return
}

//Insert will add gex_order to database
func (order *Order) Insert(caller interface{}, ctx context.Context) (err error) {

	if order.UpdateTime.Timestamp() < 1 {
		order.UpdateTime = xsql.TimeNow()
	}

	if order.CreateTime.Timestamp() < 1 {
		order.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, order, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_order to database
func (order *Order) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = order.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_order to database
func (order *Order) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = order.UpdateFilterWheref(caller, ctx, OrderFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_order to database
func (order *Order) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	order.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(order, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", order.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, order, sql, where, "and", args)
	return
}

//AddOrder will add gex_order to database
func AddOrder(ctx context.Context, order *Order) (err error) {
	err = AddOrderCall(GetQueryer, ctx, order)
	return
}

//AddOrder will add gex_order to database
func AddOrderCall(caller interface{}, ctx context.Context, order *Order) (err error) {
	err = order.Insert(caller, ctx)
	return
}

//UpdateOrderFilter will update gex_order to database
func UpdateOrderFilter(ctx context.Context, order *Order, filter string) (err error) {
	err = UpdateOrderFilterCall(GetQueryer, ctx, order, filter)
	return
}

//UpdateOrderFilterCall will update gex_order to database
func UpdateOrderFilterCall(caller interface{}, ctx context.Context, order *Order, filter string) (err error) {
	err = order.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateOrderWheref will update gex_order to database
func UpdateOrderWheref(ctx context.Context, order *Order, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateOrderWherefCall(GetQueryer, ctx, order, formats, formatArgs...)
	return
}

//UpdateOrderWherefCall will update gex_order to database
func UpdateOrderWherefCall(caller interface{}, ctx context.Context, order *Order, formats string, formatArgs ...interface{}) (err error) {
	err = order.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateOrderFilterWheref will update gex_order to database
func UpdateOrderFilterWheref(ctx context.Context, order *Order, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateOrderFilterWherefCall(GetQueryer, ctx, order, filter, formats, formatArgs...)
	return
}

//UpdateOrderFilterWherefCall will update gex_order to database
func UpdateOrderFilterWherefCall(caller interface{}, ctx context.Context, order *Order, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = order.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindOrderCall will find gex_order by id from database
func FindOrder(ctx context.Context, orderID int64) (order *Order, err error) {
	order, err = FindOrderCall(GetQueryer, ctx, orderID, false)
	return
}

//FindOrderCall will find gex_order by id from database
func FindOrderCall(caller interface{}, ctx context.Context, orderID int64, lock bool) (order *Order, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", orderID)
	order, err = FindOrderWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindOrderWhereCall will find gex_order by where from database
func FindOrderWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (order *Order, err error) {
	querySQL := crud.QuerySQL(&Order{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Order{}, "#all", querySQL, args, &order)
	return
}

//FindOrderWheref will find gex_order by where from database
func FindOrderWheref(ctx context.Context, format string, args ...interface{}) (order *Order, err error) {
	order, err = FindOrderWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindOrderWherefCall will find gex_order by where from database
func FindOrderWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (order *Order, err error) {
	order, err = FindOrderFilterWherefCall(GetQueryer, ctx, lock, "#all", format, args...)
	return
}

//FindOrderFilterWheref will find gex_order by where from database
func FindOrderFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (order *Order, err error) {
	order, err = FindOrderFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindOrderFilterWherefCall will find gex_order by where from database
func FindOrderFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (order *Order, err error) {
	querySQL := crud.QuerySQL(&Order{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Order{}, filter, querySQL, queryArgs, &order)
	return
}

//ListOrderByID will list gex_order by id from database
func ListOrderByID(ctx context.Context, orderIDs ...int64) (orderList []*Order, orderMap map[int64]*Order, err error) {
	orderList, orderMap, err = ListOrderByIDCall(GetQueryer, ctx, orderIDs...)
	return
}

//ListOrderByIDCall will list gex_order by id from database
func ListOrderByIDCall(caller interface{}, ctx context.Context, orderIDs ...int64) (orderList []*Order, orderMap map[int64]*Order, err error) {
	if len(orderIDs) < 1 {
		orderMap = map[int64]*Order{}
		return
	}
	err = ScanOrderByIDCall(caller, ctx, orderIDs, &orderList, &orderMap, "tid")
	return
}

//ListOrderFilterByID will list gex_order by id from database
func ListOrderFilterByID(ctx context.Context, filter string, orderIDs ...int64) (orderList []*Order, orderMap map[int64]*Order, err error) {
	orderList, orderMap, err = ListOrderFilterByIDCall(GetQueryer, ctx, filter, orderIDs...)
	return
}

//ListOrderFilterByIDCall will list gex_order by id from database
func ListOrderFilterByIDCall(caller interface{}, ctx context.Context, filter string, orderIDs ...int64) (orderList []*Order, orderMap map[int64]*Order, err error) {
	if len(orderIDs) < 1 {
		orderMap = map[int64]*Order{}
		return
	}
	err = ScanOrderFilterByIDCall(caller, ctx, filter, orderIDs, &orderList, &orderMap, "tid")
	return
}

//ScanOrderByID will list gex_order by id from database
func ScanOrderByID(ctx context.Context, orderIDs []int64, dest ...interface{}) (err error) {
	err = ScanOrderByIDCall(GetQueryer, ctx, orderIDs, dest...)
	return
}

//ScanOrderByIDCall will list gex_order by id from database
func ScanOrderByIDCall(caller interface{}, ctx context.Context, orderIDs []int64, dest ...interface{}) (err error) {
	err = ScanOrderFilterByIDCall(caller, ctx, "^transaction#all", orderIDs, dest...)
	return
}

//ScanOrderFilterByID will list gex_order by id from database
func ScanOrderFilterByID(ctx context.Context, filter string, orderIDs []int64, dest ...interface{}) (err error) {
	err = ScanOrderFilterByIDCall(GetQueryer, ctx, filter, orderIDs, dest...)
	return
}

//ScanOrderFilterByIDCall will list gex_order by id from database
func ScanOrderFilterByIDCall(caller interface{}, ctx context.Context, filter string, orderIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Order{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(orderIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Order{}, filter, querySQL, nil, dest...)
	return
}

//ScanOrderWherefCall will list gex_order by format from database
func ScanOrderWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanOrderWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanOrderWherefCall will list gex_order by format from database
func ScanOrderWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanOrderFilterWherefCall(caller, ctx, "^transaction#all", format, args, suffix, dest...)
	return
}

//ScanOrderFilterWheref will list gex_order by format from database
func ScanOrderFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanOrderFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanOrderFilterWherefCall will list gex_order by format from database
func ScanOrderFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Order{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &Order{}, filter, querySQL, args, dest...)
	return
}

//OrderCommFilterOptional is crud filter
const OrderCommFilterOptional = ""

//OrderCommFilterRequired is crud filter
const OrderCommFilterRequired = ""

//OrderCommFilterInsert is crud filter
const OrderCommFilterInsert = ""

//OrderCommFilterUpdate is crud filter
const OrderCommFilterUpdate = "update_time"

//OrderCommFilterFind is crud filter
const OrderCommFilterFind = "#all"

//OrderCommFilterScan is crud filter
const OrderCommFilterScan = "#all"

//EnumValid will valid value by OrderCommType
func (o *OrderCommType) EnumValid(v interface{}) (err error) {
	var target OrderCommType
	targetType := reflect.TypeOf(OrderCommType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderCommType)
	}
	for _, value := range OrderCommTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderCommTypeAll)
}

//EnumValid will valid value by OrderCommTypeArray
func (o *OrderCommTypeArray) EnumValid(v interface{}) (err error) {
	var target OrderCommType
	targetType := reflect.TypeOf(OrderCommType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderCommType)
	}
	for _, value := range OrderCommTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderCommTypeAll)
}

//DbArray will join value to database array
func (o OrderCommTypeArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o OrderCommTypeArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by OrderCommStatus
func (o *OrderCommStatus) EnumValid(v interface{}) (err error) {
	var target OrderCommStatus
	targetType := reflect.TypeOf(OrderCommStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderCommStatus)
	}
	for _, value := range OrderCommStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderCommStatusAll)
}

//EnumValid will valid value by OrderCommStatusArray
func (o *OrderCommStatusArray) EnumValid(v interface{}) (err error) {
	var target OrderCommStatus
	targetType := reflect.TypeOf(OrderCommStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(OrderCommStatus)
	}
	for _, value := range OrderCommStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", OrderCommStatusAll)
}

//DbArray will join value to database array
func (o OrderCommStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o OrderCommStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithOrderComm will return gex_order_comm meta data
func MetaWithOrderComm(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_order_comm"), fields...)
	return
}

//MetaWith will return gex_order_comm meta data
func (orderComm *OrderComm) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_order_comm"), fields...)
	return
}

//Meta will return gex_order_comm meta data
func (orderComm *OrderComm) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(orderComm, "#all")
	return
}

//Valid will valid by filter
func (orderComm *OrderComm) Valid() (err error) {
	if reflect.ValueOf(orderComm.TID).IsZero() {
		err = attrvalid.Valid(orderComm, OrderCommFilterInsert+"#all", OrderCommFilterOptional)
	} else {
		err = attrvalid.Valid(orderComm, OrderCommFilterUpdate, "")
	}
	return
}

//Insert will add gex_order_comm to database
func (orderComm *OrderComm) Insert(caller interface{}, ctx context.Context) (err error) {

	if orderComm.UpdateTime.Timestamp() < 1 {
		orderComm.UpdateTime = xsql.TimeNow()
	}

	if orderComm.CreateTime.Timestamp() < 1 {
		orderComm.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, orderComm, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_order_comm to database
func (orderComm *OrderComm) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = orderComm.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_order_comm to database
func (orderComm *OrderComm) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = orderComm.UpdateFilterWheref(caller, ctx, OrderCommFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_order_comm to database
func (orderComm *OrderComm) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	orderComm.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(orderComm, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", orderComm.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, orderComm, sql, where, "and", args)
	return
}

//AddOrderComm will add gex_order_comm to database
func AddOrderComm(ctx context.Context, orderComm *OrderComm) (err error) {
	err = AddOrderCommCall(GetQueryer, ctx, orderComm)
	return
}

//AddOrderComm will add gex_order_comm to database
func AddOrderCommCall(caller interface{}, ctx context.Context, orderComm *OrderComm) (err error) {
	err = orderComm.Insert(caller, ctx)
	return
}

//UpdateOrderCommFilter will update gex_order_comm to database
func UpdateOrderCommFilter(ctx context.Context, orderComm *OrderComm, filter string) (err error) {
	err = UpdateOrderCommFilterCall(GetQueryer, ctx, orderComm, filter)
	return
}

//UpdateOrderCommFilterCall will update gex_order_comm to database
func UpdateOrderCommFilterCall(caller interface{}, ctx context.Context, orderComm *OrderComm, filter string) (err error) {
	err = orderComm.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateOrderCommWheref will update gex_order_comm to database
func UpdateOrderCommWheref(ctx context.Context, orderComm *OrderComm, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateOrderCommWherefCall(GetQueryer, ctx, orderComm, formats, formatArgs...)
	return
}

//UpdateOrderCommWherefCall will update gex_order_comm to database
func UpdateOrderCommWherefCall(caller interface{}, ctx context.Context, orderComm *OrderComm, formats string, formatArgs ...interface{}) (err error) {
	err = orderComm.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateOrderCommFilterWheref will update gex_order_comm to database
func UpdateOrderCommFilterWheref(ctx context.Context, orderComm *OrderComm, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateOrderCommFilterWherefCall(GetQueryer, ctx, orderComm, filter, formats, formatArgs...)
	return
}

//UpdateOrderCommFilterWherefCall will update gex_order_comm to database
func UpdateOrderCommFilterWherefCall(caller interface{}, ctx context.Context, orderComm *OrderComm, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = orderComm.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindOrderCommCall will find gex_order_comm by id from database
func FindOrderComm(ctx context.Context, orderCommID int64) (orderComm *OrderComm, err error) {
	orderComm, err = FindOrderCommCall(GetQueryer, ctx, orderCommID, false)
	return
}

//FindOrderCommCall will find gex_order_comm by id from database
func FindOrderCommCall(caller interface{}, ctx context.Context, orderCommID int64, lock bool) (orderComm *OrderComm, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", orderCommID)
	orderComm, err = FindOrderCommWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindOrderCommWhereCall will find gex_order_comm by where from database
func FindOrderCommWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (orderComm *OrderComm, err error) {
	querySQL := crud.QuerySQL(&OrderComm{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &OrderComm{}, "#all", querySQL, args, &orderComm)
	return
}

//FindOrderCommWheref will find gex_order_comm by where from database
func FindOrderCommWheref(ctx context.Context, format string, args ...interface{}) (orderComm *OrderComm, err error) {
	orderComm, err = FindOrderCommWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindOrderCommWherefCall will find gex_order_comm by where from database
func FindOrderCommWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (orderComm *OrderComm, err error) {
	orderComm, err = FindOrderCommFilterWherefCall(GetQueryer, ctx, lock, "#all", format, args...)
	return
}

//FindOrderCommFilterWheref will find gex_order_comm by where from database
func FindOrderCommFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (orderComm *OrderComm, err error) {
	orderComm, err = FindOrderCommFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindOrderCommFilterWherefCall will find gex_order_comm by where from database
func FindOrderCommFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (orderComm *OrderComm, err error) {
	querySQL := crud.QuerySQL(&OrderComm{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &OrderComm{}, filter, querySQL, queryArgs, &orderComm)
	return
}

//ListOrderCommByID will list gex_order_comm by id from database
func ListOrderCommByID(ctx context.Context, orderCommIDs ...int64) (orderCommList []*OrderComm, orderCommMap map[int64]*OrderComm, err error) {
	orderCommList, orderCommMap, err = ListOrderCommByIDCall(GetQueryer, ctx, orderCommIDs...)
	return
}

//ListOrderCommByIDCall will list gex_order_comm by id from database
func ListOrderCommByIDCall(caller interface{}, ctx context.Context, orderCommIDs ...int64) (orderCommList []*OrderComm, orderCommMap map[int64]*OrderComm, err error) {
	if len(orderCommIDs) < 1 {
		orderCommMap = map[int64]*OrderComm{}
		return
	}
	err = ScanOrderCommByIDCall(caller, ctx, orderCommIDs, &orderCommList, &orderCommMap, "tid")
	return
}

//ListOrderCommFilterByID will list gex_order_comm by id from database
func ListOrderCommFilterByID(ctx context.Context, filter string, orderCommIDs ...int64) (orderCommList []*OrderComm, orderCommMap map[int64]*OrderComm, err error) {
	orderCommList, orderCommMap, err = ListOrderCommFilterByIDCall(GetQueryer, ctx, filter, orderCommIDs...)
	return
}

//ListOrderCommFilterByIDCall will list gex_order_comm by id from database
func ListOrderCommFilterByIDCall(caller interface{}, ctx context.Context, filter string, orderCommIDs ...int64) (orderCommList []*OrderComm, orderCommMap map[int64]*OrderComm, err error) {
	if len(orderCommIDs) < 1 {
		orderCommMap = map[int64]*OrderComm{}
		return
	}
	err = ScanOrderCommFilterByIDCall(caller, ctx, filter, orderCommIDs, &orderCommList, &orderCommMap, "tid")
	return
}

//ScanOrderCommByID will list gex_order_comm by id from database
func ScanOrderCommByID(ctx context.Context, orderCommIDs []int64, dest ...interface{}) (err error) {
	err = ScanOrderCommByIDCall(GetQueryer, ctx, orderCommIDs, dest...)
	return
}

//ScanOrderCommByIDCall will list gex_order_comm by id from database
func ScanOrderCommByIDCall(caller interface{}, ctx context.Context, orderCommIDs []int64, dest ...interface{}) (err error) {
	err = ScanOrderCommFilterByIDCall(caller, ctx, "#all", orderCommIDs, dest...)
	return
}

//ScanOrderCommFilterByID will list gex_order_comm by id from database
func ScanOrderCommFilterByID(ctx context.Context, filter string, orderCommIDs []int64, dest ...interface{}) (err error) {
	err = ScanOrderCommFilterByIDCall(GetQueryer, ctx, filter, orderCommIDs, dest...)
	return
}

//ScanOrderCommFilterByIDCall will list gex_order_comm by id from database
func ScanOrderCommFilterByIDCall(caller interface{}, ctx context.Context, filter string, orderCommIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&OrderComm{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(orderCommIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &OrderComm{}, filter, querySQL, nil, dest...)
	return
}

//ScanOrderCommWherefCall will list gex_order_comm by format from database
func ScanOrderCommWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanOrderCommWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanOrderCommWherefCall will list gex_order_comm by format from database
func ScanOrderCommWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanOrderCommFilterWherefCall(caller, ctx, "#all", format, args, suffix, dest...)
	return
}

//ScanOrderCommFilterWheref will list gex_order_comm by format from database
func ScanOrderCommFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanOrderCommFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanOrderCommFilterWherefCall will list gex_order_comm by format from database
func ScanOrderCommFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&OrderComm{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &OrderComm{}, filter, querySQL, args, dest...)
	return
}

//UserFilterOptional is crud filter
const UserFilterOptional = "role,name,account,phone,password,trade_pass,image,external,status"

//UserFilterRequired is crud filter
const UserFilterRequired = ""

//UserFilterInsert is crud filter
const UserFilterInsert = "role,name,account,phone,password,trade_pass,image,external,status"

//UserFilterUpdate is crud filter
const UserFilterUpdate = "update_time,role,name,account,phone,password,trade_pass,image,external,status"

//UserFilterFind is crud filter
const UserFilterFind = "^password,trade_pass#all"

//UserFilterScan is crud filter
const UserFilterScan = "^password,trade_pass#all"

//EnumValid will valid value by UserType
func (o *UserType) EnumValid(v interface{}) (err error) {
	var target UserType
	targetType := reflect.TypeOf(UserType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(UserType)
	}
	for _, value := range UserTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", UserTypeAll)
}

//EnumValid will valid value by UserTypeArray
func (o *UserTypeArray) EnumValid(v interface{}) (err error) {
	var target UserType
	targetType := reflect.TypeOf(UserType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(UserType)
	}
	for _, value := range UserTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", UserTypeAll)
}

//DbArray will join value to database array
func (o UserTypeArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o UserTypeArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by UserRole
func (o *UserRole) EnumValid(v interface{}) (err error) {
	var target UserRole
	targetType := reflect.TypeOf(UserRole(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(UserRole)
	}
	for _, value := range UserRoleAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", UserRoleAll)
}

//EnumValid will valid value by UserRoleArray
func (o *UserRoleArray) EnumValid(v interface{}) (err error) {
	var target UserRole
	targetType := reflect.TypeOf(UserRole(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(UserRole)
	}
	for _, value := range UserRoleAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", UserRoleAll)
}

//DbArray will join value to database array
func (o UserRoleArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o UserRoleArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by UserStatus
func (o *UserStatus) EnumValid(v interface{}) (err error) {
	var target UserStatus
	targetType := reflect.TypeOf(UserStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(UserStatus)
	}
	for _, value := range UserStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", UserStatusAll)
}

//EnumValid will valid value by UserStatusArray
func (o *UserStatusArray) EnumValid(v interface{}) (err error) {
	var target UserStatus
	targetType := reflect.TypeOf(UserStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(UserStatus)
	}
	for _, value := range UserStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", UserStatusAll)
}

//DbArray will join value to database array
func (o UserStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o UserStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithUser will return gex_user meta data
func MetaWithUser(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_user"), fields...)
	return
}

//MetaWith will return gex_user meta data
func (user *User) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_user"), fields...)
	return
}

//Meta will return gex_user meta data
func (user *User) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(user, "#all")
	return
}

//Valid will valid by filter
func (user *User) Valid() (err error) {
	if reflect.ValueOf(user.TID).IsZero() {
		err = attrvalid.Valid(user, UserFilterInsert+"#all", UserFilterOptional)
	} else {
		err = attrvalid.Valid(user, UserFilterUpdate, "")
	}
	return
}

//Insert will add gex_user to database
func (user *User) Insert(caller interface{}, ctx context.Context) (err error) {

	if len(user.Fee) < 1 {
		user.Fee = xsql.M{}
	}

	if len(user.External) < 1 {
		user.External = xsql.M{}
	}

	if user.UpdateTime.Timestamp() < 1 {
		user.UpdateTime = xsql.TimeNow()
	}

	if user.CreateTime.Timestamp() < 1 {
		user.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, user, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_user to database
func (user *User) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = user.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_user to database
func (user *User) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = user.UpdateFilterWheref(caller, ctx, UserFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_user to database
func (user *User) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	user.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(user, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", user.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, user, sql, where, "and", args)
	return
}

//AddUser will add gex_user to database
func AddUser(ctx context.Context, user *User) (err error) {
	err = AddUserCall(GetQueryer, ctx, user)
	return
}

//AddUser will add gex_user to database
func AddUserCall(caller interface{}, ctx context.Context, user *User) (err error) {
	err = user.Insert(caller, ctx)
	return
}

//UpdateUserFilter will update gex_user to database
func UpdateUserFilter(ctx context.Context, user *User, filter string) (err error) {
	err = UpdateUserFilterCall(GetQueryer, ctx, user, filter)
	return
}

//UpdateUserFilterCall will update gex_user to database
func UpdateUserFilterCall(caller interface{}, ctx context.Context, user *User, filter string) (err error) {
	err = user.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateUserWheref will update gex_user to database
func UpdateUserWheref(ctx context.Context, user *User, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateUserWherefCall(GetQueryer, ctx, user, formats, formatArgs...)
	return
}

//UpdateUserWherefCall will update gex_user to database
func UpdateUserWherefCall(caller interface{}, ctx context.Context, user *User, formats string, formatArgs ...interface{}) (err error) {
	err = user.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateUserFilterWheref will update gex_user to database
func UpdateUserFilterWheref(ctx context.Context, user *User, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateUserFilterWherefCall(GetQueryer, ctx, user, filter, formats, formatArgs...)
	return
}

//UpdateUserFilterWherefCall will update gex_user to database
func UpdateUserFilterWherefCall(caller interface{}, ctx context.Context, user *User, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = user.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindUserCall will find gex_user by id from database
func FindUser(ctx context.Context, userID int64) (user *User, err error) {
	user, err = FindUserCall(GetQueryer, ctx, userID, false)
	return
}

//FindUserCall will find gex_user by id from database
func FindUserCall(caller interface{}, ctx context.Context, userID int64, lock bool) (user *User, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", userID)
	user, err = FindUserWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindUserWhereCall will find gex_user by where from database
func FindUserWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (user *User, err error) {
	querySQL := crud.QuerySQL(&User{}, "^password,trade_pass#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &User{}, "^password,trade_pass#all", querySQL, args, &user)
	return
}

//FindUserWheref will find gex_user by where from database
func FindUserWheref(ctx context.Context, format string, args ...interface{}) (user *User, err error) {
	user, err = FindUserWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindUserWherefCall will find gex_user by where from database
func FindUserWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (user *User, err error) {
	user, err = FindUserFilterWherefCall(GetQueryer, ctx, lock, "^password,trade_pass#all", format, args...)
	return
}

//FindUserFilterWheref will find gex_user by where from database
func FindUserFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (user *User, err error) {
	user, err = FindUserFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindUserFilterWherefCall will find gex_user by where from database
func FindUserFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (user *User, err error) {
	querySQL := crud.QuerySQL(&User{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &User{}, filter, querySQL, queryArgs, &user)
	return
}

//ListUserByID will list gex_user by id from database
func ListUserByID(ctx context.Context, userIDs ...int64) (userList []*User, userMap map[int64]*User, err error) {
	userList, userMap, err = ListUserByIDCall(GetQueryer, ctx, userIDs...)
	return
}

//ListUserByIDCall will list gex_user by id from database
func ListUserByIDCall(caller interface{}, ctx context.Context, userIDs ...int64) (userList []*User, userMap map[int64]*User, err error) {
	if len(userIDs) < 1 {
		userMap = map[int64]*User{}
		return
	}
	err = ScanUserByIDCall(caller, ctx, userIDs, &userList, &userMap, "tid")
	return
}

//ListUserFilterByID will list gex_user by id from database
func ListUserFilterByID(ctx context.Context, filter string, userIDs ...int64) (userList []*User, userMap map[int64]*User, err error) {
	userList, userMap, err = ListUserFilterByIDCall(GetQueryer, ctx, filter, userIDs...)
	return
}

//ListUserFilterByIDCall will list gex_user by id from database
func ListUserFilterByIDCall(caller interface{}, ctx context.Context, filter string, userIDs ...int64) (userList []*User, userMap map[int64]*User, err error) {
	if len(userIDs) < 1 {
		userMap = map[int64]*User{}
		return
	}
	err = ScanUserFilterByIDCall(caller, ctx, filter, userIDs, &userList, &userMap, "tid")
	return
}

//ScanUserByID will list gex_user by id from database
func ScanUserByID(ctx context.Context, userIDs []int64, dest ...interface{}) (err error) {
	err = ScanUserByIDCall(GetQueryer, ctx, userIDs, dest...)
	return
}

//ScanUserByIDCall will list gex_user by id from database
func ScanUserByIDCall(caller interface{}, ctx context.Context, userIDs []int64, dest ...interface{}) (err error) {
	err = ScanUserFilterByIDCall(caller, ctx, "^password,trade_pass#all", userIDs, dest...)
	return
}

//ScanUserFilterByID will list gex_user by id from database
func ScanUserFilterByID(ctx context.Context, filter string, userIDs []int64, dest ...interface{}) (err error) {
	err = ScanUserFilterByIDCall(GetQueryer, ctx, filter, userIDs, dest...)
	return
}

//ScanUserFilterByIDCall will list gex_user by id from database
func ScanUserFilterByIDCall(caller interface{}, ctx context.Context, filter string, userIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&User{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(userIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &User{}, filter, querySQL, nil, dest...)
	return
}

//ScanUserWherefCall will list gex_user by format from database
func ScanUserWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanUserWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanUserWherefCall will list gex_user by format from database
func ScanUserWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanUserFilterWherefCall(caller, ctx, "^password,trade_pass#all", format, args, suffix, dest...)
	return
}

//ScanUserFilterWheref will list gex_user by format from database
func ScanUserFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanUserFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanUserFilterWherefCall will list gex_user by format from database
func ScanUserFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&User{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &User{}, filter, querySQL, args, dest...)
	return
}

//WithdrawFilterOptional is crud filter
const WithdrawFilterOptional = ""

//WithdrawFilterRequired is crud filter
const WithdrawFilterRequired = ""

//WithdrawFilterInsert is crud filter
const WithdrawFilterInsert = ""

//WithdrawFilterUpdate is crud filter
const WithdrawFilterUpdate = "update_time"

//WithdrawFilterFind is crud filter
const WithdrawFilterFind = "#all"

//WithdrawFilterScan is crud filter
const WithdrawFilterScan = "#all"

//EnumValid will valid value by WithdrawType
func (o *WithdrawType) EnumValid(v interface{}) (err error) {
	var target WithdrawType
	targetType := reflect.TypeOf(WithdrawType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(WithdrawType)
	}
	for _, value := range WithdrawTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", WithdrawTypeAll)
}

//EnumValid will valid value by WithdrawTypeArray
func (o *WithdrawTypeArray) EnumValid(v interface{}) (err error) {
	var target WithdrawType
	targetType := reflect.TypeOf(WithdrawType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(WithdrawType)
	}
	for _, value := range WithdrawTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", WithdrawTypeAll)
}

//DbArray will join value to database array
func (o WithdrawTypeArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o WithdrawTypeArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by WithdrawStatus
func (o *WithdrawStatus) EnumValid(v interface{}) (err error) {
	var target WithdrawStatus
	targetType := reflect.TypeOf(WithdrawStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(WithdrawStatus)
	}
	for _, value := range WithdrawStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", WithdrawStatusAll)
}

//EnumValid will valid value by WithdrawStatusArray
func (o *WithdrawStatusArray) EnumValid(v interface{}) (err error) {
	var target WithdrawStatus
	targetType := reflect.TypeOf(WithdrawStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(WithdrawStatus)
	}
	for _, value := range WithdrawStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", WithdrawStatusAll)
}

//DbArray will join value to database array
func (o WithdrawStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o WithdrawStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithWithdraw will return gex_withdraw meta data
func MetaWithWithdraw(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_withdraw"), fields...)
	return
}

//MetaWith will return gex_withdraw meta data
func (withdraw *Withdraw) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(string("gex_withdraw"), fields...)
	return
}

//Meta will return gex_withdraw meta data
func (withdraw *Withdraw) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(withdraw, "#all")
	return
}

//Valid will valid by filter
func (withdraw *Withdraw) Valid() (err error) {
	if reflect.ValueOf(withdraw.TID).IsZero() {
		err = attrvalid.Valid(withdraw, WithdrawFilterInsert+"#all", WithdrawFilterOptional)
	} else {
		err = attrvalid.Valid(withdraw, WithdrawFilterUpdate, "")
	}
	return
}

//Insert will add gex_withdraw to database
func (withdraw *Withdraw) Insert(caller interface{}, ctx context.Context) (err error) {

	if len(withdraw.Transaction) < 1 {
		withdraw.Transaction = xsql.M{}
	}

	if withdraw.UpdateTime.Timestamp() < 1 {
		withdraw.UpdateTime = xsql.TimeNow()
	}

	if withdraw.CreateTime.Timestamp() < 1 {
		withdraw.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, withdraw, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update gex_withdraw to database
func (withdraw *Withdraw) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = withdraw.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update gex_withdraw to database
func (withdraw *Withdraw) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = withdraw.UpdateFilterWheref(caller, ctx, WithdrawFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update gex_withdraw to database
func (withdraw *Withdraw) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	withdraw.UpdateTime = xsql.TimeNow()
	sql, args := crud.UpdateSQL(withdraw, filter, nil)
	where, args := crud.AppendWheref(nil, args, "tid=$%v", withdraw.TID)
	if len(formats) > 0 {
		where, args = crud.AppendWheref(where, args, formats, formatArgs...)
	}
	err = crud.UpdateRow(caller, ctx, withdraw, sql, where, "and", args)
	return
}

//AddWithdraw will add gex_withdraw to database
func AddWithdraw(ctx context.Context, withdraw *Withdraw) (err error) {
	err = AddWithdrawCall(GetQueryer, ctx, withdraw)
	return
}

//AddWithdraw will add gex_withdraw to database
func AddWithdrawCall(caller interface{}, ctx context.Context, withdraw *Withdraw) (err error) {
	err = withdraw.Insert(caller, ctx)
	return
}

//UpdateWithdrawFilter will update gex_withdraw to database
func UpdateWithdrawFilter(ctx context.Context, withdraw *Withdraw, filter string) (err error) {
	err = UpdateWithdrawFilterCall(GetQueryer, ctx, withdraw, filter)
	return
}

//UpdateWithdrawFilterCall will update gex_withdraw to database
func UpdateWithdrawFilterCall(caller interface{}, ctx context.Context, withdraw *Withdraw, filter string) (err error) {
	err = withdraw.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateWithdrawWheref will update gex_withdraw to database
func UpdateWithdrawWheref(ctx context.Context, withdraw *Withdraw, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateWithdrawWherefCall(GetQueryer, ctx, withdraw, formats, formatArgs...)
	return
}

//UpdateWithdrawWherefCall will update gex_withdraw to database
func UpdateWithdrawWherefCall(caller interface{}, ctx context.Context, withdraw *Withdraw, formats string, formatArgs ...interface{}) (err error) {
	err = withdraw.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateWithdrawFilterWheref will update gex_withdraw to database
func UpdateWithdrawFilterWheref(ctx context.Context, withdraw *Withdraw, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateWithdrawFilterWherefCall(GetQueryer, ctx, withdraw, filter, formats, formatArgs...)
	return
}

//UpdateWithdrawFilterWherefCall will update gex_withdraw to database
func UpdateWithdrawFilterWherefCall(caller interface{}, ctx context.Context, withdraw *Withdraw, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = withdraw.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindWithdrawCall will find gex_withdraw by id from database
func FindWithdraw(ctx context.Context, withdrawID int64) (withdraw *Withdraw, err error) {
	withdraw, err = FindWithdrawCall(GetQueryer, ctx, withdrawID, false)
	return
}

//FindWithdrawCall will find gex_withdraw by id from database
func FindWithdrawCall(caller interface{}, ctx context.Context, withdrawID int64, lock bool) (withdraw *Withdraw, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", withdrawID)
	withdraw, err = FindWithdrawWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindWithdrawWhereCall will find gex_withdraw by where from database
func FindWithdrawWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (withdraw *Withdraw, err error) {
	querySQL := crud.QuerySQL(&Withdraw{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Withdraw{}, "#all", querySQL, args, &withdraw)
	return
}

//FindWithdrawWheref will find gex_withdraw by where from database
func FindWithdrawWheref(ctx context.Context, format string, args ...interface{}) (withdraw *Withdraw, err error) {
	withdraw, err = FindWithdrawWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindWithdrawWherefCall will find gex_withdraw by where from database
func FindWithdrawWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (withdraw *Withdraw, err error) {
	withdraw, err = FindWithdrawFilterWherefCall(GetQueryer, ctx, lock, "#all", format, args...)
	return
}

//FindWithdrawFilterWheref will find gex_withdraw by where from database
func FindWithdrawFilterWheref(ctx context.Context, filter string, format string, args ...interface{}) (withdraw *Withdraw, err error) {
	withdraw, err = FindWithdrawFilterWherefCall(GetQueryer, ctx, false, filter, format, args...)
	return
}

//FindWithdrawFilterWherefCall will find gex_withdraw by where from database
func FindWithdrawFilterWherefCall(caller interface{}, ctx context.Context, lock bool, filter string, format string, args ...interface{}) (withdraw *Withdraw, err error) {
	querySQL := crud.QuerySQL(&Withdraw{}, filter)
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Withdraw{}, filter, querySQL, queryArgs, &withdraw)
	return
}

//ListWithdrawByID will list gex_withdraw by id from database
func ListWithdrawByID(ctx context.Context, withdrawIDs ...int64) (withdrawList []*Withdraw, withdrawMap map[int64]*Withdraw, err error) {
	withdrawList, withdrawMap, err = ListWithdrawByIDCall(GetQueryer, ctx, withdrawIDs...)
	return
}

//ListWithdrawByIDCall will list gex_withdraw by id from database
func ListWithdrawByIDCall(caller interface{}, ctx context.Context, withdrawIDs ...int64) (withdrawList []*Withdraw, withdrawMap map[int64]*Withdraw, err error) {
	if len(withdrawIDs) < 1 {
		withdrawMap = map[int64]*Withdraw{}
		return
	}
	err = ScanWithdrawByIDCall(caller, ctx, withdrawIDs, &withdrawList, &withdrawMap, "tid")
	return
}

//ListWithdrawFilterByID will list gex_withdraw by id from database
func ListWithdrawFilterByID(ctx context.Context, filter string, withdrawIDs ...int64) (withdrawList []*Withdraw, withdrawMap map[int64]*Withdraw, err error) {
	withdrawList, withdrawMap, err = ListWithdrawFilterByIDCall(GetQueryer, ctx, filter, withdrawIDs...)
	return
}

//ListWithdrawFilterByIDCall will list gex_withdraw by id from database
func ListWithdrawFilterByIDCall(caller interface{}, ctx context.Context, filter string, withdrawIDs ...int64) (withdrawList []*Withdraw, withdrawMap map[int64]*Withdraw, err error) {
	if len(withdrawIDs) < 1 {
		withdrawMap = map[int64]*Withdraw{}
		return
	}
	err = ScanWithdrawFilterByIDCall(caller, ctx, filter, withdrawIDs, &withdrawList, &withdrawMap, "tid")
	return
}

//ScanWithdrawByID will list gex_withdraw by id from database
func ScanWithdrawByID(ctx context.Context, withdrawIDs []int64, dest ...interface{}) (err error) {
	err = ScanWithdrawByIDCall(GetQueryer, ctx, withdrawIDs, dest...)
	return
}

//ScanWithdrawByIDCall will list gex_withdraw by id from database
func ScanWithdrawByIDCall(caller interface{}, ctx context.Context, withdrawIDs []int64, dest ...interface{}) (err error) {
	err = ScanWithdrawFilterByIDCall(caller, ctx, "#all", withdrawIDs, dest...)
	return
}

//ScanWithdrawFilterByID will list gex_withdraw by id from database
func ScanWithdrawFilterByID(ctx context.Context, filter string, withdrawIDs []int64, dest ...interface{}) (err error) {
	err = ScanWithdrawFilterByIDCall(GetQueryer, ctx, filter, withdrawIDs, dest...)
	return
}

//ScanWithdrawFilterByIDCall will list gex_withdraw by id from database
func ScanWithdrawFilterByIDCall(caller interface{}, ctx context.Context, filter string, withdrawIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Withdraw{}, filter)
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(withdrawIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Withdraw{}, filter, querySQL, nil, dest...)
	return
}

//ScanWithdrawWherefCall will list gex_withdraw by format from database
func ScanWithdrawWheref(ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanWithdrawWherefCall(GetQueryer, ctx, format, args, suffix, dest...)
	return
}

//ScanWithdrawWherefCall will list gex_withdraw by format from database
func ScanWithdrawWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanWithdrawFilterWherefCall(caller, ctx, "#all", format, args, suffix, dest...)
	return
}

//ScanWithdrawFilterWheref will list gex_withdraw by format from database
func ScanWithdrawFilterWheref(ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	err = ScanWithdrawFilterWherefCall(GetQueryer, ctx, filter, format, args, suffix, dest...)
	return
}

//ScanWithdrawFilterWherefCall will list gex_withdraw by format from database
func ScanWithdrawFilterWherefCall(caller interface{}, ctx context.Context, filter string, format string, args []interface{}, suffix string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Withdraw{}, filter)
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ", suffix)
	err = crud.Query(caller, ctx, &Withdraw{}, filter, querySQL, args, dest...)
	return
}
