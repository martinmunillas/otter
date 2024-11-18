package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

var migrations []migration

type MigrationExec = func(ctx context.Context, tx *sql.Tx) error

type migration struct {
	id   string
	up   MigrationExec
	down MigrationExec
}

func AddMigration(id string, up MigrationExec, down MigrationExec) {
	migrations = append(migrations, migration{
		id:   id,
		up:   up,
		down: down,
	})
}

func sortMigrations() {
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].id < migrations[j].id
	})
}

func RunAll(conn *sql.DB, logger *slog.Logger) error {
	sortMigrations()
	err := ensureMigrationsTable(conn, logger)
	if err != nil {
		return fmt.Errorf("error ensuring migrations table: %w", err)
	}
	records, err := getMigrationRecords(conn)
	if err != nil {
		return fmt.Errorf("error getting migrations: %w", err)
	}
	err = ensureRanMigrationsAreValid(records)
	if err != nil {
		return fmt.Errorf("error ensuring ran migrations are valid: %w", err)
	}

	err = migrateMigrations(logger, conn)
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	return nil
}

func EnsureAllMigrationsRanAndAreValid(conn *sql.DB, logger *slog.Logger) error {
	records, err := getMigrationRecords(conn)
	if err != nil {
		return err
	}

	if err := ensureRanMigrationsAreValid(records); err != nil {
		return err
	}

	if len(records) != len(migrations) {
		return fmt.Errorf("there are %d migrations left to run", len(migrations)-len(records))
	}

	logger.Info("Migrations are up to date")
	return nil
}

type migrationRecord struct {
	ID         string    `sql:"id"`
	MigratedAt time.Time `sql:"migrated_at"`
}

func getMigrationRecords(conn *sql.DB) ([]migrationRecord, error) {
	var migrations []migrationRecord
	rows, err := conn.Query("SELECT id, migrated_at FROM otter_migrations ORDER BY id")
	for rows.Next() {
		var migration migrationRecord
		err = rows.Scan(&migration.ID, &migration.MigratedAt)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	if err != nil {
		return nil, err
	}
	return migrations, nil
}

func ensureRanMigrationsAreValid(records []migrationRecord) error {
	for i, record := range records {
		migration := migrations[i]
		if migration.id != record.ID {
			return fmt.Errorf("migrations don't match, `%s` is different to `%s` though they claim to be the same migration, needs manual solving. Migrations ordering cannot be changed once migrated", migration.id, records[i].ID)
		}
	}
	return nil
}

func migrateMigrations(logger *slog.Logger, conn *sql.DB) error {
	records, err := getMigrationRecords(conn)
	if err != nil {
		return err
	}

	if len(records) == len(migrations) {
		logger.Info("All migrations already ran")
		return nil
	} else {
		logger.Info(fmt.Sprintf("%d/%d migrations already ran", len(records), len(migrations)))
	}
	logger.Info("Running migrations")

	for i, migration := range migrations {
		// Skip migrations that have already ran
		if i < len(records) {
			continue
		}

		logger.Info(fmt.Sprintf("Running migration %s", migration.id))

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		err = migration.up(ctx, tx)
		if err != nil {
			err = errors.Join(err, tx.Rollback())
			logger.Error(err.Error())
			return err
		}

		_, err = tx.Exec("INSERT INTO otter_migrations (id, migrated_at) VALUES ($1, $2)", migration.id, time.Now())
		if err != nil {
			err = errors.Join(err, tx.Rollback())
			logger.Error(err.Error())
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Migration %s ran successfully", migration.id))
	}

	logger.Info("All migrations ran successfully")
	return nil
}

func ensureMigrationsTable(conn *sql.DB, logger *slog.Logger) error {
	var exists bool
	err := conn.QueryRow(
		`SELECT EXISTS (
			SELECT 1
			FROM pg_tables
			WHERE schemaname = 'public'
			AND tablename = 'otter_migrations'
		);`,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking migrations table: %w", err)
	}

	if !exists {
		logger.Info("Creating migrations table")
		_, err := conn.Exec(
			`CREATE TABLE "otter_migrations" (
				"id" VARCHAR PRIMARY KEY NOT NULL,
				"migrated_at" TIMESTAMPTZ NOT NULL
			);
			`,
		)
		if err != nil {
			return fmt.Errorf("error creating migrations table: %w", err)
		}
		logger.Info("Migrations table created")
	}

	return nil
}
