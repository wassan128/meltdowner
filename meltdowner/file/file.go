package file

import (
	"io/ioutil"
	"os"
	"strings"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/wassan128/meltdowner/meltdowner/util"
)

func GetMarkdownPaths(sourceDir string) []string {
	dir, err := ioutil.ReadDir(sourceDir)
	util.ExitIfError(err)

	var paths []string
	for _, file := range dir {
		if strings.HasSuffix(file.Name(), ".md") {
			paths = append([]string{filepath.Join(sourceDir, file.Name())}, paths...)
		}
	}

	return paths
}
func CopyDir(srcPath, dstPath string) {
	err := copy.Copy(srcPath, dstPath)
	util.WarningIfError(err)
}
func CreateDir(dirname string) {
	err := os.Mkdir(dirname, 0777)
	util.ExitIfError(err)
}
func CreateFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
	util.WarningIfError(err)

	return file
}
func RemoveFile(filename string) {
	err := os.Remove(filename)
	util.WarningIfError(err)
}
func RemoveDir(dirname string) {
	err := os.RemoveAll(dirname)
	util.WarningIfError(err)
}
func LoadFileContents(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	util.ExitIfError(err)

	return content
}
func MoveFile(dstPath string, srcPath string) {
	err := os.Rename(dstPath, srcPath)
	util.ExitIfError(err)
}

func IsExistPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
