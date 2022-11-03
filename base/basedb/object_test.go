package basedb

import (
	"context"
	"fmt"
	"testing"

	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xsql"
)

func TestObject(t *testing.T) {
	clear()
	updated, err := UpsertObject(context.Background(), "obj", xmap.M{
		"abc": 1,
	})
	if err != nil || updated != 1 {
		t.Error(err)
		return
	}
	data, err := LoadObject(context.Background(), "obj")
	if err != nil || data.Int("abc") != 1 {
		t.Error(err)
		return
	}
	//
	//test error
	_, err = UpsertObject(context.Background(), "obj", xmap.M{
		"abc": TestObject,
	})
	if err == nil {
		t.Error(err)
		return
	}
	_, err = LoadObject(context.Background(), "xxxx")
	if err == nil {
		t.Error(err)
		return
	}
}

func TestVersionObject(t *testing.T) {
	clear()
	key := "key1"
	object := &VersionObject{
		Key: key,
		Pub: "*",
		Value: xsql.M{
			"xx": 1,
		},
		Status: VersionObjectStatusNormal,
	}
	err := UpsertVersionObject(context.Background(), object)
	if err != nil {
		t.Error(err)
		return
	}
	err = UpsertVersionObject(context.Background(), object)
	if err != nil {
		t.Error(err)
		return
	}
	searcher := VersionObjectUnifySearcher{}
	searcher.Where.Key = object.Key
	searcher.Where.Status = VersionObjectStatusArray{object.Status}
	err = searcher.Apply(context.Background())
	if err != nil || len(searcher.Query.Objects) < 1 || searcher.Count.Total < 1 {
		t.Error(err)
		return
	}
	//
	pub1 := "127.0.0.1"
	object1 := &VersionObject{
		Key: key,
		Pub: "127.0.0.1",
		Value: xsql.M{
			"xx": 2,
		},
		Status: VersionObjectStatusNormal,
	}
	err = AddVersionObject(context.Background(), object1)
	if err != nil {
		t.Error(err)
		return
	}
	loadObject, err := LoadLatestVersionObject(context.Background(), key)
	if err != nil || loadObject.TID != object1.TID {
		fmt.Println("--->", loadObject.TID)
		t.Error(err)
		return
	}
	loadObject, err = LoadLatestVersionObject(context.Background(), key, pub1)
	if err != nil || loadObject.TID != object1.TID {
		t.Error(err)
		return
	}
	pub2 := "192.168.1.1"
	loadObject, err = LoadLatestVersionObject(context.Background(), key, pub2)
	if err != nil || loadObject.TID != object.TID {
		t.Error(err)
		return
	}
}
