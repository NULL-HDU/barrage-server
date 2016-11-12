package main

import (
	"barrage-server/libs/cmdface"
	"fmt"
)

func main() {

	cmdface.AddCommand("test3", "test3", nil)
	cmdface.AddCommand("test3", "test3", nil)
	cmdface.AddCommand("test3", "test3", nil)

	for {
		if err := cmdface.InputAndRunCommand(">>> "); err != nil {
			cmdface.Show(fmt.Sprintf("%s\n", err.Error()))
		}
	}
}
