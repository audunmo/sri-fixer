package filer

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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

	err = os.Mkdir(filepath.Join(dir, "/test/dist"), 0755)
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

func TestDirNoGitignore(t *testing.T) {
	defer teardown()
	err := createFolderStructure()
	if err != nil {
		t.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	results, err := Dir(pwd, []string{})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Print(results)


	expected := []string{filepath.Join(pwd, "/test/app/static", "index.html"), filepath.Join(pwd, "/test/app/internal", "internal.html")}
  sort.Strings(expected)
  sort.Strings(results)
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("Expected to find \n%v, but got \n%v", expected, results)
	}
}

func TestDirWithGitignore(t *testing.T) {
	defer teardown()
	err := createFolderStructure()
	if err != nil {
		t.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	results, err := Dir(pwd, []string{"test/app/internal"})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Print(results)


	expected := []string{filepath.Join(pwd, "/test/app/static", "index.html")}
  sort.Strings(expected)
  sort.Strings(results)
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("Expected to find \n%v, but got \n%v", expected, results)
	}
}

