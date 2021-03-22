// impl generates method stubs for implementing an interface.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/josharian/impl/pkg/impl"
)

const usage = `impl [-dir directory] <recv> <iface>

impl generates method stubs for recv to implement iface.

Examples:

impl 'f *File' io.Reader
impl Murmur hash.Hash
impl -dir $GOPATH/src/github.com/josharian/impl Murmur hash.Hash

Don't forget the single quotes around the receiver type
to prevent shell globbing.
`

var (
	flagSrcDir = flag.String("dir", "", "package source directory, useful for vendored code")
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	recv, iface := flag.Arg(0), flag.Arg(1)
	if !impl.ValidReceiver(recv) {
		fatal(fmt.Sprintf("invalid receiver: %q", recv))
	}

	if *flagSrcDir == "" {
		if dir, err := os.Getwd(); err == nil {
			*flagSrcDir = dir
		}
	}

	fns, err := impl.Funcs(iface, *flagSrcDir)
	if err != nil {
		fatal(err)
	}

	// Get list of already implemented funcs
	implemented, err := impl.ImplementedFuncs(fns, recv, *flagSrcDir)
	if err != nil {
		fatal(err)
	}

	src := impl.GenStubs(recv, fns, implemented)
	fmt.Print(string(src))
}

func fatal(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
