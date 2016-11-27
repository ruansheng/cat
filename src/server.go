package main

import (
	"fmt"
	"proxy"
)

func main() {
	server := new(proxy.Proxy)
	err := server.NewProxy("192.168.1.104", 9090)
	if err != nil {
		fmt.Println(err)
	}
	server.Run()
}
