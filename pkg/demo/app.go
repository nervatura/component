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

const (
	httpPort         = 5000
	httpReadTimeout  = 30
	httpWriteTimeout = 30

	sessionSave = false
	sessionPath = "session"
)

type App struct {
	version    string
	infoLog    *log.Logger
	server     *http.Server
	memSession map[string]*Demo
}

func New(version string) (err error) {
	app := &App{
		version:    version,
		infoLog:    log.New(os.Stdout, "INFO: ", log.LstdFlags),
		memSession: make(map[string]*Demo),
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return app.startHttpService()
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

func (app *App) startHttpService() error {
	mux := http.NewServeMux()
	// Register API routes.
	mux.HandleFunc("/", app.HomeRoute)
	mux.HandleFunc("/save", app.HomeRoute)
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

func (app *App) SaveSession(fileName string, data any) error {
	if _, err := os.Stat(sessionPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(sessionPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filePath := fmt.Sprintf(`%s/%s.json`, sessionPath, fileName)
	sessionFile, err := os.Create(filePath)
	if err == nil {
		bin, err := json.Marshal(data)
		if err == nil {
			sessionFile.Write(bin)
		}
	}
	defer sessionFile.Close()
	return err
}

func (app *App) LoadSession(fileName string, data any) (err error) {
	filePath := fmt.Sprintf(`%s/%s.json`, sessionPath, fileName)
	sessionFile, err := os.ReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(sessionFile, &data)
	}
	return err
}

func (app *App) HomeRoute(w http.ResponseWriter, r *http.Request) {
	tokenID := ut.RandString(32)
	sessionID := base64.StdEncoding.EncodeToString([]byte(tokenID))[:24]
	dataSave := sessionSave || strings.Contains(r.URL.Path, "/save")
	demo := NewDemo("/event", "Nervatura components")
	ccApp := &ct.Application{
		Title:  "Nervatura components",
		Header: ut.SM{"X-CSRF-Token": tokenID},
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

func (app *App) AppEvent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tokenID := r.Header.Get("X-CSRF-Token")
	sessionID := base64.StdEncoding.EncodeToString([]byte(tokenID))[:24]
	dataSave := sessionSave || strings.Contains(r.Header.Get("Hx-Current-Url"), "/save")
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
	res, err = evt.Trigger.Render()
	if err != nil {
		res, _ = (&ct.Toast{
			Type: ct.ToastTypeError, Value: err.Error(),
		}).Render()
	}
	if dataSave {
		app.SaveSession(sessionID, demo)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(res))
}
