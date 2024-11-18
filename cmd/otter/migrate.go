package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/martinmunillas/otter/log"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateNewCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration utils",
	Long:  ``,
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(rune(str[i-1])) || unicode.IsDigit(rune(str[i-1]))) {
				if i != 0 || result[len(result)-1] != '_' {
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
		} else if unicode.IsSpace(r) || r == '-' || r == '.' || r == ',' {
			if i == 0 || result[len(result)-1] == '_' {
				continue
			}
			result = append(result, '_')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func toPascalCase(str string) string {

	words := strings.FieldsFunc(str, func(r rune) bool {
		return r == '_' || r == ' ' || r == '-'
	})
	for i, word := range words {
		words[i] = cases.Title(language.Und).String(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

var migrateNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Generates a new migration file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := readConfig()

		logger := log.NewLogger(true)

		if len(args) < 1 {
			logger.Error("Missing migration description, should run in the format `otter migrate new {migrateDescription}`")
			os.Exit(1)
		}
		description := ""
		for _, arg := range args {
			description += " "
			description += arg
		}
		t := time.Now().Format("20060102150405")

		err := os.MkdirAll(config.MigrationsDir, 0755)
		if err != nil {
			fatal(logger, err)
		}

		f, err := os.Create(fmt.Sprintf("%s/%s_%s.go", config.MigrationsDir, t, toSnakeCase(description)))
		if err != nil {
			fatal(logger, err)
		}
		defer f.Close()
		mainFile := makeMigrationFile(t, description)

		_, err = f.WriteString(mainFile)
		if err != nil {
			fatal(logger, err)
		}

	},
}

func makeMigrationFile(t string, description string) string {
	snake := toSnakeCase(description)
	pascal := toPascalCase(description)
	return fmt.Sprintf(`package migrations

import (
	"context"
	"database/sql"

	"github.com/martinmunillas/otter/migrate"
)

func init() {
	migrate.AddMigration("%s_%s", up%s, down%s)
}

func up%s(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "SELECT 1;")

	return err
}

func down%s(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "SELECT 1")
	return err
}

`, t, snake, pascal, pascal, pascal, pascal)
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run missing migrations",
	Long:  ``,
	Run: func(_ *cobra.Command, args []string) {
		config := readConfig()

		logger := log.NewLogger(true)

		if config.DbDriver == "" {
			logger.Error("No db driver set, make sure you have one set on your otter.json")
			os.Exit(1)

		}
		err := os.RemoveAll("./.otter")
		if err != nil {
			fatal(logger, err)
		}
		err = os.Mkdir("./.otter", 0755)
		if err != nil {
			fatal(logger, err)
		}
		err = os.Mkdir("./.otter/migrate", 0755)
		if err != nil {
			fatal(logger, err)
		}

		mainFile := makeTmpMigrationRunnerFile()

		f, err := os.Create("./.otter/migrate/main.go")
		if err != nil {
			fatal(logger, err)
		}
		defer f.Close()

		_, err = f.WriteString(mainFile)
		if err != nil {
			fatal(logger, err)
		}
		cmd := createDefaultCommand("go", "run", "./.otter/migrate/main.go")

		err = cmd.Run()
		if err != nil {
			fatal(logger, err)
		}

		f.Close()
		os.RemoveAll("./.otter/migrate/")
	},
}

func makeTmpMigrationRunnerFile() string {
	return fmt.Sprintf(`
package main

import (
	"database/sql"
	"fmt"
	"os"
	"log/slog"

	_ "%s"

	"github.com/martinmunillas/otter/migrate"
	"github.com/martinmunillas/otter/env"
	
	_ "%s/%s"
)

func main() {
	dbUser := env.RequiredStringEnvVar("DB_USER")
	dbName := env.RequiredStringEnvVar("DB_NAME")
	dbPassword := env.RequiredStringEnvVar("DB_PASSWORD")
	dbHost := env.OptionalStringEnvVar("DB_HOST", "")
	dbPort := env.OptionalStringEnvVar("DB_PORT", "")
	connStr := fmt.Sprintf("user=%%s dbname=%%s password=%%s sslmode=disable", dbUser, dbName, dbPassword)

	if dbHost != "" {
		connStr += " host=" + dbHost
	}

	if dbPort != "" {
		connStr += " port=" + dbPort
	}

	db, err := sql.Open("%s", connStr)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = migrate.RunAll(db, slog.Default())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
`, supportedDrivers[config.DbDriver], config.moduleName, config.MigrationsDir, config.DbDriver)
}
