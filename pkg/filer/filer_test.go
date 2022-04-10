package filer

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func createFolderStructure() error {
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

	err = os.Mkdir(filepath.Join(dir, "/test/app/internal"), 0755)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(dir, "/test/app/internal", "internal.html"))
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

	expected := []string{filepath.Join(pwd, "/test/app/static", "index.html"), filepath.Join(pwd, "/test/app/internal", "internal.html")}
	if !reflect.DeepEqual(results, expected) {
		fmt.Printf("Expected to find %v, but got %v", results, expected)
	}
}
