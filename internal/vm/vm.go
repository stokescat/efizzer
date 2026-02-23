package vm

import (
  "time"
  "fmt"
)

type State int

const (
  StateInit = iota
  StateRun
  StateTimeout
  StateKill
)

var (
  ErrAlreadyInit = fmt.Errorf("already init")
  ErrDied        = fmt.Errorf("machine died")
  ErrNotRun      = fmt.Errorf("machine not run")
  ErrState       = fmt.Errorf("error state")
)

type Machine interface {
  Start()                   error
  Wait()                    error
  Stop()                    error
  AppendArgs([]string)      error
  SetTimeout(time.Duration) error
  Pulse()                   error
  GetState()                State
}
