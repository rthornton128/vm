package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func load(b []byte) {
	fmt.Println()
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		load(b)
	}
}
