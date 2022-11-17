package captcha

import (
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/dchest/captcha"
	"github.com/gexservice/gexservice/base/define"
)

func CaptchaVerify(id, code string) (err error) {
	if !captcha.VerifyString(id, code) {
		err = define.ErrCodeInvalid
	}
	return
}

//NewCaptchaH is http handler
/**
 *
 * @api {GET} /pub/captcha/new New Captcha
 * @apiName NewCaptcha
 * @apiGroup Captcha
 *
 * @apiSuccess (200) {Number} code the respnose code, see the common define,
 * @apiSuccess (200) {String} captcha_id the captcha image id, to load image by /pub/captcha/${captcha_id}.png
 *
 *
 * @apiParamExample  {Query} Request-Example:
 * email=xx
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
func NewCaptchaH(s *web.Session) web.Result {
	return s.SendJSON(xmap.M{
		"code":       define.Success,
		"captcha_id": captcha.New(),
	})
}

func Hand(pre string, mux *web.SessionMux) {
	mux.HandleFunc("^"+pre+"/pub/captcha/new(\\?.*)?$", NewCaptchaH)
	mux.HandleNormal("^"+pre+"/pub/captcha/.*$", captcha.Server(captcha.StdWidth, captcha.StdHeight))
}
