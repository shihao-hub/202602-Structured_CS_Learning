package application

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// HTTPMethod HTTP方法
type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
	PATCH   HTTPMethod = "PATCH"
)

// HTTPVersion HTTP版本
type HTTPVersion string

const (
	HTTP1_0 HTTPVersion = "HTTP/1.0"
	HTTP1_1 HTTPVersion = "HTTP/1.1"
	HTTP2   HTTPVersion = "HTTP/2.0"
)

// HTTPStatusCode HTTP状态码
type HTTPStatusCode int

const (
	StatusOK                  HTTPStatusCode = 200
	StatusCreated             HTTPStatusCode = 201
	StatusAccepted            HTTPStatusCode = 202
	StatusNoContent           HTTPStatusCode = 204
	StatusMovedPermanently    HTTPStatusCode = 301
	StatusFound               HTTPStatusCode = 302
	StatusNotModified         HTTPStatusCode = 304
	StatusBadRequest          HTTPStatusCode = 400
	StatusUnauthorized        HTTPStatusCode = 401
	StatusForbidden           HTTPStatusCode = 403
	StatusNotFound            HTTPStatusCode = 404
	StatusMethodNotAllowed    HTTPStatusCode = 405
	StatusInternalServerError HTTPStatusCode = 500
	StatusNotImplemented      HTTPStatusCode = 501
	StatusBadGateway          HTTPStatusCode = 502
	StatusServiceUnavailable  HTTPStatusCode = 503
)

// String 返回状态码描述
func (code HTTPStatusCode) String() string {
	switch code {
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusAccepted:
		return "Accepted"
	case StatusNoContent:
		return "No Content"
	case StatusMovedPermanently:
		return "Moved Permanently"
	case StatusFound:
		return "Found"
	case StatusNotModified:
		return "Not Modified"
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusInternalServerError:
		return "Internal Server Error"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	default:
		return "Unknown"
	}
}

// HTTPHeader HTTP头部
type HTTPHeader struct {
	Name  string
	Value string
}

// HTTPRequest HTTP请求
type HTTPRequest struct {
	Method  HTTPMethod
	URL     string
	Version HTTPVersion
	Headers []HTTPHeader
	Body    string
}

// HTTPResponse HTTP响应
type HTTPResponse struct {
	Version    HTTPVersion
	StatusCode HTTPStatusCode
	Headers    []HTTPHeader
	Body       string
}

// NewHTTPRequest 创建HTTP请求
func NewHTTPRequest(method HTTPMethod, url string, version HTTPVersion) *HTTPRequest {
	return &HTTPRequest{
		Method:  method,
		URL:     url,
		Version: version,
		Headers: make([]HTTPHeader, 0),
		Body:    "",
	}
}

// AddHeader 添加请求头
func (req *HTTPRequest) AddHeader(name, value string) {
	req.Headers = append(req.Headers, HTTPHeader{Name: name, Value: value})
}

// SetBody 设置请求体
func (req *HTTPRequest) SetBody(body string) {
	req.Body = body
	if body != "" && req.Method == POST {
		req.AddHeader("Content-Length", strconv.Itoa(len(body)))
		req.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	}
}

// ToString 转换为字符串
func (req *HTTPRequest) ToString() string {
	var builder strings.Builder

	// 请求行
	builder.WriteString(fmt.Sprintf("%s %s %s\r\n", req.Method, req.URL, req.Version))

	// 头部
	for _, header := range req.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", header.Name, header.Value))
	}

	// 空行
	builder.WriteString("\r\n")

	// 请求体
	if req.Body != "" {
		builder.WriteString(req.Body)
	}

	return builder.String()
}

// NewHTTPResponse 创建HTTP响应
func NewHTTPResponse(version HTTPVersion, statusCode HTTPStatusCode) *HTTPResponse {
	return &HTTPResponse{
		Version:    version,
		StatusCode: statusCode,
		Headers:    make([]HTTPHeader, 0),
		Body:       "",
	}
}

// AddHeader 添加响应头
func (resp *HTTPResponse) AddHeader(name, value string) {
	resp.Headers = append(resp.Headers, HTTPHeader{Name: name, Value: value})
}

