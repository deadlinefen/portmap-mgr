package job

import (
	"os/exec"
	"sync"
)

type Process struct {
	cmd *exec.Cmd

	closed  bool
	restart chan struct{}
	wg      sync.WaitGroup
}

func (p *Process) Run() {
	p.wg.Add(1)
	defer p.wg.Done()

	if err := p.cmd.Run(); err != nil && !p.closed {
		close(p.restart)
	}
}

func (p *Process) Stop() {
	p.closed = true
	p.cmd.Process.Kill()
	p.wg.Wait()
}
