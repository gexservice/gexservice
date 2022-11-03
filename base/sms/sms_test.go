package sms

import (
	"fmt"
	"testing"

	"github.com/Centny/rediscache"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web/httptest"
	"github.com/gexservice/gexservice/base/define"
)

func TestVerifyPhone(t *testing.T) {
	func() {
		defer func() {
			recover()
		}()
		_ = SendSms(nil, "", nil)
	}()
	//
	redisURI := "redis.loc:6379?db=1"
	rediscache.InitRedisPool(redisURI)
	Redis = rediscache.C
	ts := httptest.NewMuxServer()
	Hand("", ts.Mux)
	SendSms = func(v *VerifyPhone, phoneNumber string, templateParam xmap.M) (err error) {
		return nil
	}
	//
	res, err := ts.GetMap("/pub/sendLoginSms?phone=1234567890")
	if err != nil || res.Int64("code") != 0 {
		t.Errorf("%v,%v", err, res)
		return
	}
	having, err := LoadPhoneCode(PhoneCodeTypeLogin, "1234567890")
	if err != nil || len(having) < 1 {
		t.Error(err)
		return
	}
	res, err = ts.GetMap("/pub/sendLoginSms?phone=1234567890")
	if err != nil || res.Int64("code") != define.Frequently {
		t.Errorf("%v,%v", err, res)
		return
	}
	//
	//test error
	rediscache.MockerStart()
	defer rediscache.MockerStop()
	//
	res, err = ts.GetMap("/pub/sendLoginSms?phone=")
	if err != nil || res.Int64("code") == 0 {
		t.Errorf("%v,%v", err, res)
		return
	}
	SendLoginSmsH.CalledUser = map[string]int64{}
	rediscache.MockerSet("Conn.Do", 1)
	res, err = ts.GetMap("/pub/sendLoginSms?phone=1234567810")
	if err != nil || res.Int64("code") != define.ServerError {
		t.Errorf("%v,%v", err, res)
		return
	}
	fmt.Println(converter.JSON(res))
	rediscache.MockerClear()
	//
	SendLoginSmsH.CalledUser = map[string]int64{}
	SendSms = func(v *VerifyPhone, phoneNumber string, templateParam xmap.M) (err error) {
		return fmt.Errorf("mock error")
	}
	res, err = ts.GetMap("/pub/sendLoginSms?phone=1234567810")
	if err != nil || res.Int64("code") != define.ServerError {
		t.Errorf("%v,%v", err, res)
		return
	}
}
