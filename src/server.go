package main

import (
	"fmt"
	"proxy"
)

func main() {
	server := new(proxy.Proxy)
	err := server.NewProxy("172.16.64.156", 9876)
	if err != nil {
		fmt.Println(err)
	}
	server.Run()
}
