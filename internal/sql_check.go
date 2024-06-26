package sql_check

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ledgera-io/go-sql-check/config"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"

	"github.com/jmoiron/sqlx"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = `
Validate SQL queries by running them on a database.

Requires setting DATABASE_URL environment variable with the url of the database
on which the queries are going to be run.
The SQL queries in your code are required to start with "--sql" prefix to be 
recognized by sqlchk.
`

var DatabaseEnv string

var Analyzer = &analysis.Analyzer{
	Name:      "gosqlcheck",
	Doc:       doc,
	Run:       run,
	FactTypes: []analysis.Fact{},
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
}

type SqlCheckConfig struct {
	POSTGRES_DNS string `env:"POSTGRES_DNS"`
	DATABASE_URL string `env:"DATABASE_URL"`
	DRIVER       string `env:"DRIVER"`
}

func init() {
	flag.StringVar(&DatabaseEnv, "database-env", "", "Database env")
}

func getDatabaseUrl(logger *zerolog.Logger) (string, error) {
	var databaseUrl string

	if DatabaseEnv != "" {
		if err := config.LoadDotEnv(); err != nil {
			return "", err
		}

		databaseUrl := os.Getenv(DatabaseEnv)

		if databaseUrl == "" {
			logger.Error().Msg("database url is empty")
			return "", fmt.Errorf("database url is empty")
		}

		return databaseUrl, nil
	}

	cfg := SqlCheckConfig{}

	err := config.LoadConfigFromEnv(&cfg)

	if err != nil {
		return "", err
	}

	if cfg.POSTGRES_DNS != "" {
		databaseUrl = cfg.POSTGRES_DNS
	} else {
		databaseUrl = cfg.DATABASE_URL
	}

	return databaseUrl, nil
}

func run(pass *analysis.Pass) (any, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	databaseUrl, err := getDatabaseUrl(&logger)

	if err != nil {
		return nil, err
	}

	if databaseUrl == "" {
		logger.Error().Msg("database url is empty")
		return nil, fmt.Errorf("database url is empty")
	}

	if !strings.HasSuffix(databaseUrl, "?sslmode=disable") {
		databaseUrl = databaseUrl + "?sslmode=disable"
	}

	driver := "postgres"

	db, err := sqlx.Connect(driver, databaseUrl)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	filterNodes := []ast.Node{
		(*ast.BasicLit)(nil),
	}

	inspect.Preorder(filterNodes, func(n ast.Node) {
		node := n.(*ast.BasicLit)
		if node.Kind != token.STRING {
			return
		}

		var str string
		if strings.HasPrefix(node.Value, "`") {
			str = strings.Trim(node.Value, "`")
		} else {
			str = strings.TrimRight(node.Value, "\"")
		}

		if !strings.HasPrefix(str, "--sql") {
			return
		}

		sqlStr := str
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			pass.Reportf(node.Pos(), "%s", err.Error())
			return
		}
		stmt.Close()
	})

	return nil, nil
}
