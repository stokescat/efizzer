package efi

import (
  "log"

  "github.com/tidwall/btree"
)

func (self *Module) newBtree() {

  self.hitTree = btree.NewMap[uint64, uint32](0)
  self.hitCount = 0
}

func (self *Module) flushTree(tree *btree.Map[uint64, uint32]) {

  self.fileMu.Lock()
  defer self.fileMu.Unlock()

  reader:= newBtreeReader(tree)
  _, err:= self.file.Merge(reader, nil)
  if err != nil {
    log.Fatalf("failed to merge files: %v", err)
  }

}

func (self *Module) Flush() {

  if self == nil {
    log.Panicf("method called with nil value")
  }

  // wait finish other gorutine
  self.wg.Wait()

  // flush data
  if self.hitCount > 0 {
    self.flushTree(self.hitTree)
    self.newBtree()
  }
}


func (self *Module) Close() error {

  if self == nil {
    return nil
  }

  // flush data
  self.Flush()

  self.hitTree = nil
  return self.file.Close()
}
