package main

import (
	"flag"
	"os"
	"sync"
	"syscall"

	"github.com/deadlinefen/portmap-mgr/pkg/config"
	"github.com/deadlinefen/portmap-mgr/pkg/control"
	"github.com/deadlinefen/portmap-mgr/pkg/ddns"
	"github.com/deadlinefen/portmap-mgr/pkg/job"
	"github.com/deadlinefen/portmap-mgr/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var tomlPath string

func init() {
	flag.StringVar(&tomlPath, "c", "config.toml", "toml config file path")
	flag.Parse()
}

func main() {
	cfg, err := config.NewParserFactory(tomlPath).NewParser().Parse()
	if err != nil {
		panic(errors.Wrapf(err, "parse config failed"))
	}

	// step 1: init log
	utils.InitLog(&cfg.Log)

	// step 2: init factory
	rf := ddns.NewResoluterFactory(&cfg.Resolution)
	cf := ddns.NewCheckerFactory(rf, cfg.Resolution.Ttl)
	jmf := job.NewJobManagerFactory(&cfg.Mapper)
	clf := control.NewControllerFactory(cf, jmf)

	// step 3: run controller
	var wg sync.WaitGroup
	controller := clf.NewController()
	if err := controller.Load(cfg.Jobs); err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		controller.Start()
	}()

	// Step 3: start signal mux
	signalHandler := func(signal os.Signal) bool {
		log.Infof("handle signal: %s", signal.String())
		switch signal {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			controller.Stop()
			log.Info("all service stopped.")
			return true
		default:
			return false
		}
	}
	signalMux := utils.NewSignalMux(signalHandler)
	wg.Add(1)
	go func() {
		defer wg.Done()
		signalMux.Serve()
	}()

	wg.Wait()
}
