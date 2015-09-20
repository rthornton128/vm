package main

import (
	"flag"
	"fmt"
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
	fmt.Println("text", o.SecTab[vm.TEXT])
	fmt.Println("symtab", o.SymTab)
	fmt.Println("relocs", o.RelocTab)
	return o
}

func main() {
	out := flag.String("o", "out.vm", "program name")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	o := vm.NewObject()
	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		if err := o.Merge(load(b)); err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
	fmt.Println("merged text", o.SecTab[vm.TEXT])

	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	prog := vm.NewProgram(o)
	fmt.Println("prog text", prog.SecTab[vm.TEXT])
	fmt.Println(prog.Bytes())
	n, err := f.Write(prog.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if n != len(prog.Bytes()) {
		log.Fatal("failed to write all data to output")
	}
}
