//go:generate goyacc -l -o parser.go parser.y
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	l := newLexer(bufio.NewReader(os.Stdin))
	yyParse(l)
	fmt.Printf("%v\n", l.result)
}
