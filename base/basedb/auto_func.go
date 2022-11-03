//auto gen func by autogen
package basedb

import (
	"context"
	"fmt"
	"reflect"
	"strings"

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

//AnnounceFilterOptional is crud filter
const AnnounceFilterOptional = "marked,info,content,status"

//AnnounceFilterRequired is crud filter
const AnnounceFilterRequired = "type,title"

//AnnounceFilterInsert is crud filter
const AnnounceFilterInsert = "marked,info,content,status,type,title"

//AnnounceFilterUpdate is crud filter
const AnnounceFilterUpdate = "marked,title,info,content,status,update_time"

//EnumValid will valid value by AnnounceType
func (o *AnnounceType) EnumValid(v interface{}) (err error) {
	var target AnnounceType
	targetType := reflect.TypeOf(AnnounceType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(AnnounceType)
	}
	for _, value := range AnnounceTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", AnnounceTypeAll)
}

//EnumValid will valid value by AnnounceTypeArray
func (o *AnnounceTypeArray) EnumValid(v interface{}) (err error) {
	var target AnnounceType
	targetType := reflect.TypeOf(AnnounceType(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(AnnounceType)
	}
	for _, value := range AnnounceTypeAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", AnnounceTypeAll)
}

//DbArray will join value to database array
func (o AnnounceTypeArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o AnnounceTypeArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//EnumValid will valid value by AnnounceStatus
func (o *AnnounceStatus) EnumValid(v interface{}) (err error) {
	var target AnnounceStatus
	targetType := reflect.TypeOf(AnnounceStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(AnnounceStatus)
	}
	for _, value := range AnnounceStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", AnnounceStatusAll)
}

//EnumValid will valid value by AnnounceStatusArray
func (o *AnnounceStatusArray) EnumValid(v interface{}) (err error) {
	var target AnnounceStatus
	targetType := reflect.TypeOf(AnnounceStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(AnnounceStatus)
	}
	for _, value := range AnnounceStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", AnnounceStatusAll)
}

//DbArray will join value to database array
func (o AnnounceStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o AnnounceStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithAnnounce will return announce meta data
func MetaWithAnnounce(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("announce"), fields...)
	return
}

//MetaWith will return announce meta data
func (announce *Announce) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("announce"), fields...)
	return
}

//Meta will return announce meta data
func (announce *Announce) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(announce, "#all")
	return
}

//Valid will valid by filter
func (announce *Announce) Valid() (err error) {
	if reflect.ValueOf(announce.TID).IsZero() {
		err = attrvalid.Valid(announce, AnnounceFilterInsert+"#all", AnnounceFilterOptional)
	} else {
		err = attrvalid.Valid(announce, AnnounceFilterUpdate, "")
	}
	return
}

