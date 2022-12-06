package gexapi

import (
	"sync"

	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/web"
)

type OnlineInfo struct {
	UserID int64
	Last   xsql.Time `json:"last"`
	Conn   int       `json:"conn"`
}

type OnlineHander struct {
	userAll map[int64]*OnlineInfo
	userLck sync.RWMutex
	mux     *web.SessionMux
	next    web.Handler
}

func NewOnlineHander(mux *web.SessionMux, next web.Handler) (handler *OnlineHander) {
	handler = &OnlineHander{
		userAll: map[int64]*OnlineInfo{},
		userLck: sync.RWMutex{},
		mux:     mux,
		next:    next,
	}
	return
}

func (o *OnlineHander) SrvHTTP(s *web.Session) web.Result {
	o.onConnect(s)
	o.next.SrvHTTP(s)
	o.onDisconnect(s)
	return web.Return
}

func (o *OnlineHander) onConnect(s *web.Session) {
	userID := s.Int64Def(0, "user_id")
	if userID < 1 {
		return
	}
	o.userLck.Lock()
	defer o.userLck.Unlock()
	info := o.userAll[userID]
	if info == nil {
		info = &OnlineInfo{}
		o.userAll[userID] = info
	}
	info.Last = xsql.TimeNow()
	info.Conn++
	info.UserID = userID
}

func (o *OnlineHander) onDisconnect(s *web.Session) {
	userID := s.Int64Def(0, "user_id")
	if userID < 1 {
		return
	}
	o.userLck.Lock()
	defer o.userLck.Unlock()
	info := o.userAll[userID]
	if info == nil {
		return
	}
	info.Conn--
	if info.Conn < 1 {
		delete(o.userAll, userID)
	}
}

func (o *OnlineHander) List(userIDs ...int64) (infoes map[int64]*OnlineInfo) {
	o.userLck.RLock()
	defer o.userLck.RUnlock()
	infoes = map[int64]*OnlineInfo{}
	for _, userID := range userIDs {
		info := o.userAll[userID]
		if info != nil {
			infoes[userID] = info
		}
	}
	return
}

func (o *OnlineHander) Size() (size int) {
	o.userLck.RLock()
	defer o.userLck.RUnlock()
	size = len(o.userAll)
	return
}
