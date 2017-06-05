package main

import (
	"io"
	"os"
)

func main() {
	il := &itemList{}
	io.Copy(il, os.Stdin)
	il.Flush()
	il.Sort()
	io.Copy(os.Stdout, il)
}