//Insert will add announce to database
func (announce *Announce) Insert(caller interface{}, ctx context.Context) (err error) {

	if len(announce.Info) < 1 {
		announce.Info = xsql.M{}
	}

	if len(announce.Content) < 1 {
		announce.Content = xsql.M{}
	}

	if announce.UpdateTime.Timestamp() < 1 {
		announce.UpdateTime = xsql.TimeNow()
	}

	if announce.CreateTime.Timestamp() < 1 {
		announce.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, announce, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update announce to database
func (announce *Announce) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = announce.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update announce to database
func (announce *Announce) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = announce.UpdateFilterWheref(caller, ctx, AnnounceFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update announce to database
func (announce *Announce) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	announce.UpdateTime = xsql.TimeNow()
	whereAll := []string{"tid=$%v"}
	whereArg := []interface{}{announce.TID}
	if len(formats) > 0 {
		whereAll = append(whereAll, formats)
		whereArg = append(whereArg, formatArgs...)
	}
	err = crud.UpdateRowWheref(caller, ctx, announce, filter, strings.Join(whereAll, ","), whereArg...)
	return
}

//AddAnnounce will add announce to database
func AddAnnounce(ctx context.Context, announce *Announce) (err error) {
	err = AddAnnounceCall(GetQueryer, ctx, announce)
	return
}

//AddAnnounce will add announce to database
func AddAnnounceCall(caller interface{}, ctx context.Context, announce *Announce) (err error) {
	err = announce.Insert(caller, ctx)
	return
}

//UpdateAnnounceFilter will update announce to database
func UpdateAnnounceFilter(ctx context.Context, announce *Announce, filter string) (err error) {
	err = UpdateAnnounceFilterCall(GetQueryer, ctx, announce, filter)
	return
}

//UpdateAnnounceFilterCall will update announce to database
func UpdateAnnounceFilterCall(caller interface{}, ctx context.Context, announce *Announce, filter string) (err error) {
	err = announce.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateAnnounceWheref will update announce to database
func UpdateAnnounceWheref(ctx context.Context, announce *Announce, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateAnnounceWherefCall(GetQueryer, ctx, announce, formats, formatArgs...)
	return
}

//UpdateAnnounceWherefCall will update announce to database
func UpdateAnnounceWherefCall(caller interface{}, ctx context.Context, announce *Announce, formats string, formatArgs ...interface{}) (err error) {
	err = announce.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateAnnounceFilterWheref will update announce to database
func UpdateAnnounceFilterWheref(ctx context.Context, announce *Announce, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateAnnounceFilterWherefCall(GetQueryer, ctx, announce, filter, formats, formatArgs...)
	return
}

//UpdateAnnounceFilterWherefCall will update announce to database
func UpdateAnnounceFilterWherefCall(caller interface{}, ctx context.Context, announce *Announce, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = announce.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindAnnounceCall will find announce by id from database
func FindAnnounce(ctx context.Context, announceID int64) (announce *Announce, err error) {
	announce, err = FindAnnounceCall(GetQueryer, ctx, announceID, false)
	return
}

//FindAnnounceCall will find announce by id from database
func FindAnnounceCall(caller interface{}, ctx context.Context, announceID int64, lock bool) (announce *Announce, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", announceID)
	announce, err = FindAnnounceWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindAnnounceWhereCall will find announce by where from database
func FindAnnounceWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (announce *Announce, err error) {
	querySQL := crud.QuerySQL(&Announce{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Announce{}, "#all", querySQL, args, &announce)
	return
}

//FindAnnounceWheref will find announce by where from database
func FindAnnounceWheref(ctx context.Context, format string, args ...interface{}) (announce *Announce, err error) {
	announce, err = FindAnnounceWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindAnnounceWherefCall will find announce by where from database
func FindAnnounceWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (announce *Announce, err error) {
	querySQL := crud.QuerySQL(&Announce{}, "#all")
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Announce{}, "#all", querySQL, queryArgs, &announce)
	return
}

//ListAnnounceByID will list announce by id from database
func ListAnnounceByID(ctx context.Context, announceIDs ...int64) (announceList []*Announce, announceMap map[int64]*Announce, err error) {
	announceList, announceMap, err = ListAnnounceByIDCall(GetQueryer, ctx, announceIDs...)
	return
}

//ListAnnounceByIDCall will list announce by id from database
func ListAnnounceByIDCall(caller interface{}, ctx context.Context, announceIDs ...int64) (announceList []*Announce, announceMap map[int64]*Announce, err error) {
	if len(announceIDs) < 1 {
		announceMap = map[int64]*Announce{}
		return
	}
	err = ScanAnnounceByIDCall(caller, ctx, announceIDs, &announceList, &announceMap, "tid")
	return
}

//ScanAnnounceByID will list announce by id from database
func ScanAnnounceByID(ctx context.Context, announceIDs []int64, dest ...interface{}) (err error) {
	err = ScanAnnounceByIDCall(GetQueryer, ctx, announceIDs, dest...)
	return
}

//ScanAnnounceByIDCall will list announce by id from database
func ScanAnnounceByIDCall(caller interface{}, ctx context.Context, announceIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Announce{}, "#all")
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(announceIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Announce{}, "#all", querySQL, nil, dest...)
	return
}

//ScanAnnounce will list announce by format from database
func ScanAnnounceWheref(ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	err = ScanAnnounceWherefCall(GetQueryer, ctx, format, args, dest...)
	return
}

//ScanAnnounceCall will list announce by format from database
func ScanAnnounceWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Announce{}, "#all")
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Announce{}, "#all", querySQL, args, dest...)
	return
}

//ConfigFilterOptional is crud filter
const ConfigFilterOptional = ""

//ConfigFilterRequired is crud filter
const ConfigFilterRequired = ""

//ConfigFilterInsert is crud filter
const ConfigFilterInsert = ""

//ConfigFilterUpdate is crud filter
const ConfigFilterUpdate = "update_time"

//MetaWithConfig will return config meta data
func MetaWithConfig(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("config"), fields...)
	return
}

//MetaWith will return config meta data
func (config *Config) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("config"), fields...)
	return
}

