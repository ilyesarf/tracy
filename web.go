package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ilyesarf/tracy/tracers"
)

var trace tracers.Trace

var upgrader = websocket.Upgrader{}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	trace.SendTrace(conn)

}

func RunWeb() *http.Server {
	gin.SetMode(gin.ReleaseMode)

	// Disable Gin's default logger output
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"content": "main page"})
	})

	router.GET("/ws", func(c *gin.Context) {
		handleWS(c.Writer, c.Request)
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
