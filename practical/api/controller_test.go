package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetJobs(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	expectedJobs := []*Job{}
	var actualJobs []Job

	for i := 0; i < 10; i++ {
		job, err := NewJob("", "")
		if err != nil {
			t.Error(err)
		}

		expectedJobs = append(expectedJobs, job)
	}

	c := NewController(expectedJobs, queuer)
	r.GET("/jobs", c.ListJobs)

	req := httptest.NewRequest("GET", "/jobs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	json.NewDecoder(resp.Body).Decode(&actualJobs)

	if len(actualJobs) != len(expectedJobs) {
		t.Errorf("Expected 1 job in queue, got %d", len(actualJobs))
	}
}

func TestGetJob(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	var actual Job

	job, err := NewJob("", "")
	if err != nil {
		t.Error(err)
	}

	c := NewController([]*Job{job}, queuer)
	r.GET("/jobs/:id", c.ShowJob)

	req := httptest.NewRequest("GET", fmt.Sprintf("/jobs/%s", job.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()

	json.NewDecoder(resp.Body).Decode(&actual)

	if actual.ID == "" {
		t.Error("Expected job ID to be set")
	}
}

func TestUpdateJob(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	var actualJob Job
	expectedJob := Job{}

	job, err := NewJob("", "")
	if err != nil {
		t.Error(err)
	}

	expectedJob.ID = job.ID
	expectedJob.Name = "woop woop"

	c := NewController([]*Job{job}, queuer)
	r.PATCH("/jobs/:id", c.UpdateJob)

	jb, err := json.Marshal(expectedJob)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/jobs/%s", job.ID), bytes.NewBuffer(jb))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()

	json.NewDecoder(resp.Body).Decode(&actualJob)

	if actualJob.ID != expectedJob.ID {
		t.Error("Expected job ID to be, but got", expectedJob.ID, actualJob.ID)
	}

	if actualJob.Name != expectedJob.Name {
		t.Errorf("Expected job name to be %s, got %s", expectedJob.Name, actualJob.Name)
	}
}

func TestDeleteJob(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	jobs := []*Job{}
	queueLength := 10
	lastID := ""

	for i := 0; i < queueLength; i++ {
		job, err := NewJob("", "")
		if err != nil {
			t.Error(err)
		}
		jobs = append(jobs, job)
		lastID = job.ID
	}

	c := NewController(jobs, queuer)
	r.DELETE("/jobs/:id", c.DeleteJob)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/jobs/%s", lastID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if len(c.jobs) != (queueLength - 1) {
		t.Errorf("Expected %d jobs in queue, got %d", (queueLength - 1), len(c.jobs))
	}
}

func TestGetQueueEmpty(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	c := NewController([]*Job{}, queuer)
	var queue []Job

	r.GET("/queue", c.ShowQueue)

	req := httptest.NewRequest("GET", "/queue", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	json.NewDecoder(resp.Body).Decode(&queue)

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if len(queue) != 0 {
		t.Errorf("Expected 0 jobs in queue, got %d", len(queue))
	}
}

func TestGetQueue(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	expectedJobs := []*Job{}
	var actualJobs []Job

	for i := 0; i < 10; i++ {
		job, err := NewJob("", "")
		if err != nil {
			t.Error(err)
		}

		expectedJobs = append(expectedJobs, job)
		queuer.Enqueue(job)
	}

	c := NewController(expectedJobs, queuer)
	r.GET("/queue", c.ShowQueue)

	req := httptest.NewRequest("GET", "/queue", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	json.NewDecoder(resp.Body).Decode(&actualJobs)

	if len(actualJobs) != len(expectedJobs) {
		t.Errorf("Expected %d jobs in queue, got %d", len(expectedJobs), len(actualJobs))
	}
}

func TestEnqueue(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()

	job, err := NewJob("", "")
	if err != nil {
		t.Error(err)
	}

	c := NewController([]*Job{job}, queuer)
	r.PUT("/queue", c.Enqueue)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/queue/%s", job.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	for _, j := range c.queuer.Queued() {
		if j.ID != job.ID {
			t.Errorf("Expected job ID to be %s, got %s", job.ID, j.ID)
		}
	}
}

func TestClearQueue(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	jobs := []*Job{}
	queueLength := 10

	for i := 0; i < queueLength; i++ {
		job, err := NewJob("", "")
		if err != nil {
			t.Error(err)
		}
		jobs = append(jobs, job)
		queuer.Enqueue(job)
	}

	c := NewController(jobs, queuer)
	r.DELETE("/queue", c.ClearQueue)

	req := httptest.NewRequest("DELETE", "/queue", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if len(c.queuer.Queued()) != 0 {
		t.Errorf("Expected 0 jobs in queue, got %d", len(c.queuer.Queued()))
	}
}

func TestDequeue(t *testing.T) {
	r := setupRouter()
	queuer := NewQueue()
	jobs := []*Job{}
	queueLength := 10
	lastID := ""

	for i := 0; i < queueLength; i++ {
		job, err := NewJob("", "")
		if err != nil {
			t.Error(err)
		}
		jobs = append(jobs, job)
		queuer.Enqueue(job)
		lastID = job.ID
	}

	c := NewController(jobs, queuer)
	r.DELETE("/queue/:job_id", c.Dequeue)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/queue/%s", lastID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if len(c.queuer.Queued()) != (queueLength - 1) {
		t.Errorf("Expected %d jobs in queue, got %d", (queueLength - 1), len(c.queuer.Queued()))
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}
