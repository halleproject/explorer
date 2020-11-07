package db

import (
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func (db *Database) CreateTables() error {
	for _, model := range []interface{}{&schema.TwoAuthForDCI{}, &schema.AppVersion{}} {
		// Disable pluralization
		orm.SetTableNameInflector(func(s string) string {
			return s
		})

		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Varchar:     20000, // replaces data type from `text` to `varchar(n)`
		})

		if err != nil {
			return err
		}
	}
	return nil
}

// InsertExportedData inserts exported block, transaction data
// RunInTransaction runs a function in a transaction.
// if function returns an error transaction is rollbacked, otherwise transaction is committed.
func (db *Database) InsertTwoAuthForDCI(ta *schema.TwoAuthForDCI) error {

	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if ta != nil {
			err := tx.Insert(ta)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateTwoAuthForDCI(ta *schema.TwoAuthForDCI) error {

	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if ta != nil {
			err := tx.Update(ta)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertAppVerion(av *schema.AppVersion) error {

	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if av != nil {
			err := tx.Insert(av)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateAppVerion(av *schema.AppVersion) error {

	err := db.RunInTransaction(func(tx *pg.Tx) error {
		if av != nil {
			err := tx.Update(av)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
