package main

import (
	"bytes"
	"flag"
	"go/token"
	"log"
	"os"

	vm "github.com/rthornton128/vm/lib"
)

func main() {
	flag.Parse()

	in, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	// encoding
	fset := token.NewFileSet()
	info, err := in.Stat()
	if err != nil {
		log.Fatal(err)
	}
	f := fset.AddFile(flag.Arg(0), -1, int(info.Size()))
	buf := new(bytes.Buffer)
	e := vm.NewEncoder(f, buf)
	if err := e.Encode(in); err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(flag.Arg(0) + ".o") // TODO fix extention handling
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	//fmt.Println(buf.Bytes())
	out.Write(buf.Bytes())
}
