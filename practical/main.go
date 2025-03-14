package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wnoonan/gostuff/practical/api"
	_ "github.com/wnoonan/gostuff/practical/docs"
)

// main
//
//	@title						Jobs & Queue API
//	@version					1.0
//	@description				This is a simple job and queue api.
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	var jobs []*api.Job
	for i := 0; i < 10; i++ {
		job, err := api.NewJob("", "")
		if err != nil {
			panic(err)
		}

		jobs = append(jobs, job)
	}

	queue := api.NewQueue()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	c := api.NewController(jobs, queue)

	v1 := r.Group("/api/v1")
	{
		jobs := v1.Group("/jobs")
		{
			jobs.GET("", c.ListJobs)
			jobs.GET(":id", c.ShowJob)
			jobs.POST("", c.AddJob)
			jobs.PATCH(":id", c.UpdateJob)
			jobs.DELETE(":id", c.DeleteJob)
		}
		queue := v1.Group("/queue")
		{
			queue.GET("", c.ShowQueue)
			queue.PUT(":job_id", c.Enqueue)
			queue.DELETE("", c.ClearQueue)
			queue.DELETE(":job_id", c.Dequeue)
		}
	}
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
