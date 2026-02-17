package rawcov

import (
    "os"
    "bufio"
    "fmt"
)

var (
  ErrNoRecord     = fmt.Errorf("no record")
  ErrNil          = fmt.Errorf("nil object")
  ErrArgs         = fmt.Errorf("invalid args")
  ErrInvSignature = fmt.Errorf("invalid signature")
  ErrInvFlags     = fmt.Errorf("invalid flags")
  ErrNoFile       = fmt.Errorf("no file")
)

type RawcovRecord struct {
  Value uint64
  Hit   uint32
}

type RawcovFile struct {
  file      *os.File      // linked file
  buf       *bufio.Reader // linked buffered stream
  count     uint64        // count of records
  index     uint64        // id of next record
  path      string        // path to file
  recordBuf [12]byte      // buffer for io operations with record
  headerBuf [16]byte      // buffer for io operations with header
}


type RawcovReader interface {
  IsEmpty() bool                  // return true if object is empty
  Get()     (RawcovRecord, error) // return next record value or ErrNoRecord
  Reset()   error                 // reset object to start
  Len()     uint64                // return count of items
}


// TODO TODO

