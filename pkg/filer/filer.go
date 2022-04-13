package filer

import (
	"os"
	"path/filepath"
	"strings"
)

func shouldIgnore(pwd, path string, patterns []string) (bool, error) {
  for _, p := range patterns {
    pattern := filepath.Join(pwd, p)
    m, err := filepath.Match(pattern, path)
    if err != nil {
      return false, err
    }

    if m {
      return true, nil
    }
  }

  return false, nil
}


// Reads .gitignore from pwd. If there's no gitignore, it returns an error
func ReadGitIgnore(pwd string) ([]string, error){
  data, err := Read(filepath.Join(pwd, ".gitignore")) 
  if err != nil {
    return []string{}, err
  }

  return strings.Split(data, "\n"), nil
}

// Dir gets all html files in the tree of the directory passed to it. Dir should ignore all glob ignorePatterns
// passed to ignorePatterns
func Dir(pwd string, ignorePatterns []string) ([]string, error) {
	var filepaths []string

  err := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

    ignore, err := shouldIgnore(pwd, path, ignorePatterns)
    if err != nil {
      return err
    }

    if ignore {
      // If ignore triggered on a dir, don't walk the dir, since all children should be skipped anyway
      if info.IsDir() {
        return filepath.SkipDir
      }
      return nil
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

func Read(file string) (string, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

func Write(markup, filePath string) error {
	data := []byte(markup)
	return os.WriteFile(filePath, data, 0644)
}
