package baseapi

import (
	"fmt"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
)

//UpdateSysConfigH is http handler will update sys config to database
/**
 *
 * @api {POST} /usr/updateSysConfig Update Sys Config
 * @apiName UpdateSysConfig
 * @apiGroup Sys
 *
 *
 * @apiParam  {String} key the config key to update by value, see define by Config*
 *
 * @apiSuccess (200) {Number} code the response code, see define common
 * @apiSuccess (200) {Object} config the config object
 * @apiSuccess (200) {Object} config.xxx the config value by key xxx
 *
 * @apiParamExample  {JSON} Request-Example:
 * {
 *     "key": "value",
 * }
 *
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "config": {
 *         "key": "value",
 *     }
 * }
 *
 *
 */
func UpdateSysConfigH(s *web.Session) web.Result {
	config := xmap.M{}
	_, err := s.RecvJSON(&config)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !EditSysConfigAccess(s) {
		err = fmt.Errorf("not accesss")
		return util.ReturnCodeLocalErr(s, define.NotAccess, "access-err", err)
	}
	err = basedb.UpdateConfigList(s.R.Context(), config)
	if err != nil {
		xlog.Warnf("UpdateSysConfigH update sys config by %v fail with %v", converter.JSON(config), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	xlog.Debugf("UpdateSysConfigH update sys config by %v success", converter.JSON(config))
	return s.SendJSON(xmap.M{
		"code":   0,
		"config": config,
	})
}

//LoadSysConfigH is http handler will load sys config from database
/**
 *
 * @api {GET} /usr/loadSysConfig Load Sys Config
 * @apiName LoadSysConfig
 * @apiGroup Sys
 *
 *
 * @apiSuccess (200) {Number} code the response code, see common define
 * @apiSuccess (200) {Object} config the config object
 * @apiSuccess (200) {Object} config.xxx the config value by key xxx
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "config": {
 *         "key": "value",
 *     }
 * }
 *
 *
 */
func LoadSysConfigH(s *web.Session) web.Result {
	if !EditSysConfigAccess(s) {
		err := fmt.Errorf("not accesss")
		return util.ReturnCodeLocalErr(s, define.NotAccess, "access-err", err)
	}
	config, err := basedb.LoadConfigList(s.R.Context(), basedb.ConfigAll...)
	if err != nil {
		xlog.Warnf("LoadSysConfigH load sys config fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":   0,
		"config": config,
	})
}

type ConfigLoader []string

func (c ConfigLoader) SrvHTTP(s *web.Session) web.Result {
	config, err := basedb.LoadConfigList(s.R.Context(), c...)
	if err != nil {
		xlog.Warnf("ConfigLoader load config fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":   0,
		"config": config,
	})
}