//Meta will return config meta data
func (config *Config) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(config, "#all")
	return
}

//Valid will valid by filter
func (config *Config) Valid() (err error) {
	if reflect.ValueOf(config.Key).IsZero() {
		err = attrvalid.Valid(config, ConfigFilterInsert+"#all", ConfigFilterOptional)
	} else {
		err = attrvalid.Valid(config, ConfigFilterUpdate, "")
	}
	return
}

//Insert will add config to database
func (config *Config) Insert(caller interface{}, ctx context.Context) (err error) {

	if config.UpdateTime.Timestamp() < 1 {
		config.UpdateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, config, "#all", "", "")
	return
}

//UpdateFilter will update config to database
func (config *Config) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = config.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update config to database
func (config *Config) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = config.UpdateFilterWheref(caller, ctx, ConfigFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update config to database
func (config *Config) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	config.UpdateTime = xsql.TimeNow()
	whereAll := []string{"key=$%v"}
	whereArg := []interface{}{config.Key}
	if len(formats) > 0 {
		whereAll = append(whereAll, formats)
		whereArg = append(whereArg, formatArgs...)
	}
	err = crud.UpdateRowWheref(caller, ctx, config, filter, strings.Join(whereAll, ","), whereArg...)
	return
}

//UpdateConfigFilter will update config to database
func UpdateConfigFilter(ctx context.Context, config *Config, filter string) (err error) {
	err = UpdateConfigFilterCall(GetQueryer, ctx, config, filter)
	return
}

