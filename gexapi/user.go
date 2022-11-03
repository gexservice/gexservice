package gexapi

import (
	"fmt"
	"strings"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhash"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

//LoginAccessRedirect will return the redirect url when user is not login
var LoginAccessRedirect string

//LoginAccessF is the normal user login access filter
func LoginAccessF(hs *web.Session) web.Result {
	userID, ok := hs.Value("user_id").(int64)
	if !ok || userID < 1 {
		return hs.SendJSON(xmap.M{
			"code":     define.Redirect,
			"redirect": LoginAccessRedirect,
			"message":  "not login",
		})
	}
	return web.Continue
}

func AdminAccess(s *web.Session) bool {
	userID, ok := s.Value("user_id").(int64)
	if !ok || userID < 1 {
		return false
	}
	user, err := gexdb.FindUser(s.R.Context(), userID)
	if err != nil {
		xlog.Errorf("AdminAccess find user fail with %v", err)
		return false
	}
	return user.Type == gexdb.UserTypeAdmin
}

//LoginH is http handler to login by username and password
/**
 *
 * @api {POST} /pub/login Login
 * @apiName Login
 * @apiGroup User
 *
 *
 * @apiParam  {String} [username] the login name, it can be account/phone
 * @apiParam  {String} [password] the login password
 * @apiParam  {String} [kbz_token] the kbz mp kbz_token
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Number} user.tid the user id
 * @apiSuccess (Success) {Number} user.type the user type, all type supported is <a href="#metadata-User">UserTypeAll</a>
 * @apiSuccess (Success) {Number} user.role the user role, all role supported is <a href="#metadata-User">UserRoleAll</a>
 * @apiSuccess (Success) {String} user.account the user account
 * @apiSuccess (Success) {String} user.name the user nickname
 * @apiSuccess (Success) {String} user.phone the user phone
 * @apiSuccess (Success) {Object} user.external the user external info
 * @apiSuccess (Success) {Number} user.create_time the user create time
 * @apiSuccess (Success) {Number} user.update_time the user update time
 * @apiSuccess (Success) {Number} user.status the user status
 *
 * @apiParamExample  {form} AccountPassword-Example(form):
 * username=abc&password=123
 * @apiParamExample  {form} PhonePassword-Example(form):
 * username=abc&password=123
 * @apiParamExample  {form} PhoneCode-Example(form):
 * username=abc&code=123
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "accesses": {
 *         "role": 100
 *     },
 *     "code": 0,
 *     "session_id": "613c40a7285c668490000001",
 *     "user": {
 *         "account": "abc2",
 *         "broker_id": 0,
 *         "create_time": 1631338663770,
 *         "external": {},
 *         "image": "abc2_image",
 *         "role": 100,
 *         "name": "abc2_name",
 *         "phone": "abc2_123",
 *         "status": 100,
 *         "tid": 100004,
 *         "type": 100,
 *         "update_time": 1631338663770
 *     }
 * }
 *
 */
func LoginH(s *web.Session) web.Result {
	var username, password string
	var err = s.ValidFormat(`
		username,R|S,L:0~32;
		password,R|S,L:0~16;
	`, &username, &password)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	user, err := gexdb.FindUserByUsrPwd(s.R.Context(), username, xhash.SHA1([]byte(password)))
	if err != nil {
		xlog.Errorf("LoginH user login by username to sytem from %v fail with %v", s.R.RemoteAddr, err)
		code := define.ServerError
		if strings.Contains(err.Error(), "no rows") {
			code = define.NotFound
		}
		return util.ReturnCodeLocalErr(s, code, "usr-err", err)
	}
	if user.Status != gexdb.UserStatusNormal {
		xlog.Warnf("LoginH user login to sytem from %v fail with user:%v status is %v", s.R.RemoteAddr, user.TID, user.Status)
		return util.ReturnCodeLocalErr(s, define.UserInvalid, "usr-err", fmt.Errorf("user status is %v", user.Status))
	}
	s.Clear()
	s.SetValue("user_id", user.TID)
	s.Flush()
	xlog.Infof("LoginH user login to sytem from %v success by uid:%v", s.R.RemoteAddr, user.TID)
	return sendUserInfo(s, user)
}

//LogoutH is http handler will logout current session
/**
 *
 * @api {GET} /usr/logout Logout
 * @apiName Logout
 * @apiGroup User
 *
 *
 *
 * @apiSuccess (Success) {Number} code 0 is success
 *
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code" : 0
 * }
 *
 *
 */
func LogoutH(hs *web.Session) web.Result {
	hs.Clear()
	hs.Flush()
	return util.ReturnCodeData(hs, 0, "OK")
}

//UserInfoH is http handler will return current login user
/**
 *
 * @api {GET} /usr/userInfo User Info
 * @apiName UserInfo
 * @apiGroup User
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Object} user the user info
 * @apiSuccess (Success) {Number} user.tid the user id
 * @apiSuccess (Success) {Number} user.type the user type, all type supported is <a href="#metadata-User">UserTypeAll</a>
 * @apiSuccess (Success) {Number} user.role the user role, all role supported is <a href="#metadata-User">UserRoleAll</a>
 * @apiSuccess (Success) {String} user.account the user account
 * @apiSuccess (Success) {String} user.name the user nickname
 * @apiSuccess (Success) {String} user.phone the user phone
 * @apiSuccess (Success) {Object} user.external the user external info
 * @apiSuccess (Success) {Number} user.create_time the user create time
 * @apiSuccess (Success) {Number} user.update_time the user update time
 * @apiSuccess (Success) {Number} user.status the user status, all status supported is <a href="#metadata-User">UserStatusAll</a>
 * @apiSuccess (Success) {Object} user the user info
 * @apiSuccess (Success) {Number} trade_pass_setted if user trade pass setted, 1 is setted
 * @apiSuccess (Success) {String} config.goldbar_address the goldbar withdraw address.
 * @apiSuccess (Success) {String} config.welcom_message the welcom message
 * @apiSuccess (Success) {String} config.withdraw_max the withdraw max
 *
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "accesses": {
 *         "role": 100,
 *         "type": 100
 *     },
 *     "code": 0,
 *     "session_id": "616b8ed9285c660c6a000001",
 *     "trade_pass_setted": 1,
 *     "user": {
 *         "account": "abc0",
 *         "broker_id": 100004,
 *         "create_time": 1634438873232,
 *         "external": {},
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1634438873233
 *     }
 * }
 */
func UserInfoH(s *web.Session) web.Result {
	userID := s.Value("user_id").(int64)
	user, err := gexdb.FindUser(s.R.Context(), userID)
	if err != nil {
		xlog.Warnf("UserInfoH find user info fail with %v by uid:%v", err, userID)
		return util.ReturnCodeLocalErr(s, define.ServerError, "usr-err", err)
	}
	return sendUserInfo(s, user)
}

func sendUserInfo(s *web.Session, user *gexdb.User) web.Result {
	//
	//ignore error for not return error in user info
	tradePassSetted, err := gexdb.UserHavingTradePassword(s.R.Context(), user.TID)
	if err != nil {
		xlog.Warnf("UserInfoH check user having trade password with %v by uid:%v", err, user.TID)
	}
	//config
	config, err := basedb.LoadConfigList(s.R.Context(), gexdb.ConfigWelcomeMessage, gexdb.ConfigWithdrawMax, gexdb.ConfigBrokerDesc)
	if err != nil { //ignore error
		xlog.Warnf("UserInfoH load config fail with %v by uid:%v", err, user.TID)
	}
	//
	var accesses = xmap.M{
		"type": user.Type,
		"role": user.Type,
	}
	return s.SendJSON(xmap.M{
		"code":              0,
		"user":              user,
		"trade_pass_setted": tradePassSetted,
		"config":            config,
		"accesses":          accesses,
		"session_id":        s.ID(),
	})
}

type updateUserArg struct {
	gexdb.User
	OldPassword  string `json:"old_password" valid:"old_password,o|s,r:0;"`
	OldTradePass string `json:"old_trade_pass" valid:"old_trade_pass,o|s,r:0;"`
}

//UpdateUserH is http handler to update user base info
/**
 *
 * @api {POST} /usr/updateUser Update User
 * @apiName UpdateUser
 * @apiGroup User
 *
 * @apiParam  {String} [name] will update user name
 * @apiParam  {String} [account] will update user account
 * @apiParam  {String} [phone] will update user phone
 * @apiParam  {String} [password] will update user password
 * @apiParam  {String} [old_password] old user password is required when update passowrd
 * @apiParam  {String} [trade_pass] will update user trade password
 * @apiParam  {String} [old_trade_pass] old user trade password is required when update trade password
 * @apiParam  {String} [image] will update user image
 * @apiParam  {Object} [external] will update user external
 * @apiParam  {Number} [status] will update user status, all status supported is <a href="#metadata-User">UserStatusAll</a>
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>, 9000 is old password is not correct
 *
 * @apiParamExample  {JSON} UpdateExternal-Example:
 * {
 *     "external": {
 *         "abc": 123
 *     }
 * }
 *
 * @apiParamExample  {JSON} UpdateName-Example:
 * {
 *     "name" : "abc"
 * }
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 * }
 *
 */
func UpdateUserH(s *web.Session) web.Result {
	arg := &updateUserArg{}
	err := RecvValidJSON(s, arg)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	user, err := gexdb.FindUser(s.R.Context(), userID)
	if err != nil {
		xlog.Errorf("UpdateUserH find user info fail with %v by uid:%v", err, userID)
		return util.ReturnCodeLocalErr(s, define.ServerError, "usr-err", err)
	}
	var updateUser, oldUser *gexdb.User
	updateUser = &arg.User
	if user.Type == gexdb.UserTypeAdmin {
		if updateUser.TID < 1 {
			updateUser.TID = userID
			updateUser.Role = 0 //not allow update role
		}
	} else {
		oldUser = &gexdb.User{
			Password:  converter.StringPtr(xhash.SHA1([]byte(arg.OldPassword))),
			TradePass: converter.StringPtr(xhash.SHA1([]byte(arg.OldTradePass))),
		}
		updateUser.TID = userID
		updateUser.Role = 0 //not allow update role
	}
	if updateUser.Password != nil {
		updateUser.Password = converter.StringPtr(xhash.SHA1([]byte(*updateUser.Password)))
	}
	if updateUser.TradePass != nil {
		updateUser.TradePass = converter.StringPtr(xhash.SHA1([]byte(*updateUser.TradePass)))
	}
	err = gexdb.UpdateUserByOld(s.R.Context(), updateUser, oldUser)
	if err != nil {
		xlog.Warnf("UpdateUser update user with %v fail with %v", converter.JSON(updateUser), err)
		code := define.ServerError
		if err == pgx.ErrNoRows {
			code = gexdb.CodeOldPasswordInvalid
		} else if strings.Contains(err.Error(), "duplicate") {
			code = define.Duplicate
		}
		return util.ReturnCodeLocalErr(s, code, "srv-err", err)
	}
	xlog.Debugf("update user to system from %v success with user:%+v", s.R.RemoteAddr, converter.JSON(updateUser))
	return s.SendJSON(xmap.M{
		"code": 0,
	})
}

//SearchUserH is http handler to search user base info
/**
 *
 * @api {GET} /usr/searchUser Search User
 * @apiName SearchUser
 * @apiGroup User
 *
 * @apiUse UserUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (User) {Array} users the user array
 * @apiSuccess (User) {Object} sellers the user seller info
 * @apiUse UserObject
 *
 * @apiParamExample  {Query} Search User:
 * key=x
 *
 * @apiParamExample  {Query} List Broker:
 * user_role=200
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "balances": {
 *         "100002": {
 *             "MMK": {
 *                 "asset": "MMK",
 *                 "create_time": 1632386799137,
 *                 "free": "100",
 *                 "locked": "0",
 *                 "status": 100,
 *                 "tid": 1001,
 *                 "update_time": 1632386799137,
 *                 "user_id": 100002
 *             },
 *             "YWE": {
 *                 "asset": "YWE",
 *                 "create_time": 1632386799136,
 *                 "free": "100",
 *                 "locked": "0",
 *                 "status": 100,
 *                 "tid": 1000,
 *                 "update_time": 1632386799136,
 *                 "user_id": 100002
 *             }
 *         }
 *     },
 *     "code": 0,
 *     "limit": 0,
 *     "skip": 0,
 *     "total": 1,
 *     "users": [
 *         {
 *             "account": "abc0",
 *             "broker_id": 0,
 *             "create_time": 1632386799135,
 *             "external": {},
 *             "image": "abc0_image",
 *             "name": "abc0_name",
 *             "phone": "abc0_123",
 *             "role": 100,
 *             "status": 100,
 *             "tid": 100002,
 *             "type": 100,
 *             "update_time": 1632386799135
 *         }
 *     ]
 * }
 *
 */
func SearchUserH(s *web.Session) web.Result {
	var searcher = &gexdb.UserUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !AdminAccess(s) {
		return util.ReturnCodeLocalErr(s, define.NotAccess, "srv-err", define.ErrNotAccess)
	}
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Warnf("SearchUserH search user by key:%v fail with %v", converter.JSON(searcher), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":  0,
		"users": searcher.Query.Users,
		"total": searcher.Count.Total,
		"skip":  searcher.Page.Skip,
		"limit": searcher.Page.Limit,
	})
}

