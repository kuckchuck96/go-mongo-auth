// server package is referenced from
// 'Graceful shutdown or restart' section
// of https://pkg.go.dev/github.com/gin-gonic
package serve

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func listenAndServe(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %s\n", err)
	}
}

// ListenAndServe tasks custom http server and a request wait time
// to context is used to inform the server it has x seconds to finish
func ListenAndServe(server *http.Server, waitTime time.Duration) {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go listenAndServe(server)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has x seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
