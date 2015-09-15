package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/rthornton128/vm/lib"
)

func load(b []byte) *vm.Object {
	o, err := vm.ScanObject(b)
	if err != nil {
		log.Fatal(err)
	}
	return o
}

func main() {
	out := flag.String("o", "out.vm", "program name")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	prog := vm.NewObject()
	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		if err := prog.Merge(load(b)); err != nil {
			log.Fatal(err)
		}
		f.Close()
	}

	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Write(prog.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if n != len(prog.Bytes()) {
		log.Fatal("failed to write all data to output")
	}
}
