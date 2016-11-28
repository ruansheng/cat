package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Http struct {
	Data               string              //原数据
	Method             string              //请求方法
	HttpUrl            string              //请求URL
	HttpVersion        string              //来源请求HTTP版本 example: HTTP/1.1
	SourceHeaders      map[string]string   //来源请求头
	SourceFormData     string              //来源请求体
	ResponseStatus     string              //响应状结果e.g. "200 OK"
	ResponseStatusCode int                 //响应状态码e.g. 200
	ResponseProto      string              //响应状协议版本e.g. "HTTP/1.0"
	ResponseHeaders    map[string][]string //响应头
	ResponseData       string              //响应数据
	ReturnData         string              //返回代理请求后的数据
}

func (this *Http) Send() {
	this.parseResquest()
	this.sendResquest()
}

func (this *Http) parseResquest() {
	fields := strings.Split(this.Data, "\r\n\r\n")
	if len(fields) == 2 {
		lines := strings.Split(fields[0], "\r\n")
		this.SourceHeaders = make(map[string]string)
		for index, line := range lines {
			if index == 0 {
				this.parseResquestLine(line)
			} else {
				this.parseResquestHeader(line)
			}
		}
		this.parseResquestBody(fields[1])
	}
}

func (this *Http) parseResquestLine(line string) {
	fields := strings.Split(line, " ")
	if len(fields) == 3 {
		this.Method = fields[0]
		this.HttpUrl = fields[1]
		this.HttpVersion = fields[2]
	}
}

func (this *Http) parseResquestHeader(line string) {
	fields := strings.Split(line, ":")
	if len(fields) == 2 {
		this.SourceHeaders[fields[0]] = fields[1]
	}
}

func (this *Http) parseResquestBody(line string) {
	if line != "" {
		this.SourceFormData = line
	}
}

func (this *Http) sendResquest() {
	switch this.Method {
	case "GET", "POST":
		this.Request()
	case "COMMENT":
		fmt.Println("COMMENT")
	default:
		fmt.Println(this.Method)
	}
}

func (this *Http) Request() {
	//生成client 参数为默认
	client := &http.Client{}

	//提交请求
	var req *http.Request
	var err error
	if this.SourceFormData != "" {
		req, err = http.NewRequest(this.Method, this.HttpUrl, strings.NewReader(this.SourceFormData))
	} else {
		req, err = http.NewRequest(this.Method, this.HttpUrl, nil)
	}

	if err != nil {
		panic(err)
	}

	// 设置请求头
	for key, val := range this.SourceHeaders {
		req.Header.Set(key, val)
	}

	//处理返回结果
	response, _ := client.Do(req)

	this.ResponseStatusCode = response.StatusCode
	this.ResponseStatus = response.Status
	this.ResponseProto = response.Proto
	this.ResponseHeaders = response.Header

	reponseData, _ := ioutil.ReadAll(response.Body)
	this.ResponseData = string(reponseData)
}

func (this *Http) GetReturnData() string {
	line := this.getResponseLine()
	header := this.getResponseHeader()
	return fmt.Sprintf("%s\r\n%s\r\n\r\n%s", line, header, this.ResponseData)
}

func (this *Http) getResponseLine() string {
	return fmt.Sprintf("%s %s", this.ResponseProto, this.ResponseStatus)
}

func (this *Http) getResponseHeader() string {
	var headers []string
	for key, val := range this.ResponseHeaders {
		headers = append(headers, fmt.Sprintf("%s:%s", key, val[0]))
	}
	return strings.Join(headers, "\r\n")
}

func test() {
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
}