//UpdateConfigFilterCall will update config to database
func UpdateConfigFilterCall(caller interface{}, ctx context.Context, config *Config, filter string) (err error) {
	err = config.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateConfigWheref will update config to database
func UpdateConfigWheref(ctx context.Context, config *Config, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateConfigWherefCall(GetQueryer, ctx, config, formats, formatArgs...)
	return
}

//UpdateConfigWherefCall will update config to database
func UpdateConfigWherefCall(caller interface{}, ctx context.Context, config *Config, formats string, formatArgs ...interface{}) (err error) {
	err = config.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateConfigFilterWheref will update config to database
func UpdateConfigFilterWheref(ctx context.Context, config *Config, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateConfigFilterWherefCall(GetQueryer, ctx, config, filter, formats, formatArgs...)
	return
}

//UpdateConfigFilterWherefCall will update config to database
func UpdateConfigFilterWherefCall(caller interface{}, ctx context.Context, config *Config, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = config.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindConfigCall will find config by id from database
func FindConfig(ctx context.Context, configID string) (config *Config, err error) {
	config, err = FindConfigCall(GetQueryer, ctx, configID, false)
	return
}

//FindConfigCall will find config by id from database
func FindConfigCall(caller interface{}, ctx context.Context, configID string, lock bool) (config *Config, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "key=$%v", configID)
	config, err = FindConfigWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindConfigWhereCall will find config by where from database
func FindConfigWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (config *Config, err error) {
	querySQL := crud.QuerySQL(&Config{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Config{}, "#all", querySQL, args, &config)
	return
}

//FindConfigWheref will find config by where from database
func FindConfigWheref(ctx context.Context, format string, args ...interface{}) (config *Config, err error) {
	config, err = FindConfigWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindConfigWherefCall will find config by where from database
func FindConfigWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (config *Config, err error) {
	querySQL := crud.QuerySQL(&Config{}, "#all")
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Config{}, "#all", querySQL, queryArgs, &config)
	return
}

//ListConfigByID will list config by id from database
func ListConfigByID(ctx context.Context, configIDs ...string) (configList []*Config, configMap map[string]*Config, err error) {
	configList, configMap, err = ListConfigByIDCall(GetQueryer, ctx, configIDs...)
	return
}

//ListConfigByIDCall will list config by id from database
func ListConfigByIDCall(caller interface{}, ctx context.Context, configIDs ...string) (configList []*Config, configMap map[string]*Config, err error) {
	if len(configIDs) < 1 {
		configMap = map[string]*Config{}
		return
	}
	err = ScanConfigByIDCall(caller, ctx, configIDs, &configList, &configMap, "key")
	return
}

//ScanConfigByID will list config by id from database
func ScanConfigByID(ctx context.Context, configIDs []string, dest ...interface{}) (err error) {
	err = ScanConfigByIDCall(GetQueryer, ctx, configIDs, dest...)
	return
}

//ScanConfigByIDCall will list config by id from database
func ScanConfigByIDCall(caller interface{}, ctx context.Context, configIDs []string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Config{}, "#all")
	where := append([]string{}, fmt.Sprintf("key in (%v)", xsql.StringArray(configIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Config{}, "#all", querySQL, nil, dest...)
	return
}

//ScanConfig will list config by format from database
func ScanConfigWheref(ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	err = ScanConfigWherefCall(GetQueryer, ctx, format, args, dest...)
	return
}

//ScanConfigCall will list config by format from database
func ScanConfigWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Config{}, "#all")
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Config{}, "#all", querySQL, args, dest...)
	return
}

//ObjectFilterOptional is crud filter
const ObjectFilterOptional = ""

//ObjectFilterRequired is crud filter
const ObjectFilterRequired = ""

//ObjectFilterInsert is crud filter
const ObjectFilterInsert = ""

//ObjectFilterUpdate is crud filter
const ObjectFilterUpdate = "update_time"

//EnumValid will valid value by ObjectStatus
func (o *ObjectStatus) EnumValid(v interface{}) (err error) {
	var target ObjectStatus
	targetType := reflect.TypeOf(ObjectStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(ObjectStatus)
	}
	for _, value := range ObjectStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", ObjectStatusAll)
}

//EnumValid will valid value by ObjectStatusArray
func (o *ObjectStatusArray) EnumValid(v interface{}) (err error) {
	var target ObjectStatus
	targetType := reflect.TypeOf(ObjectStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(ObjectStatus)
	}
	for _, value := range ObjectStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", ObjectStatusAll)
}

//DbArray will join value to database array
func (o ObjectStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o ObjectStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithObject will return object meta data
func MetaWithObject(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("object"), fields...)
	return
}

//MetaWith will return object meta data
func (object *Object) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("object"), fields...)
	return
}

//Meta will return object meta data
func (object *Object) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(object, "#all")
	return
}

//Valid will valid by filter
func (object *Object) Valid() (err error) {
	if reflect.ValueOf(object.Key).IsZero() {
		err = attrvalid.Valid(object, ObjectFilterInsert+"#all", ObjectFilterOptional)
	} else {
		err = attrvalid.Valid(object, ObjectFilterUpdate, "")
	}
	return
}

//Insert will add object to database
func (object *Object) Insert(caller interface{}, ctx context.Context) (err error) {

	if len(object.Value) < 1 {
		object.Value = xsql.M{}
	}

	if object.UpdateTime.Timestamp() < 1 {
		object.UpdateTime = xsql.TimeNow()
	}

	if object.CreateTime.Timestamp() < 1 {
		object.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, object, "#all", "", "")
	return
}

