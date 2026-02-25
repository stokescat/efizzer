package main

import (
  "path/filepath"
  "os"
  "net"
  "log"
  "io"
  "errors"
)

func collectorRun() {

  // signal that runner is died
  defer gWaiter.Done()

  gMmioSocketPath = filepath.Join(gWorkDir, "mmiodev.sock")

  // remove old socket file
  os.RemoveAll(gMmioSocketPath)

  // create new server
  L, err:= net.Listen("unix", gMmioSocketPath)
  if err != nil {
    log.Printf("failed to create mmio device socket (path=%v)", gMmioSocketPath)
    gEventChan <- Event{Who: EventWhoCollector, Type: EventTypeFailed}
    return
  }

  // set global var for close from other gorutine
  gCollectorListener = L
  defer os.RemoveAll(gMmioSocketPath)

  gEventChan <- Event{Who: EventWhoCollector, Type: EventTypeRun}
  for {
    // take connection
    conn, err:= L.Accept()
    if err != nil {
      if errors.Is(err, net.ErrClosed) {
        log.Printf("mmio server is closed")
        return
      }
      log.Printf("failed to accept connection: %v", err)
      continue
    }

    // set global var for close from other gorutine
    gCollectorConnection = conn
    gEventChan <- Event{Who: EventWhoCollector, Type: EventTypeConnect}

    err = collectorDispatch(conn)
    if err != nil {

      if errors.Is(err, net.ErrClosed) {
        log.Printf("mmio server is closed")
        return
      }

      if errors.Is(err, io.EOF) {
        log.Printf("mmio client is closed connect")
      } else {
        log.Printf("failed to read from socket %v", err)
      }
      
      conn.Close()
      gEventChan <- Event{Who: EventWhoCollector, Type: EventTypeBreak}
      continue 
    }

  }
}

