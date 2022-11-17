package gexupgrade

import (
	"os"
	"strings"
	"testing"

	"github.com/codingeasygo/crud/gen"
	"github.com/codingeasygo/crud/pgx"
	"github.com/codingeasygo/util/xsql"
)

func init() {
	_, err := pgx.Bootstrap("postgresql://dev:123@psql.loc:5432/gexservice")
	if err != nil {
		panic(err)
	}
	// _, _, err = pgx.Pool().Exec(context.Background(), DROP)
	// if err != nil {
	// 	panic(err)
	// }
	// _, _, err = pgx.Pool().Exec(context.Background(), LATEST)
	// if err != nil {
	// 	panic(err)
	// }
}

func nameConv(isTable bool, name string) string {
	if isTable {
		if name == "gex_kline" {
			return "KLine"
		}
		return gen.ConvCamelCase(true, strings.TrimPrefix(name, "gex_"))
	}
	if name == "tid" || name == "uuid" || name == "i18n" || name == "qq" || name == "ip" {
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
	TypeField: map[string]map[string]string{
		"gex_user": {
			"favorites": "UserFavorites",
		},
		"gex_order": {
			"transaction": "OrderTransaction",
		},
	},
	FieldFilter: map[string]map[string]string{
		"gex_user": {
			gen.FieldsOrder:    "account,phone,update_time,create_time",
			gen.FieldsOptional: "role,name,account,phone,email,password,trade_pass,image,external,status",
			gen.FieldsFind:     "^password,trade_pass,favorites#all",
			gen.FieldsScan:     "^password,trade_pass,favorites#all",
		},
		"gex_order": {
			gen.FieldsOrder:    "update_time,create_time",
			gen.FieldsOptional: "tid,quantity,price,total_price,trigger_type,trigger_price,status",
			gen.FieldsScan:     "^transaction#all",
		},
	},
	CodeAddInit:  map[string]string{},
	CodeTestInit: map[string]string{},
	CodeSlice:    gen.CodeSlicePG,
	TableRetAdd:  map[string]string{},
	TableGenAdd: xsql.StringArray{
		"gex_balance",
		"gex_balance_history",
		"gex_kline",
		"gex_order",
		"gex_order_comm",
		"gex_withdraw",
		"gex_user",
	},
	TableNotValid: xsql.StringArray{},
	TableInclude:  xsql.StringArray{},
	TableExclude: xsql.StringArray{
		"gex_config",
		"gex_object",
		"gex_version_object",
		"gex_announce",
		"gex_file",
	},
	Queryer: pgx.Pool,
	TableQueryer: func(queryer interface{}, tableSQL, columnSQL, schema string) (tables []*gen.Table, err error) {
		tables, err = gen.Query(queryer, tableSQL, columnSQL, schema)
		if err != nil {
			return
		}
		for _, table := range tables {
			for _, column := range table.Columns {
				column.Comment = strings.ReplaceAll(column.Comment, "\n", "")
			}
		}
		return
	},
	TableSQL:   gen.TableSQLPG,
	ColumnSQL:  gen.ColumnSQLPG,
	Schema:     "public",
	TypeMap:    gen.TypeMapPG,
	NameConv:   nameConv,
	GetQueryer: "Pool",
	Out:        "../gexdb/",
	OutPackage: "gexdb",
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
