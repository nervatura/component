package demo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	ut "github.com/nervatura/component/pkg/util"
)

// Saving component state in a session json file.
func (app *App) SaveFileSession(fileName string, data any) (err error) {
	if _, err = app.osStat(sessionPath); errors.Is(err, os.ErrNotExist) {
		if err = app.osMkdir(sessionPath, os.ModePerm); err != nil {
			return err
		}
	}
	filePath := fmt.Sprintf(`%s/%s.json`, sessionPath, fileName)
	sessionFile, err := app.osCreate(filePath)
	if err == nil {
		bin, err := json.Marshal(data)
		if err == nil {
			sessionFile.Write(bin)
		}
	}
	defer sessionFile.Close()
	return err
}

// Loading the state of a component from a session json file.
func (app *App) LoadFileSession(fileName string, data any) (err error) {
	filePath := fmt.Sprintf(`%s/%s.json`, sessionPath, fileName)
	sessionFile, err := app.osReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(sessionFile, &data)
	}
	return err
}

func (app *App) getSessionTableSql(driverName, query string) (sqlString string) {
	if query == "open" {
		if driverName == "sqlite3" {
			return fmt.Sprintf("select name from sqlite_master where name = '%s' ", sessionTable)
		}
		// postgres, mysql, mssql
		return fmt.Sprintf("select table_name from information_schema.tables where table_name = '%s' ", sessionTable)
	}
	if driverName == "mysql" {
		return fmt.Sprintf(
			"CREATE TABLE %s ( id VARCHAR(255) NOT NULL, value JSON, stamp VARCHAR(255), PRIMARY KEY (id) );",
			sessionTable)
	}
	jsonType := ut.SM{
		"sqlite3": "JSON", "postgres": "JSONB", "mysql": "JSON", "mssql": "NVARCHAR(MAX)", "sqltest": "JSON",
	}
	return fmt.Sprintf(
		"CREATE TABLE %s ( id VARCHAR(255) NOT NULL PRIMARY KEY, value %s, stamp VARCHAR(255) );",
		sessionTable, jsonType[driverName])
}

func (app *App) checkSessionTable() (db *sql.DB, err error) {
	if db, err = sql.Open(app.driverName, app.dataSource); err != nil {
		return db, err
	}

	var found bool
	var rows *sql.Rows
	sqlString := app.getSessionTableSql(app.driverName, "open")
	if rows, err = db.Query(sqlString); err != nil {
		return db, err
	}
	defer rows.Close()
	for rows.Next() {
		found = true
	}
	if !found {
		sqlString = app.getSessionTableSql(app.driverName, "create")
		_, err = db.Exec(sqlString)
	}

	return db, err
}

func (app *App) getSessionValue(db *sql.DB, sessionID string) (value string, err error) {
	sqlString := fmt.Sprintf("SELECT value FROM %s WHERE id='%s'", sessionTable, sessionID)
	row := db.QueryRow(sqlString)
	if err = row.Scan(&value); err == sql.ErrNoRows {
		value = ""
	}
	return value, err
}

// Saving component state in a database.
func (app *App) SaveDbSession(sessionID string, data any) (err error) {
	var db *sql.DB
	if db, err = app.checkSessionTable(); err == nil {
		var bin []byte
		if bin, err = json.Marshal(data); err == nil {
			var sqlString string = fmt.Sprintf(
				"INSERT INTO %s(id, value, stamp) VALUES('%s', '%s', '%s')",
				sessionTable, sessionID, bin, time.Now().Format("2006-01-02T15:04:05-0700"))
			value, _ := app.getSessionValue(db, sessionID)
			if value != "" {
				sqlString = fmt.Sprintf(
					"UPDATE %s SET value='%s' WHERE id='%s'", sessionTable, bin, sessionID)
			}
			_, err = db.Exec(sqlString)
		}
	}
	defer db.Close()
	return err
}

// Loading the state of a component from a database.
func (app *App) LoadDbSession(sessionID string, data any) (err error) {
	var db *sql.DB
	if db, err = app.checkSessionTable(); err == nil {
		var value string
		if value, err = app.getSessionValue(db, sessionID); value != "" {
			err = json.Unmarshal([]byte(value), &data)
		}
	}
	defer db.Close()
	return err
}
