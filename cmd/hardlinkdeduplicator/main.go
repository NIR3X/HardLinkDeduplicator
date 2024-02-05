package main

import (
	"flag"

	"github.com/NIR3X/hardlinkdeduplicator"
)

func main() {
	a := flag.Bool("a", false, "remove all duplicates (default is to keep one extra copy of file)")
	d := flag.Bool("d", false, "deduplicate files not just report duplicates")
	s := flag.Int64("s", 1024, "minimum file size to consider for deduplication (in bytes)")
	v := flag.Bool("v", false, "verbose output")
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		flag.Usage()
		return
	}
	hardlinkdeduplicator.Deduplicate(path, *a, *d, *s, *v)
}
