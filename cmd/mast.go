package main

import (
	cmd "github.com/mj37yhyy/mast/cmd/pkg"
	_ "github.com/mj37yhyy/mast/pkg/kubernetes"
)

func main() {
	fErr := cmd.Execute()
	if fErr != nil {
		panic(fErr)
	}
}
