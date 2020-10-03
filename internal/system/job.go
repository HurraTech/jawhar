package system

import "github.com/google/uuid"

var jobs map[string]*Job

type Job interface {
	SetID(string)
	GetID() string
	Run()
	Status() string
}

func RunJob(j Job) (string, error) {
	jobId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	j.SetID(jobId.String())
	jobs[jobId.String()] = &j
	go j.Run()
	return jobId.String(), nil
}

func JobStatus(uuid string) string {
	return (*jobs[uuid]).Status()
}
