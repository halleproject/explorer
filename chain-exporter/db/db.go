package db

import (
	"github.com/chain-exporter/config"
	"github.com/chain-exporter/schema"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Database implements a wrapper of golang ORM with focus on PostgreSQL
type Database struct {
	*pg.DB
}

// Connect opens a database connections with the given database connection info from config.
// It returns a database connection handle or an error if the connection fails.
func Connect(cfg config.DBConfig) *Database {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Table,
	})

	return &Database{db}
}

// Ping returns a database connection handle or an error if the connection fails.
func (db *Database) Ping() error {
	_, err := db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

// CreateTables creates database tables using object relational mapping (ORM)
func (db *Database) CreateTables() error {
	for _, model := range []interface{}{(*schema.Block)(nil), (*schema.PreCommit)(nil), (*schema.Transaction)(nil),
		(*schema.Validator)(nil), (*schema.StatAssetInfoList1H)(nil), (*schema.StatAssetInfoList24H)(nil), (*schema.TwoAuth)(nil)} {
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

	// RunInTransaction creates indexes to reduce the cost of lookup queries in case of server traffic jams.
	// If function returns an error transaction is rollbacked, otherwise transaction is committed.
	err := db.RunInTransaction(func(tx *pg.Tx) error {
		_, err := db.Model(schema.Block{}).
			Exec(`CREATE INDEX block_height_idx ON block USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.PreCommit{}).
			Exec(`CREATE INDEX pre_commit_height_idx ON pre_commit USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.PreCommit{}).
			Exec(`CREATE INDEX pre_commit_validator_address_idx ON pre_commit USING btree(validator_address);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Transaction{}).
			Exec(`CREATE INDEX transaction_height_idx ON transaction USING btree(height);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Transaction{}).
			Exec(`CREATE INDEX transaction_tx_hash_idx ON transaction USING btree(tx_hash);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Transaction{}).
			Exec(`CREATE INDEX transaction_from_address_idx ON transaction USING btree(from_address);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Transaction{}).
			Exec(`CREATE INDEX transaction_to_address_idx ON transaction USING btree(to_address);`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Transaction{}).
			Exec(`CREATE INDEX transaction_messages_symbol_idx ON validator USING btree((messages->0->'value'->>'symbol'));`)
		if err != nil {
			return err
		}
		_, err = db.Model(schema.Validator{}).
			Exec(`CREATE INDEX validator_consensus_address_idx ON validator USING btree(consensus_address);`)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
