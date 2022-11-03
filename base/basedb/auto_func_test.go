//auto gen func by autogen
package basedb

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/codingeasygo/crud"
	"github.com/codingeasygo/util/uuid"
)

func TestAutoAnnounce(t *testing.T) {
	var err error
	for _, value := range AnnounceTypeAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
		if AnnounceTypeAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if AnnounceTypeAll.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(AnnounceTypeAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(AnnounceTypeAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	for _, value := range AnnounceStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
		if AnnounceStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if AnnounceStatusAll.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(AnnounceStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(AnnounceStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithAnnounce()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	announce := &Announce{}
	announce.Valid()

	table, fields := announce.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(announce.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddAnnounce(context.Background(), announce)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(announce.TID).IsZero() {
		t.Error("not id")
		return
	}
	announce.Valid()
	err = UpdateAnnounceFilter(context.Background(), announce, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateAnnounceWheref(context.Background(), announce, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateAnnounceFilterWheref(context.Background(), announce, AnnounceFilterUpdate, "tid=$%v", announce.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findAnnounce, err := FindAnnounce(context.Background(), announce.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if announce.TID != findAnnounce.TID {
		t.Error("find id error")
		return
	}
	findAnnounce, err = FindAnnounceWheref(context.Background(), "tid=$%v", announce.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if announce.TID != findAnnounce.TID {
		t.Error("find id error")
		return
	}
	findAnnounce, err = FindAnnounceWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{announce.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if announce.TID != findAnnounce.TID {
		t.Error("find id error")
		return
	}
	findAnnounce, err = FindAnnounceWherefCall(GetQueryer, context.Background(), true, "tid=$%v", announce.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if announce.TID != findAnnounce.TID {
		t.Error("find id error")
		return
	}
	announceList, announceMap, err := ListAnnounceByID(context.Background())
	if err != nil || len(announceList) > 0 || announceMap == nil || len(announceMap) > 0 {
		t.Error(err)
		return
	}
	announceList, announceMap, err = ListAnnounceByID(context.Background(), announce.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(announceList) != 1 || announceList[0].TID != announce.TID || len(announceMap) != 1 || announceMap[announce.TID] == nil || announceMap[announce.TID].TID != announce.TID {
		t.Error("list id error")
		return
	}
	announceList = nil
	announceMap = nil
	err = ScanAnnounceByID(context.Background(), []int64{announce.TID}, &announceList, &announceMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(announceList) != 1 || announceList[0].TID != announce.TID || len(announceMap) != 1 || announceMap[announce.TID] == nil || announceMap[announce.TID].TID != announce.TID {
		t.Error("list id error")
		return
	}
	announceList = nil
	announceMap = nil
	err = ScanAnnounceWheref(context.Background(), "tid=$%v", []interface{}{announce.TID}, &announceList, &announceMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(announceList) != 1 || announceList[0].TID != announce.TID || len(announceMap) != 1 || announceMap[announce.TID] == nil || announceMap[announce.TID].TID != announce.TID {
		t.Error("list id error")
		return
	}
}

func TestAutoConfig(t *testing.T) {
	var err error
	metav := MetaWithConfig()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	config := &Config{}
	config.Valid()

	config.Key = uuid.New()

	table, fields := config.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(config.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = config.Insert(GetQueryer, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(config.Key).IsZero() {
		t.Error("not id")
		return
	}
	config.Valid()
	err = UpdateConfigFilter(context.Background(), config, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateConfigWheref(context.Background(), config, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateConfigFilterWheref(context.Background(), config, ConfigFilterUpdate, "key=$%v", config.Key)
	if err != nil {
		t.Error(err)
		return
	}
	findConfig, err := FindConfig(context.Background(), config.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if config.Key != findConfig.Key {
		t.Error("find id error")
		return
	}
	findConfig, err = FindConfigWheref(context.Background(), "key=$%v", config.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if config.Key != findConfig.Key {
		t.Error("find id error")
		return
	}
	findConfig, err = FindConfigWhereCall(GetQueryer, context.Background(), true, "and", []string{"key=$1"}, []interface{}{config.Key})
	if err != nil {
		t.Error(err)
		return
	}
	if config.Key != findConfig.Key {
		t.Error("find id error")
		return
	}
	findConfig, err = FindConfigWherefCall(GetQueryer, context.Background(), true, "key=$%v", config.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if config.Key != findConfig.Key {
		t.Error("find id error")
		return
	}
	configList, configMap, err := ListConfigByID(context.Background())
	if err != nil || len(configList) > 0 || configMap == nil || len(configMap) > 0 {
		t.Error(err)
		return
	}
	configList, configMap, err = ListConfigByID(context.Background(), config.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if len(configList) != 1 || configList[0].Key != config.Key || len(configMap) != 1 || configMap[config.Key] == nil || configMap[config.Key].Key != config.Key {
		t.Error("list id error")
		return
	}
	configList = nil
	configMap = nil
	err = ScanConfigByID(context.Background(), []string{config.Key}, &configList, &configMap, "key")
	if err != nil {
		t.Error(err)
		return
	}
	if len(configList) != 1 || configList[0].Key != config.Key || len(configMap) != 1 || configMap[config.Key] == nil || configMap[config.Key].Key != config.Key {
		t.Error("list id error")
		return
	}
	configList = nil
	configMap = nil
	err = ScanConfigWheref(context.Background(), "key=$%v", []interface{}{config.Key}, &configList, &configMap, "key")
	if err != nil {
		t.Error(err)
		return
	}
	if len(configList) != 1 || configList[0].Key != config.Key || len(configMap) != 1 || configMap[config.Key] == nil || configMap[config.Key].Key != config.Key {
		t.Error("list id error")
		return
	}
}

func TestAutoObject(t *testing.T) {
	var err error
	for _, value := range ObjectStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
		if ObjectStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if ObjectStatusAll.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(ObjectStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(ObjectStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithObject()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	object := &Object{}
	object.Valid()

	object.Key = uuid.New()

	table, fields := object.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(object.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = object.Insert(GetQueryer, context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(object.Key).IsZero() {
		t.Error("not id")
		return
	}
	object.Valid()
	err = UpdateObjectFilter(context.Background(), object, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateObjectWheref(context.Background(), object, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateObjectFilterWheref(context.Background(), object, ObjectFilterUpdate, "key=$%v", object.Key)
	if err != nil {
		t.Error(err)
		return
	}
	findObject, err := FindObject(context.Background(), object.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if object.Key != findObject.Key {
		t.Error("find id error")
		return
	}
	findObject, err = FindObjectWheref(context.Background(), "key=$%v", object.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if object.Key != findObject.Key {
		t.Error("find id error")
		return
	}
	findObject, err = FindObjectWhereCall(GetQueryer, context.Background(), true, "and", []string{"key=$1"}, []interface{}{object.Key})
	if err != nil {
		t.Error(err)
		return
	}
	if object.Key != findObject.Key {
		t.Error("find id error")
		return
	}
	findObject, err = FindObjectWherefCall(GetQueryer, context.Background(), true, "key=$%v", object.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if object.Key != findObject.Key {
		t.Error("find id error")
		return
	}
	objectList, objectMap, err := ListObjectByID(context.Background())
	if err != nil || len(objectList) > 0 || objectMap == nil || len(objectMap) > 0 {
		t.Error(err)
		return
	}
	objectList, objectMap, err = ListObjectByID(context.Background(), object.Key)
	if err != nil {
		t.Error(err)
		return
	}
	if len(objectList) != 1 || objectList[0].Key != object.Key || len(objectMap) != 1 || objectMap[object.Key] == nil || objectMap[object.Key].Key != object.Key {
		t.Error("list id error")
		return
	}
	objectList = nil
	objectMap = nil
	err = ScanObjectByID(context.Background(), []string{object.Key}, &objectList, &objectMap, "key")
	if err != nil {
		t.Error(err)
		return
	}
	if len(objectList) != 1 || objectList[0].Key != object.Key || len(objectMap) != 1 || objectMap[object.Key] == nil || objectMap[object.Key].Key != object.Key {
		t.Error("list id error")
		return
	}
	objectList = nil
	objectMap = nil
	err = ScanObjectWheref(context.Background(), "key=$%v", []interface{}{object.Key}, &objectList, &objectMap, "key")
	if err != nil {
		t.Error(err)
		return
	}
	if len(objectList) != 1 || objectList[0].Key != object.Key || len(objectMap) != 1 || objectMap[object.Key] == nil || objectMap[object.Key].Key != object.Key {
		t.Error("list id error")
		return
	}
}

func TestAutoVersionObject(t *testing.T) {
	var err error
	for _, value := range VersionObjectStatusAll {
		if value.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if value.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
		if VersionObjectStatusAll.EnumValid(int(value)) != nil {
			t.Error("not enum valid")
			return
		}
		if VersionObjectStatusAll.EnumValid(int(0)) == nil {
			t.Error("not enum valid")
			return
		}
	}
	if len(VersionObjectStatusAll.DbArray()) < 1 {
		t.Error("not array")
		return
	}
	if len(VersionObjectStatusAll.InArray()) < 1 {
		t.Error("not array")
		return
	}
	metav := MetaWithVersionObject()
	if len(metav) < 1 {
		t.Error("not meta")
		return
	}
	versionObject := &VersionObject{}
	versionObject.Valid()

	table, fields := versionObject.Meta()
	if len(table) < 1 || len(fields) < 1 {
		t.Error("not meta")
		return
	}
	fmt.Println(table, "---->", strings.Join(fields, ","))
	if table := crud.Table(versionObject.MetaWith(int64(0))); len(table) < 1 {
		t.Error("not table")
		return
	}
	err = AddVersionObject(context.Background(), versionObject)
	if err != nil {
		t.Error(err)
		return
	}
	if reflect.ValueOf(versionObject.TID).IsZero() {
		t.Error("not id")
		return
	}
	versionObject.Valid()
	err = UpdateVersionObjectFilter(context.Background(), versionObject, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateVersionObjectWheref(context.Background(), versionObject, "")
	if err != nil {
		t.Error(err)
		return
	}
	err = UpdateVersionObjectFilterWheref(context.Background(), versionObject, VersionObjectFilterUpdate, "tid=$%v", versionObject.TID)
	if err != nil {
		t.Error(err)
		return
	}
	findVersionObject, err := FindVersionObject(context.Background(), versionObject.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if versionObject.TID != findVersionObject.TID {
		t.Error("find id error")
		return
	}
	findVersionObject, err = FindVersionObjectWheref(context.Background(), "tid=$%v", versionObject.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if versionObject.TID != findVersionObject.TID {
		t.Error("find id error")
		return
	}
	findVersionObject, err = FindVersionObjectWhereCall(GetQueryer, context.Background(), true, "and", []string{"tid=$1"}, []interface{}{versionObject.TID})
	if err != nil {
		t.Error(err)
		return
	}
	if versionObject.TID != findVersionObject.TID {
		t.Error("find id error")
		return
	}
	findVersionObject, err = FindVersionObjectWherefCall(GetQueryer, context.Background(), true, "tid=$%v", versionObject.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if versionObject.TID != findVersionObject.TID {
		t.Error("find id error")
		return
	}
	versionObjectList, versionObjectMap, err := ListVersionObjectByID(context.Background())
	if err != nil || len(versionObjectList) > 0 || versionObjectMap == nil || len(versionObjectMap) > 0 {
		t.Error(err)
		return
	}
	versionObjectList, versionObjectMap, err = ListVersionObjectByID(context.Background(), versionObject.TID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(versionObjectList) != 1 || versionObjectList[0].TID != versionObject.TID || len(versionObjectMap) != 1 || versionObjectMap[versionObject.TID] == nil || versionObjectMap[versionObject.TID].TID != versionObject.TID {
		t.Error("list id error")
		return
	}
	versionObjectList = nil
	versionObjectMap = nil
	err = ScanVersionObjectByID(context.Background(), []int64{versionObject.TID}, &versionObjectList, &versionObjectMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(versionObjectList) != 1 || versionObjectList[0].TID != versionObject.TID || len(versionObjectMap) != 1 || versionObjectMap[versionObject.TID] == nil || versionObjectMap[versionObject.TID].TID != versionObject.TID {
		t.Error("list id error")
		return
	}
	versionObjectList = nil
	versionObjectMap = nil
	err = ScanVersionObjectWheref(context.Background(), "tid=$%v", []interface{}{versionObject.TID}, &versionObjectList, &versionObjectMap, "tid")
	if err != nil {
		t.Error(err)
		return
	}
	if len(versionObjectList) != 1 || versionObjectList[0].TID != versionObject.TID || len(versionObjectMap) != 1 || versionObjectMap[versionObject.TID] == nil || versionObjectMap[versionObject.TID].TID != versionObject.TID {
		t.Error("list id error")
		return
	}
}
