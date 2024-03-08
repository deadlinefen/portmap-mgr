package ddns

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type IChecker interface {
	Start()
	Stop()
}

type Checker struct {
	resoluter IResoluter
	ip        string
	ttl       int

	ipChan chan string
	stop   chan struct{}
}

func (c *Checker) checkDNSOnce() bool {
	log.Debug("Check DNS Once...")

	ipNow, err := c.resoluter.ResoluteOneIp()
	if err != nil {
		log.Warnf("Check dns failed, err: %+v", err)
		return false
	}

	if ipNow != c.ip {
		log.Debugf("Ip changed to %s", ipNow)
		c.ipChan <- ipNow
		c.ip = ipNow
	}

	return true
}

func (c *Checker) checkDNS() {
	for !c.checkDNSOnce() {
		log.Infof("check DNS failed, retry.")
	}
}

func (c *Checker) Start() {
	ticker := time.NewTicker(time.Second * time.Duration(c.ttl))
	c.checkDNS()

	for {
		select {
		case <-c.stop:
			log.Debugf("checker start() quit.")
			return
		case <-ticker.C:
			c.checkDNS()
		}
	}
}

func (c *Checker) Stop() {
	close(c.stop)
}

type ICheckerFactory interface {
	NewChecker(ipChan chan string) IChecker
}

type CheckerFactory struct {
	resoluterFactory IResoluterFactory
	ttl              int
}

func (cf *CheckerFactory) NewChecker(ipChan chan string) IChecker {
	return &Checker{
		resoluter: cf.resoluterFactory.NewResoluter(),
		ip:        "",
		ttl:       cf.ttl,
		ipChan:    ipChan,
		stop:      make(chan struct{}),
	}
}

func NewCheckerFactory(rf IResoluterFactory, ttl int) ICheckerFactory {
	return &CheckerFactory{
		resoluterFactory: rf,
		ttl:              ttl,
	}
}
