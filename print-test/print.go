package main

import (
	"fmt"
	"os"
)

func main() {
	testA()
}

func testA() { 
	fmt.Println("count: ", len(os.Args))

	for i, v := range os.Args {
		fmt.Printf("args[%d] -> %s\n", i, v)
	}
}