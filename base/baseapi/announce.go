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

//UpsertAnnounceH is http handler to add announce
/**
 *
 * @api {POST} /usr/addAnnounce Add Announce
 * @apiName AddAnnounceH
 * @apiGroup Announce
 *
 *
 * @apiParam (Announce) {Number} [tid] the announce id to update,required with update
 * @apiUse AnnounceUpdate
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Announce) {Object} announce the Announce info
 * @apiUse AnnounceObject
 *
 *
 * @apiParamExample  {JSON} AddAnnounce:
 * {
 *     "type": 100,
 *     "marked": 10,
 *     "title": "abc",
 *     "info": {
 *         "image": "xxx"
 *     },
 *     "content": {
 *         "html": "\u003ca\u003exxx\u003c/a\u003e"
 *     }
 * }
 * @apiParamExample  {JSON} UpdateAnnounce:
 * {
 *     "tid": 1,
 *     "title": "announce title",
 * }
 * @apiParamExample  {JSON} RemoveAnnounce:
 * {
 *     "tid": 1,
 *     "status": -1
 * }
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "announce": {
 *         "create_time": 1635169896607,
 *         "info": {
 *             "image": "xxx"
 *         },
 *         "marked": 10,
 *         "status": 100,
 *         "tid": 1000,
 *         "title": "abc",
 *         "type": 100,
 *         "update_time": 1635169896607
 *     },
 *     "code": 0
 * }
 *
 *
 */
func UpsertAnnounceH(s *web.Session) web.Result {
	announce := &basedb.Announce{}
	err := RecvValidJSON(s, announce)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if announce.TID < 1 {
		announce.Status = basedb.AnnounceStatusNormal
	}
	if !EditAnnounceAccess(s) {
		err := fmt.Errorf("not accesss")
		return util.ReturnCodeLocalErr(s, define.NotAccess, "access-err", err)
	}
	err = basedb.UpsertAnnounce(s.R.Context(), announce)
	if err != nil {
		xlog.Errorf("UpsertAnnounceH upsert announce with %v fail with %v", converter.JSON(announce), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     0,
		"announce": announce,
	})
}

//SearchAnnounceH is http handler will load announce from database
/**
 *
 * @api {GET} /usr/searchAnnounce Search Announce
 * @apiName SearchAnnounce
 * @apiGroup Announce
 *
 *
 * @apiUse AnnounceUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Announce) {Array} announces the Announce info
 * @apiUse AnnounceObject
 *
 *
 * @apiParamExample  {Query} Request-Example:
 * keyword=xxx
 *
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "announces": [
 *         {
 *             "content": null,
 *             "create_time": 1635170199985,
 *             "info": {
 *                 "image": "xxx"
 *             },
 *             "marked": 10,
 *             "status": 100,
 *             "tid": 1000,
 *             "title": "abc",
 *             "type": 100,
 *             "update_time": 1635170199985
 *         }
 *     ],
 *     "code": 0,
 *     "total": 1
 * }
 *
 */
func SearchAnnounceH(s *web.Session) web.Result {
	searcher := &basedb.AnnounceUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Warnf("SearchAnnounceH search announce by %v fail with %v", converter.JSON(searcher), err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":      0,
		"announces": searcher.Query.Announces,
		"total":     searcher.Count.Total,
	})
}

//LoadAnnounceH is http handler will load announce from database
/**
 *
 * @api {GET} /usr/loadAnnounce Load Announce
 * @apiName LoadAnnounce
 * @apiGroup Announce
 *
 *
 * @apiParam  {Number} announce_id the announce id
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Announce) {Object} announce the Announce info
 * @apiUse AnnounceObject
 *
 *
 * @apiParamExample  {Query} Request-Example:
 * keyword=xxx
 *
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "announces": [
 *         {
 *             "content": null,
 *             "create_time": 1635170199985,
 *             "info": {
 *                 "image": "xxx"
 *             },
 *             "marked": 10,
 *             "status": 100,
 *             "tid": 1000,
 *             "title": "abc",
 *             "type": 100,
 *             "update_time": 1635170199985
 *         }
 *     ],
 *     "code": 0,
 *     "total": 1
 * }
 *
 */
func LoadAnnounceH(s *web.Session) web.Result {
	var announceID int64
	var err = s.ValidFormat(`
		announce_id,R|I,R:0;
	`, &announceID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	announce, err := basedb.FindAnnounce(s.R.Context(), announceID)
	if err != nil {
		xlog.Errorf("LoadAnnounceH load announce fail with %v", err)
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     0,
		"announce": announce,
	})
}
