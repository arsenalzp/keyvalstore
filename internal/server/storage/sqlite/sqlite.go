// SQLite3 storage implementation.

package sqlite

import (
	"context"
	"database/sql"
	"os"

	entity "github.com/arsenalzp/keyvalstore/internal/server/storage/entity"

	_ "github.com/mattn/go-sqlite3"
)

// table schema for gokeyval storage
const schemaSQL string = `
CREATE TABLE IF NOT EXISTS "gokeyval" (
	"key"	TEXT(256) UNIQUE,
	"value"	TEXT(512) NOT NULL,
	UNIQUE(key)
);
`

// query for key searching operation
const searchSQL string = `
SELECT 
	value 
FROM gokeyval
WHERE key = ?;
`

// query for either insert key or update key operations
const insertSQL string = `
INSERT INTO
	gokeyval(key, value) 
VALUES
	(?, ?)
ON CONFLICT(key) DO UPDATE SET
	value=excluded.value;
`

// query for key delition operation
const deleteSQL string = `
DELETE FROM gokeyval
WHERE key = ?;
`

// query for selecting all rows in a table
const searchAllSQL string = `
SELECT
	*
FROM gokeyval;
`

// sqlite3 database structure
type Db struct {
	sql           *sql.DB   // sqlite3 connection
	searchStmt    *sql.Stmt // perapared statemnt for SELECT query
	insertStmt    *sql.Stmt // perapared statemnt for INSERT query
	deleteStmt    *sql.Stmt // perapared statemnt for DELETE query
	searchAllStmt *sql.Stmt // prepared statement for SELECL * query

	dbName string // database name
}

func (db *Db) Insert(ctx context.Context, k, v string) (bool, error) {
	res, err := db.insertStmt.ExecContext(ctx, k, v)
	if err != nil {
		return false, err
	}

	i, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if i == 0 {
		return false, nil
	}

	return true, nil
}

func (db *Db) Delete(ctx context.Context, k string) (bool, error) {
	res, err := db.deleteStmt.ExecContext(ctx, k)
	if err != nil {
		return false, err
	}

	i, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if i == 0 {
		return false, nil
	}

	return true, nil
}

func (db *Db) Search(ctx context.Context, k string) (string, error) {
	var value *string = new(string)

	row := db.searchStmt.QueryRowContext(ctx, k)
	//row := db.sql.QueryRowContext(ctx, "select value from gokeyval where key = ?", k)
	err := row.Scan(value)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if err := row.Err(); err != nil {
		return "", err
	}

	return *value, nil
}

func (db *Db) Import(ctx context.Context, data []entity.ImportData) (bool, error) {
	for _, item := range data {
		_, err := db.Insert(ctx, item.Key, item.Value)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (db *Db) Export(ctx context.Context) ([]entity.ExportData, error) {
	var exportRows []entity.ExportData

	rows, err := db.searchAllStmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var expRow entity.ExportData
		if err := rows.Scan(&expRow.Key, &expRow.Value); err != nil {
			return nil, err
		}
		exportRows = append(exportRows, expRow)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return exportRows, nil
}

func isDbExist(fname string) bool {
	if _, err := os.Stat(fname); err == os.ErrNotExist {
		return false
	}

	return true
}

func createDb(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		f.Close()

		return err
	}
	f.Close()

	return nil
}

func NewDb() (*Db, error) {
	var db *Db
	var fName string

	if v, ok := os.LookupEnv("SERVICE_DBNAME"); ok {
		fName = v
	} else {
		fName = "default.db"
	}

	if !isDbExist(fName) {
		if err := createDb(fName); err != nil {
			return nil, err
		}
	}

	sqlDb, err := sql.Open("sqlite3", fName)
	if err != nil {
		return nil, err
	}

	if _, err := sqlDb.Exec(schemaSQL); err != nil {
		return nil, err
	}

	searchStmt, err := sqlDb.Prepare(searchSQL)
	if err != nil {
		return nil, err
	}

	insertStmt, err := sqlDb.Prepare(insertSQL)
	if err != nil {
		return nil, err
	}

	deleteStmt, err := sqlDb.Prepare(deleteSQL)
	if err != nil {
		return nil, err
	}

	searchAllStmt, err := sqlDb.Prepare(searchAllSQL)
	if err != nil {
		return nil, err
	}

	db = &Db{
		sql:           sqlDb,
		dbName:        fName,
		searchStmt:    searchStmt,
		insertStmt:    insertStmt,
		deleteStmt:    deleteStmt,
		searchAllStmt: searchAllStmt,
	}

	return db, nil
}
