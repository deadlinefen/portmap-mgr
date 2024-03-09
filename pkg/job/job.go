package job

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/deadlinefen/portmap-mgr/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Job struct {
	info *config.Job
	name string
	bin  string
	log  *os.File
	ipv6 string

	process *Process
	restart chan struct{}

	stop chan struct{}
}

func (j *Job) Start() {
	for {
		select {
		case <-j.stop:
			j.closeAll()
			log.Infof("Job %s stopped.", j.name)
			return
		case <-j.restart:
			j.process = nil
			log.Warnf("mapper dead, restart after 3s...")
			time.Sleep(time.Second * time.Duration(3))
			j.Run(j.ipv6)
		}
	}
}

func (j *Job) Run(ipv6 string) {
	if j.process != nil {
		j.process.Stop()
	}

	j.ipv6 = ipv6
	j.process = j.createProcess()
	j.restart = j.process.restart

	go j.process.Run()
}

func (j *Job) Stop() {
	close(j.stop)
}

func (j *Job) createProcess() *Process {
	local := fmt.Sprintf("-l[%s]:%d", j.ipv6, j.info.FromPort)
	remote := fmt.Sprintf("-r%s:%d", j.info.ToIp, j.info.ToPort)
	mapType := fmt.Sprintf("-%s", j.info.Type)

	cmd := exec.Command(j.bin, local, remote, mapType)
	cmd.Stdout = j.log

	return &Process{
		cmd:     cmd,
		closed:  false,
		restart: make(chan struct{}),
	}
}

func (j *Job) closeAll() {
	if j.process != nil {
		j.process.Stop()
	}
	if j.log != nil {
		j.log.Close()
	}
}
