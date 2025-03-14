package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	jobs   []*Job
	queuer Queuer
}

// NewController creates a new Controller
func NewController(jobs []*Job, queuer Queuer) *Controller {
	return &Controller{
		jobs:   jobs,
		queuer: queuer,
	}
}

// ShowJob Shows a Job by ID
//
//	@Summary		Show a job
//	@Description	Get a job by ID
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Job ID"
//	@Success		200	{object}	Job
//	@Failure		404	{object}	map[string]string
//	@Router			/jobs/{id} [get]
func (c *Controller) ShowJob(ctx *gin.Context) {
	id := ctx.Param("id")

	var job *Job
	for _, j := range c.jobs {
		if j.ID == id {
			job = j
			break
		}
	}

	if job == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	ctx.JSON(http.StatusOK, job)
}

// ListJobs Lists all Jobs
//
//	@Summary		List jobs
//	@Description	Get all jobs
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	Job
//	@Router			/jobs [get]
func (c *Controller) ListJobs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.jobs)
}

// AddJob Adds a Job
//
//	@Summary		Add a job
//	@Description	Add a new job
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			job	body		Job	true	"Job object"
//	@Success		201	{object}	Job
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/jobs [post]
func (c *Controller) AddJob(ctx *gin.Context) {
	var job Job

	if err := ctx.ShouldBindJSON(&job); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newJob, err := NewJob(job.Name, job.Command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.jobs = append(c.jobs, newJob)

	ctx.JSON(http.StatusCreated, newJob)
}

// UpdateJob Updates a Job
//
//	@Summary		Update a job
//	@Description	Update a job by ID
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Job ID"
//	@Param			job	body		Job		true	"Add job"
//	@Success		200	{object}	Job
//	@Failure		204	{object}	map[string]string
//	@Router			/jobs/{id} [patch]
func (c *Controller) UpdateJob(ctx *gin.Context) {
	id := ctx.Param("id")
	var job Job

	if err := ctx.ShouldBindJSON(&job); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingJob Job
	for i, j := range c.jobs {
		if j.ID == id {
			existingJob = job
			c.jobs[i] = &job
			break
		}
	}

	if existingJob.ID == "" {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "job not found"})
		return
	}

	ctx.JSON(http.StatusOK, existingJob)
}

// DeleteJob Deletes a job
//
//	@Summary		Delete a job
//	@Description	Delete a job by ID
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Job ID"
//	@Success		202
//	@Failure		204	{object}	map[string]string
//	@Router			/jobs/{id} [delete]
func (c *Controller) DeleteJob(ctx *gin.Context) {
	id := ctx.Param("id")

	var job Job
	for i, j := range c.jobs {
		if j.ID == id {
			job = *j
			c.jobs = append(c.jobs[:i], c.jobs[i+1:]...)
			break
		}
	}

	if job.ID == "" {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "job not found"})
		return
	}

	ctx.JSON(http.StatusAccepted, nil)
}

// ShowQueue Shows the queue
//
//	@Summary		Show the queue
//	@Description	Get all jobs in the queue
//	@Tags			queue
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	Job
//	@Router			/queue [get]
func (c *Controller) ShowQueue(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.queuer.Queued())
}

// Enqueue Queues a job for execution
//
//	@Summary		Enqueue a job
//	@Description	Add a job to the queue
//	@Tags			queue
//	@Accept			json
//	@Produce		json
//	@Param			job_id	path		string	true	"Job ID"
//	@Success		201		{array}		Job
//	@Failure		204		{object}	map[string]string
//	@Router			/queue [put]
func (c *Controller) Enqueue(ctx *gin.Context) {
	id := ctx.Param("job_id")

	var job *Job
	for _, j := range c.jobs {
		if j.ID == id {
			job = j
			break
		}
	}

	if job == nil {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "job not found"})
		return
	}

	c.queuer.Enqueue(job)

	ctx.JSON(http.StatusCreated, c.queuer.Queued())
}

// ClearQueue Clears the queue
//
//	@Summary		Clear the queue
//	@Description	Remove all jobs from the queue
//	@Tags			queue
//	@Accept			json
//	@Produce		json
//	@Success		202
//	@Router			/queue [delete]
func (c *Controller) ClearQueue(ctx *gin.Context) {
	c.queuer.Clear()

	ctx.JSON(http.StatusAccepted, nil)
}

// Dequeue Removes a job from the queue by ID and returns the Queue
//
//	@Summary		Dequeue a job
//	@Description	Remove a job from the queue
//	@Tags			queue
//	@Accept			json
//	@Produce		json
//	@Param			job_id	path		string	true	"Job ID"
//	@Success		202		{array}		Job
//	@Failure		204		{object}	map[string]string
//	@Router			/queue/{job_id} [delete]
func (c *Controller) Dequeue(ctx *gin.Context) {
	id := ctx.Param("job_id")

	var job *Job
	for _, j := range c.queuer.Queued() {
		if j.ID == id {
			job = j
			break
		}
	}

	if job == nil {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "job not found"})
		return
	}

	c.queuer.Dequeue(job)

	ctx.JSON(http.StatusAccepted, c.queuer.Queued())
}
