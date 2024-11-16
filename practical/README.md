# Welcome

Jobs & Queues

### We'd like you to:
* Display jobs that haven't been queued
* Allow choosing jobs to queue
* Queue jobs, and display the queue
* Display the status, associated errors, created and completion times from the jobs once they've completed
* Allow us to create and queue a job

### Things to consider:
* Jobs can take some time to run, we want to display completed jobs as they're completed and not wait on the entire queue to finish.

This can be a web application, an API, or a command line tool -- It's totally up to you!

# Getting Started

## Prerequisites
* Docker

## Starting the API Server
run `./run.sh` 

# API Reference

## Endpoint
`localhost:8080`

## Jobs

### Job
To get a job, replace `{id}` with a job id and use the following `curl` command:

```sh
curl -X GET http://localhost:8080/jobs/{id} -H "Content-Type: application/json'
```

### Response Codes
`404` Job Not Found

`200` Job

### Jobs
To get all jobs, use the following `curl` command:

```sh
curl -X GET http://localhost:8080/jobs 
```
### Response Codes
`200` Jobs

### Create Job
To create a job, use the following `curl` command:

```sh
curl -X POST http://localhost:8080/jobs/new -d '{"name": "jobby", "command": "echo `Hello, World!`"}
```

### Response Codes
`200` Job

### Queue a Job

To queue a job, use the following `curl` command:

```sh
curl -X PUT http://localhost:8080/jobs/enqueue -H "Content-Type: application/json" -d '{"id": 1}'
```
### Response Codes
`404` Job Not Found
`200` Job Queued Successfully


### View Job Queue
To get the queue, use the following `curl` command:

```sh
curl -X GET http://localhost:8080/jobs/queue
```
### Response Codes
`200` Queue