//UpdateFilter will update object to database
func (object *Object) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = object.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update object to database
func (object *Object) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = object.UpdateFilterWheref(caller, ctx, ObjectFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update object to database
func (object *Object) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	object.UpdateTime = xsql.TimeNow()
	whereAll := []string{"key=$%v"}
	whereArg := []interface{}{object.Key}
	if len(formats) > 0 {
		whereAll = append(whereAll, formats)
		whereArg = append(whereArg, formatArgs...)
	}
	err = crud.UpdateRowWheref(caller, ctx, object, filter, strings.Join(whereAll, ","), whereArg...)
	return
}

//UpdateObjectFilter will update object to database
func UpdateObjectFilter(ctx context.Context, object *Object, filter string) (err error) {
	err = UpdateObjectFilterCall(GetQueryer, ctx, object, filter)
	return
}

//UpdateObjectFilterCall will update object to database
func UpdateObjectFilterCall(caller interface{}, ctx context.Context, object *Object, filter string) (err error) {
	err = object.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateObjectWheref will update object to database
func UpdateObjectWheref(ctx context.Context, object *Object, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateObjectWherefCall(GetQueryer, ctx, object, formats, formatArgs...)
	return
}

//UpdateObjectWherefCall will update object to database
func UpdateObjectWherefCall(caller interface{}, ctx context.Context, object *Object, formats string, formatArgs ...interface{}) (err error) {
	err = object.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateObjectFilterWheref will update object to database
func UpdateObjectFilterWheref(ctx context.Context, object *Object, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateObjectFilterWherefCall(GetQueryer, ctx, object, filter, formats, formatArgs...)
	return
}

//UpdateObjectFilterWherefCall will update object to database
func UpdateObjectFilterWherefCall(caller interface{}, ctx context.Context, object *Object, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = object.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindObjectCall will find object by id from database
func FindObject(ctx context.Context, objectID string) (object *Object, err error) {
	object, err = FindObjectCall(GetQueryer, ctx, objectID, false)
	return
}

//FindObjectCall will find object by id from database
func FindObjectCall(caller interface{}, ctx context.Context, objectID string, lock bool) (object *Object, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "key=$%v", objectID)
	object, err = FindObjectWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindObjectWhereCall will find object by where from database
func FindObjectWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (object *Object, err error) {
	querySQL := crud.QuerySQL(&Object{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Object{}, "#all", querySQL, args, &object)
	return
}

//FindObjectWheref will find object by where from database
func FindObjectWheref(ctx context.Context, format string, args ...interface{}) (object *Object, err error) {
	object, err = FindObjectWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindObjectWherefCall will find object by where from database
func FindObjectWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (object *Object, err error) {
	querySQL := crud.QuerySQL(&Object{}, "#all")
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &Object{}, "#all", querySQL, queryArgs, &object)
	return
}

//ListObjectByID will list object by id from database
func ListObjectByID(ctx context.Context, objectIDs ...string) (objectList []*Object, objectMap map[string]*Object, err error) {
	objectList, objectMap, err = ListObjectByIDCall(GetQueryer, ctx, objectIDs...)
	return
}

//ListObjectByIDCall will list object by id from database
func ListObjectByIDCall(caller interface{}, ctx context.Context, objectIDs ...string) (objectList []*Object, objectMap map[string]*Object, err error) {
	if len(objectIDs) < 1 {
		objectMap = map[string]*Object{}
		return
	}
	err = ScanObjectByIDCall(caller, ctx, objectIDs, &objectList, &objectMap, "key")
	return
}

//ScanObjectByID will list object by id from database
func ScanObjectByID(ctx context.Context, objectIDs []string, dest ...interface{}) (err error) {
	err = ScanObjectByIDCall(GetQueryer, ctx, objectIDs, dest...)
	return
}

//ScanObjectByIDCall will list object by id from database
func ScanObjectByIDCall(caller interface{}, ctx context.Context, objectIDs []string, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Object{}, "#all")
	where := append([]string{}, fmt.Sprintf("key in (%v)", xsql.StringArray(objectIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Object{}, "#all", querySQL, nil, dest...)
	return
}

//ScanObject will list object by format from database
func ScanObjectWheref(ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	err = ScanObjectWherefCall(GetQueryer, ctx, format, args, dest...)
	return
}

//ScanObjectCall will list object by format from database
func ScanObjectWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&Object{}, "#all")
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &Object{}, "#all", querySQL, args, dest...)
	return
}

//VersionObjectFilterOptional is crud filter
const VersionObjectFilterOptional = ""

//VersionObjectFilterRequired is crud filter
const VersionObjectFilterRequired = "key,value,pub,status"

//VersionObjectFilterInsert is crud filter
const VersionObjectFilterInsert = "key,value,pub,status"

//VersionObjectFilterUpdate is crud filter
const VersionObjectFilterUpdate = "pub,value,status,update_time"

//EnumValid will valid value by VersionObjectStatus
func (o *VersionObjectStatus) EnumValid(v interface{}) (err error) {
	var target VersionObjectStatus
	targetType := reflect.TypeOf(VersionObjectStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(VersionObjectStatus)
	}
	for _, value := range VersionObjectStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", VersionObjectStatusAll)
}

//EnumValid will valid value by VersionObjectStatusArray
func (o *VersionObjectStatusArray) EnumValid(v interface{}) (err error) {
	var target VersionObjectStatus
	targetType := reflect.TypeOf(VersionObjectStatus(0))
	targetValue := reflect.ValueOf(v)
	if targetValue.CanConvert(targetType) {
		target = targetValue.Convert(targetType).Interface().(VersionObjectStatus)
	}
	for _, value := range VersionObjectStatusAll {
		if target == value {
			return nil
		}
	}
	return fmt.Errorf("must be in %v", VersionObjectStatusAll)
}

//DbArray will join value to database array
func (o VersionObjectStatusArray) DbArray() (res string) {
	res = "{" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + "}"
	return
}

//InArray will join value to database array
func (o VersionObjectStatusArray) InArray() (res string) {
	res = "" + converter.JoinSafe(o, ",", converter.JoinPolicyDefault) + ""
	return
}

//MetaWithVersionObject will return version_object meta data
func MetaWithVersionObject(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("version_object"), fields...)
	return
}

//MetaWith will return version_object meta data
func (versionObject *VersionObject) MetaWith(fields ...interface{}) (v []interface{}) {
	v = crud.MetaWith(BaseTableName("version_object"), fields...)
	return
}

//Meta will return version_object meta data
func (versionObject *VersionObject) Meta() (table string, fileds []string) {
	table, fileds = crud.QueryField(versionObject, "#all")
	return
}

//Valid will valid by filter
func (versionObject *VersionObject) Valid() (err error) {
	if reflect.ValueOf(versionObject.TID).IsZero() {
		err = attrvalid.Valid(versionObject, VersionObjectFilterInsert+"#all", VersionObjectFilterOptional)
	} else {
		err = attrvalid.Valid(versionObject, VersionObjectFilterUpdate, "")
	}
	return
}

//Insert will add version_object to database
func (versionObject *VersionObject) Insert(caller interface{}, ctx context.Context) (err error) {

	if len(versionObject.Value) < 1 {
		versionObject.Value = xsql.M{}
	}

	if versionObject.UpdateTime.Timestamp() < 1 {
		versionObject.UpdateTime = xsql.TimeNow()
	}

	if versionObject.CreateTime.Timestamp() < 1 {
		versionObject.CreateTime = xsql.TimeNow()
	}

	_, err = crud.InsertFilter(caller, ctx, versionObject, "^tid#all", "returning", "tid#all")
	return
}

//UpdateFilter will update version_object to database
func (versionObject *VersionObject) UpdateFilter(caller interface{}, ctx context.Context, filter string) (err error) {
	err = versionObject.UpdateFilterWheref(caller, ctx, filter, "")
	return
}

//UpdateWheref will update version_object to database
func (versionObject *VersionObject) UpdateWheref(caller interface{}, ctx context.Context, formats string, formatArgs ...interface{}) (err error) {
	err = versionObject.UpdateFilterWheref(caller, ctx, VersionObjectFilterUpdate, formats, formatArgs...)
	return
}

//UpdateFilterWheref will update version_object to database
func (versionObject *VersionObject) UpdateFilterWheref(caller interface{}, ctx context.Context, filter string, formats string, formatArgs ...interface{}) (err error) {
	versionObject.UpdateTime = xsql.TimeNow()
	whereAll := []string{"tid=$%v"}
	whereArg := []interface{}{versionObject.TID}
	if len(formats) > 0 {
		whereAll = append(whereAll, formats)
		whereArg = append(whereArg, formatArgs...)
	}
	err = crud.UpdateRowWheref(caller, ctx, versionObject, filter, strings.Join(whereAll, ","), whereArg...)
	return
}

//AddVersionObject will add version_object to database
func AddVersionObject(ctx context.Context, versionObject *VersionObject) (err error) {
	err = AddVersionObjectCall(GetQueryer, ctx, versionObject)
	return
}

//AddVersionObject will add version_object to database
func AddVersionObjectCall(caller interface{}, ctx context.Context, versionObject *VersionObject) (err error) {
	err = versionObject.Insert(caller, ctx)
	return
}

//UpdateVersionObjectFilter will update version_object to database
func UpdateVersionObjectFilter(ctx context.Context, versionObject *VersionObject, filter string) (err error) {
	err = UpdateVersionObjectFilterCall(GetQueryer, ctx, versionObject, filter)
	return
}

//UpdateVersionObjectFilterCall will update version_object to database
func UpdateVersionObjectFilterCall(caller interface{}, ctx context.Context, versionObject *VersionObject, filter string) (err error) {
	err = versionObject.UpdateFilter(caller, ctx, filter)
	return
}

//UpdateVersionObjectWheref will update version_object to database
func UpdateVersionObjectWheref(ctx context.Context, versionObject *VersionObject, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateVersionObjectWherefCall(GetQueryer, ctx, versionObject, formats, formatArgs...)
	return
}

//UpdateVersionObjectWherefCall will update version_object to database
func UpdateVersionObjectWherefCall(caller interface{}, ctx context.Context, versionObject *VersionObject, formats string, formatArgs ...interface{}) (err error) {
	err = versionObject.UpdateWheref(caller, ctx, formats, formatArgs...)
	return
}

//UpdateVersionObjectFilterWheref will update version_object to database
func UpdateVersionObjectFilterWheref(ctx context.Context, versionObject *VersionObject, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = UpdateVersionObjectFilterWherefCall(GetQueryer, ctx, versionObject, filter, formats, formatArgs...)
	return
}

//UpdateVersionObjectFilterWherefCall will update version_object to database
func UpdateVersionObjectFilterWherefCall(caller interface{}, ctx context.Context, versionObject *VersionObject, filter string, formats string, formatArgs ...interface{}) (err error) {
	err = versionObject.UpdateFilterWheref(caller, ctx, filter, formats, formatArgs...)
	return
}

//FindVersionObjectCall will find version_object by id from database
func FindVersionObject(ctx context.Context, versionObjectID int64) (versionObject *VersionObject, err error) {
	versionObject, err = FindVersionObjectCall(GetQueryer, ctx, versionObjectID, false)
	return
}

//FindVersionObjectCall will find version_object by id from database
func FindVersionObjectCall(caller interface{}, ctx context.Context, versionObjectID int64, lock bool) (versionObject *VersionObject, err error) {
	where, args := crud.AppendWhere(nil, nil, true, "tid=$%v", versionObjectID)
	versionObject, err = FindVersionObjectWhereCall(caller, ctx, lock, "and", where, args)
	return
}

//FindVersionObjectWhereCall will find version_object by where from database
func FindVersionObjectWhereCall(caller interface{}, ctx context.Context, lock bool, join string, where []string, args []interface{}) (versionObject *VersionObject, err error) {
	querySQL := crud.QuerySQL(&VersionObject{}, "#all")
	querySQL = crud.JoinWhere(querySQL, where, join)
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &VersionObject{}, "#all", querySQL, args, &versionObject)
	return
}

