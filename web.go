package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyesarf/tracy/tracers"
)

func RunWeb() *http.Server {
	gin.SetMode(gin.ReleaseMode)

	// Disable Gin's default logger output
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	var trace tracers.Trace
	router.GET("/getTrace", func(c *gin.Context) {

		c.JSON(http.StatusOK, trace)
	})

	router.POST("/sendTrace", func(c *gin.Context) {
		if err := c.BindJSON(&trace); err != nil {
			return
		}

		c.Status(http.StatusOK)
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"content": "main page"})
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

	go func() {
		server := RunWeb()
		log.Println("Server listening on port 1337")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Start binary tracing in a separate goroutine
	go func() {
		args := os.Args

		var trace tracers.Trace
		if len(args) >= 2 {
			trace.Binary = args[1]
			trace.Args = args[2:]
		} else {
			trace.Binary = "tmp/a.out"
		}

		trace.TraceBin()
	}()
	select {}
}
