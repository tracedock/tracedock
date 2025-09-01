package server

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// State is used to control the current
// state of the Supervisor
type State int

// Supervisor manages the lifecycle of multiple servers
type Supervisor struct {
	state   State
	servers map[string]Server
}

const (
	Stopped State = iota
	Running
)

var (
	// ErrNotRunning is returned when is some operation that requires
	// the supervisor to be running but it isn't
	ErrNotRunning = errors.New("supervisor is not running")

	// ErrAlreadyRunning is returned when an operation is attempted
	// on the supervisor while it is already running
	ErrAlreadyRunning = errors.New("supervisor is already running")

	// ErrEmptyServerList is returned when there are no servers to start
	// but we try to start anyway
	ErrEmptyServerList = errors.New("no servers to start")
)

// NewSupervisor creates a new Supervisor instance
func NewSupervisor() *Supervisor {
	return &Supervisor{state: Stopped, servers: make(map[string]Server)}
}

// Add maps a server to run on the given address
func (o *Supervisor) Add(addr string, s Server) {
	o.servers[addr] = s
}

// Run starts all the servers managed by the Supervisor
func (o *Supervisor) Run() error {
	if o.state == Running {
		return ErrAlreadyRunning
	}

	if len(o.servers) == 0 {
		return ErrEmptyServerList
	}

	err := make(chan error)

	for addr, srv := range o.servers {
		go func() {
			err <- srv.Start(addr)
		}()
	}

	o.state = Running

	return nil
}

// Wait blocks until an interrupt signal is received and stops all servers
func (o *Supervisor) Wait() error {
	log.Print("application is running, try CTRL+C to stop")

	var sigChan = make(chan os.Signal, 1)

	if o.state != Running {
		return ErrNotRunning
	}

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	for _, srv := range o.servers {
		if err := srv.Stop(); err != nil {
			return err
		}
	}

	o.state = Stopped

	return nil
}
