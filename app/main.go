package main

import (
	"github.com/hashicorp/logutils"
	"github.com/jessevdk/go-flags"
	"github.com/theshamuel/medregistry20/app/cmd"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var version = "unknown"

type Opts struct {
	ServerCmd      cmd.ServerCommand `command:"server"`
	MedregAPIV1URL string            `long:"apiV1url" env:"MEDREG_API_V1_URL" default:"http://localhost:9000/api/v1/" description:"url to medregestry api v1 "`
	ReportsPath    string            `long:"reportsPath" env:"REPORT_PATH" required:"true" default:"./reports" description:"file system path to root report folder"`
	Debug          bool              `long:"debug" env:"DEBUG" description:"debug mode"`
}

func main() {
	setupLogLevel(false)
	log.Printf("[INFO] Starting Medregistry API v2 version:%s ...\n", version)
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLogLevel(opts.Debug)
		c := command.(cmd.CommonOptionsCommander)
		c.SetCommon(cmd.CommonOptions{
			MedregAPIV1URL: opts.MedregAPIV1URL,
			ReportsPath:    opts.ReportsPath,
		})
		err := c.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed with %v", err)
		}
		return err
	}
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func setupLogLevel(debug bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}
	log.SetFlags(log.Ldate | log.Ltime)

	if debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}
	log.SetOutput(filter)
}

func getStackTrace() string {
	maxSize := 7 * 1024 * 1024
	stacktrace := make([]byte, maxSize)
	length := runtime.Stack(stacktrace, true)
	if length > maxSize {
		length = maxSize
	}
	return string(stacktrace[:length])
}

func init() {
	sigChan := make(chan os.Signal)
	go func() {
		for range sigChan {
			log.Printf("[INFO] Singal QUITE is cought , stacktrace [\n%s", getStackTrace())
		}
	}()
	signal.Notify(sigChan, syscall.SIGQUIT)
}