//LoadUserH is http handler
/**
 *
 * @api {GET} /usr/loadUser Load User
 * @apiName LoadUser
 * @apiGroup User
 *
 *
 * @apiParam  {Number} user_id the user id
 *
 *
 * @apiParamExample  {Query} Request-Example:
 * user_id=000
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Success) {Number} user.tid the user id
 * @apiSuccess (Success) {Number} user.type the user type, all type supported is <a href="#metadata-User">UserTypeAll</a>
 * @apiSuccess (Success) {Number} user.role the user role, all role supported is <a href="#metadata-User">UserRoleAll</a>
 * @apiSuccess (Success) {String} user.account the user account
 * @apiSuccess (Success) {String} user.name the user nickname
 * @apiSuccess (Success) {String} user.phone the user phone
 * @apiSuccess (Success) {Object} user.external the user external info
 * @apiSuccess (Success) {Number} user.create_time the user create time
 * @apiSuccess (Success) {Number} user.update_time the user update time
 * @apiSuccess (Success) {Number} user.status the user status, all status supported is <a href="#metadata-User">UserStatusAll</a>
 * @apiSuccess (Success) {Object} balances the user balances info, mapping by balance asset as key.
 * @apiSuccess (Success) {String} balances.asset the user balances asset key
 * @apiSuccess (Success) {String} balances.free the user balances free amount
 * @apiSuccess (Success) {String} balances.locked the user balances locked amount
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "balances": {
 *         "MMK": {
 *             "asset": "MMK",
 *             "create_time": 1632387865995,
 *             "free": "100",
 *             "locked": "0",
 *             "status": 100,
 *             "tid": 1001,
 *             "update_time": 1632387865995,
 *             "user_id": 100002
 *         },
 *         "YWE": {
 *             "asset": "YWE",
 *             "create_time": 1632387865994,
 *             "free": "100",
 *             "locked": "0",
 *             "status": 100,
 *             "tid": 1000,
 *             "update_time": 1632387865994,
 *             "user_id": 100002
 *         }
 *     },
 *     "code": 0,
 *     "user": {
 *         "account": "abc0",
 *         "broker_id": 0,
 *         "create_time": 1632387865993,
 *         "external": {},
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1632387865993
 *     }
 * }
 *
 */
func LoadUserH(s *web.Session) web.Result {
	var (
		userID int64
		err    error
	)
	if err = s.ValidFormat(`user_id,R|I,R:-1`, &userID); err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !AdminAccess(s) {
		return util.ReturnCodeLocalErr(s, define.NotAccess, "srv-err", define.ErrNotAccess)
	}
	user, err := gexdb.FindUser(s.R.Context(), userID)
	if err != nil {
		xlog.Errorf("LoadUserH find target user(%v) err: %v", userID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": 0,
		"user": user,
	})
}
