package api

type Queue struct {
	Jobs []Job `json:"jobs"`
}

func NewQueue() *Queue {
	return &Queue{
		Jobs: []Job{},
	}
}

func (q *Queue) Enqueue(job Job) {
	q.Jobs = append(q.Jobs, job)
	go func() {
		job.Execute()
		defer q.removeJob(job)
	}()
}

func (q *Queue) removeJob(job Job) {
	for i, j := range q.Jobs {
		if j.ID == job.ID {
			q.Jobs = append((q.Jobs)[:i], (q.Jobs)[i+1:]...)
			break
		}
	}
}
