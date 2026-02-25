package main

import (
  "log"

  "github.com/stokescat/efizzer/internal/efi"
)

func handleModuleEvent(modGuid efi.Guid, modSize uint64, modAddr uint64) {

  var err error
  module, isexist:= modules[modGuid]

  // is unload event?
  if modAddr == 0 {
    if isexist {
      module.Close()
      delete(modules, modGuid)
      log.Printf("unload module event (guid=%v)", string(modGuid))
      return
    }
    log.Printf("try to unload unexisting module (guid=%v)", string(modGuid))
    return
  }

  if !isexist {
    module, err = efi.OpenModuleByGuid(gWorkDir, modGuid)
    if err != nil {
      log.Printf("failed to open module (guid=%v): %v", string(modGuid), err)
      return
    }

    module.SetAddressLayout(modAddr, modSize)
    modules[modGuid] = module
    log.Printf("load module event (guid=%v)", string(modGuid))
  } else {
    log.Printf("try to load existing module (guid=%v)", string(modGuid))
  }
}

func handleHitEvent(addr uint64) {

  isfind:= false
  for _, mod:= range modules {
    if find:= mod.Hit(addr, 1); find {
      isfind = true
      break
    }
  }

  if !isfind {
    undefModule.Hit(addr, 1)
  }
}
