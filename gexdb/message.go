package gexdb

import (
	"context"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xprop"
	"github.com/codingeasygo/util/xsql"
)

var MessageTemplate = map[string]map[string]string{}

const (
	MessageKeyWithdrawDone = "withdraw.done"
	MessageKeyWithdrawFail = "withdraw.fail"
	MessageKeyTopup        = "topup"
	MessageKeyBlowup       = "blowup"
	MessageKeyGoldbar      = "goldbar"
)

var MessageKeyAll = []string{MessageKeyWithdrawDone, MessageKeyWithdrawFail, MessageKeyTopup, MessageKeyBlowup}

func ReadMessageTemplateByConfig(config *xprop.Config) (template map[string]map[string]string) {
	template = map[string]map[string]string{}
	for _, key := range MessageKeyAll {
		titleKey := key + ".title"
		template[titleKey] = map[string]string{}
		config.Range("message."+titleKey, func(k string, val interface{}) {
			template[titleKey][k] = val.(string)
		})
		contentKey := key + ".content"
		template[contentKey] = map[string]string{}
		config.Range("message."+contentKey, func(k string, val interface{}) {
			template[contentKey][k] = val.(string)
		})
	}
	return
}

func ParseMessageTemplate(template map[string]map[string]string, env xmap.M, key string) (message xsql.M) {
	message = xsql.M{}
	for k, v := range template[key] {
		message[k] = env.ReplaceAll(v, false, true)
	}
	return
}

func AddMultiMessage(ctx context.Context, messages ...*Message) (err error) {
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
	err = AddMultiMessageCall(tx, ctx, messages...)
	return
}

func AddMultiMessageCall(caller crud.Queryer, ctx context.Context, messages ...*Message) (err error) {
	for _, message := range messages {
		err = AddMessageCall(caller, ctx, message)
		if err != nil {
			break
		}
	}
	return
}

func AddTemplateMessageCall(caller crud.Queryer, ctx context.Context, messageType MessageType, env xmap.M, key string, toUserID int64) (message *Message, err error) {
	message = &Message{
		Type:     messageType,
		Title:    ParseMessageTemplate(MessageTemplate, env, key+".title"),
		Content:  ParseMessageTemplate(MessageTemplate, env, key+".content"),
		ToUserID: toUserID,
		Status:   MessageStatusNormal,
	}
	err = AddMessageCall(caller, ctx, message)
	return
}

/**
 * @apiDefine MessageUnifySearcher
 * @apiParam  {Number} [type] the message type filter, multi with comma, all type supported is <a href="#metadata-Message">MessageTypeAll</a>
 * @apiParam  {Number} [asset] the balance asset filter, multi with comma
 * @apiParam  {Number} [start_time] the time filter
 * @apiParam  {Number} [end_time] the time filter
 * @apiParam  {Number} [status] the withdraw status filter, multi with comma, all type supported is <a href="#metadata-Message">MessageStatusAll</a>
 * @apiParam  {Number} [skip] page skip
 * @apiParam  {Number} [limit] page limit
 */
type MessageUnifySearcher struct {
	Model Message `json:"model"`
	Where struct {
		Type      MessageTypeArray   `json:"type" cmp:"type=any($%v)" valid:"type,o|i,e:;"`
		StartTime xsql.Time          `json:"start_time" cmp:"update_time>=$%v" valid:"start_time,o|i,r:-1;"`
		EndTime   xsql.Time          `json:"end_time" cmp:"update_time<$%v" valid:"end_time,o|i,r:-1;"`
		ToUserID  int64              `json:"to_user_id" cmp:"(to_user_id=$%v or to_user_id=0)" valid:"to_user_id,o|i,r:0;"`
		Status    MessageStatusArray `json:"status" cmp:"status=any($%v)" valid:"status,o|i,e:;"`
		Key       string             `json:"key" cmp:"title::text ilike $%v" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by update_time desc" valid:"order,o|s,l:0;"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Messages []*Message `json:"messages"`
	} `json:"query" filter:"#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"r.count(tid)#all"`
}

func (m *MessageUnifySearcher) Apply(ctx context.Context) (err error) {
	m.Page.Order = ""
	if len(m.Where.Key) > 0 {
		m.Where.Key = "%" + m.Where.Key + "%"
	}
	err = crud.ApplyUnify(Pool(), ctx, m)
	return
}
