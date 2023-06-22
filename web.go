package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ilyesarf/tracy/trace"

	"github.com/gin-gonic/gin"
)

func RunWeb(trace trace.Trace) *http.Server {
	router := gin.Default()

	router.GET("/getTrace", func(c *gin.Context) {

		c.JSON(http.StatusOK, trace)
	})

	server := &http.Server{
		Addr:         ":1337",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server
}

func main() {
	args := os.Args

	var path string
	if len(args) == 2 {
		path = args[1]
	} else {
		path = "tmp/a.out"
	}

	var trace trace.Trace
	trace.Binary = path
	trace.TraceBin()

	server := RunWeb(trace)
	log.Println("Server listening on port 1337")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
