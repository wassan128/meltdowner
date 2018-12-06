package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
)

func LoadMarkdown(filename string) []byte {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return md
}

func CreateDirForPublish() {
	if err := os.Mkdir("public", 0777); err != nil {
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
func MoveFile(dstPath string, srcPath string) {
	if err := os.Rename(dstPath, srcPath); err != nil {
		fmt.Println(err)
	}
}
