package gexapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/gexservice/gexservice/base/define"
	"github.com/gexservice/gexservice/gexdb"
)

func TestMessage(t *testing.T) {
	//
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	//
	ts.Should(t, "code", define.ArgsInvalid).PostJSONMap(xmap.M{"type": "xx"}, "/usr/addMessage")
	addMessage, _ := ts.Should(t, "code", define.Success).PostJSONMap(&gexdb.Message{
		Type:    gexdb.MessageTypeGlobal,
		Title:   xsql.M{"title": "test"},
		Content: xsql.M{"title": "test"},
	}, "/usr/addMessage")
	fmt.Printf("addMessage--->%v\n", converter.JSON(addMessage))
	messageID := addMessage.Int64Def(0, "/message/tid")
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/searchMessage?type=xxx")
	searchMessage, _ := ts.Should(t, "code", define.Success).GetMap("/usr/searchMessage")
	fmt.Printf("searchMessage--->%v\n", converter.JSON(searchMessage))
	ts.Should(t, "code", define.ArgsInvalid).GetMap("/usr/removeMessage?message_id=%v", "messageID")
	ts.Should(t, "code", define.Success).GetMap("/usr/removeMessage?message_id=%v", messageID)
	//
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "abc0", "123")
	ts.Should(t, "code", define.NotAccess).PostJSONMap(&gexdb.Message{
		Type:    gexdb.MessageTypeGlobal,
		Title:   xsql.M{"title": "test"},
		Content: xsql.M{"title": "test"},
	}, "/usr/addMessage")
	ts.Should(t, "code", define.NotAccess).GetMap("/usr/removeMessage?message_id=%v", messageID)
	ts.Should(t, "code", define.Success).GetMap("/usr/searchMessage")
	//
	//test error
	pgx.MockerStart()
	defer pgx.MockerStop()
	clearCookie()
	ts.Should(t, "code", define.Success).GetMap("/pub/login?username=%v&password=%v", "admin", "123")
	pgx.MockerClear()

	pgx.MockerSetCall("Rows.Scan", 2).Should(t, "code", define.ServerError).PostJSONMap(&gexdb.Message{
		Type:    gexdb.MessageTypeGlobal,
		Title:   xsql.M{"title": "test"},
		Content: xsql.M{"title": "test"},
	}, "/usr/addMessage")
	pgx.MockerSetCall("Pool.Exec", 1).Should(t, "code", define.ServerError).GetMap("/usr/removeMessage?message_id=%v", messageID)
	pgx.MockerSetCall("Pool.Query", 1).Should(t, "code", define.ServerError).GetMap("/usr/searchMessage")
}
