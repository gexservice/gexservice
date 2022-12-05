package gexapi

import (
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/base/util"
	"github.com/gexservice/gexservice/base/xlog"
	"github.com/gexservice/gexservice/gexdb"
)

//AddMessageH is http handler
/**
 *
 * @api {GET} /usr/addMessage Add Message
 * @apiName AddMessage
 * @apiGroup Message
 *
 * @apiUse MessageUpdate
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Message) {Object} message the message records
 * @apiUse MessageObject
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0,
 *     "message": {
 *         "content": {
 *             "title": "test"
 *         },
 *         "create_time": 1670244390094,
 *         "status": 100,
 *         "tid": 1000,
 *         "title": {
 *             "title": "test"
 *         },
 *         "type": 200,
 *         "update_time": 1670244390094
 *     }
 * }
 *
 */
func AddMessageH(s *web.Session) web.Result {
	var message = &gexdb.Message{}
	_, err := s.RecvJSON(message)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	message.Status = gexdb.MessageStatusNormal
	err = gexdb.AddMessage(s.R.Context(), message)
	if err != nil {
		xlog.Errorf("AddMessageH add message fail with %v by %v", err, converter.JSON(message))
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":    define.Success,
		"message": message,
	})
}

//RemoveMessageH is http handler
/**
 *
 * @api {GET} /usr/removeMessage Remove Message
 * @apiName RemoveMessage
 * @apiGroup Message
 *
 * @apiParam  {String} symbol the symbol to add
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 *
 *
 * @apiSuccessExample {JSON} Success-Response:
 * {
 *     "code": 0
 * }
 *
 */
func RemoveMessageH(s *web.Session) web.Result {
	var messageID int64
	err := s.ValidFormat(`
		message_id,r|i,r:0;
	`, &messageID)
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	if !AdminAccess(s) {
		return s.SendJSON(xmap.M{
			"code":    define.NotAccess,
			"message": define.ErrNotAccess.String(),
		})
	}
	message := &gexdb.Message{
		TID:    messageID,
		Status: gexdb.MessageStatusRemoved,
	}
	err = gexdb.UpdateMessageFilter(s.R.Context(), message, "status")
	if err != nil {
		xlog.Errorf("RemoveMessageH update message fail with %v by %v", err, converter.JSON(message))
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code": define.Success,
	})
}

//SearchMessageH is http handler
/**
 *
 * @api {GET} /usr/searchMessage Search Message
 * @apiName SearchMessage
 * @apiGroup Message
 *
 *
 * @apiUse MessageUnifySearcher
 *
 * @apiSuccess (Success) {Number} code the result code, see the common define <a href="#metadata-ReturnCode">ReturnCode</a>
 * @apiSuccess (Message) {Array} messages the message records
 * @apiUse MessageObject
 *
 * @apiSuccessExample {type} Success-Response:
 * {
 *     "code": 0,
 *     "messages": [
 *         {
 *             "content": {
 *                 "title": "test"
 *             },
 *             "create_time": 1670244390094,
 *             "status": 100,
 *             "tid": 1000,
 *             "title": {
 *                 "title": "test"
 *             },
 *             "type": 200,
 *             "update_time": 1670244390094
 *         }
 *     ],
 *     "total": 1
 * }
 */
func SearchMessageH(s *web.Session) web.Result {
	searcher := &gexdb.MessageUnifySearcher{}
	err := s.Valid(searcher, "#all")
	if err != nil {
		return util.ReturnCodeLocalErr(s, define.ArgsInvalid, "arg-err", err)
	}
	userID := s.Int64("user_id")
	if !AdminAccess(s) {
		searcher.Where.ToUserID = userID
	}
	err = searcher.Apply(s.R.Context())
	if err != nil {
		xlog.Errorf("ListMessageH list message fail with %v by %v", err, converter.JSON(searcher))
		return util.ReturnCodeLocalErr(s, define.ServerError, "srv-err", err)
	}
	return s.SendJSON(xmap.M{
		"code":     define.Success,
		"messages": searcher.Query.Messages,
		"total":    searcher.Count.Total,
	})
}
