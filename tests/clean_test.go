package tests

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/amarjeetanandsingh/tgcon/files"
)

func TestCleanCmd(t *testing.T) {
	// test each flags
	t.Run("dir", cleanCmd_Dir)
	t.Run("verbose", cleanCmd_Verbose)
	t.Run("recursive", cleanCmd_Recursive)
	t.Run("verbose recursive dir", cleanCmd_Verbose_Dir_Recursive)
}

func cleanCmd_Verbose(t *testing.T) {

	// generate a tmp file to be cleaned by tgcon clean
	generatedFile, err := ioutil.TempFile(".", "*_tgconst_gen.go")
	if err != nil {
		t.Error(err)
	}
	generatedFile.Close()
	defer os.Remove(generatedFile.Name())

	// run clean command for current dir
	cleanCmd, out := setupRootCmd("clean", "--verbose")
	if err := cleanCmd.Execute(); err != nil {
		t.Error(err)
	}

	if strings.Count(out.String(), "_tgconst_gen.go") == 0 {
		t.Errorf("no _tgconst_gen.go deleted in current dir")
	}
}

func cleanCmd_Recursive(t *testing.T) {

	// copy files to tmp dir
	dir, cleanup, err := copyDirToTmp("testdata/onlyGenerated")
	if err != nil {
		t.Error(err)
	}
	defer cleanup()

	// run clean command for dir
	cleanCmd, _ := setupRootCmd("clean", "--recursive", "--dir="+dir)
	if err := cleanCmd.Execute(); err != nil {
		t.Error(err)
	}

	// search generated file in dir
	generatedFiles, err := files.ListFilesInDirRecursive(dir, func(fileName string) bool {
		return strings.HasSuffix(fileName, "_tgconst_gen.go")
	})
	if err != nil {
		t.Error(err)
	}
	if len(generatedFiles) != 0 {
		t.Error("generated files found after clean")
	}
}

func cleanCmd_Dir(t *testing.T) {

	// copy files to tmp dir
	dir, cleanup, err := copyDirToTmp("testdata/onlyGenerated")
	if err != nil {
		t.Error(err)
	}
	defer cleanup()

	// run clean command for dir
	cleanCmd, _ := setupRootCmd("clean", "--dir="+dir)
	if err := cleanCmd.Execute(); err != nil {
		t.Error(err)
	}

	// search generated file in dir
	generatedFiles, err := files.ListFilesInDir(dir, func(fileName string) bool {
		return strings.HasSuffix(fileName, "_tgconst_gen.go")
	})
	if err != nil {
		t.Error(err)
	}
	if len(generatedFiles) != 0 {
		t.Error("generated files found after clean")
	}
}

func cleanCmd_Verbose_Dir_Recursive(t *testing.T) {

	// copy files to tmp dir
	dir, cleanup, err := copyDirToTmp("testdata/onlyGenerated")
	if err != nil {
		t.Error(err)
	}
	defer cleanup()

	// run clean command for dir
	cleanCmd, out := setupRootCmd("clean", "--verbose", "--dir="+dir, "--recursive")
	if err := cleanCmd.Execute(); err != nil {
		t.Error(err)
	}

	// search generated file in dir
	generatedFiles, err := files.ListFilesInDir(dir, func(fileName string) bool {
		return strings.HasSuffix(fileName, "_tgconst_gen.go")
	})
	if err != nil {
		t.Error(err)
	}
	if len(generatedFiles) != 0 {
		t.Error("generated files found after clean")
	}

	if strings.Count(out.String(), "_tgconst_gen.go") == 0 {
		t.Error("cleaned filed not logged.")
	}
}
