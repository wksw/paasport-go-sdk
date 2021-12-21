package paasport

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// signHeaderList 需要签名的头域列表
var signHeaderList = []string{
	HTTPHeaderSubToken,
	HTTPHeaderAppId,
	HTTPHeaderSubAppId,
	HTTPHeaderDeviceId,
	HTTPHeaderSubDeviceId,
	HTTPHeaderRegion,
	HTTPHeaderTerminalType,
	HTTPHeaderTenant,
	HTTPHeaderSubTenant,
}

// Client paasport client define
type Client struct {
	conf   *Config
	client *http.Client
}

// NewClient create a new paasport client
func NewClient(ak, sk, endpoint string, configures ...configurer) (*Client, error) {
	conf := &Config{
		endpoint: endpoint,
		ak:       ak,
		sk:       sk,
		headers:  make(map[string]string),
	}
	for _, configure := range configures {
		configure(conf)
	}
	conf.withDefault()
	client := Client{
		conf: conf,
	}
	client.withLogger()
	return client.init()
}

func (c *Client) init() (*Client, error) {
	transport, err := c.newTransport()
	if err != nil {
		return c, err
	}
	c.client = &http.Client{
		Transport: transport,
	}
	return c, nil
}

func (c *Client) newTransport() (*http.Transport, error) {
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, c.conf.HTTPTimeout.ConnectTimeout)
			if err != nil {
				return nil, err
			}
			return newConn(conn, c.conf.HTTPTimeout.ReadWriteTimeout,
				c.conf.HTTPTimeout.LongTimeout), nil

		},
		MaxIdleConns:          c.conf.HTTPMaxConns.MaxIdleConns,
		MaxIdleConnsPerHost:   c.conf.HTTPMaxConns.MaxIdleConnsPerHost,
		IdleConnTimeout:       c.conf.HTTPTimeout.IdleConnTimeout,
		ResponseHeaderTimeout: c.conf.HTTPTimeout.HeaderTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: c.conf.sslVerify},
	}
	if c.conf.proxyHost != "" {
		proxyUrl, err := url.Parse(c.conf.proxyHost)
		if err != nil {
			return nil, err
		}
		if c.conf.proxyUser != "" {
			if c.conf.proxyPwd != "" {
				proxyUrl.User = url.UserPassword(c.conf.proxyUser, c.conf.proxyPwd)
			} else {
				proxyUrl.User = url.User(c.conf.proxyUser)
			}
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	return transport, nil
}

func (c *Client) withLogger() {
	if c.conf.logger == nil {
		c.conf.logger = defaultLogger
	}
}

// WithAuthToken set x-auth-token in header
func (c *Client) WithAuthToken(token string) {
	c.conf.headers[HTTPHeaderToken] = token
}

// AddHTTPHeader add http request header
func (c *Client) AddHTTPHeader(key, value string) {
	c.conf.headers[key] = value
}

// WithApiVersion set api version
func (c *Client) WithApiVersion(apiVersion string) {
	if apiVersion != "" {
		c.conf.apiVersion = apiVersion
	}
}

// Do send http request
// replace path with request arguments
func (c Client) Do(method, path string, in interface{}, out interface{}) *Error {
	requestPath, requestBody, rerr := c.request(method, path, in)
	if rerr != nil {
		return rerr
	}
	req, err := http.NewRequest(method, requestPath, bytes.NewBuffer(requestBody))
	if err != nil {
		return &Error{
			Code:    -1,
			Message: err.Error(),
		}
	}
	c.setDefaultHeader(req)
	for k, v := range c.conf.headers {
		req.Header.Set(k, v)
	}
	c.sign(requestBody, req)
	return c.do(req, requestBody, out)
}

// send request
func (c Client) do(req *http.Request, requestBody []byte, out interface{}) *Error {
	resp, err := c.client.Do(req)
	if err != nil {
		return &Error{
			Code:       -1,
			StatusCode: resp.StatusCode,
			Message:    err.Error(),
		}
	}
	// parse error
	if req.Method != http.MethodHead {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return &Error{
				Code:       -1,
				StatusCode: resp.StatusCode,
				Message:    err.Error(),
			}
		}
		defer resp.Body.Close()
		c.debug(requestBody, body, req, resp)
		rerr := c.response(body, out)
		if rerr != nil {
			rerr.StatusCode = resp.StatusCode
		}
		return rerr
	}
	return nil
}

func (c Client) debug(requestBody, responseBody []byte, req *http.Request, resp *http.Response) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("> %s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)
		fmt.Printf("> Host: %s\n", req.Host)
		for key, header := range req.Header {
			for _, value := range header {
				fmt.Printf("> %s: %s\n", key, value)
			}
		}
		fmt.Println(">")
		fmt.Println(string(requestBody))
		fmt.Println(">")
		fmt.Printf("< %s %s\n", resp.Proto, resp.Status)
		for key, header := range resp.Header {
			for _, value := range header {
				fmt.Printf("< %s: %s\n", key, value)
			}
		}

		fmt.Println("< ")
		fmt.Println(string(responseBody))
		fmt.Println("< ")
	}
}

