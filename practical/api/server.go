package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gorilla/mux"
)

type Server struct {
	jobs  []Job
	queue *Queue
}

func NewServer(jobs []Job, queue *Queue) *Server {
	return &Server{
		jobs:  jobs,
		queue: queue,
	}
}

func (srvr *Server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/jobs/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$}", srvr.JobHandler).Methods("GET")
	r.HandleFunc("/jobs", srvr.JobHandler).Methods("GET")
	r.HandleFunc("/jobs/new", srvr.NewJobHandler).Methods("POST")
	r.HandleFunc("/jobs/queue", srvr.QueueHandler).Methods("GET")
	r.HandleFunc("/jobs/enqueue", srvr.EnqueueHandler).Methods("PUT")

	http.ListenAndServe(":8080", r)
}

func (srvr *Server) JobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(srvr.jobs)
		return
	}

	var job *Job
	for _, j := range srvr.jobs {
		if j.ID == vars["id"] {
			job = &j
			break
		}
	}

	if job == nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (srvr *Server) NewJobHandler(w http.ResponseWriter, r *http.Request) {
	var job Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	job.ID = gofakeit.UUID()
	job.Created = time.Now()
	job.Success = false
	job.Completed = nil

	srvr.jobs = append(srvr.jobs, job)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (srvr *Server) QueueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(srvr.queue)
}

func (srvr *Server) EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	var req Job
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var job *Job
	for _, j := range srvr.jobs {
		if j.ID == req.ID {
			job = &j
			break
		}
	}

	if job == nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	srvr.queue.Enqueue(*job)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(srvr.queue)
}
