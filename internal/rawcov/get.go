package rawcov

import (
  "io"
  "encoding/binary"
  "fmt"
)

func (self *File) getRecord() (Record, error) {

  // return ErrNoRecord if object is empty
  if self.IsEmpty() {
    return Record{}, ErrNoRecord
  }

  // return ErrNoRecord if object have no record
  if self.index >= self.count {
    return Record{}, ErrNoRecord
  }

  // read record from file to buffer
  if _, err:= io.ReadFull(self.buf, self.inpRecordBuf[:]); err != nil {
    return Record{}, fmt.Errorf("failed to read record: %w", err)
  }

  // get values from buffer
  recordValue:= binary.LittleEndian.Uint64(self.inpRecordBuf[:8])
  recordHit:= binary.LittleEndian.Uint32(self.inpRecordBuf[8:])
  self.index++

  // return record value
  return Record{
    Value: recordValue,
    Hit: recordHit,
  }, nil
}

func (self *File) Get() (Record, error) {

  if self == nil {
    return Record{}, ErrNil
  }

  return self.getRecord()
}
