package vm

import (
  "time"
  "os/exec"
  "sync"
  "log"
)

type QemuMachine struct {
  bin       string
  args      []string

  timeout   time.Duration
  timer     *time.Timer

  cmd       *exec.Cmd

  mu        sync.Mutex
  state     State
}

// Накапливает аргументы
func (self *QemuMachine) AppendArgs(args []string) error {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  self.mu.Lock()
  defer self.mu.Unlock()

  if self.state != StateInit {
    return ErrAlreadyInit
  }

  self.args = append(self.args, args...)
  return nil
}


// Устанавливает тайм-аут
func (self *QemuMachine) SetTimeout(t time.Duration) error {

  if self == nil {
   log.Panicf("method called with nil value")
  }

  self.mu.Lock()
  defer self.mu.Unlock()

  if self.state != StateInit {
    return ErrAlreadyInit
  }

  self.timeout = t
  return nil
}


// heart beat
func (self *QemuMachine) Pulse() error {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  self.mu.Lock()
  defer self.mu.Unlock()

  switch self.state {
    case StateInit:
      return ErrNotRun

    case StateTimeout, StateKill:
      return ErrDied

    case StateRun:
      if self.timeout > 0 {
        self.timer.Reset(self.timeout)
      }
      return nil

    default:
      return ErrState
  }
}

// возвращает состояние
func (self *QemuMachine) GetState() State {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  self.mu.Lock()
  defer self.mu.Unlock()

  return self.state
}
