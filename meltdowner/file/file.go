package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"path/filepath"

	"github.com/otiai10/copy"
)

func GetMarkdownPaths(sourceDir string) []string {
	dir, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		fmt.Println(err)
	}

	var paths []string
	for _, file := range dir {
		paths = append(paths, filepath.Join(sourceDir, file.Name()))
	}

	return paths
}
func CopyDir(srcPath, dstPath string) {
	if err := copy.Copy(srcPath, dstPath); err != nil {
		fmt.Println(err)
	}
}
func CreateDir(dirname string) {
	if err := os.Mkdir(dirname, 0777); err != nil {
		log.Fatal(err)
	}
}
func CreateFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
func RemoveFile(filename string) {
	if err := os.Remove(filename); err != nil {
		fmt.Println(err)
	}
}
func RemoveDir(dirname string) {
	if err := os.RemoveAll(dirname); err != nil {
		log.Fatal(err)
	}
}
func LoadFileContents(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return content
}
func MoveFile(dstPath string, srcPath string) {
	if err := os.Rename(dstPath, srcPath); err != nil {
		fmt.Println(err)
	}
}

func IsExistPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
