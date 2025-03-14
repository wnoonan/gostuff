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
	Jobs() ([]*Job, error)
	// Job returns a job by id
	Job(id string) (*Job, error)
	// CreateJob creates a new job
	CreateJob(job *Job) error
	// Queue returns all jobs in the queue
	Queue() ([]*Job, error)
	// Enqueue adds a job to the queue
	Enqueue(job *Job) ([]*Job, error)
	// Dequeue removes a job from the queue
	Dequeue(job *Job) ([]*Job, error)
	// ClearQueue clears the queue
	ClearQueue() error
	// DeleteJob deletes a job
	DeleteJob(job *Job) error
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

type client struct {
	env *EnvConfig
}

func NewJobbr(env *EnvConfig) Jobber {
	return &client{env: env}
}

func (c *client) Jobs() ([]*Job, error) {
	resp, err := http.Get(fmt.Sprintf("%s:%d/api/v1/jobs", c.env.URL, c.env.Port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jobs []*Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *client) Job(id string) (*Job, error) {
	resp, err := http.Get(fmt.Sprintf("%s:%d/api/v1/jobs/%s", c.env.URL, c.env.Port, id))
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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s:%d/api/v1/jobs", c.env.URL, c.env.Port), bytes.NewBuffer(j))
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

func (c *client) Queue() ([]*Job, error) {
	resp, err := http.Get(fmt.Sprintf("%s:%d/api/v1/queue", c.env.URL, c.env.Port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var queue []*Job
	if err := json.NewDecoder(resp.Body).Decode(&queue); err != nil {
		return nil, err
	}

	return queue, nil
}

func (c *client) Enqueue(job *Job) ([]*Job, error) {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s:%d/api/v1/queue/%s", c.env.URL, c.env.Port, job.ID), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var queue []*Job
	if err := json.NewDecoder(resp.Body).Decode(&queue); err != nil {
		return nil, err
	}

	return queue, nil

}

func (c *client) Dequeue(job *Job) ([]*Job, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s:%d/api/v1/queue/%s", c.env.URL, c.env.Port, job.ID), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var queue []*Job
	if err := json.NewDecoder(resp.Body).Decode(&queue); err != nil {
		return nil, err
	}

	return queue, nil
}

func (c *client) ClearQueue() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s:%d/api/v1/queue", c.env.URL, c.env.Port), nil)
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

func (c *client) DeleteJob(job *Job) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s:%d/api/v1/jobs/%s", c.env.URL, c.env.Port, job.ID), nil)
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
