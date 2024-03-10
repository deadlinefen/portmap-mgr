package job

import (
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Process struct {
	name string
	cmd  *exec.Cmd

	closed  bool
	restart chan struct{}
	wg      sync.WaitGroup
}

func (p *Process) needRestart(err error) bool {
	if p.closed {
		log.Debugf("process %s close", p.name)
		return false
	}

	if err != nil {
		log.Warnf("Process %s run failed, err: %+v", p.name, err)
	} else if !p.cmd.ProcessState.Success() {
		log.Warnf("Process %s unexpectedly exited, status: %+v", p.name, p.cmd.ProcessState)
	}

	return true
}

func (p *Process) Run() {
	p.wg.Add(1)
	defer p.wg.Done()

	if err := p.cmd.Run(); p.needRestart(err) {
		p.restart <- struct{}{}
	}
}

func (p *Process) Stop() {
	p.closed = true
	p.cmd.Process.Kill()
	p.wg.Wait()
}
