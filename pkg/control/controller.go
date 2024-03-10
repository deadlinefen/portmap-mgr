package control

import (
	"sync"

	"github.com/deadlinefen/tinyPortMapper-manager-ipv6/pkg/config"
	"github.com/deadlinefen/tinyPortMapper-manager-ipv6/pkg/ddns"
	"github.com/deadlinefen/tinyPortMapper-manager-ipv6/pkg/job"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type IController interface {
	Start()
	Load(jobs map[string]config.Job) error
	Stop()
}

type Controller struct {
	checker    ddns.IChecker
	jobManager job.IJobManager
}

type IControllerFactory interface {
	NewController() IController
}

type ControllerFactory struct {
	checkerFactory    ddns.ICheckerFactory
	jobManagerFactory job.IJobManagerFactory
}

func (clf *ControllerFactory) NewController() IController {
	ipChan := make(chan string)
	return &Controller{
		checker:    clf.checkerFactory.NewChecker(ipChan),
		jobManager: clf.jobManagerFactory.NewJobManager(ipChan),
	}
}

func NewControllerFactory(cf ddns.ICheckerFactory, jmf job.IJobManagerFactory) IControllerFactory {
	return &ControllerFactory{checkerFactory: cf, jobManagerFactory: jmf}
}

func (c *Controller) Start() {
	var wg sync.WaitGroup

	log.Infof("Job manager starting...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.jobManager.Start()
	}()

	log.Infof("Checker starting...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.checker.Start()
	}()

	wg.Wait()
}

func (c *Controller) Load(jobs map[string]config.Job) error {
	if err := c.jobManager.AddJobs(jobs); err != nil {
		return errors.Wrapf(err, "load jobs failed")
	}
	return nil
}

func (c *Controller) Stop() {
	c.jobManager.Stop()
	log.Infof("Job manager stopped.")
	c.checker.Stop()
	log.Infof("Checker stopped.")
}
