package efi_test

import (
  "testing"

  "github.com/stokescat/efizzer/internal/efi"
)

func TestBytesToGuid(t *testing.T) {

  guid1_str:= "500DF8D1-CA05-3042-959E-8872D94F6BC6"
  guid1_bytes:= [16]byte{0xD1,0xF8,0x0D,0x50,0x05,0xCA,0x42,0x30,0x95,0x9E,0x88,0x72,0xD9,0x4F,0x6B,0xC6}
  guid2_str:= "336F5A24-E9F9-B643-972A-14C8738E3F5B"
  guid2_bytes:= [16]byte{0x24,0x5A,0x6F,0x33,0xF9,0xE9,0x43,0xB6,0x97,0x2A,0x14,0xC8,0x73,0x8E,0x3F,0x5B}

  if guid:= efi.BytesToGuid(guid1_bytes[:]); guid != efi.Guid(guid1_str) {
      t.Fatalf("invalid EFI GUID: expected %s, got %s", guid1_str, guid)
  }

  if guid:= efi.BytesToGuid(guid2_bytes[:]); guid != efi.Guid(guid2_str) {
      t.Fatalf("invalid EFI GUID: expected %s, got %s", guid1_str, guid)
  }
}
