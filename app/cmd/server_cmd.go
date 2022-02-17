package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/theshamuel/medregistry20/app/rest"
	"github.com/theshamuel/medregistry20/app/service"
	"github.com/theshamuel/medregistry20/app/store"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerCommand struct {
	StoreEngine StoreGroup `group:"store" namespace:"store" env-namespace:"STORE"`
	Version     string
	Port        int `long:"port" env:"SERVER_PORT" default:"9002" description:"application port"`
	CommonOptions
}

type StoreGroup struct {
	Type   string      `long:"type" env:"TYPE" description:"type of storage" choice:"Remote" default:"Remote"`
	Remote RemoteGroup `group:"Remote" namespace:"Remote" env-namespace:"Remote"`
}

type RemoteGroup struct {
	API          string        `long:"api" env:"API" description:"Remote extension api url"`
	Timeout      time.Duration `long:"timeout" env:"TIMEOUT" default:"5s" description:"http timeout"`
	AuthUser     string        `long:"auth_user" env:"AUTH_USER" description:"basic auth user name"`
	AuthPassword string        `long:"auth_passwd" env:"AUTH_PASSWD" description:"basic auth user password"`
}

type application struct {
	*ServerCommand
	rest        *rest.Rest
	dataService *service.DataStore
	terminated  chan struct{}
}

//Execute is the entry point for server command
func (sc *ServerCommand) Execute(_ []string) error {
	log.Printf("[INFO] start app server")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] Get interrupt signal")
		cancel()
	}()
	app, err := sc.bootstrapApp()
	if err != nil {
		log.Printf("[PANIC] Failed to setup application, %v", err)
		return err
	}
	if err = app.run(ctx); err != nil {
		log.Printf("[ERROR] Server terminated with error %v", err)
		return err
	}
	log.Printf("[INFO] Server terminated")
	return nil
}

func (app *application) run(ctx context.Context) error {

	go func() {
		<-ctx.Done()
		app.rest.Shutdown()
		log.Print("[INFO] shutdown is completed")
		if e := app.dataService.Close()
			e != nil {
			log.Printf("[WARN] failed to close data store, %s", e)
		}
	}()

	app.rest.Run(app.Port)
	close(app.terminated)
	return nil
}

func (sc *ServerCommand) bootstrapApp() (*application, error) {

	storeEngine, err := sc.buildDataEngine()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build data store engine")
	}

	rest := &rest.Rest{
		Version:     sc.Version,
		URI:         sc.MedregAPIV1URL,
		ReportPath:  sc.ReportsPath,
		DataService: &service.DataStore{Engine: storeEngine, ReportPath: sc.ReportsPath},
	}

	return &application{
		ServerCommand: sc,
		rest:          rest,
		dataService:   &service.DataStore{Engine: storeEngine},
		terminated:    make(chan struct{}),
	}, nil
}

func (sc *ServerCommand) buildDataEngine() (result store.EngineInterface, err error) {
	log.Printf("[INFO] build data engine store. Type=%s", sc.StoreEngine.Type)

	switch sc.StoreEngine.Type {
	case "Remote":
		r := &store.Remote{URI: sc.MedregAPIV1URL}
		return r, nil
	default:
		return nil, errors.Errorf("can't initialize data store, unsupported store type %s", sc.StoreEngine.Type)
	}
}

// Wait for application completion (termination)
func (app *application) Wait() {
	<-app.terminated
}
