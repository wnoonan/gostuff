package api

import (
	"os/exec"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Jobber interface {
	Execute() error
}

type Job struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `fake:"{sentence:3}" json:"name"`
	Created   time.Time  `json:"created"`
	Command   string     `json:"command,omitempty"`
	Output    string     `json:"output,omitempty"`
	Completed *time.Time `json:"completed,omitempty"`
	Success   bool       `json:"success,omitempty"`
	Executed  bool       `json:"executed,omitempty"`
	Error     string     `json:"error,omitempty"`
}

func NewJob() (*Job, error) {
	var job Job

	err := gofakeit.Struct(&job)
	if err != nil {
		return nil, err
	}
	job.Created = time.Now()
	job.ID = gofakeit.UUID()

	return &job, nil
}

func (j *Job) Execute() {
	time.Sleep(time.Duration(gofakeit.Number(5, 60)) * time.Second)

	j.Executed = true

	if j.Command != "" {
		cmd := exec.Command(j.Command)
		out, err := cmd.Output()
		if err != nil {
			j.Error = err.Error()
		}
		j.Output = string(out)
		return
	}

	completionTime := time.Now()
	j.Success = gofakeit.Bool()
	j.Completed = &completionTime
	j.Output = gofakeit.HackerPhrase()

	if !j.Success {
		j.Error = gofakeit.Error().Error()
	}

}
