package rawcov

import (
  "os"
  "fmt"
  "errors"
)

func Open(path string) (*File, error) {

  obj:= new(File)
  obj.path = path

  newFile, err:= os.Open(path)
  if err != nil {
    if errors.Is(err, os.ErrNotExist) {
      obj.initFile(nil)
      return obj, nil
    }
    return nil, fmt.Errorf("failed to open .rawcov | %w", err)
  }

  if err = obj.initFile(newFile); err != nil {
    newFile.Close()
    return nil, fmt.Errorf("failed to init file | %w", err)
  }

  return obj, nil
}
