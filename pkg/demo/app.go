/*
Component demo application

1. 💻 Ensure that you have Golang installed on your system. If not, please follow the
https://golang.org/doc/install.

2. 📦 Clone the repository:

	git clone https://github.com/nervatura/component.git

3. 📂 Change into the project directory:

	cd component

4. 🔨 Build the demo project:

	go build -ldflags="-w -s -X main.version=demo" -o ./component main.go

5. 🌍 Run the demo application:

	./component 5000

The demo application can store session data in memory and as
session files:
  - open the http://localhost:5000/ (memory session)
  - or http://localhost:5000/session/ (file session)
*/
package demo

import (
	"database/sql"
	"embed"
	"errors"
	"strings"

	"encoding/base64"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	// _ "github.com/mattn/go-sqlite3"
	// _ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/denisenkom/go-mssqldb"

	ct "github.com/nervatura/component/pkg/component"
	st "github.com/nervatura/component/pkg/static"
	ut "github.com/nervatura/component/pkg/util"
)

// Demo [App] constants
const (
	httpReadTimeout  = 30
	httpWriteTimeout = 30
	sessionPath      = "session"
	sessionTable     = "session"
	// sqlite3
	sessionSource = "file:./session.db?cache=shared&mode=rwc"
	// postgres
	// sessionSource = "postgres://postgres:password@172.18.0.1:5432/session?sslmode=disable"
	// mysql
	// sessionSource = "mysql://root:password@tcp(localhost:3306)/session"
	// mssql
	// sessionSource = "mssql://sa:Password1234_1@localhost:1433?database=session"
)

//go:embed static
var Public embed.FS

// Demo application
type App struct {
	version     string
	infoLog     *log.Logger
	memSession  map[string]*Demo
	driverName  string
	dataSource  string
	saveSession func(name string, data any) (err error)
	loadSession func(name string, data any) (err error)
	osStat      func(name string) (fs.FileInfo, error)
	osMkdir     func(name string, perm fs.FileMode) error
	osCreate    func(name string) (*os.File, error)
	osReadFile  func(name string) ([]byte, error)
}

// It creates a new application and starts an http server.
func New(version string, httpPort int64) {
	app := &App{
		version:    version,
		infoLog:    log.New(os.Stdout, "INFO: ", log.LstdFlags),
		memSession: make(map[string]*Demo),
		dataSource: sessionSource,
		osStat:     os.Stat,
		osMkdir:    os.Mkdir,
		osCreate:   os.Create,
		osReadFile: os.ReadFile,
	}
	app.saveSession = app.SaveFileSession
	app.loadSession = app.LoadFileSession
	if len(sql.Drivers()) > 0 {
		app.driverName = sql.Drivers()[0]
		app.saveSession = app.SaveDbSession
		app.loadSession = app.LoadDbSession
	}
	mux := http.NewServeMux()
	// Register API routes.
	mux.HandleFunc("/", app.HomeRoute)
	mux.HandleFunc("/session", app.HomeRoute)
	mux.HandleFunc("POST /event", app.AppEvent)

	// Register static dirs.
	// app (demo component) css files
	var publicFS, _ = fs.Sub(Public, "static")
	// components css files
	var staticFS, _ = fs.Sub(st.Static, ".")
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(publicFS))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	server := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf(":%d", httpPort),
		ReadTimeout:  time.Duration(httpReadTimeout) * time.Second,
		WriteTimeout: time.Duration(httpWriteTimeout) * time.Second,
	}

	app.infoLog.Printf("HTTP server serving at: %d. \n", httpPort)
	if err := server.ListenAndServe(); err != nil {
		app.infoLog.Printf("server error: %s\n", err)
	}
}

func (app *App) respondMessage(w http.ResponseWriter, res string, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(res))
}

// Creates and returns an Application/[Demo] component.
// It stores the state of the component in memory or in a session file
func (app *App) HomeRoute(w http.ResponseWriter, r *http.Request) {
	tokenID := ut.RandString(24)
	sessionID := base64.StdEncoding.EncodeToString([]byte(tokenID))
	dataSave := strings.Contains(r.URL.Path, "/session")
	demo := NewDemo("/event", "Nervatura components")
	ccApp := &ct.Application{
		Title:  "Nervatura components",
		Header: ut.SM{"X-Session-Token": tokenID},
		HeadLink: []ct.HeadLink{
			{Rel: "icon", Href: "/static/favicon.svg", Type: "image/svg+xml"},
			{Rel: "stylesheet", Href: "/public/demo.css"},
			{Rel: "stylesheet", Href: "/static/css/index.css"},
		},
		MainComponent: demo,
	}
	var err error
	var res string
	if res, err = ccApp.Render(); err == nil {
		if dataSave {
			err = app.saveSession(sessionID, demo)
		} else {
			app.memSession[sessionID] = demo
		}
	}
	app.respondMessage(w, res, err)
}

// Receive the component event request.
// Loads the Demo component based on the X-Session-Token identifier.
func (app *App) AppEvent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tokenID := r.Header.Get("X-Session-Token")
	sessionID := base64.StdEncoding.EncodeToString([]byte(tokenID))
	dataSave := strings.Contains(r.Header.Get("Hx-Current-Url"), "/session")
	te := ct.TriggerEvent{
		Id:     r.Header.Get("HX-Trigger"),
		Name:   r.Header.Get("HX-Trigger-Name"),
		Target: r.Header.Get("HX-Target"),
		Values: r.Form,
	}
	var err error
	var evt ct.ResponseEvent
	var demo *Demo
	if mem, found := app.memSession[sessionID]; found {
		evt = mem.OnRequest(te)
	} else if dataSave {
		if err = app.loadSession(sessionID, &demo); err == nil {
			demo.DemoMap = DemoMap
			demo.RequestMap = map[string]ct.ClientComponent{}
			demo.InitDemoMap()
			_, err = demo.Render()
			if err == nil {
				evt = demo.OnRequest(te)
			}
		}
	}
	for key, value := range evt.Header {
		w.Header().Set(key, value)
	}
	var res string
	if evt.Trigger != nil {
		res, err = evt.Trigger.Render()
	} else {
		err = errors.New("missing component")
	}
	if err != nil {
		res, _ = (&ct.Toast{
			Type: ct.ToastTypeError, Value: err.Error(),
		}).Render()
	}
	if dataSave && (err == nil) {
		app.saveSession(sessionID, demo)
	}
	app.respondMessage(w, res, nil)
}
