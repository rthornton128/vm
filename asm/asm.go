package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	vm "github.com/rthornton128/vm/lib"
)

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	src, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// encoding
	buf := new(bytes.Buffer)
	e := vm.NewEncoder(buf)
	if err := e.Encode(src); err != nil {
		log.Fatal(err)
	}

	f, err = os.Create(flag.Arg(0) + ".vm") // TODO fix extention handling
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Println(buf.Bytes())
	f.Write(buf.Bytes())
	// linking
	//if err := vm.Link(buf.Bytes(), f); err != nil {
	///log.Fatal(err)
	//}
}
