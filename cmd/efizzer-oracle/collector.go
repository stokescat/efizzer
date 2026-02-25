package main

import (
  "log"
  "net"
  "sync"

  "github.com/stokescat/efizzer/internal/efi"
  "github.com/stokescat/efizzer/internal/mmep"
)

var (
  modules       = make(map[efi.Guid]*efi.Module)

  messagePool   = sync.Pool{
    New: func() any {
      return mmep.NewMessage()
    },
  }
)

var (
  undefModule   *efi.Module

)

func getMessage() *mmep.Message {
  return messagePool.Get().(*mmep.Message)
}

func putMessage(msg *mmep.Message) {
  msg.Clear()
  messagePool.Put(msg)
}


func collectorStart() {

  var err error
  undefModule, err = efi.OpenModuleByName(gWorkDir, "undefined")
  if err != nil {
    log.Printf("failed to open undefined module")
    gEventChan <- Event{Who: EventWhoCollector, Type: EventTypeFailed}
    return
  }

  undefModule.SetAddressLayout(uint64(0), 0xffffffffffffffff)

  gWaiter.Add(1)
  go collectorRun()
}


func collectorStop() {

  if gCollectorConnection != nil {
    gCollectorConnection.Close()
    gCollectorConnection = nil
  }

  if gCollectorListener != nil {
    gCollectorListener.Close()
    gCollectorListener = nil
  }
}

func collectorFlushModules() {

  for _, module:= range modules {
    module.Flush()
  }

  undefModule.Flush()
}

func collectorCloseModules() {

  for guid, module:= range modules {
    module.Close()
    delete(modules, guid)
  }

  undefModule.Close()

}

func collectorDispatch(conn net.Conn) error {

  reader:= mmep.NewReader(conn)

  for {
    msg:= getMessage()

    err:= reader.Read(msg)
    if err != nil {
      return err
    }

    switch tag:= msg.GetTag(); tag {
      case mmep.MsgTagEvent:
        switch eventType:= msg.GetEventType(); eventType {
          case mmep.MsgEventHit:
            //hit event
            addrHit:= msg.GetHit()
            handleHitEvent(addrHit)
            

          case mmep.MsgEventMod:
            // module load/unload event
            modGuid:= msg.GetModuleGuid()
            modSize:= msg.GetModuleSize()
            modAddr:= msg.GetModuleAddr()

            handleModuleEvent(modGuid, modSize, modAddr)

          default:
            log.Printf("undefined event type (type=%v)", eventType)
        }

      default:
        log.Printf("undefined message tag")
    }

    err = gMachine.Pulse()
    if err != nil {
      log.Printf("failed to pulse vm: %v", err)
    }
    putMessage(msg)
  }
}

