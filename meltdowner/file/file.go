package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
)

func CreateDir(dirname string) {
	if err := os.Mkdir(dirname, 0777); err != nil {
		fmt.Println(err)
	}
}
func CreateFile(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
func LoadFileContents(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
