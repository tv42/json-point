package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	jp "github.com/dustin/go-jsonpointer"
	"io"
	"log"
	"os"
	"path"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s [-pretty] QUERY.. <FILE\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -list <FILE\n", os.Args[0])
	flag.PrintDefaults()
}

var list = flag.Bool("list", false, "list all possible query strings")
var pretty = flag.Bool("pretty", false, "output in non-JSON, pretty format; not machine readable")

func main() {
	prog := path.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = Usage
	flag.Parse()

	var all_matched = true

	dec := json.NewDecoder(os.Stdin)

	switch {

	case flag.NArg() == 0 && !*list:
		Usage()

	case *list && flag.NArg() > 0:
		log.Fatal("cannot combine -list with query")

	case *list && *pretty:
		log.Fatal("cannot combine -list with -pretty")

	case *list:
		for {
			var raw json.RawMessage

			err := dec.Decode(&raw)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("reading json: %v", err)
			}
			p, err := jp.ListPointers(raw)
			if err != nil {
				log.Fatalf("listing pointers: %v", err)
			}
			for _, s := range p {
				fmt.Println(s)
			}
		}

	case !*list:
		for {
			var raw json.RawMessage

			err := dec.Decode(&raw)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("reading json: %v", err)
			}

			for _, query := range flag.Args() {
				match, err := jp.Find(raw, query)
				if err != nil {
					log.Fatalf("matching: %v", err)
				}

				if match == nil {
					// did not match anything;
					// signal with an empty line
					all_matched = false
					fmt.Println()
					continue
				}

				if *pretty {
					var data interface{}
					err = json.Unmarshal(match, &data)
					if err != nil {
						log.Fatalf("parsing json: %v", err)
					}
					switch data.(type) {
					case string, int, float64:
						fmt.Println(data)
						continue
					default:
						// handle difficult types below
					}
				}

				// canonicalize whitespace
				var buf bytes.Buffer
				err = json.Compact(&buf, match)
				if err != nil {
					log.Fatalf("serializing json: %v", err)
				}
				buf.WriteByte('\n')
				_, err = buf.WriteTo(os.Stdout)
				if err != nil {
					log.Fatalf("writing json: %v", err)
				}
			}
		}

	default:
		panic("unreachable")
	}

	if !all_matched {
		os.Exit(1)
	}
}
