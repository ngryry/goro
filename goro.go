package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/ngryry/goro/internal"
	"golang.org/x/tools/imports"
)

var (
	src        = flag.String("s", "", "path to input Go source file")
	dst        = flag.String("d", "", "path to output file")
	enabledTag = flag.Bool("t", false, "enable tag mode")
)

func main() {
	flag.Parse()

	if *src == "" || *dst == "" {
		log.Fatal("[goro]: must specify -s and -d")
	}

	p := internal.Parser{EnabledTag: *enabledTag}
	w := internal.Writer{}

	log.Printf("[goro]: parse %s", *src)
	f, err := p.Parse(*src)
	if err != nil {
		log.Fatalf("[goro]: failed to parse src: %s", err.Error())
	}

	log.Printf("[goro]: generate %s", *dst)
	if err := w.Write(f, *dst); err != nil {
		log.Fatalf("[goro]: failed to write file: %s", err.Error())
	}

	log.Printf("[goro]: format %s", *dst)
	if err := format(); err != nil {
		log.Fatalf("[goro]: failed to format file: %s", err.Error())
	}

	log.Print("[goro]: complete")
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
