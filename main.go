package main

import (
	sql_check "github.com/ledgera-io/go-sql-check/internal"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(sql_check.Analyzer)
}
