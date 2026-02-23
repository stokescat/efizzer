package efi

import (
  "sync"

  "github.com/stokescat/efizzer/internal/rawcov"

  "github.com/tidwall/btree"
)

const (
  maxHitCapacity = 1024
  maxHitValue = 0xffffffff
)

type Module struct {

  Name string
  Guid Guid

  modAddr uint64
  modSize uint64
  endAddr uint64

  file *rawcov.File
  fileMu sync.Mutex
  wg  sync.WaitGroup

  hitTree *btree.Map[uint64, uint32]
  hitCount int
}



