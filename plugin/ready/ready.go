// Package ready is used to signal readiness of the CoreDNS binary (it's a process wide thing).
package ready

import (
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin("ready")

var r = &m{m: make(map[string]int)}

// RegisterPlugin registers plugin to signal readiness. A plugin may be registered multiple times
// each of the instances needs to call Signal() is signal that it's ready.
func RegisterPlugin(plugin string) { r.Store(plugin) }

// Signal is called by each plugin that is ready. Once all plugins have called in the plugin will signal readiness
// by return 200 OK on the http handler. Every 10 seconds we log which plugins haven't called back yet.
func Signal(plugin string) { r.Delete(plugin) }

type ready struct {
	Addr string
	Ok   bool
	sync.RWMutex

	ln   net.Listener
	done bool
	mux  *http.ServeMux
}

// new returns a new initialized ready.
func new(addr string) *ready { return &ready{Addr: addr} }

// wait waits for all plugins that have been registered to call Ready.
func (rd *ready) wait() error {
	total, _ := r.Len()
	for {
		l, list := r.Len()
		if l == 0 {
			log.Infof("All plugins signal readiness")
			break
		}
		log.Infof("Waiting on %d/%d: %s", l, total, list)

		time.Sleep(1 * time.Second)
	}

	rd.setOK()
	return nil
}

func (rd *ready) setOK() {
	rd.Lock()
	defer rd.Unlock()
	rd.Ok = true
}

func (rd *ready) ok() bool {
	rd.RLock()
	defer rd.RUnlock()
	return rd.Ok
}

func (rd *ready) onStartup() error {
	if rd.Addr == "" {
		rd.Addr = defAddr
	}

	ln, err := net.Listen("tcp", rd.Addr)
	if err != nil {
		return err
	}

	rd.Lock()
	rd.ln = ln
	rd.mux = http.NewServeMux()
	rd.done = true
	rd.Unlock()

	rd.mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		if rd.ok() {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "OK")
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	go func() { http.Serve(rd.ln, rd.mux) }()

	return nil
}

func (rd *ready) onRestart() error { return rd.onFinalShutdown() }

func (rd *ready) onFinalShutdown() error {
	rd.Lock()
	defer rd.Unlock()
	if !rd.done {
		return nil
	}

	rd.ln.Close()
	rd.done = false
	return nil
}

const defAddr = ":8181"
