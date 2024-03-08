package job

import (
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Process struct {
	name    string
	cmd     *exec.Cmd
	errChan chan error
	stop    chan struct{}

	wg *sync.WaitGroup
}

func (p *Process) run() {
	p.errChan <- p.cmd.Run()
}

func (p *Process) Stop() {
	close(p.stop)
}

func (p *Process) Begin() {
	p.wg.Add(1)
	defer p.wg.Done()

	for {
		select {
		case <-p.stop:
			p.cmd.Process.Kill()
			<-p.errChan
			return
		case err := <-p.errChan:
			log.Warnf("job %s process exited abnormally, err: %+v", p.name, err)
			go p.run()
		}
	}
}
