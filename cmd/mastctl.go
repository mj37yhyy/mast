package main

import cmd "github.com/mj37yhyy/mast/cmd/pkg"

func main() {

	fErr := cmd.Execute()
	if fErr != nil {
		panic(fErr)
	}
}
