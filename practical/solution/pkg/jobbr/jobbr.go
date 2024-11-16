package jobbr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Jobber interface {
	// Jobs returns all jobs
	Jobs() (*[]Job, error)
	// Job returns a job by id
	Job(id string) (*Job, error)
	// CreateJob creates a new job
	CreateJob(job *Job) error
	// Queue returns all jobs in the queue
	Queue() (*Queue, error)
	// Enqueue adds a job to the queue
	Enqueue(job *Job) error
}

type Job struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Created   time.Time  `json:"created"`
	Command   string     `json:"command,omitempty"`
	Output    string     `json:"output,omitempty"`
	Completed *time.Time `json:"completed,omitempty"`
	Success   bool       `json:"success,omitempty"`
	Executed  bool       `json:"executed,omitempty"`
	Error     string     `json:"error,omitempty"`
}

type Queue struct {
	Jobs []Job `json:"jobs"`
}

type client struct {
	env *EnvConfig
}

func NewJobbr(env *EnvConfig) Jobber {
	return &client{env: env}
}

func (c *client) Jobs() (*[]Job, error) {
	resp, err := http.Get(fmt.Sprintf("%s:%d/jobs", c.env.URL, c.env.Port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jobs []Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return &jobs, nil
}

func (c *client) Job(id string) (*Job, error) {
	resp, err := http.Get(fmt.Sprintf("%s:%d/jobs/%s", c.env.URL, c.env.Port, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var job Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *client) CreateJob(job *Job) error {
	j, err := json.Marshal(job)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s:%d/jobs/new", c.env.URL, c.env.Port), bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *client) Queue() (*Queue, error) {
	resp, err := http.Get(fmt.Sprintf("%s:%d/jobs/queue", c.env.URL, c.env.Port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var queue Queue
	if err := json.NewDecoder(resp.Body).Decode(&queue); err != nil {
		return nil, err
	}

	return &queue, nil
}

func (c *client) Enqueue(job *Job) error {
	j, err := json.Marshal(job)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s:%d/jobs/enqueue", c.env.URL, c.env.Port), bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
