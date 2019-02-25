package ready

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	clog "github.com/coredns/coredns/plugin/pkg/log"
)

func init() { clog.Discard() }

func TestReady(t *testing.T) {
	rd := new(":0")
	RegisterPlugin("test")
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err := rd.onStartup(); err != nil {
			t.Fatalf("Unable to startup the readiness server: %v", err)
		}
		wg.Done()
	}()
	wg.Wait()

	defer rd.onFinalShutdown()

	address := fmt.Sprintf("http://%s/ready", rd.ln.Addr().String())

	wg.Add(1)
	go func() {
		response, err := http.Get(address)
		if err != nil {
			t.Fatalf("Unable to query %s: %v", address, err)
		}
		if response.StatusCode != 503 {
			t.Errorf("Invalid status code: expecting %q, got %q'", 503, response.StatusCode)
		}
		response.Body.Close()
		wg.Done()
	}()

	wg.Wait()

	// Now, make it ready.
	Signal("test")
	rd.wait()

	response, err := http.Get(address)
	if err != nil {
		t.Fatalf("Unable to query %s: %v", address, err)
	}
	if response.StatusCode != 200 {
		t.Errorf("Invalid status code: expecting '200', got '%d'", response.StatusCode)
	}
	response.Body.Close()
}
