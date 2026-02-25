package main

type EventWho uint16
type EventType uint16

type Event struct {
  Who     EventWho
  Type    EventType
  Data    any
}

const (
  EventWhoMain = iota
  EventWhoCollector
  EventWhoMachine
)

const (
  EventTypePing = iota
  EventTypeRun
  EventTypeConnect
  EventTypeError
  EventTypeFailed
  EventTypeBreak
  EventTypeDied
)

func (self *Event) Check(w EventWho, t EventType) bool {

  return ((self.Who == w) && (self.Type == t))
}

