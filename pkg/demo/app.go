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
	"context"
	"errors"
	"strings"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ct "github.com/nervatura/component/pkg/component"
	st "github.com/nervatura/component/pkg/static"
	ut "github.com/nervatura/component/pkg/util"
	"golang.org/x/sync/errgroup"
)

// Demo [App] constants
const (
	httpReadTimeout  = 30
	httpWriteTimeout = 30

	sessionSave = false
	sessionPath = "session"
)

// Demo application
type App struct {
	version    string
	infoLog    *log.Logger
	server     *http.Server
	memSession map[string]*Demo
	osStat     func(name string) (fs.FileInfo, error)
	osMkdir    func(name string, perm fs.FileMode) error
	osCreate   func(name string) (*os.File, error)
	osReadFile func(name string) ([]byte, error)
}

// It creates a new application and starts an http server.
func New(version string, httpPort int64) (err error) {
	app := &App{
		version:    version,
		infoLog:    log.New(os.Stdout, "INFO: ", log.LstdFlags),
		memSession: make(map[string]*Demo),
		osStat:     os.Stat,
		osMkdir:    os.Mkdir,
		osCreate:   os.Create,
		osReadFile: os.ReadFile,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return app.startHttpService(httpPort)
	})

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}
	app.infoLog.Println("received shut down signal")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if app.server != nil {
		app.infoLog.Println("stopping HTTP server")
		_ = app.server.Shutdown(shutdownCtx)
	}

	return g.Wait()
}

// It sets the http routes and starts the server.
func (app *App) startHttpService(httpPort int64) error {
	mux := http.NewServeMux()
	// Register API routes.
	mux.HandleFunc("/", app.HomeRoute)
	mux.HandleFunc("/session", app.HomeRoute)
	mux.HandleFunc("POST /event", app.AppEvent)

	// Register static dirs.
	var publicFS, _ = fs.Sub(st.Static, ".")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(publicFS))))

	app.server = &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf(":%d", httpPort),
		ReadTimeout:  time.Duration(httpReadTimeout) * time.Second,
		WriteTimeout: time.Duration(httpWriteTimeout) * time.Second,
	}

	app.infoLog.Printf("HTTP server serving at: %d. \n", httpPort)
	return app.server.ListenAndServe()
}

// Saving component state in a session json file.
func (app *App) SaveSession(fileName string, data any) error {
	if _, err := app.osStat(sessionPath); errors.Is(err, os.ErrNotExist) {
		err := app.osMkdir(sessionPath, os.ModePerm)
		if err != nil {
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
func (app *App) LoadSession(fileName string, data any) (err error) {
	filePath := fmt.Sprintf(`%s/%s.json`, sessionPath, fileName)
	sessionFile, err := app.osReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(sessionFile, &data)
	}
	return err
}

// Creates and returns an Application/[Demo] component.
// It stores the state of the component in memory or in a session file
func (app *App) HomeRoute(w http.ResponseWriter, r *http.Request) {
	tokenID := ut.RandString(24)
	sessionID := base64.StdEncoding.EncodeToString([]byte(tokenID))
	dataSave := sessionSave || strings.Contains(r.URL.Path, "/session")
	demo := NewDemo("/event", "Nervatura components")
	ccApp := &ct.Application{
		Title:  "Nervatura components",
		Header: ut.SM{"X-Session-Token": tokenID},
		HeadLink: []ct.HeadLink{
			{Rel: "icon", Href: "/static/favicon.svg", Type: "image/svg+xml"},
			{Rel: "stylesheet", Href: "/static/css/index.css"},
		},
		MainComponent: demo,
	}
	var err error
	var res string
	res, err = ccApp.Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if dataSave {
		err = app.SaveSession(sessionID, demo)
	}
	if (err != nil) || !dataSave {
		app.memSession[sessionID] = demo
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(res))
}

// Receive the component event request.
// Loads the Demo component based on the X-Session-Token identifier.
func (app *App) AppEvent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tokenID := r.Header.Get("X-Session-Token")
	sessionID := base64.StdEncoding.EncodeToString([]byte(tokenID))
	dataSave := sessionSave || strings.Contains(r.Header.Get("Hx-Current-Url"), "/session")
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
		if err = app.LoadSession(sessionID, &demo); err == nil {
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
		app.SaveSession(sessionID, demo)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(res))
}
