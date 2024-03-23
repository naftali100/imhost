package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	listenPort    = "80" // Listen port
	srvTimeoutSec = 20   // Server timeout
	FAILURE       = 1    // Exit code for failure
	SUCCESS       = 0    // Exit code for success
)

// Monkey patching for testing.
var (
	// osExit is a copy of os.Exit that can be mocked in tests.
	osExit = os.Exit
	// osHostname is a copy of os.Hostname that can be mocked in tests.
	osHostname = os.Hostname
)

// ----------------------------------------------------------------------------
//  Main
// ----------------------------------------------------------------------------

//nolint:forbidigo // allow fmt.Print* due to its nature
func main() {
	addr := "0.0.0.0:" + listenPort
	url := "http://" + addr
	handler := new(defaultHandler)

	// If the args contain --ping, do a health check and exit
	healthCheck(url)

	// Start server
	fmt.Println("* Starging up server ...")
	fmt.Printf("* Listening to: %s\n", url)

	log.Fatal(spawnServer(addr, handler))
}

// ----------------------------------------------------------------------------
//  Type: defaultHandler (implements http.Handler)
// ----------------------------------------------------------------------------

// defaultHandler is a "catch all" handler that responds with a message
// containing the hostname of the server.
type defaultHandler struct{}

// ServeHTTP writes a message containing the hostname of the server to the
// response writer.
// This is an implementation of the http.Handler interface.
func (h *defaultHandler) ServeHTTP(respWriter http.ResponseWriter, _ *http.Request) {
	var respMsg string

	status := http.StatusOK

	hostname, err := osHostname()
	if err != nil {
		respMsg = fmt.Sprintf("failed to get hostname. error: %s", err)
		status = http.StatusInternalServerError
	} else {
		respMsg = fmt.Sprintf("Hello from host: %s\n", hostname)
	}

	respWriter.WriteHeader(status)

	if _, err := respWriter.Write([]byte(respMsg)); err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// healthCheck exits if the --ping flag is set. Otherwise it will do nothing.
//
// If the --ping flag is set, it will:
// Make a request to the local server and check if it is running. It will exit
// with a status code of 0 if the server is running, and 1 if it is not.
func healthCheck(rawURL string) {
	if !strings.Contains(strings.Join(os.Args, " "), "--ping") {
		return
	}

	ctx, doCancel := context.WithTimeout(context.Background(), srvTimeoutSec*time.Second)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		doCancel()
		fmt.Fprintln(os.Stderr, err)
		osExit(FAILURE)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		doCancel()
		fmt.Fprintln(os.Stderr, err)
		osExit(FAILURE)
	}

	resp.Body.Close()
	doCancel()

	osExit(SUCCESS)
}

func spawnServer(addr string, handler http.Handler) error {
	srv := new(http.Server)

	srv.Addr = addr
	srv.Handler = handler
	srv.ReadHeaderTimeout = srvTimeoutSec * time.Second

	return errors.Wrap(srv.ListenAndServe(), "failed to start server")
}
