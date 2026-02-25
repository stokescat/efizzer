package mmep

import (
  "io"
  "fmt"
)

var (
  ErrInvalidLength = fmt.Errorf("invalid message size")
)

type Reader struct {
  src io.Reader
}

func NewReader(src io.Reader) *Reader {

  return &Reader{
      src: src,
    }

}

func (self *Reader) Read(msg *Message) error {

  _, err:= io.ReadFull(self.src, msg.buf[:8])
  if err != nil {
    return err
  }

  msg_length:= msg.GetLength()
  if msg_length < 8 {
    return ErrInvalidLength
  }

  _, err = io.ReadFull(self.src, msg.buf[8:msg_length])
  if err != nil {
    return err
  }

  return err
}