// parse response
func (c Client) response(body []byte, out interface{}) *Error {
	// onebox response parse
	if c.conf.onebox {
		return c.oneboxResponse(body, out)
	}
	var respError Error
	if err := json.Unmarshal(body, &respError); err != nil {
		return &Error{
			Code:    -1,
			Message: fmt.Sprintf("can not parse response body '%s'", string(body)),
		}
	}
	if respError.ErrCode != 0 {
		return &respError
	}
	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return &Error{
				Code:    -1,
				Message: fmt.Sprintf("unmarshal response body '%s' to output fail[%s]", string(body), err.Error()),
			}
		}
	}
	return nil
}

// parse onebox response
func (c Client) oneboxResponse(body []byte, out interface{}) *Error {
	var respError OneboxError
	if err := json.Unmarshal(body, &respError); err != nil {
		return &Error{
			Code:    -1,
			Message: err.Error(),
		}
	}
	if respError.ErrCode != 0 {
		return &Error{
			Code:    respError.Code,
			ErrCode: respError.ErrCode,
			Message: respError.Message,
			RequestMeta: RequestMeta{
				RequestId:     respError.RequestId,
				RequestMethod: respError.RequestMethod,
			},
		}
	}
	if out != nil {
		if err := json.Unmarshal(respError.Data, out); err != nil {
			return &Error{
				Code:    -1,
				Message: err.Error(),
			}
		}
	}
	return nil
}

// if ak and sk not empty then sign request
// and send request with aksk sign
// http://192.168.6.35:6001/api/%E7%AD%BE%E5%90%8D%E7%AE%97%E6%B3%95/
func (c Client) sign(requestBody []byte, req *http.Request) {
	if c.conf.ak != "" && c.conf.sk != "" {
		var signMap = make(map[string]interface{})
		// 获取请求头
		for inkey := range req.Header {
			for _, key := range signHeaderList {
				if strings.ToLower(inkey) == strings.ToLower(key) {
					signMap[strings.ToLower(inkey)] = req.Header.Get(inkey)
					break
				}
			}
		}
		// 获取请求query参数
		query := req.URL.Query()
		// 添加time, sign_nonce参数
		query.Add(HTTPQueryTime, fmt.Sprintf("%d", time.Now().Unix()))
		query.Add(HTTPQuerySignNonce, randStr(6))
		req.URL.RawQuery = query.Encode()
		for key := range query {
			signMap[key] = query.Get(key)
		}
		var keys []string
		for k := range signMap {
			if signMap[k] != "" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		var signStr string
		for _, k := range keys {
			signStr += fmt.Sprintf("&%s=%v", k, signMap[k])
		}
		// 去除首尾连接符号
		signStr = strings.Trim(signStr, "&")
		if hexencode, err := hexEncodeSHA256Hash(requestBody); err == nil {
			signStr += "&" + hexencode
		}
		sign := hmac256(fmt.Sprintf("%s%s?%s",
			req.Method,
			req.URL.Path,
			signStr), c.conf.sk)
		req.Header.Set(HTTPHeaderToken, fmt.Sprintf("%s %s:%s", PROJECT_NAME, c.conf.ak, sign))
	}
}

// set default headers
func (c Client) setDefaultHeader(req *http.Request) {
	req.Header.Add(HTTPHeaderAppId, fmt.Sprintf("%d", c.conf.appId))
	if c.conf.deviceId != "" {
		req.Header.Add(HTTPHeaderDeviceId, c.conf.deviceId)
	}
	if c.conf.tenantName != "" {
		req.Header.Add(HTTPHeaderTenant, c.conf.tenantName)
	}
	if c.conf.region != "" {
		req.Header.Add(HTTPHeaderRegion, c.conf.region)
	}
	req.Header.Add(HTTPHeaderUserAgent, userAgent())
	req.Header.Add(HTTPHeaderContentType, "application/json")
}

// get request path
// if is onebox request add onebox param in query
// if ignore http code then add ihc param in query
// requestPath = apiVerison + requestPath + queryParams
func (c Client) request(method, path string, in interface{}) (string, []byte, *Error) {
	var requestBody []byte

	requestPath := fmt.Sprintf("%s/%s/%s",
		strings.TrimRight(c.conf.endpoint, "/"),
		c.conf.apiVersion,
		strings.TrimLeft(path, "/"))

	requestUrl, err := url.Parse(requestPath)
	if err != nil {
		return requestPath, requestBody, &Error{
			Code:    -1,
			Message: err.Error(),
		}
	}

	requestParams, err := queryMap(in)
	if err != nil {
		return path, requestBody, &Error{
			Code:    -1,
			Message: err.Error(),
		}
	}

	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if in != nil {
			body, err := json.Marshal(in)
			if err != nil {
				return requestPath, requestBody, &Error{
					Code:    -1,
					Message: err.Error(),
				}
			}
			requestBody = body
		} else {
			requestBody = []byte("{}")
		}
	default:
		values := requestUrl.Query()
		for k, v := range requestParams {
			for _, v1 := range v {
				values.Add(k, v1)
			}
		}
		if c.conf.ihc {
			values.Set("ihc", "true")
		}
		if c.conf.onebox {
			values.Set("onebox", "true")
		}
		requestUrl.RawQuery = values.Encode()

	}

	// replace path
	for value := range requestParams {
		requestUrl.Path = strings.Replace(requestUrl.Path, fmt.Sprintf("{%s}", value), requestParams.Get(value), -1)
	}
	if len(requestBody) == 0 {
		requestBody = []byte("{}")
	}
	return requestUrl.String(), requestBody, nil
}
