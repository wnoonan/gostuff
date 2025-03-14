package api

type Queuer interface {
	// Enqueue adds a job to the queue for execution
	Enqueue(*Job)
	// Queued returns all jobs in the queue
	Queued() []*Job
	// Dequeue removes a job from the queue
	Dequeue(*Job)
	// Clear removes all jobs from the queue
	Clear()
}

type queue struct {
	jobs []*Job
}

func NewQueue() Queuer {
	return &queue{}
}

func (q *queue) Enqueue(job *Job) {
	q.jobs = append(q.jobs, job)
	go job.Execute()
}

func (q *queue) Queued() []*Job {
	return q.jobs
}

func (q *queue) Dequeue(job *Job) {
	for i, j := range q.jobs {
		if j.ID == job.ID {
			job.Executed = false
			q.jobs = append(q.jobs[:i], q.jobs[i+1:]...)
			break
		}
	}
}

func (q *queue) Clear() {
	for _, job := range q.jobs {
		job.Executed = false
	}

	q.jobs = []*Job{}
}
