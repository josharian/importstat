package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/josharian/importstat/github.com/kr/fs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: importstat <dir> | sort | uniq -c | sort -n -r")
		os.Exit(1)
	}
	walker := fs.Walk(os.Args[1])

	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Printf("Error during filesystem walk: %v\n", err)
			continue
		}

		if walker.Stat().IsDir() || !strings.HasSuffix(walker.Path(), ".go") {
			continue
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, walker.Path(), nil, parser.ImportsOnly)
		if err != nil {
			// don't print err here; it is too chatty, due to (un?)surprising
			// amounts of broken code in the wild
			continue
		}

		for _, imp := range f.Imports {
			fmt.Println(imp.Name, imp.Path.Value)
		}
	}
}
