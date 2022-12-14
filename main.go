package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	lflag    = flag.Bool("l", false, "print long format")
	maxdepth = flag.Int("max", -1, "max depth")
	uflag    = flag.Bool("unix", false, "with -l, print unix time")
	iflag    = flag.Bool("i", false, "ignore dotfiles")

	stdout = bufio.NewWriter(os.Stdout)
)

func main() {
	log.SetPrefix("walk: ")
	log.SetFlags(0)
	flag.Parse()

	exitCode := 0
	doWalk := func(d string) {
		if err := walk(d); err != nil {
			log.Println(err)
			exitCode = 1
		}
	}

	if flag.NArg() == 0 {
		doWalk(".")
	} else {
		for _, dir := range flag.Args() {
			doWalk(dir)
		}
	}

	stdout.Flush()
	os.Exit(exitCode)
}
