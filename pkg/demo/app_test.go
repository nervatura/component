package demo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/nervatura/component/pkg/demo/sqltest"
)

func TestApp_HomeRoute(t *testing.T) {
	type fields struct {
		version     string
		infoLog     *log.Logger
		memSession  map[string]*Demo
		saveSession func(name string, data any) (err error)
		loadSession func(name string, data any) (err error)
		osStat      func(name string) (fs.FileInfo, error)
		osMkdir     func(name string, perm fs.FileMode) error
		osCreate    func(name string) (*os.File, error)
		osReadFile  func(name string) ([]byte, error)
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
				loadSession: func(name string, data any) (err error) {
					return nil
				},
				saveSession: func(name string, data any) (err error) {
					return nil
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
				version:     tt.fields.version,
				infoLog:     tt.fields.infoLog,
				memSession:  tt.fields.memSession,
				loadSession: tt.fields.loadSession,
				saveSession: tt.fields.saveSession,
				osStat:      tt.fields.osStat,
				osMkdir:     tt.fields.osMkdir,
				osCreate:    tt.fields.osCreate,
				osReadFile:  tt.fields.osReadFile,
			}
			app.HomeRoute(tt.args.w, tt.args.r)
		})
	}
}

func TestApp_AppEvent(t *testing.T) {
	type fields struct {
		version     string
		infoLog     *log.Logger
		memSession  map[string]*Demo
		saveSession func(name string, data any) (err error)
		loadSession func(name string, data any) (err error)
		osStat      func(name string) (fs.FileInfo, error)
		osMkdir     func(name string, perm fs.FileMode) error
		osCreate    func(name string) (*os.File, error)
		osReadFile  func(name string) ([]byte, error)
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	upload := func() (body io.Reader) {
		bodyBuffer := new(bytes.Buffer)
		mw := multipart.NewWriter(bodyBuffer)
		if part, err := mw.CreateFormFile("file", "myfile.txt"); err == nil {
			file := strings.NewReader(`------WebKitFormBoundaryePkpFF7tjBAqx29L--`)
			io.Copy(part, file)
		}
		return bodyBuffer
	}

	sessionID := base64.StdEncoding.EncodeToString([]byte("SessionID"))
	demoApp := &Demo{}
	tests := []struct {
		name    string
		fields  fields
		args    args
		content string
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
				loadSession: func(name string, data any) (err error) {
					return nil
				},
				saveSession: func(name string, data any) (err error) {
					return nil
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", nil),
			},
			content: "application/x-www-form-urlencoded",
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
				loadSession: func(name string, data any) (err error) {
					return nil
				},
				saveSession: func(name string, data any) (err error) {
					return nil
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", nil),
			},
			content: "application/x-www-form-urlencoded",
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
				loadSession: func(name string, data any) (err error) {
					return nil
				},
				saveSession: func(name string, data any) (err error) {
					return nil
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", nil),
			},
			content: "application/x-www-form-urlencoded",
		},
		{
			name: "file",
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
				loadSession: func(name string, data any) (err error) {
					return nil
				},
				saveSession: func(name string, data any) (err error) {
					return nil
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", upload()),
			},
			content: "multipart/form-data; boundary=----WebKitFormBoundaryePkpFF7tjBAqx29L",
		},
		{
			name: "json",
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
				loadSession: func(name string, data any) (err error) {
					return nil
				},
				saveSession: func(name string, data any) (err error) {
					return nil
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/event", nil),
			},
			content: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				version:     tt.fields.version,
				infoLog:     tt.fields.infoLog,
				memSession:  tt.fields.memSession,
				loadSession: tt.fields.loadSession,
				saveSession: tt.fields.saveSession,
				osStat:      tt.fields.osStat,
				osMkdir:     tt.fields.osMkdir,
				osCreate:    tt.fields.osCreate,
				osReadFile:  tt.fields.osReadFile,
			}
			app.loadSession = app.LoadFileSession
			app.saveSession = app.SaveFileSession
			tt.args.r.Header.Set("X-Session-Token", "SessionID")
			tt.args.r.Header.Set("Hx-Current-Url", "/session")
			tt.args.r.Header.Set("Content-Type", tt.content)
			app.AppEvent(tt.args.w, tt.args.r)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		version  string
		httpPort int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "start",
			args: args{
				version:  "test",
				httpPort: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New(tt.args.version, tt.args.httpPort)
		})
	}
}

func TestApp_respondMessage(t *testing.T) {
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
		w    http.ResponseWriter
		html template.HTML
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "error",
			args: args{
				w:   httptest.NewRecorder(),
				err: errors.New("error"),
			},
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
			app.respondMessage(tt.args.w, tt.args.html, tt.args.err)
		})
	}
}
