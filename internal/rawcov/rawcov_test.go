package rawcov_test

import (
  "testing"
  "log"
  "os"
  "path/filepath"

  "github.com/tidwall/btree"
  "github.com/stokescat/efizzer/internal/rawcov"
)

type SimpleReader struct {

  Tree    *btree.Map[uint64, uint32]
  Index  uint64
}

func (self *SimpleReader) IsEmpty() bool {

  return (self.Tree.Len() <= 0)
}

func (self *SimpleReader) Len() uint64 {

  return uint64(self.Tree.Len())
}

func (self *SimpleReader) Reset() error {

  self.Index = 0
  return nil
}

func (self *SimpleReader) Get() (rawcov.Record, error) {

  count:= uint64(self.Tree.Len())

  if self.Index >= count {
    return rawcov.Record{}, rawcov.ErrNoRecord
  }

  value, hit, _:= self.Tree.GetAt(int(self.Index))
  self.Index++
  return rawcov.Record{
    Value: value,
    Hit: hit,
  }, nil
}

var (

  addrOddTree   btree.Map[uint64, uint32]
  addrEvenTree  btree.Map[uint64, uint32]
  addrTree      btree.Map[uint64, uint32]

  addrOddFileName string
  addrEvenFileName string
  addrFileName string

)

func TestMain(m *testing.M) {

// inits btree
  for I:= 0; I < 2048; I++ {
    if I % 2 == 0 {
      addrEvenTree.Set(uint64(I), uint32(I))
      addrTree.Set(uint64(I), uint32(I))
    } else {
      addrOddTree.Set(uint64(I), uint32(I))
      addrTree.Set(uint64(I), uint32(I))
    }
  }

  DirName, err:= os.MkdirTemp("", "rawcov_testdir_*")
  if err != nil {
    log.Fatalf("[rawcov test] failed to created temp dir: %v", err)
  }

  addrOddFileName = filepath.Join(DirName, "addrodd.rawcov")
  addrEvenFileName = filepath.Join(DirName, "addreven.rawcov")
  addrFileName = filepath.Join(DirName, "addr.rawcov")


  exitCode:= m.Run()

  os.RemoveAll(DirName)
  os.Exit(exitCode)
}

