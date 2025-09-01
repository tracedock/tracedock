package server

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewSupervisor(t *testing.T) {
	t.Run("should correctly create a Supervisor", func(t *testing.T) {
		o := NewSupervisor()

		assert.NotNil(t, o)
		assert.Equal(t, Stopped, o.state)
		assert.NotNil(t, o.servers)
		assert.Empty(t, o.servers)
	})
}

func Test_Supervisor_Add(t *testing.T) {
	t.Run("should add servers to Supervisor", func(t *testing.T) {
		var superv = NewSupervisor()
		var addrs = []string{":7070", ":8080", ":9090"}

		for i, addr := range addrs {
			superv.Add(addr, NewMockServer(t))

			assert.Len(t, superv.servers, i+1)
		}

		for key, val := range superv.servers {
			assert.Equal(t, val, superv.servers[key])
		}
	})

}

func Test_Supervisor_Run(t *testing.T) {
	tests := []struct {
		name          string
		currentState  State
		desiredState  State
		servers       []string
		expectedError error
	}{
		{
			name:          "should return no error when isn't running",
			currentState:  Stopped,
			desiredState:  Running,
			servers:       []string{":8080"},
			expectedError: nil,
		},
		{
			name:          "should return error when already is running",
			currentState:  Running,
			desiredState:  Running,
			servers:       []string{":8080", ":9090"},
			expectedError: ErrAlreadyRunning,
		},
		{
			name:          "should return error when have no servers to start",
			currentState:  Stopped,
			desiredState:  Stopped,
			servers:       []string{},
			expectedError: ErrEmptyServerList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewSupervisor()

			for _, addr := range tt.servers {
				mockServer := NewMockServer(t)
				mockServer.EXPECT().Start(addr).Return(nil).Maybe()
				o.Add(addr, mockServer)
			}

			o.state = tt.currentState

			err := o.Run()

			assert.Equal(t, tt.desiredState, o.state)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func Test_Supervisor_Wait(t *testing.T) {
	t.Run("returns error when not running", func(t *testing.T) {
		o := NewSupervisor()

		err := o.Wait()

		assert.Error(t, err)
		assert.Equal(t, ErrNotRunning, err)
		assert.Equal(t, Stopped, o.state)
	})

	t.Run("stops servers and sets state to stopped", func(t *testing.T) {
		o := NewSupervisor()
		mockServer1 := NewMockServer(t)
		mockServer2 := NewMockServer(t)
		addr1 := ":8080"
		addr2 := ":9090"

		mockServer1.EXPECT().Start(addr1).Return(nil).Maybe()
		mockServer2.EXPECT().Start(addr2).Return(nil).Maybe()
		mockServer1.EXPECT().Stop().Return(nil)
		mockServer2.EXPECT().Stop().Return(nil)

		o.Add(addr1, mockServer1)
		o.Add(addr2, mockServer2)
		o.Run()

		done := make(chan error, 1)
		go func() {
			done <- o.Wait()
		}()

		time.Sleep(10 * time.Millisecond)

		// Send SIGINT to trigger the Wait() method to complete
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGINT)

		err := <-done

		assert.NoError(t, err)
		assert.Equal(t, Stopped, o.state)
	})

	t.Run("returns error when server stop fails", func(t *testing.T) {
		o := NewSupervisor()
		mockServer := NewMockServer(t)
		addr := ":8080"
		expectedError := errors.New("stop error")

		mockServer.EXPECT().Start(addr).Return(nil).Maybe()
		mockServer.EXPECT().Stop().Return(expectedError)

		o.Add(addr, mockServer)
		o.Run()

		done := make(chan error, 1)
		go func() {
			done <- o.Wait()
		}()

		time.Sleep(10 * time.Millisecond)

		// Send SIGINT to trigger the Wait() method to complete
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGINT)

		err := <-done

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}
