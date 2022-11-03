package main

import (
	"context"
	"os"
	"strings"

	"github.com/codingeasygo/crud/pgx"
	"github.com/gexservice/gexservice/base/basedb"
	"github.com/gexservice/gexservice/base/baseupgrade"
)

func main() {
	basedb.SYS = "exs"
	_, err := pgx.Bootstrap(os.Args[1])
	if err != nil {
		panic(err)
	}
	_, _, err = pgx.Exec(context.Background(), strings.ReplaceAll(baseupgrade.DROP, "_sys_", basedb.SYS+"_"))
	if err != nil {
		panic(err)
	}
}
