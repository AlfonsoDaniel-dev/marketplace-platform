package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"shopperia/src/db"
)

var querys = []string{sqlMigrateUuidExtension, sqlMigrateUserTable, sqlMigrateUserAddressesTable, sqlMigrateCollectionsTable, sqlMigrateImagesTable}

var specialQuerys = []string{sqlAddConstraintForUserWithAddress}

type Migrator struct {
	Db *sql.DB
}

func NewMigrator(db *sql.DB) Migrator {
	return Migrator{
		Db: db,
	}
}

func (m *Migrator) Migrate() error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}

	for _, query := range querys {
		_, err := db.ExecQuery(tx, query)
		if err != nil {
			fmt.Println(query)
			return err
		}
	}

	tx.Commit()

	log.Println("Migraciones Realizadas")

	return nil
}
