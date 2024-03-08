package job

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/deadlinefen/portmap-mgr/pkg/config"
)

type Job struct {
	info *config.Job
	name string
	file *os.File

	process *Process
	wg      *sync.WaitGroup
}

func (j *Job) Run(mapper *config.Mapper, ipv6 string) {
	if j.process != nil {
		j.process.Stop()
	}

	local := fmt.Sprintf("-l[%s]:%d", ipv6, j.info.FromPort)
	remote := fmt.Sprintf("-r%s:%d", j.info.ToIp, j.info.ToPort)
	mapType := fmt.Sprintf("-%s", j.info.Type)

	j.process = &Process{
		name:    j.name,
		cmd:     exec.Command(mapper.Bin, local, remote, mapType),
		file:    j.file,
		errChan: make(chan error),
		stop:    make(chan struct{}),
		wg:      j.wg,
	}

	j.process.cmd.Stdout = j.file
	go j.process.Begin()
}

func (j *Job) Stop() {
	if j.process != nil {
		j.process.Stop()
	}
	j.wg.Wait()
}