//FindVersionObjectWheref will find version_object by where from database
func FindVersionObjectWheref(ctx context.Context, format string, args ...interface{}) (versionObject *VersionObject, err error) {
	versionObject, err = FindVersionObjectWherefCall(GetQueryer, ctx, false, format, args...)
	return
}

//FindVersionObjectWherefCall will find version_object by where from database
func FindVersionObjectWherefCall(caller interface{}, ctx context.Context, lock bool, format string, args ...interface{}) (versionObject *VersionObject, err error) {
	querySQL := crud.QuerySQL(&VersionObject{}, "#all")
	where, queryArgs := crud.AppendWheref(nil, nil, format, args...)
	querySQL = crud.JoinWhere(querySQL, where, "and")
	if lock {
		querySQL += " for update "
	}
	err = crud.QueryRow(caller, ctx, &VersionObject{}, "#all", querySQL, queryArgs, &versionObject)
	return
}

//ListVersionObjectByID will list version_object by id from database
func ListVersionObjectByID(ctx context.Context, versionObjectIDs ...int64) (versionObjectList []*VersionObject, versionObjectMap map[int64]*VersionObject, err error) {
	versionObjectList, versionObjectMap, err = ListVersionObjectByIDCall(GetQueryer, ctx, versionObjectIDs...)
	return
}

