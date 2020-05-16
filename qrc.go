package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"rsc.io/qr"
)

func exitError(msg string) {
	fmt.Fprintf(flag.CommandLine.Output(), "qrc: %s\n", msg)
	os.Exit(1)
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options] [text]\n", os.Args[0],
		)
		fmt.Fprintln(flag.CommandLine.Output(),
			"If text is not given as an argument, input is taken from stdin until EOF.",
		)
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
	}
}

func main() {
	var o, e string

	flag.StringVar(&o, "o", "-",
		"output file",
	)
	flag.StringVar(&e, "e", "L",
		"error correction level: L, M, Q, or H",
	)
	flag.Parse()

	level := map[string]qr.Level{"L": qr.L, "M": qr.M, "Q": qr.Q, "H": qr.H}[e]
	if level == 0 && e != "L" {
		exitError("invalid error correction level")
	}

	var data string
	switch flag.NArg() {
	case 0:
		// A QR code can contain at most 7089 bytes (version 40, L error
		// correction, and numeric encoding), so a 8192 byte buffer is
		// certainly sufficient.
		buf := make([]byte, 8192)
		n := 0
		for n < 8192 {
			r, err := os.Stdin.Read(buf[n:])
			if err == io.EOF {
				break
			} else if err != nil {
				exitError(err.Error())
			}
			n += r
		}
		// If we exit the loop and we didn't reach EOF, qr.Encode will return
		// an error (text too long), so we don't need to check for that case.
		data = string(buf[:n])
	case 1:
		data = flag.Arg(0)
	default:
		exitError("too many arguments")
	}

	code, err := qr.Encode(data, level)
	if err != nil {
		exitError(err.Error())
	}

	out := os.Stdout
	if o != "-" {
		f, err := os.Create(o)
		if err != nil {
			exitError(err.Error())
		}
		defer f.Close()
		out = f
	}

	_, err = out.Write(code.PNG())
	if err != nil {
		exitError(err.Error())
	}
}
