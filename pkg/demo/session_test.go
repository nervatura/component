package demo

import (
	"database/sql"
	"errors"
	"io/fs"
	"log"
	"os"
	"testing"
)

func TestApp_SaveFileSession(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		memSession map[string]*Demo
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		fileName string
		data     any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "mk_error",
			fields: fields{
				version: "1.0.0",
				osStat: func(name string) (fs.FileInfo, error) {
					return nil, os.ErrNotExist
				},
				osMkdir: func(name string, perm fs.FileMode) error {
					return errors.New("error")
				},
			},
			args: args{
				fileName: "filename",
				data:     &Demo{},
			},
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				version: "1.0.0",
				osStat: func(name string) (fs.FileInfo, error) {
					return nil, nil
				},
				osMkdir: func(name string, perm fs.FileMode) error {
					return nil
				},
				osCreate: func(name string) (*os.File, error) {
					return os.NewFile(0, name), nil
				},
			},
			args: args{
				fileName: "filename",
				data:     &Demo{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				memSession: tt.fields.memSession,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			if err := app.SaveFileSession(tt.args.fileName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("App.SaveSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_LoadFileSession(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		memSession map[string]*Demo
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		fileName string
		data     any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				osReadFile: func(name string) ([]byte, error) {
					return []byte{}, nil
				},
			},
			args: args{
				fileName: "filename",
				data:     map[string]interface{}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				memSession: tt.fields.memSession,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			if err := app.LoadFileSession(tt.args.fileName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("App.LoadSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_checkSessionTable(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		memSession map[string]*Demo
		driverName string
		dataSource string
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "open_error",
			fields: fields{
				driverName: "missing",
				dataSource: "test",
			},
			wantErr: true,
		},
		{
			name: "query_error",
			fields: fields{
				driverName: "sqltest",
				dataSource: "query_error",
			},
			wantErr: true,
		},
		{
			name: "not_found",
			fields: fields{
				driverName: "sqltest",
				dataSource: "not_found",
			},
			wantErr: false,
		},
		{
			name: "found",
			fields: fields{
				driverName: "sqltest",
				dataSource: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				memSession: tt.fields.memSession,
				driverName: tt.fields.driverName,
				dataSource: tt.fields.dataSource,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			_, err := app.checkSessionTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("App.checkSessionTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestApp_getSessionValue(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		memSession map[string]*Demo
		driverName string
		dataSource string
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		db        *sql.DB
		sessionID string
	}
	db, _ := sql.Open("sqltest", "not_found")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "not_found",
			args: args{
				db:        db,
				sessionID: "SESID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				memSession: tt.fields.memSession,
				driverName: tt.fields.driverName,
				dataSource: tt.fields.dataSource,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			_, err := app.getSessionValue(tt.args.db, tt.args.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.getSessionValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestApp_SaveDbSession(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		memSession map[string]*Demo
		driverName string
		dataSource string
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		sessionID string
		data      any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update",
			fields: fields{
				driverName: "sqltest",
				dataSource: "test",
			},
			args: args{
				sessionID: "SESID",
				data:      &Demo{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				memSession: tt.fields.memSession,
				driverName: tt.fields.driverName,
				dataSource: tt.fields.dataSource,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			if err := app.SaveDbSession(tt.args.sessionID, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("App.SaveDbSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_LoadDbSession(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		memSession map[string]*Demo
		driverName string
		dataSource string
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		sessionID string
		data      any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "load",
			fields: fields{
				driverName: "sqltest",
				dataSource: "test",
			},
			args: args{
				sessionID: "SESID",
				data:      &Demo{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				memSession: tt.fields.memSession,
				driverName: tt.fields.driverName,
				dataSource: tt.fields.dataSource,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			if err := app.LoadDbSession(tt.args.sessionID, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("App.LoadDbSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_getSessionTableSql(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		memSession map[string]*Demo
		driverName string
		dataSource string
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		driverName string
		query      string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantSqlString string
	}{
		{
			name: "open_sqlite3",
			args: args{
				driverName: "sqlite3",
				query:      "open",
			},
			wantSqlString: "select name from sqlite_master where name = 'session' ",
		},
		{
			name: "open_postgres",
			args: args{
				driverName: "postgres",
				query:      "open",
			},
			wantSqlString: "select table_name from information_schema.tables where table_name = 'session' ",
		},
		{
			name: "create_mysql",
			args: args{
				driverName: "mysql",
				query:      "create",
			},
			wantSqlString: "CREATE TABLE session ( id VARCHAR(255) NOT NULL, value JSON, stamp VARCHAR(255), PRIMARY KEY (id) );",
		},
		{
			name: "create_postgres",
			args: args{
				driverName: "postgres",
				query:      "create",
			},
			wantSqlString: "CREATE TABLE session ( id VARCHAR(255) NOT NULL PRIMARY KEY, value JSONB, stamp VARCHAR(255) );",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				memSession: tt.fields.memSession,
				driverName: tt.fields.driverName,
				dataSource: tt.fields.dataSource,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			if gotSqlString := app.getSessionTableSql(tt.args.driverName, tt.args.query); gotSqlString != tt.wantSqlString {
				t.Errorf("App.getSessionTableSql() = %v, want %v", gotSqlString, tt.wantSqlString)
			}
		})
	}
}
