package basedb

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xmap"
)

var ConfigAll = []string{}
var ConfigKey = map[string][]string{}

//StoreConf will store config valu to database
func StoreConf(ctx context.Context, key string, val string) (err error) {
	err = StoreConfCall(Pool(), ctx, key, val)
	return
}

//StoreConfTx will store config valu to database
func StoreConfCall(caller crud.Queryer, ctx context.Context, key string, val string) (err error) {
	sql := fmt.Sprintf(`
		insert into %v_config(key,value,update_time) values($1,$2,$3) 
		on conflict(key) 
		do update set value=$2, update_time=$3`, SYS)
	args := []interface{}{key, val, time.Now()}
	_, _, err = caller.Exec(ctx, sql, args...)
	return
}

//LoadConf will return config by key from data
func LoadConf(ctx context.Context, key string, val interface{}) (err error) {
	err = LoadConfCall(Pool(), ctx, key, val)
	return
}

func LoadConfCall(caller crud.Queryer, ctx context.Context, key string, val interface{}) (err error) {
	returnType := "text"
	kind := reflect.Indirect(reflect.ValueOf(val)).Type()
	if kind.Kind() == reflect.Ptr {
		kind = kind.Elem()
	}
	if kind.Kind().String() == "map" {
		data := ""
		err = caller.QueryRow(ctx, `select value::text from `+SYS+`_config where key=$1`, key).Scan(&data)
		if err == nil {
			err = json.Unmarshal([]byte(data), val)
		}
	} else {
		switch kind.Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			returnType = "int8"
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			returnType = "float8"
		default:
			returnType = "text"
		}
		err = caller.QueryRow(ctx, `select value::`+returnType+` from `+SYS+`_config where key=$1`, key).Scan(val)
	}
	return
}

//UpdateConfigList update all configue
func UpdateConfigList(ctx context.Context, config xmap.M) (err error) {
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
	for key, val := range config {
		_, _, err = tx.Exec(ctx, `update `+SYS+`_config set value=$1,update_time=$2 where key=$3`, val, time.Now(), key)
		if err != nil {
			err = fmt.Errorf("update config by key(%v).Value(%+v) fail with %v", key, val, err)
			return
		}
	}
	return
}

//LoadConfigList load config by keys
func LoadConfigList(ctx context.Context, keys ...string) (config xmap.M, err error) {
	rows, err := Pool().Query(ctx, `select key,value from `+SYS+`_config where key=any($1)`, keys)
	if err != nil {
		return
	}
	defer rows.Close()
	config = xmap.M{}
	for rows.Next() {
		var key, value string
		err = rows.Scan(&key, &value)
		if err != nil {
			return
		}
		config[key] = value
	}
	return
}
