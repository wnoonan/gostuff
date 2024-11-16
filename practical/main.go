package main

import "github.com/wnoonan/gostuff/practical/api"

func main() {
	var jobs []api.Job

	for i := 0; i < 10; i++ {
		job, err := api.NewJob()
		if err != nil {
			panic(err)
		}

		jobs = append(jobs, *job)
	}

	queue := api.NewQueue()
	srvr := api.NewServer(jobs, queue)

	srvr.Start()
}
