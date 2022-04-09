package filer

import (
	"os"
	"path/filepath"
)

func Dir(pwd string) ([]string, error) {
  var filepaths []string
  err := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error{
    if err != nil {
      return err
    }

    ext := filepath.Ext(path)
    if ext == ".html" {
      filepaths = append(filepaths, path)
    }

    return nil
  })

  if err != nil {
    return []string{}, err
  }

  return filepaths, nil
}

func Read(file string) ([]byte, error){
  return os.ReadFile(file)
}

func Write(markup, filePath string) error {
  data := []byte(markup)
  return os.WriteFile(filePath, data, 0644)
}
