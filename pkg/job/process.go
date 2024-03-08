package job

import (
	"os"
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Process struct {
	name    string
	cmd     *exec.Cmd
	file    *os.File
	errChan chan error
	stop    chan struct{}

	wg *sync.WaitGroup
}

func (p *Process) run() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		p.file.Truncate(0)
		
		p.errChan <- p.cmd.Run()
	}()
}

func (p *Process) Stop() {
	close(p.stop)
}

func (p *Process) Begin() {
	p.run()

	for {
		select {
		case <-p.stop:
			p.cmd.Process.Kill()
			<-p.errChan
			return
		case err := <-p.errChan:
			log.Warnf("job %s process exited abnormally, err: %+v", p.name, err)
			p.run()
		}
	}
}
