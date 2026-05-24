package main

import (
	"fmt"
	"os"

	"github.com/yunluoicu/code-flow/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
