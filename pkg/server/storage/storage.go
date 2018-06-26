package storage

import (
	"github.com/yiqinguo/armyant/pkg/models"
	//"github.com/influxdata/influxdb/client/v2"
)

type Storage interface {
	AddJob(job *models.Job) error
	GetJob(id string) (*models.Job, error)
	ListJobs() ([]models.Job, error)
	DeleteJob(id string) error
	UpdateJob(job *models.Job) error
	AddResult(stats *models.Stats) error
	GetJobResult(jobId, instanceId string) (*models.JobResultResponse, error)
}
