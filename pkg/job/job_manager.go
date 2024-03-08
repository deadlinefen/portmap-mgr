package job

import (
	"fmt"
	"os"
	"sync"

	"github.com/deadlinefen/portmap-mgr/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type IJobManager interface {
	AddJobs(jobs map[string]config.Job) error
	Start()
	Stop()
}

type JobManager struct {
	mapper *config.Mapper
	jobs   []Job
	ipChan chan string
	stop   chan struct{}
}

type IJobManagerFactory interface {
	NewJobManager(ipChan chan string) IJobManager
}

type JobManagerFactory struct {
	mapper *config.Mapper
}

func (jmf *JobManagerFactory) NewJobManager(ipChan chan string) IJobManager {
	return &JobManager{
		mapper: jmf.mapper,
		jobs:   []Job{},
		ipChan: ipChan,
		stop:   make(chan struct{}),
	}
}

func NewJobManagerFactory(mapper *config.Mapper) IJobManagerFactory {
	return &JobManagerFactory{mapper: mapper}
}

func (jm *JobManager) AddJobs(jobs map[string]config.Job) error {
	for name, info := range jobs {
		log.Infof("Add job %s", name)
		filePath := fmt.Sprintf("%s/%d-%d.log", jm.mapper.FileDir, info.FromPort, info.ToPort)
		file, err := jm.createLogFile(filePath)
		if err != nil {
			return errors.Wrapf(err, "create log file failed")
		}
		jm.jobs = append(jm.jobs, Job{name: name, info: &info, file: file,wg: &sync.WaitGroup{}})
	}

	return nil
}

func (jm *JobManager) createLogFile(logFile string) (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "create log file %s failed", logFile)
	}

	return file, nil
}

func (jm *JobManager) stopAll() {
	for i := range jm.jobs {
		jm.jobs[i].Stop()
		log.Infof("job %s stopped.", jm.jobs[i].name)
	}
}

func (jm *JobManager) runAll(ipv6 string) {
	for i := range jm.jobs {
		log.Infof("job %s run with ip: %s", jm.jobs[i].name, ipv6)
		jm.jobs[i].Run(jm.mapper, ipv6)
	}
}

func (jm *JobManager) Start() {
	for {
		select {
		case <-jm.stop:
			jm.stopAll()
			return
		case newIpv6 := <-jm.ipChan:
			log.Infof("Get a new ipv6 addr: %s", newIpv6)
			jm.runAll(newIpv6)
		}
	}
}

func (jm *JobManager) Stop() {
	close(jm.stop)
}
