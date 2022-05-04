package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ngryry/goro/internal"
	"golang.org/x/tools/imports"
)

var (
	src        = flag.String("s", "", "Path to input Go source file")
	dst        = flag.String("d", "", "Path to output file")
	enabledTag = flag.Bool("t", false, "Enable tag mode")
)

const usage = `Goro (go-readonly) is a Go code generator to automate the creation of constructor, getter and setter.
See: https://github.com/ngryry/goro

USAGE:
	goro -s <PATH> -d <PATH> [-t]

OPTIONS:
	-s Path to input Go source file (required)
	-d Path to output file (reqired)
	-t Enable tag mode 
`

func main() {
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	if *src == "" || *dst == "" {
		log.Fatal("[goro]: Must specify -s and -d")
	}

	p := internal.Parser{EnabledTag: *enabledTag}
	w := internal.Writer{}

	f, err := p.Parse(*src)
	if err != nil {
		log.Fatalf("[goro]: Failed to parse src: %s", err.Error())
	}

	if err := w.Write(f, *dst); err != nil {
		log.Fatalf("[goro]: Failed to write file: %s", err.Error())
	}

	if err := format(); err != nil {
		log.Fatalf("[goro]: Failed to format file: %s", err.Error())
	}
}

func format() error {
	b, err := ioutil.ReadFile(*dst)
	if err != nil {
		return err
	}

	res, err := imports.Process(*dst, b, nil)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(*dst, res, 0777); err != nil {
		return err
	}
	return nil
}
