package infrastructure

import (
	"io/ioutil"
	"os"

	"github.com/otiai10/copy"
	"github.com/wassan128/meltdowner/meltdowner/util"
)

func CopyDir(srcPath, dstPath string) {
	err := copy.Copy(srcPath, dstPath)
	util.WarningIfError(err)
}

func CreateDir(dirname string) {
	err := os.Mkdir(dirname, 0777)
	util.ExitIfError(err)
}

func DeleteDir(dirname string) {
	err := os.RemoveAll(dirname)
	util.WarningIfError(err)
}

func CreateFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	util.WarningIfError(err)

	return file
}

func ReadFile(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	util.ExitIfError(err)

	return content
}

func MoveFile(dstPath string, srcPath string) {
	err := os.Rename(dstPath, srcPath)
	util.ExitIfError(err)
}

func DeleteFile(filename string) {
	err := os.Remove(filename)
	util.WarningIfError(err)
}

func IsExistPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
