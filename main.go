package main

import (
	gosqltest "github.com/chalfel/sql-check-test/internal"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(gosqltest.Analyzer)
}
