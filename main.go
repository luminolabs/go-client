package main

import (
	"fmt"
	"lumino/cmd"
)

func main() {
	fmt.Println("Initializing interfaces")
	cmd.InitializeInterfaces()
	fmt.Println("Executing command")
	cmd.Execute()
}
