package baseapi

import (
	"fmt"
	"strings"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
)

//UpsertVersionObjectH is http handler
/**
 *
 * @api {POST} /usr/upsertVersionObject Upsert Version Object
 * @apiName UpsertVersionObject
 * @apiGroup Object
 *
 *
 * @apiParam (VersionObject) {String} key the object key
 * @apiUse VersionObjectUpdate
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (VersionObject) {Object} object the Version Object info
 * @apiUse VersionObjectObject
 *
 * @apiParamExample  {JSON} Request-Example:
 * {
 *     "key": "key1",
 *     "value": {
 *         "xx": 1
 *     },
 *     "pub": "*",
 *     "status": 100
 * }
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "object": {
 *         "create_time": 1558586942626,
 *         "key": "key1",
 *         "pub": "*",
 *         "status": 100,
 *         "tid": 1000,
 *         "update_time": 1558586942626,
 *         "value": {
 *             "xx": 1
 *         }
 *     }
 * }
 *
 */
func UpsertVersionObjectH(s *web.Session) web.Result {
	object := &basedb.VersionObject{}
	err := RecvValidJSON(s, object)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !EditVersionObjectAccess(s) {
		err = fmt.Errorf("not accesss")
		return util.ReturnCodeLocalErr(s, define.NotAccess, "access-err", err)
	}
	err = basedb.UpsertVersionObject(s.R.Context(), object)
	if err != nil {
		xlog.Warnf("UpdateVersionObjectH update version object with %v fail with %v", converter.JSON(object), err)
		return util.ReturnCodeLocalErr(s, 20, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":   0,
		"object": object,
	})
}

//FindVersionObjectH is http handler
/**
 *
 * @api {GET} /pub/findVersionObject Find Version Object
 * @apiName FindVersionObject
 * @apiGroup Object
 *
 *
 * @apiParam  {Number} object_id the object id
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (VersionObject) {Object} object the Version Object info
 * @apiUse VersionObjectObject
 *
 * @apiParamExample  {Query} Request-Example:
 * object_id=100
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "object": {
 *         "create_time": 1558586942626,
 *         "key": "key1",
 *         "pub": "*",
 *         "status": 100,
 *         "tid": 1000,
 *         "update_time": 1558586942626,
 *         "value": {
 *             "xx": 1
 *         }
 *     }
 * }
 *
 */
func FindVersionObjectH(s *web.Session) web.Result {
	var objectID int64
	err := s.ValidFormat(`
			object_id,R|I,R:0;
		`, &objectID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	object, err := basedb.FindVersionObject(s.R.Context(), objectID)
	if err != nil {
		xlog.Warnf("FindVersionObjectH find object with %v fail with %v", objectID, err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":   0,
		"object": object,
	})
}

//SearchVersionObjectH is http handler
/**
 *
 * @api {GET} /pub/searchVersionObject Search Version Object
 * @apiName SearchVersionObject
 * @apiGroup Object
 *
 *
 * @apiUse VersionObjectUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (VersionObject) {Array} objects the version object info
 * @apiUse VersionObjectObject
 *
 * @apiParamExample  {Query} Request-Example:
 * key=123
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "limit": 30,
 *     "objects": [
 *         {
 *             "create_time": 1558586942626,
 *             "keys": "key1",
 *             "pub": "*",
 *             "status": 100,
 *             "tid": 1000,
 *             "update_time": 1558586942628,
 *             "value": {
 *                 "xx": 1
 *             }
 *         }
 *     ],
 *     "skip": 0,
 *     "total": 1
 * }
 */
func SearchVersionObjectH(s *web.Session) web.Result {
	searcher := &basedb.VersionObjectUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Warnf("SearchVersionObjectH search version object by %v fail with %v", converter.JSON(searcher), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":    0,
		"objects": searcher.Query.Objects,
		"total":   searcher.Count.Total,
	})
}

//LoadLatestVersionObjectH is http handler
/**
 *
 * @api {GET} /pub/vobject/<key>.json Load Latest Version Object
 * @apiName LoadLatestVersionObject
 * @apiGroup Object
 *
 *
 * @apiParam  {String} key the key, eg: /pub/vobject/xxx.json, see define VersionObjectKey* and VersionObjectApp
 *
 * @apiSuccess (200) {Number} code the response code, see common define
 *
 * @apiParamExample  {Query} Request-Example:
 * /pub/vobject/xxx.json
 *
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "xxx": 2
 * }
 */
func LoadLatestVersionObjectH(s *web.Session) web.Result {
	key := strings.SplitN(s.R.URL.Path, "vobject", 2)[1]
	key = strings.Trim(key, "/ ")
	key = strings.TrimSuffix(key, ".json")
	//
	pubs := []string{}
	pub := s.R.Header.Get("X-Forwarded-For")
	if len(pub) < 1 {
		pub = s.R.Header.Get("X-Real-IP")
	}
	if len(pub) < 1 {
		pub = s.R.RemoteAddr
	}
	pub = strings.TrimSpace(strings.Split(pub, ",")[0])
	ip := s.R.URL.Query().Get("ip")
	pubLIke := "%" + strings.SplitN(pub, ":", 2)[0] + "%"
	pubs = append(pubs, pubLIke)
	if len(ip) > 0 {
		ipLkie := "%" + pub + "-" + ip + "%"
		pubs = append(pubs, ipLkie)
	}
	object, err := basedb.LoadLatestVersionObject(s.R.Context(), key, pubs...)
	if err != nil {
		xlog.Warnf("LoadLatestVersionObjectH load latest version object by key:%v fail with %v", converter.IndirectString(key), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	var res = xsql.M{}
	if object.Value != nil {
		res = object.Value
	}
	res["code"] = 0
	return s.SendJSON(res)
}
