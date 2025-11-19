package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"path/filepath"
	"slices"

	"github.com/pressly/goose"
)

type Migration interface {
	CreateFile(path string, migrationName string, migrationType string) error
	Apply(path string, version int64, fake bool) error
	Rollback(path string, version int64) error
}

type MigrationConfig struct {
	db *sql.DB
}

func NewMigration(db *sql.DB) Migration {
	return &MigrationConfig{
		db: db,
	}
}

func (m *MigrationConfig) CreateFile(path string, migrationName string, migrationType string) error {
	migrationDir, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := goose.Create(nil, migrationDir, migrationName, migrationType); err != nil {
		return err
	}

	log.Printf("Migration file created successfully: %s.%s", migrationName, migrationType)
	return nil
}

func (m *MigrationConfig) Apply(path string, version int64, fake bool) error {
	if m.db == nil {
		return errors.New("connection to database is not established")
	}
	db := m.db
	defer db.Close()

	migrationDir, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if fake {
		return fakeApply(db, path, version)
	}
	if version == 0 {
		return goose.Up(db, migrationDir)
	}
	return goose.UpTo(db, migrationDir, version)
}

func (m *MigrationConfig) Rollback(path string, version int64) error {
	if m.db == nil {
		return errors.New("connection to database is not established")
	}
	db := m.db
	defer db.Close()

	migrationDir, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if version == 0 {
		return goose.Down(db, migrationDir)
	}
	return goose.DownTo(db, path, version)
}

func fakeApply(db *sql.DB, path string, version int64) error {
	// Retrieves the current version for this DB.
	// Creates and initializes the DB version table if it doesn't exist.
	if _, err := goose.EnsureDBVersion(db); err != nil {
		return err
	}

	ctx := context.Background()
	allAppliedMigrations, err := listAllDBVersions(ctx, db)
	if err != nil {
		return err
	}

	// If the fake apply has a version
	if version != 0 {
		if !slices.Contains(allAppliedMigrations, version) {
			_, err = db.Exec(`
			INSERT INTO public.goose_db_version (version_id, is_applied) 
			VALUES ($1, true)
			`, version)
			if err != nil {
				return err
			}

			log.Printf("Migration %d applied successfully!", version)
			return nil
		}
		log.Printf("Migration %d is already applied.", version)
		return nil
	}

	// If the fake apply does not have a version
	foundMigrations, err := goose.CollectMigrations(path, 0, goose.MaxVersion)
	if err != nil {
		return err
	}

	var toApply []int64
	for _, migration := range foundMigrations {
		if !slices.Contains(allAppliedMigrations, migration.Version) {
			toApply = append(toApply, migration.Version)
		}
	}

	if len(toApply) == 0 {
		log.Println("No unapplied migration found.")
		return nil
	}

	for _, migrationVer := range toApply {
		_, err := db.Exec(`
		INSERT INTO public.goose_db_version (version_id, is_applied) 
		VALUES ($1, true)
		`, migrationVer)
		if err != nil {
			return err
		}
	}

	log.Println("Migrations fake applied successfully!")
	return nil
}

// List all version_id values in the goose_db_version table
func listAllDBVersions(ctx context.Context, db *sql.DB) ([]int64, error) {
	var versionIDs []int64

	query := `SELECT version_id FROM goose_db_version`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var versionID sql.NullInt64
		if err := rows.Scan(&versionID); err != nil {
			return nil, err
		}

		if versionID.Valid {
			versionIDs = append(versionIDs, versionID.Int64)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versionIDs, nil
}
