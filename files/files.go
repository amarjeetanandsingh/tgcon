package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type files struct {
	RecursiveOp bool
	CurDir      string
}

func New(isRecursive bool, curDir string) *files {
	f := &files{
		RecursiveOp: isRecursive,
		CurDir:      curDir,
	}
	return f
}

func (f files) DeleteFilesWithSuffix(suffix string) error {
	return deleteFilesWithSuffix(f.RecursiveOp, f.CurDir, suffix)
}

func deleteFilesWithSuffix(isRecursive bool, dir, suffix string) error {
	files, err := listFilesInDir(dir, func(fileName string) bool {
		return strings.HasSuffix(fileName, suffix)
	})
	if err != nil {
		return fmt.Errorf("files.deleteFilesWithSuffix():: error geting files. dir = %s, suffix=%s :: %w", dir, suffix, err)
	}

	// delete files in current directory
	for _, file := range files {
		filePath := path.Join(dir, file.Name())
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("files.deleteFilesWithSuffix():: error deleting %s file :: %w", filePath, err)
		}
	}

	// Stop processing if it's not recursive
	if !isRecursive {
		return nil
	}

	// ******* process recursive call *********

	subDirs, err := listDirs(dir, func(dirName string) bool {
		// ignore hidden dirs
		return len(dirName) > 0 && dirName[0] != '.'
	})
	if err != nil {
		return fmt.Errorf("files.deleteFilesWithSuffix():: error getting dirs. dir = %s :: %w", dir, err)
	}

	for _, subDir := range subDirs {
		subDirPath := path.Join(dir, subDir.Name())
		if err := deleteFilesWithSuffix(isRecursive, subDirPath, suffix); err != nil {
			return err
		}
	}
	return nil
}

// List returns the list of files and directories in the given dir
func (f files) ListFiles(checkName func(string) bool) ([]os.FileInfo, error) {
	return listFilesInDir(f.CurDir, checkName)
}

func listFilesInDir(dir string, checkName func(string) bool) ([]os.FileInfo, error) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading path: %s :: %w", dir, err)
	}

	files := []os.FileInfo{}

	for _, f := range ff {
		if !f.IsDir() && checkName(f.Name()) {
			files = append(files, f)
		}
	}
	return files, nil
}

func (f files) ListDirs(checkName func(string) bool) ([]os.FileInfo, error) {
	return listDirs(f.CurDir, checkName)
}

func listDirs(dir string, checkName func(string) bool) ([]os.FileInfo, error) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("files.listDirs():: error reading path: %s :: %w", dir, err)
	}

	dirs := []os.FileInfo{}

	for _, f := range ff {
		if f.IsDir() && checkName(f.Name()) {
			dirs = append(dirs, f)
		}
	}
	return dirs, nil
}
