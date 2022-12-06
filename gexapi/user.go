package gexapi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhash"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/email"
	"github.com/gexservice/gexservice/base/sms"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

//LoginAccessRedirect will return the redirect url when user is not login
var LoginAccessRedirect string

//LoginAccessF is the normal user login access filter
func LoginAccessF(s *web.Session) web.Result {
	userID, ok := s.Value("user_id").(int64)
	if !ok || userID < 1 {
		return s.SendJSON(xmap.M{
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

//AdminAccessF is the admin user login access filter
func AdminAccessF(s *web.Session) web.Result {
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	return web.Continue
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
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (User) {Object} user the user info
 * @apiUse UserObject
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
 *         "role": 100,
 *         "type": 100
 *     },
 *     "code": 0,
 *     "coin_rate": [
 *         {
 *             "key": "cn",
 *             "name": "人民币",
 *             "rate": 7.2
 *         },
 *         {
 *             "key": "en",
 *             "name": "美金",
 *             "rate": 1
 *         },
 *         {
 *             "key": "xx",
 *             "name": "缅甸币",
 *             "rate": 100
 *         }
 *     ],
 *     "config": {
 *         "welcome_message": "welcom",
 *         "withdraw_max": "50000"
 *     },
 *     "session_id": "637614b4285c666ec7000001",
 *     "trade_pass_setted": 1,
 *     "user": {
 *         "account": "abc0",
 *         "create_time": 1668682932081,
 *         "favorites": {},
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1668682932081
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
	remoteAddr := s.R.RemoteAddr
	realIP := s.R.Header.Get("X-Real-IP")
	if len(realIP) > 0 {
		remoteAddr = realIP
	}
	err = gexdb.AddUserRecord(s.R.Context(), &gexdb.UserRecord{
		UserID:   user.TID,
		Type:     gexdb.UserRecordTypeLogin,
		FromAddr: remoteAddr,
		Status:   gexdb.UserRecordStatusNormal,
	})
	if err != nil {
		xlog.Errorf("LoginH add user record to sytem from %v fail with %v", s.R.RemoteAddr, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	s.Clear()
	s.SetValue("user_id", user.TID)
	s.Flush()
	xlog.Infof("LoginH user login to sytem from %v success by uid:%v", s.R.RemoteAddr, user.TID)
	return sendUserInfo(s, user)
}

//RegisterUserH is http handler to update user base info
/**
 *
 * @api {POST} /pub/registerUser Register User
 * @apiName RegisterUser
 * @apiGroup User
 *
 *
 * @apiParam  {String} [phone] the user phone
 * @apiParam  {String} [email] the user email
 * @apiParam  {String} [name] the user nick name
 * @apiParam  {String} [account] the user account，若没有，则填入phone
 * @apiParam  {String} [password] the user password
 * @apiParam  {String} code the phone verify code
 *
 * @apiSuccess (User) {Object} user the user info
 * @apiUse UserObject
 *
 * @apiParamExample  {JSON} PhoneRegister-Example:
 * {
 *     "phone" : "12345678901",
 *     "account" : "12345678901",
 *     "code" : "123",
 *     "password" : "123",
 *     "name" : "用户昵称"
 * }
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "accesses": {
 *         "role": 100,
 *         "type": 100
 *     },
 *     "code": 0,
 *     "coin_rate": [
 *         {
 *             "key": "cn",
 *             "name": "人民币",
 *             "rate": 7.2
 *         },
 *         {
 *             "key": "en",
 *             "name": "美金",
 *             "rate": 1
 *         },
 *         {
 *             "key": "xx",
 *             "name": "缅甸币",
 *             "rate": 100
 *         }
 *     ],
 *     "config": {
 *         "welcome_message": "welcom",
 *         "withdraw_max": "50000"
 *     },
 *     "session_id": "637614b4285c666ec7000001",
 *     "trade_pass_setted": 1,
 *     "user": {
 *         "account": "abc0",
 *         "create_time": 1668682932081,
 *         "favorites": {},
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1668682932081
 *     }
 * }
 */
func RegisterUserH(s *web.Session) web.Result {
	var user struct {
		Code string `json:"code"`
		gexdb.User
	}
	if err := RecvValidJSON(s, &user); err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if user.Password != nil {
		user.Password = converter.StringPtr(gexdb.EncryptionUserPassword(*user.Password))
	}
	if user.Phone != nil && len(*user.Phone) > 0 {
		expected, err := sms.LoadPhoneCode(sms.PhoneCodeTypeLogin, *user.Phone)
		if err != nil {
			xlog.Warnf("RegisterUser load phone code with %v fail with %v", *user.Phone, err)
			return util.ReturnCodeLocalErr(s, define.CodeInvalid, "srv-err", err)
		}
		if len(user.Code) < 1 || len(expected) < 1 || expected != user.Code {
			xlog.Warnf("RegisterUser verify phone code with phone %v fail with expect:%v,having:%v", converter.IndirectString(user.Phone), expected, user.Code)
			return util.ReturnCodeLocalErr(s, define.CodeInvalid, "srv-err", define.ErrCodeInvalid)
		}
	} else if user.Email != nil && len(*user.Email) > 0 {
		if len(*user.Email) > 200 || !regexp.MustCompile("^.*@.*$").MatchString(*user.Email) {
			return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", fmt.Errorf("user email is invalid or empty "))
		}
		expected, err := email.LoadEmailCode(email.EmailCodeTypeLogin, *user.Email)
		if err != nil {
			xlog.Warnf("RegisterUser load email code with %v fail with %v", *user.Email, err)
			return util.ReturnCodeLocalErr(s, define.CodeInvalid, "srv-err", err)
		}
		if len(user.Code) < 1 || len(expected) < 1 || expected != user.Code {
			xlog.Warnf("RegisterUser verify email code with email %v fail with expect:%v,having:%v", converter.IndirectString(user.Email), expected, user.Code)
			return util.ReturnCodeLocalErr(s, define.CodeInvalid, "srv-err", define.ErrCodeInvalid)
		}
	} else {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", fmt.Errorf("user phone or email is required"))
	}
	user.Type = gexdb.UserTypeNormal
	user.Role = gexdb.UserRoleNormal
	user.Status = gexdb.UserStatusNormal
	err := gexdb.AddUser(s.R.Context(), &user.User)
	if err != nil {
		xlog.Warnf("RegisterUser add user with %v fail with %v", converter.JSON(user), err)
		code := define.ServerError
		if strings.Contains(err.Error(), "duplicate") {
			code = define.Duplicate
		}
		return util.ReturnCodeLocalErr(s, code, "srv-err", err)
	}
	s.Clear()
	xlog.Infof("RegisterUserH add user to system from %v success with user:%+v", s.R.RemoteAddr, converter.JSON(user))
	s.SetValue("user_id", user.TID)
	s.Flush()
	//返回时不应该把password的内容返回
	user.Password = nil
	return sendUserInfo(s, &user.User)
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
 * @apiSuccess (Success) {Number} trade_pass_setted if user trade pass setted, 1 is setted
 * @apiSuccess (Config) {Object} config the golbal config info
 * @apiSuccess (Config) {String} config.goldbar_address the goldbar withdraw address.
 * @apiSuccess (Config) {String} config.welcom_message the welcom message
 * @apiSuccess (Config) {String} config.withdraw_max the withdraw max
 * @apiSuccess (User) {Object} user the user info
 * @apiUse UserObject
 * @apiSuccess (CoinRate) {Object} coin_rate the coin rate info, mapping by coin as key
 *
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "accesses": {
 *         "role": 100,
 *         "type": 100
 *     },
 *     "code": 0,
 *     "coin_rate": [
 *         {
 *             "key": "cn",
 *             "name": "人民币",
 *             "rate": 7.2
 *         },
 *         {
 *             "key": "en",
 *             "name": "美金",
 *             "rate": 1
 *         },
 *         {
 *             "key": "xx",
 *             "name": "缅甸币",
 *             "rate": 100
 *         }
 *     ],
 *     "config": {
 *         "welcome_message": "welcom",
 *         "withdraw_max": "50000"
 *     },
 *     "session_id": "637614b4285c666ec7000001",
 *     "trade_pass_setted": 1,
 *     "user": {
 *         "account": "abc0",
 *         "create_time": 1668682932081,
 *         "favorites": {},
 *         "image": "abc0_image",
 *         "name": "abc0_name",
 *         "phone": "abc0_123",
 *         "role": 100,
 *         "status": 100,
 *         "tid": 100002,
 *         "type": 100,
 *         "update_time": 1668682932081
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
	config, err := basedb.LoadConfigList(s.R.Context(), gexdb.ConfigWelcomeMessage, gexdb.ConfigWithdrawMax)
	if err != nil { //ignore error
		xlog.Warnf("UserInfoH load config fail with %v by uid:%v", err, user.TID)
	}
	coinRate, err := gexdb.LoadCoinRate(s.R.Context())
	if err != nil { //ignore error
		xlog.Warnf("UserInfoH load coin rate fail with %v by uid:%v", err, user.TID)
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
		"coin_rate":         coinRate,
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

//UpdateUserConfigH is http handler to update user base info
/**
 *
 * @api {POST} /usr/updateUserConfig Update User Config
 * @apiName UpdateUserConfig
 * @apiGroup User
 *
 * @apiParam  {String} [price_show_coin] will update user price show coin
 *
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 * @apiParamExample  {JSON} UpdateConfig-Example:
 * {
 *     "price_show_coin" : "abc"
 * }
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 * }
 *
 */
func UpdateUserConfigH(s *web.Session) web.Result {
	args := xmap.M{}
	_, err := s.RecvJSON(&args)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	err = gexdb.UpdateUserConfig(s.R.Context(), userID, func(config xmap.M) {
		for k, v := range args {
			if v == nil {
				delete(config, k)
			} else {
				config[k] = v
			}
		}
	})
	if err != nil {
		xlog.Errorf("UpdateUserConfigH update user config with %v fail with %v", converter.JSON(args), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	xlog.Infof("UpdateUserConfigH update user %v config with %v success", userID, converter.JSON(args))
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
 * @apiSuccess (User) {Object} user the user info
 * @apiUse UserObject
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
