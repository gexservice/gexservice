package baseapi

import (
	"fmt"
	"testing"

	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xsql"
	"github.com/codingeasygo/web"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/define"
)

func TestVersionObject(t *testing.T) {
	clearCookie()
	//
	key := "key1"
	pub := "*"
	objectArgs := &basedb.VersionObject{
		Key: key,
		Pub: pub,
		Value: xsql.M{
			"xx": 1,
		},
		Status: basedb.VersionObjectStatusNormal,
	}
	//
	fmt.Printf("addVersionObjectArgs->%v\n\n", converter.JSON(objectArgs))
	addVersionObject, err := ts.PostJSONMap(objectArgs, "/usr/upsertVersionObject")
	if err != nil || addVersionObject.IntDef(-1, "code") != 0 || addVersionObject.Int64("/object/tid") < 1 || addVersionObject.Int64("/object/value/xx") != 1 {
		t.Errorf("err:%v,addVersionObject:%v", err, converter.JSON(addVersionObject))
		return
	}
	objectID := addVersionObject.Int64("/object/tid")
	objectArgs.TID = objectID
	fmt.Printf("addVersionObject->%v\n\n", converter.JSON(addVersionObject))
	//
	fmt.Printf("updateVersionObjectArgs->%v\n\n", converter.JSON(objectArgs))
	updateVersionObject, err := ts.PostJSONMap(objectArgs, "/usr/upsertVersionObject")
	if err != nil || updateVersionObject.IntDef(-1, "code") != 0 {
		t.Errorf("err:%v,updateVersionObject:%v", err, converter.JSON(updateVersionObject))
		return
	}
	//
	findVersionObject, err := ts.GetMap("/pub/findVersionObject?object_id=%v", objectID)
	if err != nil || findVersionObject.IntDef(-1, "code") != 0 || findVersionObject.Int64("/object/value/xx") != 1 {
		t.Errorf("err:%v,findVersionObject:%v", err, converter.JSON(findVersionObject))
		return
	}
	fmt.Printf("findVersionObject->%v\n\n", converter.JSON(findVersionObject))
	//
	searchVersionObject, err := ts.GetMap("/pub/searchVersionObject?key=%v", key)
	if err != nil || searchVersionObject.IntDef(-1, "code") != 0 || len(searchVersionObject.ArrayMapDef(nil, "/objects")) != 1 {
		t.Errorf("err:%v,searchVersionObject:%v", err, converter.JSON(searchVersionObject))
		return
	}
	fmt.Printf("searchVersionObject->%v\n\n", converter.JSON(searchVersionObject))
	//
	vobject, err := ts.GetMap("/pub/vobject/%v.json", key)
	if err != nil || vobject.IntDef(-1, "code") != 0 || vobject.Str("xx") != "1" {
		t.Errorf("err:%v,vobject:%v", err, converter.JSON(vobject))
		return
	}
	fmt.Printf("\n vobject->%v\n\n", converter.JSON(vobject))
	//
	pub1 := "127.0.0.1"
	objectArg1 := &basedb.VersionObject{
		Key: key,
		Pub: pub1,
		Value: xsql.M{
			"xx": 2,
		},
		Status: basedb.VersionObjectStatusNormal,
	}
	addVersionObject1, err := ts.PostJSONMap(objectArg1, "/usr/upsertVersionObject")
	if err != nil || addVersionObject1.IntDef(-1, "code") != 0 ||
		addVersionObject1.Int64("/object/tid") < 1 || addVersionObject1.Int64("/object/value/xx") != 2 {
		t.Errorf("err:%v,upsertVersionObject:%v", err, converter.JSON(addVersionObject1))
		return
	}
	objectID1 := addVersionObject1.Int64("/object/tid")
	objectArg1.TID = objectID1
	//
	vobject, err = ts.GetMap("/pub/vobject/%v.json", key)
	if err != nil || vobject.IntDef(-1, "code") != 0 || vobject.Str("xx") != "2" {
		t.Errorf("err:%v,vobject:%v", err, converter.JSON(vobject))
		return
	}
	fmt.Printf("\n vobject->%v\n\n", converter.JSON(vobject))
	//
	pub2 := "127.0.0.1-xxx"
	objectArg2 := &basedb.VersionObject{
		Key: key,
		Pub: pub2,
		Value: xsql.M{
			"xx": 3,
		},
		Status: basedb.VersionObjectStatusNormal,
	}
	addVersionObject2, err := ts.PostJSONMap(objectArg2, "/usr/upsertVersionObject")
	if err != nil || addVersionObject2.IntDef(-1, "code") != 0 ||
		addVersionObject2.Int64("/object/tid") < 1 || addVersionObject2.Int64("/object/value/xx") != 3 {
		t.Errorf("err:%v,upsertVersionObject:%v", err, converter.JSON(addVersionObject2))
		return
	}
	objectID2 := addVersionObject2.Int64("/object/tid")
	objectArg2.TID = objectID2
	//
	vobject, err = ts.GetMap("/pub/vobject/%v.json?ip=xxx", key)
	if err != nil || vobject.IntDef(-1, "code") != 0 || vobject.Str("xx") != "3" {
		t.Errorf("err:%v,vobject:%v", err, converter.JSON(vobject))
		return
	}
	//
	//
	pgx.MockerStart()
	defer pgx.MockerStop()
	//add object
	res, err := ts.PostJSONMap("xx", "/usr/upsertVersionObject")
	if err != nil || res.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
	EditVersionObjectAccess = func(s *web.Session) bool { return false }
	res, err = ts.PostJSONMap(objectArgs, "/usr/upsertVersionObject")
	if err != nil || res.IntDef(-1, "code") != define.NotAccess {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
	pgx.MockerClear()
	EditVersionObjectAccess = func(s *web.Session) bool { return true }
	pgx.MockerClear()
	pgx.MockerSet("Pool.Exec", 1)
	res, err = ts.PostJSONMap(objectArgs, "/usr/upsertVersionObject")
	if err != nil || res.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
	pgx.MockerClear()
	//find object
	res, err = ts.GetMap("/pub/findVersionObject?object_id=%v", "xxx")
	if err != nil || res.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
	pgx.MockerSet("Row.Scan", 1)
	res, err = ts.GetMap("/pub/findVersionObject?object_id=%v", "1")
	if err != nil || res.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
	pgx.MockerClear()
	//search object
	res, err = ts.GetMap("/pub/searchVersionObject?skip=xxx")
	if err != nil || res.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
	pgx.MockerSet("Pool.Query", 1)
	res, err = ts.GetMap("/pub/searchVersionObject")
	if err != nil || res.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
	pgx.MockerClear()
	//vobject
	res, err = ts.GetMap("/pub/vobject/xxxx")
	if err != nil || res.IntDef(-1, "code") == 0 {
		t.Errorf("err:%v,res:%v", err, res)
		return
	}
}
