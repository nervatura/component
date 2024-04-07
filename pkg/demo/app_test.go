package demo

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		version  string
		httpPort int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "start",
			args: args{
				version:  "test",
				httpPort: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New(tt.args.version, tt.args.httpPort); (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_SaveSession(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		server     *http.Server
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
				server:     tt.fields.server,
				memSession: tt.fields.memSession,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			if err := app.SaveSession(tt.args.fileName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("App.SaveSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_LoadSession(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		server     *http.Server
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
				server:     tt.fields.server,
				memSession: tt.fields.memSession,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			if err := app.LoadSession(tt.args.fileName, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("App.LoadSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_HomeRoute(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		server     *http.Server
		memSession map[string]*Demo
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok_save",
			fields: fields{
				memSession: map[string]*Demo{},
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
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/session", nil),
			},
		},
		{
			name: "ok_mem",
			fields: fields{
				memSession: map[string]*Demo{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				server:     tt.fields.server,
				memSession: tt.fields.memSession,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			app.HomeRoute(tt.args.w, tt.args.r)
		})
	}
}

func TestApp_AppEvent(t *testing.T) {
	type fields struct {
		version    string
		infoLog    *log.Logger
		server     *http.Server
		memSession map[string]*Demo
		osStat     func(name string) (fs.FileInfo, error)
		osMkdir    func(name string, perm fs.FileMode) error
		osCreate   func(name string) (*os.File, error)
		osReadFile func(name string) ([]byte, error)
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	sessionID := base64.StdEncoding.EncodeToString([]byte("SessionID"))
	demoApp := &Demo{}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "mem",
			fields: fields{
				memSession: map[string]*Demo{
					sessionID: demoApp,
				},
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
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", nil),
			},
		},
		{
			name: "file",
			fields: fields{
				memSession: map[string]*Demo{},
				osStat: func(name string) (fs.FileInfo, error) {
					return nil, nil
				},
				osMkdir: func(name string, perm fs.FileMode) error {
					return nil
				},
				osCreate: func(name string) (*os.File, error) {
					return os.NewFile(0, name), nil
				},
				osReadFile: func(name string) ([]byte, error) {
					app, _ := json.Marshal(demoApp)
					return []byte(app), nil
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", nil),
			},
		},
		{
			name: "missing",
			fields: fields{
				memSession: map[string]*Demo{},
				osStat: func(name string) (fs.FileInfo, error) {
					return nil, nil
				},
				osMkdir: func(name string, perm fs.FileMode) error {
					return nil
				},
				osCreate: func(name string) (*os.File, error) {
					return os.NewFile(0, name), nil
				},
				osReadFile: func(name string) ([]byte, error) {
					return nil, errors.New("error")
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:    tt.fields.version,
				infoLog:    tt.fields.infoLog,
				server:     tt.fields.server,
				memSession: tt.fields.memSession,
				osStat:     tt.fields.osStat,
				osMkdir:    tt.fields.osMkdir,
				osCreate:   tt.fields.osCreate,
				osReadFile: tt.fields.osReadFile,
			}
			tt.args.r.Header.Set("X-Session-Token", "SessionID")
			tt.args.r.Header.Set("Hx-Current-Url", "/session")
			app.AppEvent(tt.args.w, tt.args.r)
		})
	}
}
