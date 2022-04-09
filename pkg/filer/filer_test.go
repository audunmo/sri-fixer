package filer

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func createFolderStructure() error{
  dir, err := os.Getwd()
  if err != nil {
    return err
  }

  err = os.Mkdir(filepath.Join(dir, "/test/"), 0755)
  if err != nil {
    return err
  }

  err = os.Mkdir(filepath.Join(dir, "/test/app"), 0755)
  if err != nil {
    return err
  }

  err = os.Mkdir(filepath.Join(dir, "/test/app/static"), 0755)
  if err != nil {
    return err
  }

  _, err = os.Create(filepath.Join(dir, "/test/app/static", "index.html"))
  if err != nil {
    return err
  }

  return nil
}

func teardown() error {
  pwd, err := os.Getwd()
  if err != nil {
    return err
  }

  err = os.RemoveAll(filepath.Join(pwd, "/test"))
  if err != nil {
    return err
  }

  return nil
}

func TestDir(t *testing.T) {
  defer teardown()
  err := createFolderStructure()
  if err != nil {
    t.Fatal(err)
  }

  pwd, err := os.Getwd()
  if err != nil {
    t.Fatal(err)
  }

  results, err := Dir(pwd)
  if err != nil {
    t.Fatal(err)
  }

  fmt.Print(results)

  indexNestedPath := filepath.Join(pwd, "/test/app/static", "index.html")
  if !reflect.DeepEqual(results, []string{indexNestedPath}) {
    t.Fatal("couldn't find nested index.html")
  }
}
