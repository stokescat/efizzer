package main

import (
  "log"
  "os"

  "context"
  "os/signal"
  "syscall"

  "github.com/stokescat/efizzer/internal/vm"
)

func main() {

// заглушка

  gWorkDir = ""
  gMachinePath = "/home/pavel/efizzer/qemu/build/qemu-system-x86_64"
  firmwarePath:= "/home/pavel/efizzer/audk1/Build/OvmfX64/DEBUG_CLANGDWARF/FV/OVMF.fd"

  SigtermEvent, stopSigterm:= signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
  defer stopSigterm()

  collectorStart()

  event:= <-gEventChan
  if !event.Check(EventWhoCollector, EventTypeRun) {
    log.Fatalf("failed to start mmio device server")
  }

  machine, err:= vm.NewQemu(gMachinePath)
  if err != nil {
    log.Fatalf("failed to create vm")
  }

  machineArgs:= []string{
      "-machine", "q35,efizzer-addr=0xfeb00000,efizzer-sock=unix:"+gMmioSocketPath,
      "-m",       "5G",
      "-bios",    firmwarePath,
      "-nographic",
    }

  gMachine = machine
  gMachine.AppendArgs(machineArgs)

  //run machine
  err = gMachine.Start()
  if err != nil {
    log.Fatalf("failed to run vm")
  }

  go func() {
    err:= gMachine.Wait()
    gEventChan <- Event{Who: EventWhoMachine, Type: EventTypeDied, Data: err}
  } ()

  log.Printf("vm started") 

  for {
    var ev Event
    select {
      case ev = <-gEventChan:
        log.Printf("Something event... %v", ev)

      case <-SigtermEvent.Done():
        gMachine.Stop()
        collectorStop()
        log.Printf("Waiting to stop all...")
        gWaiter.Wait()
        collectorCloseModules()
        log.Printf("goodbye! :)")
        return
    }
  }
}
