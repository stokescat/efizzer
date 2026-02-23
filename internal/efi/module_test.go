package efi_test

import (
  "testing"
  "log"
  "path/filepath"
  "os"

  "github.com/stokescat/efizzer/internal/efi"
  "github.com/stokescat/efizzer/internal/rawcov"

  "github.com/tidwall/btree"
)

const (
  addrCount = 65536
)

var (
  addrTree btree.Map[uint64, uint32]

  DirPath string
)

func TestMain(m *testing.M) {

  //init addr tree
  for I:= 0; I < addrCount; I++ {
    addrTree.Set(uint64(I), uint32(I))
  }

  var err error

  DirPath, err = os.MkdirTemp("", "efi_module_testdir_*")
  if err != nil {
    log.Fatalf("[rawcov test] failed to created temp dir: %v", err)
  }

  exitCode:= m.Run()

  os.RemoveAll(DirPath)
  os.Exit(exitCode)
}


func TestCreateNewModule(t *testing.T) {

  module, err:= efi.OpenModuleByGuid(DirPath, efi.UndefinedGuid)
  if err != nil {
    t.Fatalf("failed to create module: %v", err)
  }

  module.SetAddressLayout(0x0,0xffffffffffffffff)

  addrTree.Scan(func(value uint64, hit uint32) bool {
    module.Hit(value, hit)
    return true
  })

  err = module.Close()
  if err != nil {
    t.Fatalf("failed to close new module: %v", err)
  }
}

func TestNewModuleFileContent(t *testing.T) {

  rawcovFilePath:= filepath.Join(DirPath, string(efi.UndefinedGuid) + ".rawcov")
  file, err:= rawcov.Open(rawcovFilePath)
  if err != nil {
    t.Fatalf("failed to open new .rawcov file: %v", err)
  }

  for I:= 0; I < addrCount; I++ {
    value, hit, _:= addrTree.GetAt(I)
    rec, err:= file.Get()
    if err != nil {
      t.Fatalf("failed to get record from file: %v", err)
    }

    if rec.Value != value {
      t.Fatalf("value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }

    if rec.Hit != hit {
      t.Fatalf("hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  err = file.Close()
  if err != nil {
    t.Fatalf("failed to close .rawcov file: %v", err)
  }
}

func TestAppendModule(t *testing.T) {

  module, err:= efi.OpenModuleByGuid(DirPath, efi.UndefinedGuid)
  if err != nil {
    t.Fatalf("failed to open module: %v", err)
  }

  module.SetAddressLayout(0x0,0xffffffffffffffff)

  addrTree.Scan(func(value uint64, hit uint32) bool {
    module.Hit(value, hit)
    return true
  })

  err = module.Close()
  if err != nil {
    t.Fatalf("failed to close module: %v", err)
  }
}

func TestApendModuleFileContent(t *testing.T) {

  rawcovFilePath:= filepath.Join(DirPath, string(efi.UndefinedGuid) + ".rawcov")
  file, err:= rawcov.Open(rawcovFilePath)
  if err != nil {
    t.Fatalf("failed to open new .rawcov file: %v", err)
  }

  for I:= 0; I < addrCount; I++ {
    value, hit, _:= addrTree.GetAt(I)
    hit = 2*hit
    rec, err:= file.Get()
    if err != nil {
      t.Fatalf("failed to get record from file: %v", err)
    }

    if rec.Value != value {
      t.Fatalf("value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }

    if rec.Hit != hit {
      t.Fatalf("hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  err = file.Close()
  if err != nil {
    t.Fatalf("failed to close .rawcov file: %v", err)
  }
}
