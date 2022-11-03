package sms

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"

	"github.com/gomodule/redigo/redis"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const (
	//PhoneCodeTypeVerify is phone code type to verify
	PhoneCodeTypeVerify = "verify"
	//PhoneCodeTypeLogin is phone code type to login
	PhoneCodeTypeLogin   = "login"
	VerifyPhoneTypeUser  = "user"
	VerifyPhoneTypePhone = "phone"
)

//Redis will return redis connection
var Redis = func() redis.Conn {
	panic("redis is not initial")
}

//SendSms will send message to phone
var SendSms = func(v *VerifyPhone, phoneNumber string, templateParam xmap.M) (err error) {
	panic("send sms is not initial")
}

//LoadPhoneCode will return send code
func LoadPhoneCode(key, phone string) (having string, err error) {
	conn := Redis()
	defer conn.Close()
	val, err := conn.Do("get", key+"_sms_"+phone)
	if err != nil {
		return
	}
	if val != nil {
		having, err = redis.String(val, nil)
	}
	return
}

//VerifyPhone is verify phone impl
type VerifyPhone struct {
	Key           string
	Type          string
	Limit         int64
	CalledUser    map[string]int64
	CalledUserLck sync.RWMutex
}

//NewVerifyPhone will craete new VerifyPhone
func NewVerifyPhone(key, verifyType string, limit int64) (v *VerifyPhone) {
	v = &VerifyPhone{
		Key:           key,
		Type:          verifyType,
		Limit:         limit,
		CalledUser:    map[string]int64{},
		CalledUserLck: sync.RWMutex{},
	}
	return
}

//SrvHTTP is http handler
func (v *VerifyPhone) SrvHTTP(hs *web.Session) web.Result {
	var phone string
	err := hs.ValidFormat(`
		phone,R|S,L:0;
	`, &phone)
	if err != nil {
		return util.ReturnCodeLocalErr(hs, define.ArgsInvalid, "arg-err", err)
	}
	unique := ""
	if v.Type == "user" {
		unique = hs.Str("uid")
	} else if v.Type == "phone" {
		unique = phone
	} else {
		unique = strings.Split(hs.R.RemoteAddr, ":")[0]
	}
	v.CalledUserLck.Lock()
	now := xsql.TimeNow().Timestamp()
	last := v.CalledUser[unique]
	if now-last < v.Limit {
		v.CalledUserLck.Unlock()
		// return util.ReturnCodeLocalErr(hs, define.Frequently, "srv-err", err)
		return hs.SendJSON(map[string]interface{}{
			"code":    define.Frequently,
			"after":   v.Limit - (now - last),
			"message": "call too frequently",
		})
	}
	v.CalledUser[unique] = now
	v.CalledUserLck.Unlock()
	number := rand.Intn(1000000)
	err = SendSms(v, phone, xmap.M{
		"code": number,
	})
	if err != nil {
		xlog.Warnf("VerifyPhone send sms by %v fail with %v", phone, err)
		return util.ReturnCodeLocalErr(hs, define.ServerError, "srv-err", err)
	}
	conn := Redis()
	defer conn.Close()
	_, err = conn.Do("setex", v.Key+"_sms_"+phone, 1800, number)
	if err != nil {
		xlog.Warnf("VerifyPhone save sened sms by %v fail with %v", phone, err)
		return util.ReturnCodeLocalErr(hs, define.ServerError, "srv-err", err)
	}
	return util.ReturnCodeData(hs, 0, "OK")
}

//SendVerifySmsH is http handler
/**
 *
 * @api {GET} /usr/sendVerifySms Send Verify Sms
 * @apiName SendVerifySms
 * @apiGroup User
 *
 *
 * @apiParam  {String} phone the phone number
 * @apiSuccess (200) {Number} code the respnose code, see the common define,
 * @apiSuccess (200) {Number} after the after time when call frequently
 *
 * @apiParamExample  {Query} Request-Example:
 * phone=xx
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "data": "OK"
 * }
 * @apiSuccessExample {JSON} Frequently-Response:
 * {
 *     "code": 1600,
 *     "after": 10000,
 * }
 *
 */
var SendVerifySmsH = NewVerifyPhone(PhoneCodeTypeVerify, VerifyPhoneTypeUser, 10000)

//SendLoginSmsH is http handler
/**
 *
 * @api {GET} /pub/sendLoginSms Send Login Sms
 * @apiName SendLoginSms
 * @apiGroup User
 *
 *
 * @apiParam  {String} phone the phone number
 * @apiSuccess (200) {Number} code the respnose code, see the common define,
 * @apiSuccess (200) {Number} after the after time when call frequently
 *
 *
 * @apiParamExample  {Query} Request-Example:
 * phone=xx
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "data": "OK"
 * }
 * @apiSuccessExample {JSON} Frequently-Response:
 * {
 *     "code": 1600,
 *     "after": 10000,
 * }
 *
 */
var SendLoginSmsH = NewVerifyPhone(PhoneCodeTypeLogin, VerifyPhoneTypePhone, 10000)

func LoadPhoneCodeH(s *web.Session) web.Result {
	var key, phone string
	var err = s.ValidFormat(`
		key,R|S,L:0;
		phone,R|S,L:0;
	`, &key, &phone)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	xlog.Warnf("debug api is enabled for load phone code")
	having, err := LoadPhoneCode(key, phone)
	if err != nil {
		xlog.Warnf("DebugLoadPhoneCodeH load %v sended sms by %v fail with %v", key, phone, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(map[string]interface{}{
		"code":      0,
		"phoneCode": having,
	})
}

func Hand(pre string, mux *web.SessionMux) {
	mux.Handle("^"+pre+"/usr/sendVerifySms(\\?.*)?$", SendVerifySmsH)
	mux.Handle("^"+pre+"/pub/sendLoginSms(\\?.*)?$", SendLoginSmsH)
}

func HandDebug(pre string, mux *web.SessionMux) {
	xlog.Warnf("debug api is enabled for load phone code")
	mux.HandleFunc("^"+pre+"/pub/loadPhoneCode(\\?.*)?$", LoadPhoneCodeH)
}
