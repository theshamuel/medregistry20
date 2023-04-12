package cmd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/theshamuel/medregistry20/app/rest"
	"github.com/theshamuel/medregistry20/app/service"
	"github.com/theshamuel/medregistry20/app/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerCommand struct {
	StoreEngine    StoreGroup `group:"store" namespace:"store" env-namespace:"STORE"`
	MongoURL       string     `long:"mongo-url" env:"MONGO_URL" default:"mongodb://medregdb:27017" description:"url to connect to mongo db server"`
	MongoUsername  string     `long:"mongo-username" env:"MONGO_USERNAME" default:"admin" description:"username to connect to mongo db server"`
	MongoPassword  string     `long:"mongo-password" env:"MONGO_PASSWORD" default:"admin" description:"password to connect to mongo db server "`
	MedregAPIV1URL string     `long:"apiV1url" env:"MEDREG_API_V1_URL" default:"http://localhost:9000/api/v1/" description:"url to medregestry api v1 "`
	ReportsPath    string     `long:"reportsPath" env:"REPORT_PATH" required:"true" default:"./reports" description:"file system path to root report folder"`
	Port           int        `long:"port" env:"SERVER_PORT" default:"9002" description:"application port"`
	Version        string
}

type StoreGroup struct {
	//nolint: staticcheck
	Type   string      `long:"type" env:"STORE_TYPE" description:"type of storage" choice:"Remote" choice:"Mongo" choice:"Mix" default:"Remote"`
	Remote RemoteGroup `group:"Remote" namespace:"Remote" env-namespace:"Remote"`
	Mongo  MongoGroup  `group:"Mongo" namespace:"Mongo" env-namespace:"Mongo"`
	Mix    MixGroup    `group:"Mix" namespace:"Mix" env-namespace:"Mix"`
}

type RemoteGroup struct {
	API          string        `long:"api" env:"API" description:"Remote extension api url"`
	Timeout      time.Duration `long:"timeout" env:"TIMEOUT" default:"5s" description:"http timeout"`
	AuthUser     string        `long:"auth_user" env:"AUTH_USER" description:"basic auth user name"`
	AuthPassword string        `long:"auth_passwd" env:"AUTH_PASSWD" description:"basic auth user password"`
}

type MongoGroup struct {
	URL string `long:"mongo-url" env:"MONGO_URL" description:"MongoDB URL"`
}

type MixGroup struct {
	RemoteGroup
	MongoGroup
}

type application struct {
	*ServerCommand
	rest        *rest.Rest
	dataService *service.DataStore
	terminated  chan struct{}
}

// Execute is the entry point for server command
func (sc *ServerCommand) Execute(_ []string) error {
	log.Printf("[INFO] start app server")
	log.Printf("[INFO] server args:\n"+
		"		port: %d;\n"+
		"		report path: %s;\n"+
		"		mongoURL %s;\n"+
		"		medregAPIV1URL: %s;\n"+
		"		store engine type: %s;\n",
		sc.Port, sc.ReportsPath, sc.MongoURL, sc.MedregAPIV1URL, sc.StoreEngine.Type)

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
		if e := app.dataService.Close(); e != nil {
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
	log.Printf("[DEBUG] build data engine store. Type=%s", sc.StoreEngine.Type)

	switch sc.StoreEngine.Type {
	case "Remote":
		r := &store.Remote{URI: sc.MedregAPIV1URL}
		return r, nil
	case "Mongo":
		return nil, errors.Errorf("not implemented yet")
	case "Mix":
		clientOpts := options.Client().ApplyURI(sc.MongoURL)
		credential := options.Credential{
			AuthSource: "medregDB",
			Username:   sc.MongoUsername,
			Password:   sc.MongoPassword,
		}

		//Go with assumption that password cannot be empty
		if credential.Password != "" && credential.Username != "" {
			clientOpts = clientOpts.SetAuth(credential)
		}

		client, err := mongo.Connect(context.Background(), clientOpts)
		if err != nil {
			return nil, errors.Errorf("can't initialize data store because failed to establish mongo connection: %s", sc.StoreEngine.Type)
		}

		var result bson.M
		if err := client.Database("medregDB").RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
			log.Printf("[ERROR] cannot ping medregDB: %#v", err)
			return nil, errors.New("cannot connect to MongoDB")
		}

		log.Printf("[INFO] ping medregDB successfully")
		r := &store.Mix{
			URI:         sc.MedregAPIV1URL,
			MongoClient: client,
		}
		return r, nil
	default:
		return nil, errors.Errorf("can't initialize data store, unsupported store type %s", sc.StoreEngine.Type)
	}
}

// Wait for application completion (termination)
func (app *application) Wait() {
	<-app.terminated
}
