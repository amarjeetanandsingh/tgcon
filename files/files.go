package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//TODO think again, if this package is needed

type files struct {
	RecursiveOp bool
	CurDir      string
	Verbose     bool
}

func New(verbose, isRecursive bool, curDir string) *files {
	return &files{
		RecursiveOp: isRecursive,
		Verbose:     verbose,
		CurDir:      curDir,
	}
}

// todo:v2: think of making it transactional
func (f files) DeleteFilesWithSuffix(suffix string) error {
	return f.deleteFilesWithSuffix(f.CurDir, suffix)
}

// todo:v2: optimisation.
// Do not call listFilesInDir & listDirs twice.
// We can get the list of file & folders in a single call.
// todo:v2: we can reuse deleteFileWithSuffix() in case of recursive call
func (f files) deleteFilesWithSuffix(dir, suffix string) error {
	fileNames, err := ListFilesInDir(dir, func(fileName string) bool {
		return strings.HasSuffix(fileName, suffix)
	})
	if err != nil {
		return fmt.Errorf("files.deleteFilesWithSuffix():: error geting files. dir = %s, suffix=%s :: %w", dir, suffix, err)
	}

	// delete files in current directory
	for _, fileName := range fileNames {
		filePath := path.Join(dir, fileName)
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("files.deleteFilesWithSuffix():: error deleting %s file :: %w", filePath, err)
		}

		// log
		if f.Verbose {
			// todo:v2: use buffered out writer
			fmt.Println("Deleted:", filePath)
		}
	}

	// Stop processing if it's not recursive
	if !f.RecursiveOp {
		return nil
	}

	// ******* process recursive call *********

	subDirs, err := ListDirs(dir, func(dirName string) bool {
		// ignore hidden dirs
		return len(dirName) > 0 && dirName[0] != '.'
	})
	if err != nil {
		return fmt.Errorf("files.deleteFilesWithSuffix():: error getting dirs. dir = %s :: %w", dir, err)
	}

	for _, subDir := range subDirs {
		subDirPath := path.Join(dir, subDir)
		if err := f.deleteFilesWithSuffix(subDirPath, suffix); err != nil {
			return err
		}
	}
	return nil
}

// List returns the list of files and directories in the given dir
func (f files) ListFiles(checkName func(string) bool) ([]string, error) {
	return ListFilesInDir(f.CurDir, checkName)
}

func ListFilesInDir(dir string, checkName func(string) bool) ([]string, error) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading path: %s :: %w", dir, err)
	}

	var fileNames []string

	for _, f := range ff {
		if !f.IsDir() && checkName(f.Name()) {
			fileNames = append(fileNames, f.Name())
		}
	}
	return fileNames, nil
}

func (f files) ListDirs(checkName func(string) bool) ([]string, error) {
	return ListDirs(f.CurDir, checkName)
}

func ListDirs(dir string, checkName func(string) bool) ([]string, error) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("files.listDirs():: error reading path: %s :: %w", dir, err)
	}

	var dirNames []string

	for _, f := range ff {
		if f.IsDir() && checkName(f.Name()) {
			dirNames = append(dirNames, f.Name())
		}
	}
	return dirNames, nil
}