// SetBody 设置响应体
func (resp *HTTPResponse) SetBody(body string) {
	resp.Body = body
	resp.AddHeader("Content-Length", strconv.Itoa(len(body)))
	resp.AddHeader("Content-Type", "text/html; charset=utf-8")
}

// ToString 转换为字符串
func (resp *HTTPResponse) ToString() string {
	var builder strings.Builder

	// 状态行
	builder.WriteString(fmt.Sprintf("%s %d %s\r\n", resp.Version, resp.StatusCode, resp.StatusCode.String()))

	// 头部
	for _, header := range resp.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", header.Name, header.Value))
	}

	// 空行
	builder.WriteString("\r\n")

	// 响应体
	if resp.Body != "" {
		builder.WriteString(resp.Body)
	}

	return builder.String()
}

// ParseHTTPRequest 解析HTTP请求
func ParseHTTPRequest(data string) (*HTTPRequest, error) {
	lines := strings.Split(data, "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("无效的HTTP请求")
	}

	// 解析请求行
	requestLine := strings.Fields(lines[0])
	if len(requestLine) < 3 {
		return nil, fmt.Errorf("无效的请求行")
	}

	method := HTTPMethod(requestLine[0])
	url := requestLine[1]
	version := HTTPVersion(requestLine[2])

	req := NewHTTPRequest(method, url, version)

	// 解析头部
	bodyStart := -1
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			bodyStart = i + 1
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			req.AddHeader(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// 解析请求体
	if bodyStart > 0 && bodyStart < len(lines) {
		req.Body = strings.Join(lines[bodyStart:], "\r\n")
	}

	return req, nil
}

// ParseHTTPResponse 解析HTTP响应
func ParseHTTPResponse(data string) (*HTTPResponse, error) {
	lines := strings.Split(data, "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("无效的HTTP响应")
	}

	// 解析状态行
	statusLine := strings.Fields(lines[0])
	if len(statusLine) < 2 {
		return nil, fmt.Errorf("无效的状态行")
	}

	version := HTTPVersion(statusLine[0])
	statusCode, err := strconv.Atoi(statusLine[1])
	if err != nil {
		return nil, fmt.Errorf("无效的状态码: %v", err)
	}

	resp := NewHTTPResponse(version, HTTPStatusCode(statusCode))

	// 解析头部
	bodyStart := -1
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			bodyStart = i + 1
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			resp.AddHeader(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	// 解析响应体
	if bodyStart > 0 && bodyStart < len(lines) {
		resp.Body = strings.Join(lines[bodyStart:], "\r\n")
	}

	return resp, nil
}

// URLParse 简单的URL解析
type URLInfo struct {
	Scheme   string
	Host     string
	Port     string
	Path     string
	Query    string
	Fragment string
}

func ParseURL(rawURL string) (*URLInfo, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	return &URLInfo{
		Scheme:   u.Scheme,
		Host:     u.Hostname(),
		Port:     port,
		Path:     u.Path,
		Query:    u.RawQuery,
		Fragment: u.Fragment,
	}, nil
}

// SimpleHTTPClient 简单HTTP客户端
type SimpleHTTPClient struct {
	Timeout time.Duration
}

func NewSimpleHTTPClient() *SimpleHTTPClient {
	return &SimpleHTTPClient{
		Timeout: 30 * time.Second,
	}
}

func (client *SimpleHTTPClient) Get(url string) (*HTTPResponse, error) {
	return client.Do(GET, url, "")
}

func (client *SimpleHTTPClient) Post(url, body string) (*HTTPResponse, error) {
	return client.Do(POST, url, body)
}

func (client *SimpleHTTPClient) Do(method HTTPMethod, url, body string) (*HTTPResponse, error) {
	urlInfo, err := ParseURL(url)
	if err != nil {
		return nil, err
	}

	// 创建TCP连接
	address := fmt.Sprintf("%s:%s", urlInfo.Host, urlInfo.Port)
	conn, err := net.DialTimeout("tcp", address, client.Timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 构建HTTP请求
	path := urlInfo.Path
	if path == "" {
		path = "/"
	}
	if urlInfo.Query != "" {
		path += "?" + urlInfo.Query
	}

	req := NewHTTPRequest(method, path, HTTP1_1)
	req.AddHeader("Host", urlInfo.Host)
	req.AddHeader("User-Agent", "SimpleHTTPClient/1.0")
	req.AddHeader("Accept", "*/*")
	req.AddHeader("Connection", "close")

	if body != "" {
		req.SetBody(body)
	}

	// 发送请求
	requestStr := req.ToString()
	_, err = conn.Write([]byte(requestStr))
	if err != nil {
		return nil, err
	}

	// 读取响应
	reader := bufio.NewReader(conn)
	responseBuilder := strings.Builder{}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		responseBuilder.WriteString(line)
	}

	// 解析响应
	responseData := responseBuilder.String()
	resp, err := ParseHTTPResponse(responseData)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 示例函数
func HTTPExample() {
	fmt.Println("=== HTTP 协议示例 ===")

	fmt.Println("1. HTTP 请求构建:")
	// 构建GET请求
	getReq := NewHTTPRequest(GET, "/index.html", HTTP1_1)
	getReq.AddHeader("Host", "example.com")
	getReq.AddHeader("User-Agent", "Mozilla/5.0")
	getReq.AddHeader("Accept", "text/html")
	fmt.Println("GET 请求:")
	fmt.Println(getReq.ToString())

	// 构建POST请求
	postReq := NewHTTPRequest(POST, "/api/data", HTTP1_1)
	postReq.AddHeader("Host", "api.example.com")
	postReq.AddHeader("Content-Type", "application/json")
	postReq.SetBody(`{"name": "John", "age": 30}`)
	fmt.Println("\nPOST 请求:")
	fmt.Println(postReq.ToString())

	fmt.Println("\n2. HTTP 响应构建:")
	// 构建成功响应
	successResp := NewHTTPResponse(HTTP1_1, StatusOK)
	successResp.AddHeader("Server", "ExampleServer/1.0")
	successResp.AddHeader("Content-Type", "text/html")
	successResp.SetBody("<html><body>Hello, World!</body></html>")
	fmt.Println("成功响应:")
	fmt.Println(successResp.ToString())

	// 构建错误响应
	errorResp := NewHTTPResponse(HTTP1_1, StatusNotFound)
	errorResp.AddHeader("Server", "ExampleServer/1.0")
	errorResp.SetBody("<html><body>404 Not Found</body></html>")
	fmt.Println("\n错误响应:")
	fmt.Println(errorResp.ToString())

	fmt.Println("\n3. URL 解析:")
	// 解析URL
	testURLs := []string{
		"http://example.com/path/to/resource",
		"https://api.example.com:8080/v1/users?id=123&name=john",
		"http://localhost:3000/about#section1",
	}

	for _, u := range testURLs {
		urlInfo, err := ParseURL(u)
		if err != nil {
			fmt.Printf("解析URL失败: %v\n", err)
			continue
		}
		fmt.Printf("URL: %s\n", u)
		fmt.Printf("  Scheme: %s\n", urlInfo.Scheme)
		fmt.Printf("  Host: %s\n", urlInfo.Host)
		fmt.Printf("  Port: %s\n", urlInfo.Port)
		fmt.Printf("  Path: %s\n", urlInfo.Path)
		fmt.Printf("  Query: %s\n", urlInfo.Query)
		fmt.Printf("  Fragment: %s\n", urlInfo.Fragment)
		fmt.Println()
	}

	fmt.Println("4. HTTP 客户端示例:")
	// 创建HTTP客户端
	client := NewSimpleHTTPClient()

	// 注意：由于网络连接可能失败，这里只演示客户端的使用方法
	fmt.Printf("创建HTTP客户端成功 (超时: %v)\n", client.Timeout)
	fmt.Println("示例方法:")
	fmt.Println("  client.Get('http://example.com')")
	fmt.Println("  client.Post('http://api.example.com', 'data')")

	fmt.Println("\n5. HTTP 状态码说明:")
	statusCodes := []HTTPStatusCode{
		StatusOK, StatusCreated, StatusNoContent,
		StatusMovedPermanently, StatusFound,
		StatusBadRequest, StatusUnauthorized, StatusNotFound,
		StatusInternalServerError, StatusServiceUnavailable,
	}

	for _, code := range statusCodes {
		fmt.Printf("  %d - %s\n", code, code.String())
	}

	fmt.Println()
}
