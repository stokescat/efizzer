package main

import (
  "sync"
  "net"

  "github.com/stokescat/efizzer/internal/vm"
)

const (
  eventChanSize = 50
)

var (
  gWaiter               sync.WaitGroup

  gWorkDir              string
  gMmioSocketPath       string
  gMachinePath          string

  gMachine              vm.Machine
  gCollectorListener    net.Listener
  gCollectorConnection  net.Conn
)

var (
  gEventChan = make(chan Event, eventChanSize)
)

