package main

import (
	"fmt"
	"gar"
)

func main() {
	fmt.Println("gar a ar command line app.")
	target := "/home/jesse/work/artest/lib/gar.a"
	src1 := "abc.o"
	src2 := "abc.a"
	arf, _ := gar.Open(target)

	arf.List()
	arf.Append(src1)
	arf.Append(src2)
	arf.Close()
}
