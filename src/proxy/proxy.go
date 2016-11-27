package proxy

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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
	//tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)
	//conn, err := net.DialTCP("tcp", nil, tcpAddr)
	for {
		buf := make([]byte, 10240)
		n, err := client.Read(buf)
		if err != nil {
			break
		}

		fmt.Println(string(buf[:n]))
		this.parseResquest(string(buf[:n]))
		client.Write(buf[:n])
	}
	client.Close()
}

func (this *Proxy) parseResquest(cmd string) {
	lines := strings.Split(cmd, "\r\n")
	method := ""
	httpurl := ""
	httpversion := ""
	body := ""
	headers := make(map[string]string)
	for index, line := range lines {
		fmt.Println("index:", index, "line:", line)
		if index == 0 {
			method, httpurl, httpversion = this.parseMethodUrl(line)
		} else {
			key, val := this.parseHeader(line)
			if key != "" {
				headers[key] = val
			}
		}
	}
	fmt.Println("method:", method)
	fmt.Println("httpurl:", httpurl)
	fmt.Println("httpversion:", httpversion)
	fmt.Println("headers:", headers)
	this.sendResquest(method, httpurl, httpversion, headers, body)
}

func (this *Proxy) parseMethodUrl(line string) (method string, url string, httpversion string) {
	fields := strings.Split(line, " ")
	if len(fields) == 0 {
		return "", "", ""
	}
	fmt.Println("type:", reflect.TypeOf(fields))
	if len(fields) != 3 {
		return "", "", ""
	}
	return fields[0], fields[1], fields[2]
}

func (this *Proxy) parseHeader(line string) (key string, val string) {
	fields := strings.Split(line, ":")
	if len(fields) == 0 {
		return "", ""
	}
	if len(fields) != 2 {
		return "", ""
	}
	return fields[0], fields[1]
}

func (this *Proxy) sendResquest(method string, httpurl string, httpversion string, headers map[string]string, body string) {
	switch method {
	case "GET":
		this.RequestGet(httpurl, headers)
	case "POST":
		this.RequestPost(httpurl, headers)
	case "COMMENT":
		fmt.Println("COMMENT")
	}
}

func (this *Proxy) RequestGet(httpurl string, headers map[string]string) {
	//生成client 参数为默认
	client := &http.Client{}
	//提交请求
	req, err := http.NewRequest("GET", httpurl, nil)

	if err != nil {
		panic(err)
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	//处理返回结果
	response, _ := client.Do(req)
	body := response.Body
	respheader := response.Header

	fmt.Println("body", body)
	fmt.Println("respheader", respheader)

	reponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(reponseData), err)
}

func (this *Proxy) RequestPost(httpurl string, headers map[string]string) {
	//生成client 参数为默认
	client := &http.Client{}

	//请求数据
	formData := url.Values{}
	formData.Add("username", "ruansheng")
	formData.Add("password", "123")
	data := formData.Encode()

	//提交请求
	req, err := http.NewRequest("POST", httpurl, strings.NewReader(data))

	if err != nil {
		panic(err)
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	//处理返回结果
	response, _ := client.Do(req)
	body := response.Body
	respheader := response.Header

	fmt.Println("body", body)
	fmt.Println("respheader", respheader)

	reponseData, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(reponseData), err)
}
