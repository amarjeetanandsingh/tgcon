package tests

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/amarjeetanandsingh/tgconst/cmd"
	"github.com/spf13/cobra"
)

func emptyRun(*cobra.Command, []string) {}

func setupCleanCmd(args ...string) (*cobra.Command, *bytes.Buffer) {
	cmd := cmd.NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs(args)
	return cmd, buf
}

func setupGenCmd(args ...string) *cobra.Command {
	genCmd := cmd.NewGenCmd()
	genCmd.SetArgs(args)
	return genCmd
}

func copyDirToTmp(src string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "tgconst")
	if err != nil {
		return "", err
	}

	if err := copyDirRecursive(src, tmpDir); err != nil {
		return "", err
	}
	return tmpDir, nil
}

func copyDirRecursive(src, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDirRecursive(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
