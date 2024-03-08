package job

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/deadlinefen/portmap-mgr/pkg/config"
)

type Job struct {
	info *config.Job
	name string

	process *Process
	wg      *sync.WaitGroup
}

func (j *Job) Run(mapper *config.Mapper, ipv6 string) {
	if j.process != nil {
		j.process.Stop()
	}

	cmdLine := fmt.Sprintf("%s -r[%s]:%d -l%s:%d -%s > %s/%d->%d.log",
		mapper.Bin, ipv6, j.info.FromPort,
		j.info.ToIp, j.info.ToPort, j.info.Type,
		mapper.FileDir, j.info.FromPort, j.info.ToPort)

	j.process = &Process{
		name:    j.name,
		cmd:     exec.Command(cmdLine),
		errChan: make(chan error),
		stop:    make(chan struct{}),
		wg:      j.wg,
	}

	go j.process.Begin()
}

func (j *Job) Stop() {
	if j.process != nil {
		j.process.Stop()
	}
	j.wg.Wait()
}
