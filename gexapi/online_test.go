package gexapi

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/codingeasygo/web"
	"github.com/codingeasygo/web/httptest"
	"golang.org/x/net/websocket"
)

func TestOnline(t *testing.T) {
	ts := httptest.NewMuxServer()
	ts.Mux.FilterFunc("/test", func(s *web.Session) web.Result {
		s.SetValue("user_id", 100)
		return web.Continue
	})
	wsServer := websocket.Server{Handler: func(c *websocket.Conn) {
		io.Copy(io.Discard, c)
	}}
	handler := NewOnlineHander(ts.Mux, web.HandlerFunc(func(s *web.Session) web.Result {
		wsServer.ServeHTTP(s.W, s.R)
		return web.Return
	}))
	ts.Mux.Handle("/", handler)
	conn, err := websocket.Dial(strings.ReplaceAll(ts.URL, "http://", "ws://")+"/test", "", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	infoes := handler.List(100)
	if err != nil || len(infoes) != 1 {
		t.Error(err)
		return
	}
	conn.Close()
	time.Sleep(100 * time.Millisecond)
	infoes = handler.List(100)
	if err != nil || len(infoes) != 0 {
		t.Error(err)
		return
	}

	//
	ts.Mux.FilterFunc("/error", func(s *web.Session) web.Result {
		s.SetValue("user_id", 0)
		handler.onConnect(s)
		handler.onDisconnect(s)
		s.SetValue("user_id", 1000)
		handler.onDisconnect(s)
		return web.Continue
	})
	conn, err = websocket.Dial(strings.ReplaceAll(ts.URL, "http://", "ws://")+"/error", "", ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	conn.Close()
}
