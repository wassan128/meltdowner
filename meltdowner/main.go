package main

import (
	"fmt"
	"flag"
	"net/http"

	"github.com/wassan128/meltdowner/meltdowner/build"
)

func main() {
	var (
		genFlag = flag.Bool("g", false, "generate before serve")
	)
	flag.Parse()

	if *genFlag {
		build.Run()
	}

	http.Handle("/", http.FileServer(http.Dir("public")))

	fmt.Println("[*] public/ is being served on http://localhost:5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Println(err)
		return
	}
}

