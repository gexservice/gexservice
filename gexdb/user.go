package gexdb

import (
	"context"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/converter"
)

//FindUserByUsrPwd will return user by match account/email/phone=username and passowrd=matched
func FindUserByUsrPwd(ctx context.Context, username, password string) (*User, error) {
	return FindUserWheref(ctx, "(account=$%v or phone=$%v),password=$%v", username, password)
}

//FindUserByAccount will return user by match account/email/phone=username and passowrd=matched
func FindUserByAccount(ctx context.Context, account string) (*User, error) {
	return FindUserWheref(ctx, "account=$%v", account)
}

//FindUserByPhone will return user by match phone
func FindUserByPhone(ctx context.Context, phone string) (*User, error) {
	return FindUserWheref(ctx, "phone=$%v", phone)
}

//UpdateUserCaller will update user to database
func UpdateUserCaller(caller crud.Queryer, ctx context.Context, user, old *User) (err error) {
	sql, args := crud.UpdateSQL(user, UserFilterUpdate, nil)
	where, args := crud.AppendWhere(nil, args, true, "tid=$%v", user.TID)
	if old != nil && user.Password != nil && old.Password != nil {
		where, args = crud.AppendWhere(where, args, true, "(password is null or password='' or password=$%v)", old.Password)
	}
	if old != nil && user.TradePass != nil && old.TradePass != nil {
		where, args = crud.AppendWhere(where, args, true, "(trade_pass is null or trade_pass='' or trade_pass=$%v)", old.TradePass)
	}
	err = crud.UpdateRow(caller, ctx, user, sql, where, "and", args)
	return
}

//UpdateUser will update user to database
func UpdateUser(ctx context.Context, user *User) (err error) {
	return UpdateUserCaller(Pool(), ctx, user, nil)
}

func UpdateUserByOld(ctx context.Context, user, old *User) (err error) {
	return UpdateUserCaller(Pool(), ctx, user, old)
}

func UserHavingTradePassword(ctx context.Context, userID int64) (having int, err error) {
	var tradePass *string
	err = Pool().QueryRow(ctx, `select trade_pass from exs_user where tid=$1`, userID).Scan(&tradePass)
	if err == nil && tradePass != nil && len(*tradePass) > 0 {
		having = 1
	}
	return
}

func UserVerifyTradePassword(ctx context.Context, userID int64, password string) (err error) {
	err = Pool().QueryRow(ctx, `select tid from exs_user where tid=$1 and trade_pass=$2`, userID, password).Scan(converter.Int64Ptr(0))
	return
}

/**
 * @apiDefine UserUnifySearcher
 * @apiParam  {Number} [type] the type filter, multi with comma, all type supported is <a href="#metadata-User">UserTypeAll</a>
 * @apiParam  {Number} [role] the role filter, multi with comma, all type supported is <a href="#metadata-User">UserRoleAll</a>
 * @apiParam  {Number} [status] the status filter, multi with comma, all status supported is <a href="#metadata-User">UserStatusAll</a>
 * @apiParam  {String} [key] search key
 * @apiParam  {Number} [offset] page offset
 * @apiParam  {Number} [limit] page limit
 */
type UserUnifySearcher struct {
	Model User `json:"model"`
	Where struct {
		Type   UserTypeArray   `json:"type" cmp:"type=any($%v)" valid:"type,o|i,e:;"`
		Role   UserRoleArray   `json:"role" cmp:"role=any($%v)" valid:"role,o|i,e:;"`
		Status UserStatusArray `json:"status" cmp:"status=any($%v)" valid:"status,o|i,e:;"`
		Key    string          `json:"key" cmp:"(tid::text like $%v or name like $%v or phone like $%v or account like $%v)" valid:"key,o|s,l:0;"`
	} `json:"where" join:"and" valid:"inline"`
	Page struct {
		Order string `json:"order" default:"order by update_time desc" valid:"order,o|s,l:0;"`
		Skip  int    `json:"skip" valid:"skip,o|i,r:-1;"`
		Limit int    `json:"limit" valid:"limit,o|i,r:0;"`
	} `json:"page" valid:"inline"`
	Query struct {
		Users   []*User `json:"users"`
		UserIDs []int64 `json:"user_ids" scan:"tid"`
	} `json:"query" filter:"^password,external#all"`
	Count struct {
		Total int64 `json:"total" scan:"tid"`
	} `json:"count" filter:"count(tid)#all"`
}

func (t *UserUnifySearcher) Apply(ctx context.Context) (err error) {
	if len(t.Where.Key) > 0 {
		t.Where.Key = "%" + t.Where.Key + "%"
	}
	t.Page.Order = crud.BuildOrderby(UserOrderbyAll, t.Page.Order)
	err = crud.ApplyUnify(Pool(), ctx, t)
	return
}
