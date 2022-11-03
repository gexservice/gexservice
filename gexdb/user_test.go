package gexdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xsql"
)

func testAddUser(prefix string) (user *User) {
	account, phone, password := prefix+"_acc", prefix+"_123", "123"
	image := prefix + "_image"
	user = &User{
		Type:      UserTypeNormal,
		Role:      UserRoleNormal,
		Name:      &prefix,
		Account:   &account,
		Phone:     &phone,
		Image:     &image,
		Password:  &password,
		TradePass: &password,
		External:  xsql.M{"abc": 1},
		Status:    UserStatusNormal,
	}
	err := AddUser(ctx, user)
	if err != nil {
		panic(err)
	}
	return
}

func TestUser(t *testing.T) {
	clear()
	user := testAddUser("abc")
	user2, err := FindUserByUsrPwd(ctx, "abc_acc", "123")
	if err != nil || user.TID != user2.TID {
		t.Errorf("err:%v,user:%v,user2:%v", err, user.TID, user2.TID)
		return
	}
	user3, err := FindUserByUsrPwd(ctx, "abc_123", "123")
	if err != nil || user.TID != user3.TID {
		t.Error(err)
		return
	}
	user4, err := FindUserByAccount(ctx, *user.Account)
	if err != nil || user.TID != user4.TID {
		t.Error(err)
		return
	}
	user.TradePass = converter.StringPtr("123")
	err = UpdateUser(ctx, user)
	if err != nil {
		t.Error(err)
		return
	}
	user.Password = converter.StringPtr("abc")
	user.TradePass = converter.StringPtr("abc")
	old := &User{Password: converter.StringPtr("123"), TradePass: converter.StringPtr("123")}
	err = UpdateUserByOld(ctx, user, old)
	if err != nil {
		t.Error(err)
		return
	}

	findUser, err := FindUser(ctx, user.TID)
	if err != nil || user.TID != findUser.TID {
		t.Error(err)
		fmt.Printf("%v\n", user.TID)
		return
	}

	findPhoneUser, err := FindUserByPhone(ctx, *user.Phone)
	if err != nil || user.TID != findPhoneUser.TID {
		t.Error(err)
		return
	}

	having, err := UserHavingTradePassword(ctx, user.TID)
	if err != nil || having != 1 {
		t.Error(err)
		return
	}

	err = UserVerifyTradePassword(ctx, user.TID, *user.TradePass)
	if err != nil {
		t.Error(err)
		return
	}

	searcher := &UserUnifySearcher{}
	searcher.Where.Type = UserTypeAll
	searcher.Where.Key = *user.Name
	searcher.Where.Status = UserStatusArray{user.Status}
	err = searcher.Apply(context.Background())
	if err != nil || len(searcher.Query.Users) < 1 || searcher.Count.Total < 1 {
		t.Error(err)
		return
	}
}
