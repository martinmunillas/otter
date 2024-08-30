package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration utils",
	Long:  ``,
	Run: func(_ *cobra.Command, args []string) {
		config := readConfig()
		if config.DbDriver == "" {
			panic(fmt.Errorf("no db driver set, make sure you have one set on your otter.json"))
		}
		err := os.RemoveAll("./.otter")
		if err != nil {
			panic(err)
		}
		err = os.Mkdir("./.otter", 0755)
		if err != nil {
			panic(err)
		}
		err = os.Mkdir("./.otter/migrate", 0755)
		if err != nil {
			panic(err)
		}

		mainFile := makeTmpMigrationFile()

		f, err := os.Create("./.otter/migrate/main.go")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = f.WriteString(mainFile)
		if err != nil {
			panic(err)
		}
		cmd := createDefaultCommand("go", "run", "./.otter/migrate/main.go")

		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	},
}

func makeTmpMigrationFile() string {
	return fmt.Sprintf(`
package main

import (
	"database/sql"
	"fmt"

	_ "%s"

	"github.com/martinmunillas/otter/migrate"
	"github.com/martinmunillas/otter/log"
	"github.com/martinmunillas/otter/env"
	
	_ "%s/%s"
)

func main() {
	dbUser := env.RequiredStringEnvVar("DB_USER")
	dbName := env.RequiredStringEnvVar("DB_NAME")
	dbPassword := env.RequiredStringEnvVar("DB_PASSWORD")
	connStr := fmt.Sprintf("user=%%s dbname=%%s password=%%s sslmode=disable", dbUser, dbName, dbPassword)
	db, err := sql.Open("%s", connStr)
	if err != nil {
		panic(err)
	}

	logger := log.NewLogger(false)
	err = migrate.RunAll(db, logger)
	if err != nil {
		panic(err)
	}
}
`, supportedDrivers[config.DbDriver], config.moduleName, config.MigrationsDir, config.DbDriver)
}
