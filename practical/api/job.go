package api

import (
	"os/exec"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Executor interface {
	// Execute runs the job's command
	Execute()
}

type Job struct {
	ID        string     `json:"id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string     `fake:"{sentence:3}" json:"name" example:"Job Name"`
	Created   time.Time  `json:"created" example:"2021-07-01T00:00:00Z"`
	Command   string     `json:"command,omitempty" example:"echo 'Hello, World!'"`
	Output    string     `json:"output,omitempty" example:"Hello, World!"`
	Completed *time.Time `json:"completed,omitempty" example:"2021-07-01T00:00:00Z"`
	Success   bool       `json:"success,omitempty" example:"true"`
	Executed  bool       `json:"executed,omitempty" example:"true"`
	Error     string     `json:"error,omitempty" example:"error message"`
}

// NewJob creates a new Job
func NewJob(name string, command string) (*Job, error) {
	var job Job

	err := gofakeit.Struct(&job)
	if err != nil {
		return nil, err
	}

	job.Created = time.Now()
	job.ID = gofakeit.UUID()
	job.Command = command
	job.Executed = false

	if name != "" {
		job.Name = name
	}

	return &job, nil
}

func (j *Job) Execute() {
	time.Sleep(time.Duration(gofakeit.Number(1, 20)) * time.Second)

	completionTime := time.Now()

	if j.Command != "" {
		j.Success = true
		j.Executed = true
		cmd := exec.Command(j.Command)
		out, err := cmd.CombinedOutput()
		if err != nil {
			j.Error = err.Error()
			j.Success = false
		}

		j.Output = string(out)
		j.Completed = &completionTime
		return
	}

	j.Success = gofakeit.Bool()
	j.Completed = &completionTime
	j.Output = gofakeit.HackerPhrase()
	j.Executed = true

	if !j.Success {
		err := gofakeit.Error()
		j.Error = err.Error()
	}
}
