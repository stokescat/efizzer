package rawcov

import (
  "os"
  "bufio"
  "fmt"
  "errors"
)

func (dst *File) mergeToNewFile(src Reader, filename string, tr func(Record) Record) (uint64, error) {

  countRecords:= uint64(0) // count of records in new file
  countMatches:= uint64(0) // count matches value in dst

  file, err:= os.Create(filename)
  if err != nil {
    return 0, fmt.Errorf("failed to open new file: %w", err)
  }
  defer file.Close()

  // reset dst file
  if err = dst.resetFile(); err != nil {
    return 0, fmt.Errorf("failed to reset destination file: %w", err)
  }

  // reset src file
  if err = src.Reset(); err != nil {
    return 0, fmt.Errorf("failed to reset source file: %w", err)
  }

  // create writeable buffer
  writer:= bufio.NewWriter(file)

  // make null header
  dst.fillSignature()
  dst.fillHeaderBuf(0, 0)
  if _, err = writer.Write(dst.headerBuf[:]); err != nil {
    return 0, fmt.Errorf("failed to write initial header: %w", err)
  }

  // function get and transform record from src file
  getSrcRecord:= func() (Record, error) {
    if rec, err:= src.Get(); err == nil {
       return tr(rec), nil
    } else {
       return rec, err
    }
  }

  // init values
  dRec, dErr:= dst.getRecord()
  sRec, sErr:= getSrcRecord()
  for {
    // exit condition
    if errors.Is(dErr, ErrNoRecord) && errors.Is(sErr, ErrNoRecord) {
      break
    }
    // handle errors
    if (dErr != nil) && !errors.Is(dErr, ErrNoRecord) {
      return 0, fmt.Errorf("failed to read record from destination file: %w", dErr)
    }
    if (sErr != nil) && !errors.Is(sErr, ErrNoRecord) {
      return 0, fmt.Errorf("failed to read record from source file: %w", sErr)
    }

    // if dst is end
    if errors.Is(dErr, ErrNoRecord) {
       dst.fillRecordBuf(sRec.Value, sRec.Hit)
       sRec, sErr = getSrcRecord()
       countRecords++
       goto LOOP_WRITE_RECORD
    }

    // if src is end
    if errors.Is(sErr, ErrNoRecord) {
       dst.fillRecordBuf(dRec.Value, dRec.Hit)
       dRec, dErr = dst.getRecord()
       countRecords++
       goto LOOP_WRITE_RECORD
    }

    switch {
      case dRec.Value < sRec.Value:
        dst.fillRecordBuf(dRec.Value, dRec.Hit)
        dRec, dErr = dst.getRecord()
        countRecords++
        goto LOOP_WRITE_RECORD

      case sRec.Value < dRec.Value:
        dst.fillRecordBuf(sRec.Value, sRec.Hit)
        sRec, sErr = getSrcRecord()
        countRecords++
        goto LOOP_WRITE_RECORD

      default:
        dst.fillAndSumRecordBuf(dRec.Value, dRec.Hit, sRec.Hit)
        dRec, dErr = dst.getRecord()
        sRec, sErr = getSrcRecord()
        countRecords++
        countMatches++
        goto LOOP_WRITE_RECORD
    }

    continue

LOOP_WRITE_RECORD:

    if _, err = writer.Write(dst.outRecordBuf[:]); err != nil {
      return 0, fmt.Errorf("failed to write record: %w", err)
    }
  }

  if err = writer.Flush(); err != nil {
      return 0, fmt.Errorf("failed to flush data: %w", err)
  }

  // make right header
  dst.fillHeaderBuf(0, countRecords)
  // write right header
  if _, err = file.Seek(0, 0); err != nil {
      return 0, fmt.Errorf("failed to seek file: %w", err)
  }
  if _, err = file.Write(dst.headerBuf[:]); err != nil {
      return 0, fmt.Errorf("failed to write header: %w", err)
  }

  countUniqueValues:= src.Len() - countMatches
  return countUniqueValues, nil
}


func (dst *File) Merge(src Reader, tr func(Record) Record) (uint64, error) {

  if dst == nil {
    return 0, ErrNil
  }

  if src == nil {
    return 0, ErrArgs
  }

  if dst == src {
    return 0, ErrArgs
  }

  if tr == nil {
    tr = func (r Record) Record {return r}
  }

  if src.IsEmpty() {
    // do nothing
    return 0, nil
  }

  // create temporary file name
  newFileName:= dst.path + ".new"
  countNewRecords, err:= dst.mergeToNewFile(src, newFileName, tr)
  if err != nil {
    return 0, fmt.Errorf("failed to merge files | %w", err)
  }

  // now we need to close old destination file
  if err = dst.closeFile(); err != nil {
    return 0, fmt.Errorf("failed to close old file | %w", err)
  }

  // now we need to replace old file
  if err = os.Rename(newFileName, dst.path); err != nil {
    return 0, fmt.Errorf("failed to replace old file | %w", err)
  }

  // now we need to open new file
  newFile, err:= os.Open(dst.path)
  if err != nil {
    return 0, fmt.Errorf("failed to open new file | %w", err)
  }

  if err = dst.initFile(newFile); err != nil {
    return 0, fmt.Errorf("failed to init new file | %w", err)
  }

  return countNewRecords, nil
}