//ListVersionObjectByIDCall will list version_object by id from database
func ListVersionObjectByIDCall(caller interface{}, ctx context.Context, versionObjectIDs ...int64) (versionObjectList []*VersionObject, versionObjectMap map[int64]*VersionObject, err error) {
	if len(versionObjectIDs) < 1 {
		versionObjectMap = map[int64]*VersionObject{}
		return
	}
	err = ScanVersionObjectByIDCall(caller, ctx, versionObjectIDs, &versionObjectList, &versionObjectMap, "tid")
	return
}

//ScanVersionObjectByID will list version_object by id from database
func ScanVersionObjectByID(ctx context.Context, versionObjectIDs []int64, dest ...interface{}) (err error) {
	err = ScanVersionObjectByIDCall(GetQueryer, ctx, versionObjectIDs, dest...)
	return
}

//ScanVersionObjectByIDCall will list version_object by id from database
func ScanVersionObjectByIDCall(caller interface{}, ctx context.Context, versionObjectIDs []int64, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&VersionObject{}, "#all")
	where := append([]string{}, fmt.Sprintf("tid in (%v)", xsql.Int64Array(versionObjectIDs).InArray()))
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &VersionObject{}, "#all", querySQL, nil, dest...)
	return
}

//ScanVersionObject will list version_object by format from database
func ScanVersionObjectWheref(ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	err = ScanVersionObjectWherefCall(GetQueryer, ctx, format, args, dest...)
	return
}

//ScanVersionObjectCall will list version_object by format from database
func ScanVersionObjectWherefCall(caller interface{}, ctx context.Context, format string, args []interface{}, dest ...interface{}) (err error) {
	querySQL := crud.QuerySQL(&VersionObject{}, "#all")
	var where []string
	if len(format) > 0 {
		where, args = crud.AppendWheref(nil, nil, format, args...)
	}
	querySQL = crud.JoinWhere(querySQL, where, " and ")
	err = crud.Query(caller, ctx, &VersionObject{}, "#all", querySQL, args, dest...)
	return
}