func TestWriteAndReadOddRecords(t *testing.T) {

  reader:= &SimpleReader {
    Tree: &addrOddTree,
    Index: 0, 
  }

  rawcovFile, err:= rawcov.Open(addrOddFileName)
  if err != nil {
    t.Fatalf("failed to create file: %v", err)
  }

  count, err:= rawcovFile.Merge(reader, nil)
  if err != nil {
    t.Fatalf("failed to merge file: %v", err)
  }

  if count != 1024 {
    t.Fatalf("unexpected number of new records: expected %v, got %v", 1024, count)
  }

  if err = rawcovFile.Reset(); err != nil {
    t.Fatalf("failed to reset file: %v", err)
  }

  for I:= 0; I < addrOddTree.Len(); I++ {
    value, hit, _:= addrOddTree.GetAt(I)
    rec, err:= rawcovFile.Get()
    if err != nil {
      t.Fatalf("failed to get record: %v", err)
    }
    if rec.Value != value {
      t.Fatalf("value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }
    if rec.Hit != hit {
      t.Fatalf("hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  if err = rawcovFile.Close(); err != nil {
    t.Fatalf("failed to close file: %v", err)
  }

}

func TestWriteAndReadEvenRecords(t *testing.T) {

  reader:= &SimpleReader {
    Tree: &addrEvenTree,
    Index: 0, 
  }

  rawcovFile, err:= rawcov.Open(addrEvenFileName)
  if err != nil {
    t.Fatalf("failed to create file: %v", err)
  }

  count, err:= rawcovFile.Merge(reader, nil)
  if err != nil {
    t.Fatalf("failed to merge file: %v", err)
  }

  if count != 1024 {
    t.Fatalf("unexpected number of new records: expected %v, got %v", 1024, count)
  }

  if err = rawcovFile.Reset(); err != nil {
    t.Fatalf("failed to reset file: %v", err)
  }

  for I:= 0; I < addrEvenTree.Len(); I++ {
    value, hit, _:= addrEvenTree.GetAt(I)
    rec, err:= rawcovFile.Get()
    if err != nil {
      t.Fatalf("failed to get record: %v", err)
    }
    if rec.Value != value {
      t.Fatalf("value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }
    if rec.Hit != hit {
      t.Fatalf("hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  if err = rawcovFile.Close(); err != nil {
    t.Fatalf("failed to close file: %v", err)
  }

}



func TestOddAndEvenFilesContent(t *testing.T) {

  oddFile, err:= rawcov.Open(addrOddFileName)
  if err != nil {
    t.Fatalf("[odd] failed to open odd file: %v", err)
  }

  evenFile, err:= rawcov.Open(addrEvenFileName)
  if err != nil {
    t.Fatalf("[even] failed to open even file: %v", err)
  }

  for I:= 0; I < addrOddTree.Len(); I++ {
    value, hit, _:= addrOddTree.GetAt(I)
    rec, err:= oddFile.Get()
    if err != nil {
      t.Fatalf("[odd] failed to get record: %v", err)
    }
    if rec.Value != value {
      t.Fatalf("[odd] value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }
    if rec.Hit != hit {
      t.Fatalf("[odd] hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  for I:= 0; I < addrEvenTree.Len(); I++ {
    value, hit, _:= addrEvenTree.GetAt(I)
    rec, err:= evenFile.Get()
    if err != nil {
      t.Fatalf("[even] failed to get record: %v", err)
    }
    if rec.Value != value {
      t.Fatalf("[even] value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }
    if rec.Hit != hit {
      t.Fatalf("[even] hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  if err = oddFile.Close(); err != nil {
    t.Fatalf("[odd] failed to close file: %v", err)
  }

  if err = evenFile.Close(); err != nil {
    t.Fatalf("[even] failed to close file: %v", err)
  }
}


// TODO TODO
func TestMergeTwoFiles(t *testing.T) {

  oddFile, err:= rawcov.Open(addrOddFileName)
  if err != nil {
    t.Fatalf("[odd] failed to open odd file: %v", err)
  }

  evenFile, err:= rawcov.Open(addrEvenFileName)
  if err != nil {
    t.Fatalf("[even] failed to open even file: %v", err)
  }

  addrFile, err:= rawcov.Open(addrFileName)
  if err != nil {
    t.Fatalf("[merge] failed to create file: %v", err)
  }

  count, err:= addrFile.Merge(oddFile, nil)
  if err != nil {
    t.Fatalf("failed to merge file [odd -> merge]: %v", err)
  }

  if count != 1024 {
    t.Fatalf("unexpected number of new records [odd -> merge]: expected %v, got %v", 1024, count)
  }

  count, err = addrFile.Merge(evenFile, nil)
  if err != nil {
    t.Fatalf("failed to merge file [even -> merge]: %v", err)
  }

  if count != 1024 {
    t.Fatalf("unexpected number of new records [even -> merge]: expected %v, got %v", 1024, count)
  }

  if addrFile.Len() != 2048 {
    t.Fatalf("unexpected number of new records [merge]: expected %v, got %v", 2048, addrFile.Len())
  }

  if err = oddFile.Close(); err != nil {
    t.Fatalf("[odd] failed to close file: %v", err)
  }

  if err = evenFile.Close(); err != nil {
    t.Fatalf("[even] failed to close file: %v", err)
  }

  if err = addrFile.Close(); err != nil {
    t.Fatalf("[merge] failed to close file: %v", err)
  }
}


func TestMergedFileContent(t *testing.T) {

  addrFile, err:= rawcov.Open(addrFileName)
  if err != nil {
    t.Fatalf("failed to open file: %v", err)
  }

  if err = addrFile.Reset(); err != nil {
    t.Fatalf("failed to reset file: %v", err)
  }

  for I:= 0; I < addrTree.Len(); I++ {
    value, hit, _:= addrTree.GetAt(I)
    rec, err:= addrFile.Get()
    if err != nil {
      t.Fatalf("failed to get record: %v", err)
    }
    if rec.Value != value {
      t.Fatalf("value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }
    if rec.Hit != hit {
      t.Fatalf("hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  if err = addrFile.Close(); err != nil {
    t.Fatalf("failed to close file: %v", err)
  }
}

func TestMergeSameDataTwice(t *testing.T) {

  reader:= &SimpleReader {
    Tree: &addrOddTree,
    Index: 0, 
  }

  rawcovFile, err:= rawcov.Open(addrOddFileName)
  if err != nil {
    t.Fatalf("failed to open file: %v", err)
  }

  count, err:= rawcovFile.Merge(reader, nil)
  if err != nil {
    t.Fatalf("failed to merge file: %v", err)
  }

  if count != 0 {
    t.Fatalf("unexpected number of new records [merge]: expected %v, got %v", 0, count)
  }

  if err = rawcovFile.Reset(); err != nil {
    t.Fatalf("failed to reset file: %v", err)
  }

  for I:= 0; I < addrOddTree.Len(); I++ {
    value, hit, _:= addrOddTree.GetAt(I)
    hit = 2*hit
    rec, err:= rawcovFile.Get()
    if err != nil {
      t.Fatalf("failed to get record: %v", err)
    }
    if rec.Value != value {
      t.Fatalf("value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }
    if rec.Hit != hit {
      t.Fatalf("hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  if err = rawcovFile.Close(); err != nil {
    t.Fatalf("failed to close file: %v", err)
  }

}


func TestMergeExistingData(t *testing.T) {

  oddFile, err:= rawcov.Open(addrOddFileName)
  if err != nil {
    t.Fatalf("[odd] failed to open file: %v", err)
  }

  addrFile, err:= rawcov.Open(addrFileName)
  if err != nil {
    t.Fatalf("[merge] failed to open file: %v", err)
  }

  count, err:= addrFile.Merge(oddFile, nil)
  if err != nil {
    t.Fatalf("failed to merge file [odd -> merge]: %v", err)
  }

  if count != 0 {
    t.Fatalf("unexpected number of new records [merge]: expected %v, got %v", 0, count)
  }

  if addrFile.Len() != 2048 {
    t.Fatalf("unexpected number of records [merge]: expected %v, got %v", 2048, addrFile.Len())
  }

  if err = oddFile.Close(); err != nil {
    t.Fatalf("[odd] failed to close file: %v", err)
  }

  if err = addrFile.Close(); err != nil {
    t.Fatalf("[merge] failed to close file: %v", err)
  }
}

func TestMergedFileAfterAddingOddAgain(t *testing.T) {

  addrFile, err:= rawcov.Open(addrFileName)
  if err != nil {
    t.Fatalf("failed to open file: %v", err)
  }

  if err = addrFile.Reset(); err != nil {
    t.Fatalf("failed to reset file: %v", err)
  }

  for I:= 0; I < addrTree.Len(); I++ {
    value, hit, _:= addrTree.GetAt(I)
    if (value % 2 != 0) {
      hit = 3*hit
    }
    rec, err:= addrFile.Get()
    if err != nil {
      t.Fatalf("failed to get record: %v", err)
    }
    if rec.Value != value {
      t.Fatalf("value mismatch at index %d: expected %v, got %v", I, value, rec.Value)
    }
    if rec.Hit != hit {
      t.Fatalf("hit count mismatch for value %v: expected %v, got %v", value, hit, rec.Hit)
    }
  }

  if err = addrFile.Close(); err != nil {
    t.Fatalf("failed to close file: %v", err)
  }
}

