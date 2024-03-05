package main

import (
	"os"

	sql_check "github.com/ledgera-io/go-sql-check/internal"
	"github.com/rs/zerolog"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting gosqlcheck")
	singlechecker.Main(sql_check.Analyzer)
}
