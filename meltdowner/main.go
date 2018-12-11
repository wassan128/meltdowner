package main

import (
	"fmt"
	"net/http"

	"github.com/wassan128/meltdowner/meltdowner/build"
)

func main() {
	build.Run()

	http.Handle("/", http.FileServer(http.Dir("public")))

	fmt.Println("[*] public/ is being served on http://localhost:5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Println(err)
		return
	}
}

