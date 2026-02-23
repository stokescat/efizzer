package rawcov

import (
  "os"
  "bufio"
  "fmt"
)

var (
  ErrNoRecord     = fmt.Errorf("no record")
  ErrArgs         = fmt.Errorf("invalid args")
  ErrInvSignature = fmt.Errorf("invalid signature")
  ErrInvFlags     = fmt.Errorf("invalid flags")
  ErrNoFile       = fmt.Errorf("no file")
)


type Record struct {
  Value uint64
  Hit   uint32
}

type File struct {
  file         *os.File      // linked file
  buf          *bufio.Reader // linked buffered stream
  count        uint64        // count of records
  index        uint64        // id of next record
  path         string        // path to file
  outRecordBuf [12]byte      // buffer for io operations with record
  inpRecordBuf [12]byte
  headerBuf    [16]byte      // buffer for io operations with header
}


type Reader interface {
  IsEmpty() bool                  // return true if object is empty
  Get()     (Record, error) // return next record value or ErrNoRecord
  Reset()   error                 // reset object to start
  Len()     uint64                // return count of items
}


// TODO: add panic 

