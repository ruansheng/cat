package proxy

import (
	"fmt"
	"net"
)

type Proxy struct {
	ip       string
	port     int
	listener *net.TCPListener
}

func (this *Proxy) NewProxy(ip string, port int) error {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
	if err != nil {
		fmt.Println("listen error: ", err.Error())
		return err
	}
	fmt.Println("init done...")

	this.ip = ip
	this.port = port
	this.listener = listener
	return nil
}

func (this *Proxy) Run() {
	for {
		client, err := this.listener.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err.Error())
			continue
		}
		go this.Channal(client)
	}
}

func (this *Proxy) Channal(client *net.TCPConn) {
	for {
		buf := make([]byte, 10240)
		n, err := client.Read(buf)
		if err != nil {
			break
		}

		http := new(Http)
		http.Data = string(buf[:n])
		http.Send()
		client.Write([]byte(http.ReturnData))
		break
	}
	client.Close()
}
