//auto gen models by autogen
package basedb

import (
	"github.com/codingeasygo/util/xsql"
)

/***** metadata:Announce *****/
type AnnounceType int
type AnnounceTypeArray []AnnounceType

const (
	AnnounceTypeNormal AnnounceType = 100 //is normal type
)

//AnnounceTypeAll is the announce type
var AnnounceTypeAll = AnnounceTypeArray{AnnounceTypeNormal}

//AnnounceTypeShow is the announce type
var AnnounceTypeShow = AnnounceTypeArray{AnnounceTypeNormal}

type AnnounceStatus int
type AnnounceStatusArray []AnnounceStatus

const (
	AnnounceStatusNormal  AnnounceStatus = 100 //is normal status
	AnnounceStatusRemoved AnnounceStatus = -1  //is removed status
)

//AnnounceStatusAll is the announce status
var AnnounceStatusAll = AnnounceStatusArray{AnnounceStatusNormal, AnnounceStatusRemoved}

//AnnounceStatusShow is the announce status
var AnnounceStatusShow = AnnounceStatusArray{AnnounceStatusNormal}

/*Announce  represents announce */
type Announce struct {
	T          BaseTableName  `table:"announce"`                                          /* the table name tag */
	TID        int64          `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the announce id */
	Type       AnnounceType   `json:"type,omitempty" valid:"type,r|i,e:0;"`               /* the announce type, Normal=100:is normal type */
	Marked     int            `json:"marked,omitempty" valid:"marked,r|i,r:0;"`           /* the announce marked */
	Title      string         `json:"title,omitempty" valid:"title,r|s,l:0;"`             /* the announce title */
	Info       xsql.M         `json:"info,omitempty" valid:"info,r|s,l:0;"`               /* the announce external info */
	Content    xsql.M         `json:"content,omitempty" valid:"content,r|s,l:0;"`         /* the announce content */
	UpdateTime xsql.Time      `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the announce update time */
	CreateTime xsql.Time      `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the announce create time */
	Status     AnnounceStatus `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the announce status, Normal=100:is normal status, Removed=-1:is removed status */
}

/***** metadata:Config *****/

/*Config  represents config */
type Config struct {
	T          BaseTableName `table:"config"`                                            /* the table name tag */
	Key        string        `json:"key,omitempty" valid:"key,r|s,l:0;"`                 /*  */
	Value      string        `json:"value,omitempty" valid:"value,r|s,l:0;"`             /*  */
	UpdateTime xsql.Time     `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /*  */
}

/***** metadata:Object *****/
type ObjectStatus int
type ObjectStatusArray []ObjectStatus

const (
	ObjectStatusNormal  ObjectStatus = 100 //is normal
	ObjectStatusRemoved ObjectStatus = -1  //is removed
)

//ObjectStatusAll is the status
var ObjectStatusAll = ObjectStatusArray{ObjectStatusNormal, ObjectStatusRemoved}

//ObjectStatusShow is the status
var ObjectStatusShow = ObjectStatusArray{ObjectStatusNormal}

/*Object  represents object */
type Object struct {
	T          BaseTableName `table:"object"`                                            /* the table name tag */
	Key        string        `json:"key,omitempty" valid:"key,r|s,l:0;"`                 /* the object key */
	Value      xsql.M        `json:"value,omitempty" valid:"value,r|s,l:0;"`             /* the object value */
	UpdateTime xsql.Time     `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /*  */
	CreateTime xsql.Time     `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the create time */
	Status     ObjectStatus  `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the status, Normal=100:is normal, Removed=-1:is removed */
}

/***** metadata:VersionObject *****/
type VersionObjectStatus int
type VersionObjectStatusArray []VersionObjectStatus

const (
	VersionObjectStatusNormal   VersionObjectStatus = 100 //is normal
	VersionObjectStatusDisabled VersionObjectStatus = 200 //is disabled
	VersionObjectStatusRemoved  VersionObjectStatus = -1  //is removed
)

//VersionObjectStatusAll is the status
var VersionObjectStatusAll = VersionObjectStatusArray{VersionObjectStatusNormal, VersionObjectStatusDisabled, VersionObjectStatusRemoved}

//VersionObjectStatusShow is the status
var VersionObjectStatusShow = VersionObjectStatusArray{VersionObjectStatusNormal, VersionObjectStatusDisabled}

/*VersionObject  represents version_object */
type VersionObject struct {
	T          BaseTableName       `table:"version_object"`                                    /* the table name tag */
	TID        int64               `json:"tid,omitempty" valid:"tid,r|i,r:0;"`                 /* the primary key */
	Key        string              `json:"key,omitempty" valid:"key,r|s,l:0;"`                 /* the name of key */
	Pub        string              `json:"pub,omitempty" valid:"pub,r|s,l:0;"`                 /* the publish scoe of version object, split multi by comma, * to all, x.x.x.x for ip */
	Value      xsql.M              `json:"value,omitempty" valid:"value,r|s,l:0;"`             /* the version of key */
	UpdateTime xsql.Time           `json:"update_time,omitempty" valid:"update_time,r|i,r:1;"` /* the update time */
	CreateTime xsql.Time           `json:"create_time,omitempty" valid:"create_time,r|i,r:1;"` /* the create time */
	Status     VersionObjectStatus `json:"status,omitempty" valid:"status,r|i,e:0;"`           /* the status, Normal=100:is normal, Disabled=200:is disabled, Removed=-1:is removed */
}
