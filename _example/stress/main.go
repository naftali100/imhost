package main

import (
	"context"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	targetHost     = "127.0.0.1" // Localhost
	listenPort     = "8080"      // Listen port
	reqTimeoutSec  = 20          // Request timeout
	testTimeoutSec = 30          // Test timeout
	numScaledHost  = 3           // Number of hosts expected (default)
)

func main() {
	numScale := flag.Int("s", numScaledHost, "number to scale")
	flag.Parse()

	found := map[uint32]string{}
	timeBegin := time.Now()

	fmt.Printf("\tSearching for %d hosts.\n", *numScale)

	for {
		resp, err := doRequest()
		panicOnError(err)

		hashedHost := crc32.ChecksumIEEE(resp)
		if _, ok := found[hashedHost]; !ok {
			found[hashedHost] = string(resp)
		}

		numHostsFound := len(found)
		dots := strings.Repeat(".", numHostsFound)

		fmt.Printf("\tFound %d hosts %s\r", numHostsFound, dots)

		if len(found) == *numScale {
			fmt.Println("\n\tOK:found all hosts")
			break
		}

		if time.Since(timeBegin) > testTimeoutSec*time.Second {
			panicOnError(errors.New("timeout: exceeded to find all hosts"))
		}

		//time.Sleep(1 * time.Second)
	}

	fmt.Println("\tList responses:")
	index := 0
	for _, nameHost := range found {
		fmt.Printf("\t  #%02d: \"%s\"\n", index+1, strings.TrimSpace(nameHost))
		index++
	}
}

func doRequest() ([]byte, error) {
	req, doCancel, err := newRequest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	defer doCancel()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return response, nil
}

func newRequest() (*http.Request, func(), error) {
	rURL := new(url.URL)

	rURL.Scheme = "http"
	rURL.Host = targetHost + ":" + listenPort

	rawURL := rURL.String()
	ctx, doCancel := context.WithTimeout(context.Background(), reqTimeoutSec*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)

	return req, doCancel, errors.Wrap(err, "failed to create request")
}

// panicOnError panics if an error is not nil. Otherwise, it does nothing.
func panicOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
