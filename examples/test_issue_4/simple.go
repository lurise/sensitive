package main

import (
	"fmt"

	"github.com/lurise/sensitive"
)

func main() {
	filter := sensitive.New()
	filter.LoadWordDict("../../dict/dict.txt")
	fmt.Println(filter.Replace("xC4x", '*'))
}
