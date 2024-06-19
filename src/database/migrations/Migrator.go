package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

var querys = []string{sqlMigrateUuidExtension, sqlMigrateUserTable, sqlMigrateUserAddressesTable}

// var specialQuerys = []string{sqlAddConstraintForUserWithAddress, sqlAddConstraintForAddressWithUser}

type Migrator struct {
	Db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		Db: db,
	}
}

func (m *Migrator) migrateTable(tx *sql.Tx, query string) error {
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) RunSpecialQuery(tx *sql.Tx, query string) error {
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrator) Migrate() error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}

	for _, query := range querys {
		err = m.migrateTable(tx, query)
		if err != nil {
			fmt.Println(query)
			tx.Rollback()
			return err
		}
	}

	defer tx.Commit()

	log.Println("Migraciones Realizadas")

	return nil
}
