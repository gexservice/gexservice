package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/codingeasygo/web/httptest"

	"github.com/codingeasygo/crud/pgx"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
)

func TestLoginByUsername(t *testing.T) {
	//login
	login, _ := ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	fmt.Printf("login--->%v\n", converter.JSON(login))
	ts.Should(t, "code", define.NotFound).GetMap("/pub/login?username=%v&password=%v", "abc0", "1x23")
	ts.Should(t, "code", define.NotFound).GetMap("/pub/login?username=%v&password=%v", "abcxx", "123")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/pub/login?username=%v&password=%v", "", "")
	ts.Should(t, "code", define.UserInvalid).GetMap("/pub/login?username=%v&password=%v", "abc1", "123")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/pub/login?username=%v&passwordx=%v", "abc1", "123")
	//user info
	userInfo, _ := ts.Should(t, "code", define.Success, "user", xmap.ShouldIsNoNil).GetMap("/usr/userInfo")
	fmt.Printf("userInfo--->%v\n", converter.JSON(userInfo))
	//update
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap("xxx", "/usr/updateUser")
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap(&gexdb.User{Status: -100}, "/usr/updateUser")
	ts.Should(t, "code", gexdb.CodeOldPasswordInvalid).PostJSONMap(&gexdb.User{
		Password: converter.StringPtr("123"),
	}, "/usr/updateUser")
	ts.Should(t, "code", define.Duplicate).PostJSONMap(&gexdb.User{
		Account: converter.StringPtr("abc1"),
	}, "/usr/updateUser")
	updateUser, _ := ts.Should(t, "code", define.Success).PostJSONMap(&updateUserArg{
		OldPassword:  "123",
		OldTradePass: "123",
		User: gexdb.User{
			Phone:     converter.StringPtr("12345678901"),
			Password:  converter.StringPtr("123"),
			TradePass: converter.StringPtr("123"),
		},
	}, "/usr/updateUser")
	fmt.Printf("updateUser->%v\n", converter.JSON(updateUser))
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	pgx.MockerClear()
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).GetMap("/usr/userInfo")
	pgx.MockerSetRangeCall("Rows.Scan", 2, 4).Should(t, "code", define.Success).GetMap("/usr/userInfo")
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.ServerError).PostJSONMap(&gexdb.User{
		Password: converter.StringPtr("123"),
	}, "/usr/updateUser")
	//logout
	ts.Should(t, "code", define.Success).GetMap("/usr/logout")
	ts.Should(t, "code", define.Redirect).GetMap("/usr/userInfo")
}

func TestManageUser(t *testing.T) {
	//not access
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	ts.Should(t, "code", define.NotAccess).GetMap("/usr/searchUser")
	ts.Should(t, "code", define.NotAccess).GetMap("/usr/loadUser?user_id=%v", userabc0.TID)

	//admin access
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	//searcher user
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/searchUser?limit=xx")
	searchUser, _ := ts.Should(t, "code", define.Success, "users", xmap.ShouldIsNoEmpty).GetMap("/usr/searchUser?key=abc0&ret_balance=1")
	fmt.Printf("searchUser-->%v\n", converter.JSON(searchUser))
	userID := searchUser.Int64Def(0, "/users/0/tid")
	//load user
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/loadUser?user_id=%v", "xxx")
	loadUser, _ := ts.Should(t, "code", define.Success).GetMap("/usr/loadUser?user_id=%v", userID)
	fmt.Printf("loadUser-->%v\n", converter.JSON(loadUser))
	//update user
	ts.Should(t, "code", define.Success).PostJSONMap(&gexdb.User{
		Password:  converter.StringPtr("123"),
		TradePass: converter.StringPtr("123"),
	}, "/usr/updateUser")
	ts.Should(t, "code", define.Success).PostJSONMap(&gexdb.User{
		TID:       userID,
		Password:  converter.StringPtr("123"),
		TradePass: converter.StringPtr("123"),
	}, "/usr/updateUser")

	pgx.MockerStart()
	defer pgx.MockerStop()
	//
	pgx.MockerSetCall("Rows.Scan", 2).Should(t, "code", define.ServerError).GetMap("/usr/searchUser?key=abc0")
	pgx.MockerSetRangeCall("Pool.Query", 1, 2).Should(t, "code", define.ServerError).GetMap("/usr/searchUser?key=abc0&ret_balance=1")
	pgx.MockerSetCall("Rows.Scan", 2).Should(t, "code", define.ServerError).GetMap("/usr/loadUser?user_id=%v", userID)
}

func TestAdminAccess(t *testing.T) {
	clearCookie()
	ts := httptest.NewMuxServer()
	ts.Mux.HandleFunc("/pub/login", LoginH)
	ts.Mux.HandleFunc("/testAccess", func(s *web.Session) web.Result {
		if AdminAccess(s) {
			return s.SendJSON(xmap.M{"code": define.Success})
		} else {
			return s.SendJSON(xmap.M{"code": define.NotAccess})
		}
	})
	ts.Should(t, "code", define.NotAccess).GetMap("/testAccess")
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	ts.Should(t, "code", define.Success).GetMap("/testAccess")

	pgx.MockerStart()
	defer pgx.MockerStop()
	//
	pgx.MockerSetCall("Rows.Scan", 1).Should(t, "code", define.NotAccess).Call(func(trigger int) (res xmap.M, err error) {
		return ts.GetMap("/testAccess")
	})
}
