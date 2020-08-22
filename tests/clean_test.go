package tests

import (
	"strings"
	"testing"
)

func TestCleanCmd(t *testing.T) {
	// test each flags
	//t.Run("dir", cleanCmd_Dir)
	t.Run("verbose", cleanCmd_Verbose) // todo: use out
	//t.Run("recursive", cleanCmd_Recursive)
	//t.Run("verbose recursive dir", cleanCmd_Verbose_Dir_Recursive)
}

func cleanCmd_Verbose(t *testing.T) {
	if err := copyDirRecursive("./tests/testdata/onlyGenerated", "."); err != nil {
		t.Error(err)
	}
	cleanCmd, out := setupCleanCmd("clean", "--verbose")

	// run clean command for current dir
	if err := cleanCmd.Execute(); err != nil {
		t.Error(err)
	}

	if count := strings.Count(out.String(), "_tgconst_gen.go"); count == 0 {
		t.Errorf("no _tgconst_gen.go deleted in cur dir")
	}
}

func cleanCmd_Recursive(t *testing.T) {
	// TODO
}

func cleanCmd_Dir(t *testing.T) {
	// TODO
}

func cleanCmd_Verbose_Dir_Recursive(t *testing.T) {
	// todo
}
