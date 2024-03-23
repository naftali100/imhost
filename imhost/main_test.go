package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zenizh/go-capturer"
)

//nolint:paralleltest // disable parallel test due to changing global state
func Test_main_server_not_running(t *testing.T) {
	// Backup and defer restore the original osExit function
	oldOsExit := osExit

	t.Cleanup(func() {
		osExit = oldOsExit
	})

	// Backup and defer restore args
	oldOsArgs := os.Args

	t.Cleanup(func() {
		os.Args = oldOsArgs
	})

	var exitCode int // capture the exit code

	// Mock osExit to panic instead of exiting
	osExit = func(code int) {
		exitCode = code
		fmt.Fprintf(os.Stderr, "os.Exit called with code: %d", code)
		panic("paniced instead of exiting")
	}

	// Mock os.Args
	os.Args = []string{
		t.Name(),
		"--ping", // This will cause the healthCheck function to exit
	}

	out := capturer.CaptureOutput(func() {
		require.Panics(t, func() {
			main()
		}, "Expected os.Exit to be called")
	})

	expectCode := FAILURE
	require.Equal(t, expectCode, exitCode, "incase of --ping flag, expected exit code to be 1")

	expectContain := "os.Exit called with code: 1"
	require.Contains(t, out, expectContain, "Expected output to contain: %s", expectContain)
}

//nolint:paralleltest // disable parallel test due to changing global state
func Test_healthCheck_golden(t *testing.T) {
	// Backup and defer restore args
	oldOsArgs := os.Args

	t.Cleanup(func() {
		os.Args = oldOsArgs
	})

	// Mock os.Args
	os.Args = []string{
		t.Name(),
		"", // This will cause the healthCheck function to do nothing (just return)
	}

	out := capturer.CaptureOutput(func() {
		require.NotPanics(t, func() {
			healthCheck("")
		}, "Expected os.Exit not to be called")
	})

	require.Empty(t, out, "Expected no output")
}

//nolint:paralleltest // disable parallel test due to changing global state
func Test_healthCheck_url_contains_control_character(t *testing.T) {
	// Backup and defer restore the original osExit function
	oldOsExit := osExit

	t.Cleanup(func() {
		osExit = oldOsExit
	})

	// Backup and defer restore args
	oldOsArgs := os.Args

	t.Cleanup(func() {
		os.Args = oldOsArgs
	})

	// Mock os.Args
	os.Args = []string{
		t.Name(),
		"--ping", // This will cause the healthCheck function to exit
	}

	var exitCode int // capture the exit code

	// Mock osExit to panic instead of exiting
	osExit = func(code int) {
		exitCode = code
		fmt.Fprintf(os.Stderr, "os.Exit called with code: %d", code)
		panic("paniced instead of exiting")
	}

	out := capturer.CaptureOutput(func() {
		require.Panics(t, func() {
			invalidURL := "http://localhost:8080/\x7f"
			healthCheck(invalidURL)
		}, "Expected os.Exit to be called")
	})

	expectCode := FAILURE
	require.Equal(t, expectCode, exitCode,
		"incase of --ping flag, if the URL contains control character, expected exit code to be 1")

	expectContain := "invalid control character in URL"
	require.Contains(t, out, expectContain,
		"Expected output to contain: \"%s\"", expectContain)
}

//nolint:paralleltest // disable parallel test due to changing global state
func Test_healthCheck_ping_golden(t *testing.T) {
	// Backup and defer restore the original osExit function
	oldOsExit := osExit

	t.Cleanup(func() {
		osExit = oldOsExit
	})

	// Backup and defer restore args
	oldOsArgs := os.Args

	t.Cleanup(func() {
		os.Args = oldOsArgs
	})

	// Start a test server
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer testSrv.Close()

	var exitCode int // capture the exit code

	// Mock osExit to panic instead of exiting
	osExit = func(code int) {
		exitCode = code
		fmt.Fprintf(os.Stderr, "os.Exit called with code: %d", code)
		panic("paniced instead of exiting")
	}

	// Mock os.Args
	os.Args = []string{
		t.Name(),
		"--ping", // This will cause the healthCheck function to exit
	}

	out := capturer.CaptureOutput(func() {
		require.Panics(t, func() {
			healthCheck(testSrv.URL)
		}, "Expected os.Exit to be called")
	})

	expectCode := SUCCESS
	require.Equal(t, expectCode, exitCode,
		"incase of --ping flag and server running, expected exit code to be 0")

	expectContain := "os.Exit called with code: 0"
	require.Contains(t, out, expectContain,
		"Expected output to contain: %s", expectContain)
}

//nolint:paralleltest // disable parallel test due to changing global state
func Test_defaultHandler_golden(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err, "Failed to create a new request")

	response := httptest.NewRecorder()
	handler := new(defaultHandler)

	// Test the handler
	require.NotPanics(t, func() {
		handler.ServeHTTP(response, request)
	}, "Expected ServeHTTP to not panic")

	expectContain := "Hello from host:"
	require.Contains(t, response.Body.String(), expectContain,
		"Expected response body to contain: \"%s\"", expectContain)
}

func Test_spawnServer_fail_to_start(t *testing.T) {
	t.Parallel()

	err := spawnServer("unknown://foo", nil)
	require.Error(t, err, "malformed address should return an error")

	expectContain := "unknown port"
	assert.Contains(t, err.Error(), expectContain, "Expected error message to contain: \"address\"")
}

//nolint:paralleltest // disable parallel test due to changing global state
func Test_defaultHandler_fail_to_get_hostname(t *testing.T) {
	// Backup and defer restore the original osHostname function
	oldOsHostname := osHostname

	t.Cleanup(func() {
		osHostname = oldOsHostname
	})

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err, "Failed to create a new request")

	response := httptest.NewRecorder()
	handler := new(defaultHandler)

	// Mock osHostname to return an error
	osHostname = func() (string, error) {
		return "", errors.New("forced error")
	}

	// Test the handler
	require.NotPanics(t, func() {
		handler.ServeHTTP(response, request)
	}, "Expected ServeHTTP to not panic")

	respBody := response.Body.String()

	assert.Contains(t, respBody, "failed to get hostname. error: forced error",
		"Expected response body to contain error message")
	assert.Contains(t, respBody, "forced error",
		"message should contain the error message")
	require.Equal(t, http.StatusInternalServerError, response.Code,
		"Expected response code to be 500")
}

func Test_defaultHandler_fail_to_write_response(t *testing.T) {
	t.Parallel()

	response := new(DumyResponseWriter)
	handler := new(defaultHandler)

	// Mock osHostname to return an error
	osHostname = func() (string, error) {
		return "", errors.New("forced error")
	}

	// Test the handler
	require.Panics(t, func() {
		handler.ServeHTTP(response, nil)
	}, "Expected ServeHTTP to panic")
}

// ----------------------------------------------------------------------------
//  Helper Functions
// ----------------------------------------------------------------------------

type DumyResponseWriter struct{}

// Header returns a nil Header.
func (DumyResponseWriter) Header() http.Header {
	return nil
}

// Write always returns an error with the message "forced error".
func (DumyResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("forced error")
}

// WriteHeader does nothing.
func (DumyResponseWriter) WriteHeader(int) {}
