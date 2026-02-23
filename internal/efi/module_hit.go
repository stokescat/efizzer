package efi

import (
  "log"

  "github.com/tidwall/btree"
)


func (self *Module) SetAddressLayout(addr, size uint64) {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  if self.modSize > 0 {
    return
  }

  self.modAddr = addr
  self.modSize = size
  self.endAddr = addr + size
  self.newBtree()
}

func (self *Module) Hit(addr uint64, hit uint32) bool {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  if !((self.modAddr <= addr) && (addr < self.endAddr)) {
    return false
  }

  addr = addr - self.modAddr

  hitVal, isexist:= self.hitTree.Get(addr)
  if isexist {
    hitNewVal:= uint64(hitVal) + uint64(hit)
    if hitNewVal > maxHitValue {
      hitNewVal = maxHitValue
    }
    self.hitTree.Set(addr, uint32(hitNewVal))
    return true
  }

  if self.hitCount >= maxHitCapacity {
    oldTree:= self.hitTree
    self.newBtree()

    self.wg.Add(1)
    go func(tree *btree.Map[uint64, uint32]) {
        defer self.wg.Done()
        self.flushTree(tree)
    } (oldTree)
  }

  self.hitTree.Set(addr, hit)
  self.hitCount++
  return true
}

