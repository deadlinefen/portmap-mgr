package job

import (
	"sync"

	"github.com/deadlinefen/portmap-mgr/pkg/config"
)

type IJobManager interface {
	AddJobs(jobs map[string]config.Job)
	Run()
	Stop()
}

type JobManager struct {
	mapper *config.Mapper
	jobs   []Job
	ipChan chan string
	stop   chan struct{}
}

type IJobManagerFactory interface {
	NewJobManager(mapper *config.Mapper, ipChan chan string) IJobManager
}

type JobManagerFactory struct {
}

func (jmf *JobManagerFactory) NewJobManager(mapper *config.Mapper, ipChan chan string) IJobManager {
	return &JobManager{
		mapper: mapper,
		jobs:   []Job{},
		ipChan: ipChan,
		stop:   make(chan struct{}),
	}
}

func NewJobManagerFactory() IJobManagerFactory {
	return &JobManagerFactory{}
}

func (jm *JobManager) AddJobs(jobs map[string]config.Job) {
	for name, info := range jobs {
		jm.jobs = append(jm.jobs, Job{name: name, info: &info, wg: &sync.WaitGroup{}})
	}
}

func (jm *JobManager) stopAll() {
	for i := range jm.jobs {
		jm.jobs[i].Stop()
	}
}

func (jm *JobManager) runAll(ipv6 string) {
	for i := range jm.jobs {
		jm.jobs[i].Run(jm.mapper, ipv6)
	}
}

func (jm *JobManager) Run() {
	for {
		select {
		case <-jm.stop:
			jm.stopAll()
			return
		case newIpv6 := <-jm.ipChan:
			jm.runAll(newIpv6)
		}
	}
}

func (jm *JobManager) Stop() {
	close(jm.stop)
}
