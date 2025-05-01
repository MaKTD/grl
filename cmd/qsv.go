package main

import (
	"fmt"
	"grl/internal/qsv"
	"os"
)

func main() {
	if err := qsv.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
