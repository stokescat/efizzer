package rawcov

import (
  "io"
  "encoding/binary"
  "fmt"
)

func (self *RawcovFile) getRecord() (RawcovRecord, error) {

  // return ErrNoRecord if object is empty
  if self.IsEmpty() {
    return RawcovRecord{}, ErrNoRecord
  }

  // return ErrNoRecord if object have no record
  if self.index >= self.count {
    return RawcovRecord{}, ErrNoRecord
  }

  // read record from file to buffer
  if _, err:= io.ReadFull(self.buf, self.recordBuf[:]); err != nil {
    return RawcovRecord{}, fmt.Errorf("failed to read record: %w", err)
  }

  // get values from buffer
  recordValue:= binary.LittleEndian.Uint64(self.recordBuf[:8])
  recordHit:= binary.LittleEndian.Uint32(self.recordBuf[8:])
  self.index++

  // return record value
  return RawcovRecord{
    Value: recordValue,
    Hit: recordHit,
  }, nil
}

func (self *RawcovFile) Get() (RawcovRecord, error) {

  if self == nil {
    return RawcovRecord{}, ErrNil
  }

  return self.getRecord()
}
