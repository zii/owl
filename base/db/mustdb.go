package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func raise(err error) {
	if err != nil {
		panic(err)
	}
}

type MustDB struct {
	*sqlx.DB
}

func (db *MustDB) Exec(query string, args ...interface{}) *Result {
	r, err := db.DB.Exec(query, args...)
	raise(err)
	return FromSqlResult(r)
}

func (db *MustDB) Get(dest interface{}, query string, args ...interface{}) {
	err := db.DB.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return
	}
	raise(err)
}

func (db *MustDB) Query(query string, args ...interface{}) *MustRows {
	r, err := db.DB.Query(query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	raise(err)
	return &MustRows{Rows: r}
}

func (db *MustDB) Queryx(query string, args ...interface{}) *MustXRows {
	r, err := db.DB.Queryx(query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	raise(err)
	return &MustXRows{Rows: r}
}

func (db *MustDB) QueryRow(query string, args ...interface{}) *MustRow {
	r := db.DB.QueryRow(query, args...)
	return &MustRow{Row: r}
}

func (db *MustDB) QueryRowx(query string, args ...interface{}) *MustRowx {
	r := db.DB.QueryRowx(query, args...)
	return &MustRowx{Row: r}
}

func (db *MustDB) Begin() *MustTx {
	tx, err := db.DB.Beginx()
	raise(err)
	return &MustTx{Tx: tx}
}

type Result struct {
	affectedRows int64
	insertId     int64
}

func (r *Result) LastInsertId() int64 {
	return r.insertId
}

func (r *Result) RowsAffected() int64 {
	return r.affectedRows
}

func (r *Result) OK() bool {
	return r.insertId != 0 || r.affectedRows > 0
}

func FromSqlResult(r sql.Result) *Result {
	insertId, err := r.LastInsertId()
	raise(err)
	affectedRows, err := r.RowsAffected()
	raise(err)
	out := &Result{
		insertId:     insertId,
		affectedRows: affectedRows,
	}
	return out
}

type MustRow struct {
	*sql.Row
}

func (rs *MustRow) Scan(dest ...interface{}) bool {
	err := rs.Row.Scan(dest...)
	if err == sql.ErrNoRows {
		return false
	}
	raise(err)
	return true
}

type MustRowx struct {
	*sqlx.Row
}

func (r *MustRowx) Scan(dest ...interface{}) bool {
	err := r.Row.Scan(dest...)
	if err == sql.ErrNoRows {
		return false
	}
	raise(err)
	return true
}

func (r *MustRowx) StructScan(dest interface{}) bool {
	err := r.Row.StructScan(dest)
	if err == sql.ErrNoRows {
		return false
	}
	raise(err)
	return true
}

type MustRows struct {
	*sql.Rows
}

func (rs *MustRows) Scan(dest ...interface{}) {
	err := rs.Rows.Scan(dest...)
	raise(err)
}

type MustXRows struct {
	*sqlx.Rows
}

func (rs *MustXRows) Scan(dest ...interface{}) {
	err := rs.Rows.Scan(dest...)
	raise(err)
}

type MustTx struct {
	*sqlx.Tx
}

func (tx *MustTx) Exec(query string, args ...interface{}) *Result {
	r, err := tx.Tx.ExecContext(context.Background(), query, args...)
	raise(err)
	return FromSqlResult(r)
}

func (tx *MustTx) Get(dest interface{}, query string, args ...interface{}) bool {
	err := tx.Tx.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return false
	}
	raise(err)
	return true
}

func (tx *MustTx) Query(query string, args ...interface{}) *MustRows {
	r, err := tx.Tx.Query(query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	raise(err)
	return &MustRows{Rows: r}
}

func (tx *MustTx) Queryx(query string, args ...interface{}) *MustXRows {
	r, err := tx.Tx.Queryx(query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	raise(err)
	return &MustXRows{Rows: r}
}

func (tx *MustTx) QueryRow(query string, args ...interface{}) *MustRow {
	r := tx.Tx.QueryRow(query, args...)
	return &MustRow{Row: r}
}

func (tx *MustTx) QueryRowx(query string, args ...interface{}) *MustRowx {
	r := tx.Tx.QueryRowx(query, args...)
	return &MustRowx{Row: r}
}

func (tx *MustTx) Rollback() {
	err := tx.Tx.Rollback()
	if err == sql.ErrTxDone {
		return
	}
	raise(err)
}

func (tx *MustTx) Commit() {
	err := tx.Tx.Commit()
	if err == sql.ErrTxDone {
		return
	}
	raise(err)
}
