package efi

import (

  "github.com/stokescat/efizzer/internal/rawcov"

  "github.com/tidwall/btree"
)


type btreeReader struct {

  tree *btree.Map[uint64, uint32]
  index uint64
  count uint64
}

func (self *btreeReader) Get() (rawcov.Record, error) {

  if self.index >= self.count {
    return rawcov.Record{}, rawcov.ErrNoRecord
  }

  value, hit, _:= self.tree.GetAt(int(self.index))
  self.index++

  return rawcov.Record {
    Value: value,
    Hit: hit,
  }, nil
}

func (self *btreeReader) IsEmpty() bool {

  return (self.count == 0)
}

func (self *btreeReader) Reset() error {

  self.index = 0
  return nil
}

func (self *btreeReader) Len() uint64 {

  return (self.count)
}

func newBtreeReader(tree *btree.Map[uint64, uint32]) *btreeReader {

  obj:= &btreeReader{
    tree: tree,
    index: 0,
    count: uint64(tree.Len()),
  }

  return obj
}
