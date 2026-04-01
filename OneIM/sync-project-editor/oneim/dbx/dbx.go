package dbx

import (
	//    "database/sql"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/jimsmart/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"

	"pso.oneidentity.com/sped/oneim"
)

const MAX_RESULTS = 1000

type DBConfig struct {
	UserName       string
	Password       string
	HostName       string
	Port           int
	DatabaseName   string
	MaxConnections int
}

func CreateCtxFromStruct(config *DBConfig) (*sqlx.DB, error) {
	return CreateCtx(config.UserName, config.Password, config.HostName, config.Port, config.DatabaseName, config.MaxConnections)
}

func CreateCtx(
	username string, password string,
	host string, port int, database string,
	maxConnections int) (*sqlx.DB, error) {

	query := url.Values{}
	query.Add("database", database)

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%d", host, port),
		RawQuery: query.Encode(),
	}
	db, err := sqlx.Open("sqlserver", u.String())
	if err != nil {
		log.Println(err)
		return db, err
	}
	if maxConnections > 0 {
		db.SetMaxOpenConns(maxConnections)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return db, err
	}

	return db, nil
}

// ===========================================

type DBID struct {
	ID string `db:"id"`
}

func GetNewId(db *sqlx.DB) (string, error) {
	var id string
	err := db.Get(&id, "SELECT id=REPLACE(CONVERT(varchar(255), newid()), '-', '')")
	if err != nil {
		return "", err
	}
	return "CCC-" + id, nil
}

func GetTableValue(db *sqlx.DB, table string, column string, whereClause string) (string, error) {
	var val string
	stmt := fmt.Sprintf(`select %s from %s where %s`, column, table, whereClause)
	err := db.Get(&val, stmt)
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetTableCount(db *sqlx.DB, table string, clause string) (int, error) {
	count := -1
	q := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", table, clause)
	row := db.QueryRow(q)
	row.Scan(&count)
	return count, nil
}

func getTableRows(db *sqlx.DB, table string, clause string) (*sqlx.Rows, error) {
	//fmt.Println(fmt.Sprintf("SELECT * FROM %s WHERE %s", table, clause))
	return db.Queryx(fmt.Sprintf("SELECT * FROM %s WHERE %s", table, clause))
}

// use maxRows=-1 to fetch all
func GetBufferedTableData(db *sqlx.DB, table string, clause string, maxRows int) ([]map[string]interface{}, error) {
	bufSize := maxRows
	if bufSize <= 0 {
		bufSize = MAX_RESULTS
	}

	// array of results
	var r = make([]map[string]interface{}, 0, bufSize)
	count := 0

	rows, err := getTableRows(db, table, clause)
	if err != nil {
		log.Println(err)
		return r, err
	}
	defer rows.Close()

	for (count < bufSize) && rows.Next() {
		buf := make(map[string]interface{})
		err = rows.MapScan(buf)
		if err != nil {
			log.Println(err)
			return r, err
		}

		r = append(r, buf)
		count++
	}

	return r, nil
}

func GetSingletonTableData(db *sqlx.DB, table string, wc string) (map[string]interface{}, error) {

	buffs, err := GetBufferedTableData(db, table, wc, 1)
	if err != nil {
		fmt.Println("failed to fetch " + table + " records")
		return nil, err
	}

	if len(buffs) == 1 {
		return buffs[0], nil
	}

	return nil, errors.New("expected singleton " + table + " for " + wc)
}

func GetForeignSingleton(db *sqlx.DB, table string, columnName string, columnValue interface{}) (map[string]interface{}, error) {

	wc := fmt.Sprintf(`%s = '%v'`, columnName, columnValue)
	return GetSingletonTableData(db, table, wc)
}

// returns table name, table data, error
func GetSingletonTableDataByKey(db *sqlx.DB, objectKey interface{}) (string, map[string]interface{}, error) {

	keyTable, keyIDs := oneim.GetKeyParts(fmt.Sprintf("%v", objectKey))
	if len(keyTable) == 0 || len(keyIDs) != 1 {
		return "", nil, errors.New("failed to parse object key")
	}

	tdata, err := GetForeignSingleton(db, keyTable, "UID_"+keyTable, keyIDs[0])
	return keyTable, tdata, err
}

// return foreign key where clause
func GetFKWC(row *map[string]interface{}, attrName string) (string, error) {
	return GetCRWC(attrName, row, attrName)
}

// return CR relation where clause
func GetCRWC(foreignAttrName string, row *map[string]interface{}, attrName string) (string, error) {
	rval, ok := (*row)[attrName]
	if ok && rval != nil {
		return fmt.Sprintf("%s = '%v'", foreignAttrName, rval), nil
	}
	return "", errors.New("unable to parse foreign key value")
}

func GetStructData[T any](db *sqlx.DB, table string, clause string) ([]T, error) {

	tt := make([]T, 0)

	udb := db.Unsafe()
	rows, err := udb.Queryx(fmt.Sprintf("select * from %s where %s", table, clause))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		buf := new(T)
		err = rows.StructScan(buf)
		if err != nil {
			return nil, err
		}
		tt = append(tt, *buf)
	}

	return tt, nil
}

