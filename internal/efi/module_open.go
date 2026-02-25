package efi

import (
  "path/filepath"

  "github.com/stokescat/efizzer/internal/rawcov"
)

func OpenModuleByGuid(dir string, g Guid) (*Module, error) {

  obj:= new(Module)
  obj.Name = string(g)
  obj.Guid = g

  rawcovFileName:= string(g) + ".rawcov"
  rawcovFilePath:= filepath.Join(dir, rawcovFileName)

  rawcovFile, err:= rawcov.Open(rawcovFilePath)
  if err != nil {
    return nil, err
  }

  obj.file = rawcovFile

  return obj, nil
}

func OpenModuleByName(dir, name string) (*Module, error) {

  obj:= new(Module)
  obj.Name = name

  rawcovFileName:= name + ".rawcov"
  rawcovFilePath:= filepath.Join(dir, rawcovFileName)

  rawcovFile, err:= rawcov.Open(rawcovFilePath)
  if err != nil {
    return nil, err
  }

  obj.file = rawcovFile

  return obj, nil

}
