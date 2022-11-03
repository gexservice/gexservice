package baseapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
)

func TestAnnounce(t *testing.T) {
	addAnnounceArgs := &basedb.Announce{
		Type:    basedb.AnnounceTypeNormal,
		Marked:  10,
		Title:   "abc",
		Info:    xsql.M{"image": "xxx"},
		Content: xsql.M{"html": "<a>xxx</a>"},
	}
	//
	fmt.Printf("addAnnounceArgs->%v\n\n", converter.JSON(addAnnounceArgs))
	addAnnounce, err := ts.PostJSONMap(addAnnounceArgs, "/usr/upsertAnnounce")
	if err != nil || addAnnounce.IntDef(-1, "code") != 0 || addAnnounce.Int64("/announce/tid") < 1 {
		t.Errorf("err:%v,addAnnounce:%v", err, converter.JSON(addAnnounce))
		return
	}
	addAnnounceArgs.TID = addAnnounce.Int64("/announce/tid")
	fmt.Printf("addAnnounce->%v\n\n", converter.JSON(addAnnounce))

	fmt.Printf("updateAnnounceArgs->%v\n\n", converter.JSON(addAnnounceArgs))
	updateAnnounce, err := ts.PostJSONMap(addAnnounceArgs, "/usr/upsertAnnounce")
	if err != nil || updateAnnounce.IntDef(-1, "code") != 0 {
		t.Errorf("err:%v,updateAnnounce:%v", err, converter.JSON(updateAnnounce))
		return
	}
	fmt.Printf("updateAnnounce->%v\n\n", converter.JSON(updateAnnounce))

	searchAnnounce, err := ts.GetMap("/pub/searchAnnounce?keyword=abc")
	if err != nil || searchAnnounce.IntDef(-1, "code") != 0 {
		t.Errorf("err:%v,searchAnnounce:%v", err, converter.JSON(searchAnnounce))
		return
	}
	fmt.Printf("searchAnnounce->%v\n\n", converter.JSON(searchAnnounce))

	loadAnnounce, err := ts.GetMap("/pub/loadAnnounce?announce_id=%v", addAnnounceArgs.TID)
	if err != nil || loadAnnounce.IntDef(-1, "code") != 0 {
		t.Errorf("err:%v,loadAnnounce:%v", err, converter.JSON(loadAnnounce))
		return
	}
	fmt.Printf("loadAnnounce->%v\n\n", converter.JSON(loadAnnounce))

	pgx.MockerStart()
	defer pgx.MockerStop()
	var res xmap.M

	//update error
	res, err = ts.PostJSONMap("", "/usr/upsertAnnounce")
	if err != nil || res.IntDef(-1, "code") != define.ArgsInvalid {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
	res, err = ts.PostJSONMap(xmap.M{}, "/usr/upsertAnnounce")
	if err != nil || res.IntDef(-1, "code") != define.ArgsInvalid {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
	pgx.MockerClear()
	pgx.MockerSet("Pool.Exec", 1)
	res, err = ts.PostJSONMap(addAnnounceArgs, "/usr/upsertAnnounce")
	if err != nil || res.IntDef(-1, "code") != define.ServerError {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
	pgx.MockerClear()

	//search error
	res, err = ts.GetMap("/pub/searchAnnounce?type=xxx")
	if err != nil || res.IntDef(-1, "code") != define.ArgsInvalid {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
	pgx.MockerClear()
	pgx.MockerSet("Pool.Query", 1)
	res, err = ts.GetMap("/pub/searchAnnounce")
	if err != nil || res.IntDef(-1, "code") != define.ServerError {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
	pgx.MockerClear()

	//load error
	res, err = ts.GetMap("/pub/loadAnnounce?announce_id=xxx")
	if err != nil || res.IntDef(-1, "code") != define.ArgsInvalid {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
	pgx.MockerClear()
	pgx.MockerSet("Rows.Scan", 1)
	res, err = ts.GetMap("/pub/loadAnnounce?announce_id=%v", addAnnounceArgs.TID)
	if err != nil || res.IntDef(-1, "code") != define.ServerError {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
	pgx.MockerClear()

	//not access
	EditAnnounceAccess = func(s *web.Session) bool {
		return false
	}
	res, err = ts.PostJSONMap(addAnnounceArgs, "/usr/upsertAnnounce")
	if err != nil || res.IntDef(-1, "code") != define.NotAccess {
		t.Errorf("err:%v,res:%v", err, converter.JSON(res))
		return
	}
}
