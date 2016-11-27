package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Http struct {
	Data            string            //原数据
	Method          string            //请求方法
	HttpUrl         string            //请求URL
	HttpVersion     string            //来源请求HTTP版本 example: HTTP/1.1
	SHeaders        map[string]string //来源请求头
	SFormData       string            //来源请求体
	ResponseHeaders map[string]string //响应头
	ResponseData    string            //响应数据
	ReturnData      string            //返回代理请求后的数据
}

func (this *Http) Send() {
	this.parseResquest()
	this.sendResquest()
}

func (this *Http) parseResquest() {
	lines := strings.Split(this.Data, "\r\n")
	this.SHeaders = make(map[string]string)
	for index, line := range lines {
		if index == 0 {
			this.parseResquestLine(line)
		} else {
			this.parseResquestHeader(line)
		}
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
		this.SHeaders[fields[0]] = fields[1]
	}
}

func (this *Http) parseResquestBody(line string) {
	//请求数据
	formData := url.Values{}
	formData.Add("username", "ruansheng")
	formData.Add("password", "123")
	data := formData.Encode()

	this.SFormData = line
	this.SFormData = data
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
	if this.SFormData != "" {
		req, err = http.NewRequest(this.Method, this.HttpUrl, strings.NewReader(this.SFormData))
	} else {
		req, err = http.NewRequest(this.Method, this.HttpUrl, nil)
	}

	if err != nil {
		panic(err)
	}

	// 设置请求头
	for key, val := range this.SHeaders {
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
