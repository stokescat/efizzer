package rawcov

import (
  "encoding/binary"
  "fmt"
  "os"
  "bufio"
  "io"
)


func (self *File) fillRecordBuf(value uint64, hit uint32) {

  binary.LittleEndian.PutUint64(self.outRecordBuf[:8], value)
  binary.LittleEndian.PutUint32(self.outRecordBuf[8:], hit)
}

func (self *File) fillAndSumRecordBuf(value uint64, hit1 uint32, hit2 uint32) {

  hitSum:= uint64(hit1) + uint64(hit2)

  if hitSum > 0xffffffff {
    self.fillRecordBuf(value, 0xffffffff)
  } else {
    self.fillRecordBuf(value, uint32(hitSum))
  }

}

func (self *File) fillHeaderBuf(flags uint16, count uint64) {

  binary.LittleEndian.PutUint16(self.headerBuf[6:8], flags)
  binary.LittleEndian.PutUint64(self.headerBuf[8:], count)
}

func (self *File) fillSignature() {

  copy(self.headerBuf[:6], "RAWCOV")
}


func (self *File) IsEmpty() bool {

  return (self.count == 0)
}

func (self *File) closeFile() error {

  // clear object
  self.buf = nil
  self.index = 0
  self.count = 0
  self.fillSignature()
  self.fillHeaderBuf(0, 0)

  if self.file != nil {
    // close file if there is
    err:= self.file.Close()
    self.file = nil
    return err
  }
  return nil
}


func (self *File) resetFile() error {

  // if we have not file, then return nil
  if self.file == nil {
    return nil
  }

  if _, err:= self.file.Seek(16, 0); err != nil {
    return fmt.Errorf("failed to seek file: %w", err)
  }

  self.buf.Reset(self.file)
  self.index = 0

  return nil
}

func (self *File) Reset() error {

  if self == nil {
    return ErrNil
  }

  return self.resetFile()
}

func (self *File) initFile(file *os.File) error {

  if file == nil {
    // init empty file
    self.file = nil
    self.buf = nil
    self.count = 0
    self.index = 0
    self.fillSignature()
    self.fillHeaderBuf(0, 0)
    return nil
  }

  self.file = file
  self.buf = bufio.NewReader(file)

  var err error

  //read and check signature
  if _, err = io.ReadFull(self.buf, self.headerBuf[:]); err != nil {
    return fmt.Errorf("failed to read header: %w", err)
  }

  if string(self.headerBuf[:6]) != "RAWCOV" {
    return ErrInvSignature
  }

  if hFlags:= binary.LittleEndian.Uint16(self.headerBuf[6:8]); hFlags != 0 {
    return ErrInvFlags
  }

  self.count = binary.LittleEndian.Uint64(self.headerBuf[8:])
  self.index = 0
  return nil
}

func (self *File) Len() uint64 {
  if self == nil {
    return 0
  }
  return self.count
}

func (self *File) Close() error {

  if self == nil {
    return nil
  }

  return self.closeFile()
}
