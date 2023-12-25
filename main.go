package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("foi")
	vals, err := os.ReadDir("/home")
	if err != nil {
		fmt.Printf("Invalid path received on readfile, %v", err)
		return
	}
	net.IPv4bcast

	fmt.Print(vals)

}
