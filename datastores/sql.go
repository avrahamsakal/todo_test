package datastores

import (
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/jordan-borges-lark/todo_test/models"
)

const migrationDirName = "./datastores/migrations"

var db *sqlx.DB

func GetSqlDB(driverName, dataSourceName string) (*sqlx.DB, error) {
	if db != nil {
		return db, nil
	}

	db, err := sqlx.Open(driverName, dataSourceName)

	return db, err
}

// RunSqlMigrations runs all migrations higher version
// than what the DB reports we're on.
func RunSqlMigrations(db *sqlx.DB) error {
	// Ask DB the current migration version
	currentVersionStr, err := models.GetMetadataByKey(db, "migration_version")
	if err != nil {
		return err
	}
	currentVersion, _ := strconv.Atoi(currentVersionStr)

	// Loop over migration SQL files in alphanumeric order (hence big-endian dated filenames)
	entries, err := os.ReadDir(migrationDirName)
	if err != nil { return err }
	numIgnored := 0
	for i, entry := range entries {
		if entry.IsDir() {
			numIgnored++
			continue
		}
		// If we haven't yet run the migration in the current SQL file...
		if version := i + 1 - numIgnored; version > currentVersion {
			if err = RunSqlMigration(db, entry.Name()); err != nil {
				return err
			}
			currentVersion = version // Increment currentVersion with each migration we run
		}
	}

	return nil

	// @TODO: If using gorm, we could run gorm.RunMigrations()
	// and it would autogenerate/execute them based on diffs
	// of the model definitions vs. what's in the database
}


func RunSqlMigration(db *sqlx.DB, entryName string) error {
	byt, err := os.ReadFile(migrationDirName + "/" + entryName)
	if err != nil {
		return err
	}

	// @TODO: actually check the sql result that comes back,
	// rather than throwing it away
	if _, err = db.Exec(string(byt)); err != nil {
		return err
	}

	return nil
}