func GetStructSingleton[T any](db *sqlx.DB, id string) (T, error) {
	var t T
	return GetStructSingletonByWC[T](db, fmt.Sprintf(`UID_%s='%s'`, reflect.TypeOf(t).Name(), id))
}

func GetStructSingletonByWC[T any](db *sqlx.DB, wc string) (T, error) {

	var t T
	table := reflect.TypeOf(t).Name()

	ts, err := GetStructData[T](db, table, wc)
	if err != nil {
		return t, err
	}
	if len(ts) == 0 {
		return t, errors.New("object not found: " + table)
	} else if len(ts) > 1 {
		return t, errors.New(fmt.Sprintf("Too many rows returned for %s", table))
	}

	return ts[0], nil
}

func GenerateInsertStmt(s interface{}) (string, error) {

	table := reflect.TypeOf(s).Name()

	names, err := oneim.GetNonNullFieldNames(s)
	if err != nil {
		return "", nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`INSERT INTO %s (`, table))

	// column names
	for i, n := range names {
		b.WriteString(n)
		if i < len(names)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString(") VALUES (")

	// values
	for i, n := range names {
		b.WriteString(":")
		b.WriteString(n)
		if i < len(names)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString(")")

	return b.String(), nil
}

func GenerateUpdateStmt(s interface{}, id string) (string, error) {

	table := reflect.TypeOf(s).Name()

	names, err := oneim.GetNonNullFieldNames(s)
	if err != nil {
		return "", nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`UPDATE %s SET `, table))

	for i, n := range names {
		b.WriteString(fmt.Sprintf(`%s = :%s`, n, n))
		if i < len(names)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString(fmt.Sprintf(` WHERE UID_%s='%s'`, table, id))

	return b.String(), nil
}

func StructExists[T any](db *sqlx.DB, id string) (bool, error) {

	var t_empty T
	t, err := GetStructSingleton[T](db, id)
	if err != nil {
		return false, err
	}

	t_j, _ := json.Marshal(t)
	t_empty_j, _ := json.Marshal(t_empty)

	return len(t_j) != len(t_empty_j), nil
}

func GetTableColumns(db *sqlx.DB, table string) ([][2]string, error) {

	tcols, err := schema.ColumnTypes(db.DB, "", table)
	if err != nil {
		return nil, err
	}

	cols := make([][2]string, len(tcols))
	for i := range tcols {
		cols[i][0] = tcols[i].Name()
		cols[i][1] = tcols[i].DatabaseTypeName()
	}

	return cols, nil
}

func WaitForDBResult(db *sqlx.DB, ctx context.Context, dbFunc func(db *sqlx.DB) (bool, error), interval time.Duration) error {

	timer := time.NewTimer(interval)
	defer timer.Stop()

	// Start polling until the context is done
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout waiting for SQL task")
		case <-timer.C:
			result, err := dbFunc(db)
			if err != nil {
				return nil
			}
			if result {
				return nil
			} else {
			}
			timer.Reset(interval)
		}
	}

	return nil
}
