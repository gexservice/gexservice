package basedb

import (
	"context"

	"github.com/codingeasygo/crud"
)

func UpsertAnnounce(ctx context.Context, announce *Announce) (err error) {
	if announce.TID > 0 {
		err = UpdateAnnounceWheref(ctx, announce, "")
	} else {
		err = AddAnnounce(ctx, announce)
	}
	return
}

/**
 * @apiDefine AnnounceUnifySearcher
 * @apiParam  {Number} [type] the type filter, multi with comma, all type supported is <a href="#metadata-Announce">AnnounceTypeAll</a>
 * @apiParam  {Number} [marked] the marked filter, multi with comma
 * @apiParam  {Number} [status] the status filter, multi with comma, all status supported is <a href="#metadata-Announce">AnnounceStatusAll</a>
 * @apiParam  {String} [key] search key
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type AnnounceUnifySearcher struct {
	Model Announce `json:"model"`
	Where struct {
		Type   AnnounceTypeArray   `json:"type" cmp:"type=any($%v)" valid:"type,o|i,e:;"`
		Marked []int               `json:"marked" cmp:"marked=any($%v)" valid:"marked,o|i,r:0;"`
		Status AnnounceStatusArray `json:"status" cmp:"status=any($%v)" valid:"status,o|i,e:;"`
		Key    string              `json:"key" cmp:"title like $%v or info::text like $%v or content::text like $%v" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by update_time desc"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Return struct {
		Content int `json:"ret_content" valid:"ret_content,o|i,o:0~1;"`
	} `json:"return" valid:"inline"`
	Query struct {
		Filter    crud.FilterValue `json:"filter"`
		Announces []*Announce      `json:"announces"`
	} `json:"query" filter:"^content#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"count(tid)#all"`
}

func (t *AnnounceUnifySearcher) Apply(ctx context.Context) (err error) {
	if len(t.Where.Key) > 0 {
		t.Where.Key = "%" + t.Where.Key + "%"
	}
	if t.Return.Content > 0 {
		t.Query.Filter = crud.FilterValue("#all")
	}
	err = crud.ApplyUnify(Pool(), ctx, t)
	return
}
