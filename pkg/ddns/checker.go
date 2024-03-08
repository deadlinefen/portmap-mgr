package ddns

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type IChecker interface {
	Begin()
	CheckDNSOnce() bool
	Stop()
}

type Checker struct {
	resoluter IResoluter
	ip        string
	ttl       int

	ipChan chan string
	stop   chan struct{}
}

func (c *Checker) CheckDNSOnce() bool {
	log.Debug("Check DNS Once...")

	ipNow, err := c.resoluter.ResoluteOneIp()
	if err != nil {
		log.Warnf("Check dns failed, err: %+v", err)
		return false
	}

	if ipNow != c.ip {
		log.Debugf("Ip changed to %s", ipNow)
		c.ipChan <- ipNow
	}

	return true
}

func (c *Checker) Begin() {
	ticker := time.NewTicker(time.Second * time.Duration(c.ttl))

	for {
		select {
		case <-c.stop:
			break
		case <-ticker.C:
			for !c.CheckDNSOnce() {
			}
		}
	}
}

func (c *Checker) Stop() {
	close(c.stop)
}

type ICheckerFactory interface {
	NewChecker(resoluter IResoluter, ttl int, ipChan chan string) IChecker
}

type CheckerFactory struct {
}

func (cf *CheckerFactory) NewChecker(resoluter IResoluter, ttl int, ipChan chan string) IChecker {
	return &Checker{
		resoluter: resoluter,
		ip:        "",
		ttl:       ttl,
		ipChan:    ipChan,
		stop:      make(chan struct{}),
	}
}

func NewCheckerFactory() ICheckerFactory {
	return &CheckerFactory{}
}
