package demo

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	ct "github.com/nervatura/component/component"
	fm "github.com/nervatura/component/component/atom"
	bc "github.com/nervatura/component/component/base"
	pg "github.com/nervatura/component/component/page"
	"golang.org/x/sync/errgroup"
)

const (
	httpPort         = 5000
	httpReadTimeout  = 30
	httpWriteTimeout = 30
)

type App struct {
	version    string
	infoLog    *log.Logger
	server     *http.Server
	appSession map[string]bc.ClientComponent
}

func New(version string) (err error) {
	app := &App{
		version:    version,
		infoLog:    log.New(os.Stdout, "INFO: ", log.LstdFlags),
		appSession: make(map[string]bc.ClientComponent),
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
	r := mux.NewRouter()

	r.HandleFunc("/", app.HomeRoute)
	r.HandleFunc("/event", app.AppEvent).Methods("POST")
	var publicFS, _ = fs.Sub(ct.Style, "style")
	r.PathPrefix("/style/").Handler(http.StripPrefix("/style/", http.FileServer(http.FS(publicFS))))

	app.server = &http.Server{
		Handler: csrf.Protect(
			[]byte(bc.RandString(32)),
			csrf.Secure(app.version != "dev"),
		)(r),
		Addr:         fmt.Sprintf(":%d", httpPort),
		ReadTimeout:  time.Duration(httpReadTimeout) * time.Second,
		WriteTimeout: time.Duration(httpWriteTimeout) * time.Second,
	}

	app.infoLog.Printf("HTTP server serving at: %d. \n", httpPort)
	return app.server.ListenAndServe()
}

func (app *App) HomeRoute(w http.ResponseWriter, r *http.Request) {
	tokenID := csrf.Token(r)
	app.appSession[tokenID] = &pg.Application{
		Title:  "Nervatura components",
		Header: bc.SM{"X-CSRF-Token": tokenID},
		HeadLink: []pg.HeadLink{
			{Rel: "icon", Href: "/style/static/favicon.svg", Type: "image/svg+xml"},
			{Rel: "stylesheet", Href: "/style/index.css"},
		},
		MainComponent: pg.NewDemo("/event", "Nervatura components"),
	}
	res, err := app.appSession[tokenID].Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(res))
}

func (app *App) AppEvent(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("app_event_start")
	r.ParseForm()
	sessionID := r.Header.Get("X-CSRF-Token")
	te := bc.TriggerEvent{
		Id:     r.Header.Get("HX-Trigger"),
		Name:   r.Header.Get("HX-Trigger-Name"),
		Target: r.Header.Get("HX-Target"),
		Values: r.Form,
	}
	app.infoLog.Println(te.Name)
	ccApp := app.appSession[sessionID]
	evt := ccApp.OnRequest(te)
	for key, value := range evt.Header {
		w.Header().Set(key, value)
	}
	res, err := evt.Trigger.Render()
	if err == nil {
		ccApp.SetProperty("request_map", evt.Trigger.GetProperty("request_map"))
	} else {
		res, _ = (&fm.Toast{
			Type: fm.ToastTypeError, Value: err.Error(),
		}).Render()
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(res))
}
