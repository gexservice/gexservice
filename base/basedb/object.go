package basedb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
)

//LoadObject will return the object data
func LoadObject(ctx context.Context, key string) (data xmap.M, err error) {
	object, err := FindObject(ctx, key)
	if err == nil && object.Value != nil {
		data = object.Value.AsMap()
	}
	return
}

//UpsertObject will update object data
func UpsertObject(ctx context.Context, key string, data interface{}) (effected int64, err error) {
	effected, err = UpsertObjectCall(Pool(), ctx, key, data)
	return
}

func UpsertObjectCall(caller crud.Queryer, ctx context.Context, key string, data interface{}) (effected int64, err error) {
	bys, err := json.Marshal(data)
	if err != nil {
		return
	}
	sql, args := crud.InsertSQL(
		MetaWithObject(key, string(bys), time.Now(), time.Now(), ObjectStatusNormal),
		"key,value,create_time,update_time,status#all",
		"on conflict (key) do update set value=$2",
	)
	_, effected, err = caller.Exec(ctx, sql, args...)
	return
}

func UpsertVersionObject(ctx context.Context, object *VersionObject) (err error) {
	if object.TID > 0 {
		err = UpdateVersionObjectWheref(ctx, object, "")
	} else {
		err = AddVersionObject(ctx, object)
	}
	return
}

//LoadLatestVersionObject will return latest version object by key
func LoadLatestVersionObject(ctx context.Context, key string, pubs ...string) (object *VersionObject, err error) {
	object = &VersionObject{}
	sql := crud.QuerySQL(object, "#all")
	sql, args := crud.JoinWheref(
		sql, nil,
		"key=$%v,status=$%v,(pub='*' or pub like any($%v))",
		key, VersionObjectStatusNormal, xsql.StringArray(pubs),
	)
	sql = crud.JoinPage(sql, "order by tid desc", 0, 1)
	err = crud.QueryRow(Pool, ctx, object, "#all", sql, args, &object)
	return
}

/**
 * @apiDefine VersionObjectUnifySearcher
 * @apiParam  {Number} [status] the status filter, multi with comma, all status supported is <a href="#metadata-Announce">AnnounceStatusAll</a>
 * @apiParam  {String} [key] list by key
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type VersionObjectUnifySearcher struct {
	Model VersionObject `json:"model"`
	Where struct {
		Status VersionObjectStatusArray `json:"status" cmp:"status=any($%v)" valid:"status,o|i,e:;"`
		Key    string                   `json:"key" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by update_time desc"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Objects []*VersionObject `json:"objects"`
	} `json:"query" filter:"^content#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"count(tid)#all"`
}

func (v *VersionObjectUnifySearcher) Apply(ctx context.Context) (err error) {
	err = crud.ApplyUnify(Pool(), ctx, v)
	return
}
