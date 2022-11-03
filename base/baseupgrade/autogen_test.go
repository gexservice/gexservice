package baseupgrade

import (
	"os"
	"strings"
	"testing"

	"github.com/codingeasygo/crud/gen"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xsql"
)

func init() {
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/base")
	if err != nil {
		panic(err)
	}
}

func nameConv(isTable bool, name string) string {
	if isTable {
		return gen.ConvCamelCase(true, strings.TrimPrefix(name, "_sys_"))
	}
	if name == "tid" || name == "uuid" || name == "i18n" || name == "qq" {
		return strings.ToUpper(name)
	} else if strings.HasSuffix(name, "_id") {
		return gen.ConvCamelCase(false, strings.TrimSuffix(name, "_id")+"_ID")
	} else if strings.HasSuffix(name, "_ids") {
		return gen.ConvCamelCase(false, strings.TrimSuffix(name, "_ids")+"_IDs")
	} else {
		return gen.ConvCamelCase(false, name)
	}
}

var PgGen = gen.AutoGen{
	TypeField: map[string]map[string]string{},
	FieldFilter: map[string]map[string]string{
		"version_object": {
			gen.FieldsRequired: "key,value,pub,status",
			gen.FieldsUpdate:   "pub,value,status",
		},
		"announce": {
			gen.FieldsRequired: "type,title",
			gen.FieldsOptional: "marked,info,content,status",
			gen.FieldsUpdate:   "marked,title,info,content,status",
		},
	},
	CodeAddInit: map[string]string{},
	CodeTestInit: map[string]string{
		"config": `
			ARG.Key = uuid.New()
		`,
		"object": `
			ARG.Key = uuid.New()
		`,
	},
	CodeSlice: gen.CodeSlicePG,
	TableRetAdd: map[string]string{
		"config": "",
		"object": "",
	},
	TableGenAdd: xsql.StringArray{
		"version_object",
		"announce",
	},
	TableInclude: xsql.StringArray{},
	TableExclude: xsql.StringArray{},
	Queryer:      pgx.Pool,
	TableQueryer: func(queryer interface{}, tableSQL, columnSQL, schema string) (tables []*gen.Table, err error) {
		tables, err = gen.Query(queryer, tableSQL, columnSQL, schema)
		if err != nil {
			return
		}
		for _, table := range tables {
			table.Name = strings.TrimPrefix(table.Name, "_sys_")
		}
		return
	},
	TableNameType: "BaseTableName",
	TableSQL:      gen.TableSQLPG,
	ColumnSQL:     gen.ColumnSQLPG,
	Schema:        "public",
	TypeMap:       gen.TypeMapPG,
	NameConv:      nameConv,
	GetQueryer:    "Pool",
	Out:           "../basedb/",
	OutPackage:    "basedb",
	OutStructPre: `
		//auto gen models by autogen
		package %v
		import (
			"github.com/codingeasygo/util/xsql"
		)
	`,
	OutTestPre: `
		//auto gen func by autogen
		package %v
		import (
			"context"
			"fmt"
			"reflect"
			"strings"
			"testing"

			"github.com/codingeasygo/crud"
			"github.com/codingeasygo/util/uuid"
		)
	`,
}

func TestPgGen(t *testing.T) {
	// defer os.RemoveAll(PgGen.Out)
	os.MkdirAll(PgGen.Out, os.ModePerm)
	err := PgGen.Generate()
	if err != nil {
		t.Error(err)
		return
	}
}
