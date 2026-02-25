package vm

import (
  "os/exec"
  "errors"
  "syscall"
  "time"
  "fmt"
  "log"
)

const (
  stopGracePeriod = 3*time.Second
)


func NewQemu(binPath string) (*QemuMachine, error) {

  bin, err:= exec.LookPath(binPath)
  if err != nil {
    return nil, fmt.Errorf("failed to get binary file: %w", err) 
  }

  obj:= new(QemuMachine)

  obj.bin = bin // set full path to qemu binary
  obj.state = StateInit // set initial state

  return obj, nil
}


func (self *QemuMachine) Start() error {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  self.mu.Lock()
  defer self.mu.Unlock()

  if self.state != StateInit {
    // run machine we can in Init state
    return ErrState
  }

  // make command
  cmd:= exec.Command(self.bin, self.args...)
  // make process group for command
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Setpgid: true,
  }

  if err:= cmd.Start(); err != nil {
    return err
  }

  self.cmd = cmd
  self.state = StateRun

  if self.timeout > 0 {
    self.timer = time.AfterFunc(self.timeout, self.timeoutHandler)
  }

  return nil
}

func (self *QemuMachine) Wait() error {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  self.mu.Lock()
  if self.state != StateRun {
    self.mu.Unlock()
    return ErrNotRun
  }

  cmd:= self.cmd
  self.mu.Unlock()

  err:= cmd.Wait()

  self.mu.Lock()
  if self.state == StateRun {
    self.state = StateKill
  }
  self.mu.Unlock()

  return err
}

func (self *QemuMachine) Stop() error {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  return self.kill(StateKill)

}

func (self *QemuMachine) timeoutHandler() {

  _ = self.kill(StateTimeout)

}

func (self *QemuMachine) kill(newState State) error {

  self.mu.Lock()
  if self.state != StateRun {
    // do nothing if machine already killed
    self.mu.Unlock()
    return nil
  }

  // get process id
  pid:= self.cmd.Process.Pid
  // change state
  self.state = newState
  // free mutex
  self.mu.Unlock()

  // kill process

  // kill timer
  if self.timeout > 0 {
    self.timer.Stop()
  }

  // get process id
  pgid, err:= syscall.Getpgid(pid)
  if err != nil {
    if errors.Is(err, syscall.ESRCH) {
      return nil
    }
    return fmt.Errorf("failed to find process group: %w", err)
  }

  // soft kill process group
  _ = syscall.Kill(-pgid, syscall.SIGTERM)

  // delay
  time.Sleep(stopGracePeriod)

  if err = syscall.Kill(pid, 0); err == nil {
    // if we here then pid already exist
    // hard kill process group
    _ = syscall.Kill(-pgid, syscall.SIGKILL)
  }

  return nil
}
