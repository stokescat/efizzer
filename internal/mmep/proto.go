package mmep

import (
  "encoding/binary"

  "github.com/stokescat/efizzer/internal/efi"
)

const (
  MaxMessageSize = 4096

  maskHeaderLength    = 0x0000000000000fff
  maskHeaderTag       = 0xf000000000000000
  maskHeaderEventType = 0x000000000000f000
)

type Message struct {
  buf [MaxMessageSize]byte
}

type MsgTag uint
const (
  MsgTagEvent = MsgTag(0)
  MsgTagBegin = MsgTag(1)
  MsgTagEnd   = MsgTag(2)
  MsgTagData  = MsgTag(3)
)

type MsgEventType uint
const (
  MsgEventHit         = MsgEventType(1)
  MsgEventMod         = MsgEventType(2)
  MsgEventExecutorHi  = MsgEventType(3)
  MsgEventStart       = MsgEventType(4)
  MsgEventStop        = MsgEventType(5)
)

func NewMessage() *Message {

  return new(Message)
}

func (self *Message) GetLength() uint {

  header:= binary.LittleEndian.Uint64(self.buf[:8])
  header = header & maskHeaderLength
  return uint(header)
}

func (self *Message) GetTag() MsgTag {

  header:= binary.LittleEndian.Uint64(self.buf[:8])
  header = (header & maskHeaderTag) >> 60
  return MsgTag(header)
}

func (self *Message) GetEventType() MsgEventType {

  header:= binary.LittleEndian.Uint64(self.buf[:8])
  header = (header & maskHeaderEventType) >> 12
  return MsgEventType(header)
}

func (self *Message) GetHit() uint64 {

  addr:= binary.LittleEndian.Uint64(self.buf[8:16])
  return addr
}

func (self *Message) GetModuleGuid() efi.Guid {

  guid:= efi.BytesToGuid(self.buf[8:24])
  return guid
}

func (self *Message) GetModuleSize() uint64 {
  size:= binary.LittleEndian.Uint64(self.buf[24:32])
  return size
}

func (self *Message) GetModuleAddr() uint64 {
  addr:= binary.LittleEndian.Uint64(self.buf[32:40])
  return addr
}

func (self *Message) Clear() {

  msg_length:= self.GetLength()
  clear(self.buf[:msg_length])
}

