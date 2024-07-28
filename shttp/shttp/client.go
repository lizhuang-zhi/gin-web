package shttp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	defaultClient *Client
	defaultOnce   sync.Once
	defaultMutex  sync.Mutex
)

type Client struct {
	*http.Client
}

func NewClient(opts ...ClientOptionFunc) *Client {
	// 通过options的写法, 扩展了NewClient的参数(options写法)
	cliOptions := newClientOptions(opts...)

	c := &Client{
		Client: &http.Client{
			Timeout: cliOptions.TimeOut,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: cliOptions.DialTimeout,
				}).Dial,
				MaxIdleConnsPerHost: cliOptions.MaxIdleConnsPerHost,
			},
		},
	}

	return c
}

func DefaultNewClient() *Client {
	defaultMutex.Lock()
	defer defaultMutex.Unlock()

	defaultOnce.Do(func() {
		defaultClient = NewClient()
	})

	return defaultClient
}

func WithTimeout(timeout time.Duration) ClientOptionFunc {
	return func(cliOpts *ClientOptions) {
		cliOpts.TimeOut = timeout
	}
}

func WithDialTimeout(dailTimeout time.Duration) ClientOptionFunc {
	return func(cliOpts *ClientOptions) {
		cliOpts.DialTimeout = dailTimeout
	}
}

func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) ClientOptionFunc {
	return func(cliOpts *ClientOptions) {
		cliOpts.MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

type HttpRequest struct {
	options  *HttpOptions
	ctx      context.Context
	cancelFn context.CancelFunc
	url      string
	method   string
	body     []byte
}

func NewHttpRequest(ctx context.Context, url, method string, body []byte, opts ...HttpOptionFunc) *HttpRequest {
	httpOpts := newHttpOptions(opts...)

	return &HttpRequest{
		options: httpOpts,
		ctx:     ctx,
		url:     url,
		method:  method,
		body:    body,
	}
}

func (req *HttpRequest) Do() *HTTPResponse {
	return req.DoWithClient(DefaultNewClient())
}

func (req *HttpRequest) DoWithClient(c *Client) *HTTPResponse {
	res := req.do(c)
	if res.Err == nil {
		return res
	}

	for i := 0; i < req.options.Retry; i++ { // 重试
		res.Close() // 关闭前一个Response

		res = req.do(c)
		if res.Err == nil {
			return res
		}
	}

	return res
}

func (req *HttpRequest) do(c *Client) *HTTPResponse {
	request, err := req.BuildRequest()
	if err != nil {
		return NewHTTTPResponse(nil, req.cancelFn, err)
	}

	resp, err := c.Do(request)
	if err != nil {
		return NewHTTTPResponse(resp, req.cancelFn, err)
	}

	return NewHTTTPResponse(resp, req.cancelFn, err)
}

func (req *HttpRequest) BuildRequest() (*http.Request, error) {
	requestCtx, cancelFn := context.WithDeadline(req.ctx, time.Now().Add(req.options.Timeout))

	req.cancelFn = cancelFn

	httpRequest, err := http.NewRequestWithContext(requestCtx, req.method, req.url, bytes.NewBuffer(req.body))
	if err != nil {
		return nil, err
	}

	for k, v := range req.options.Header {
		httpRequest.Header.Set(k, v)
	}

	return httpRequest, nil
}

func Post(ctx context.Context, url string, body []byte, opts ...HttpOptionFunc) *HTTPResponse {
	return NewHttpRequest(ctx, url, http.MethodPost, body, opts...).Do()
}

func (c *Client) Post(ctx context.Context, url string, body []byte, opts ...HttpOptionFunc) *HTTPResponse {
	return NewHttpRequest(ctx, url, http.MethodPost, body, opts...).DoWithClient(c)
}

func Get(ctx context.Context, url string, body []byte, opts ...HttpOptionFunc) *HTTPResponse {
	return NewHttpRequest(ctx, url, http.MethodGet, body, opts...).Do()
}

func (c *Client) Get(ctx context.Context, url string, body []byte, opts ...HttpOptionFunc) *HTTPResponse {
	return NewHttpRequest(ctx, url, http.MethodGet, body, opts...).DoWithClient(c)
}

type HTTPResponse struct {
	*http.Response
	cancelFn context.CancelFunc
	Err      error
}

func NewHTTTPResponse(resp *http.Response, cancelFn context.CancelFunc, err error) *HTTPResponse {
	res := &HTTPResponse{
		Response: resp,
		cancelFn: cancelFn,
		Err:      err,
	}

	if err != nil {
		cancelFn()
	}

	return res
}

func (res *HTTPResponse) ReadAll() ([]byte, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	if res.Response == nil {
		return nil, errors.New("res is nil")
	}

	if res.Response.Body == nil {
		return nil, errors.New("body is nil")
	}

	defer func() {
		res.Response.Body.Close()
		res.Response.Body = nil
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func (res *HTTPResponse) Status() int {
	if res.Response == nil {
		return -1
	}

	return res.Response.StatusCode
}

// 关闭Body
func (res *HTTPResponse) Close() {
	if res.Response == nil {
		return
	}

	if res.Response.Body == nil {
		return
	}

	if res.cancelFn != nil {
		defer res.cancelFn()
	}

	// 读完所有的数据，用于HTTP的连接重用
	_, err := io.ReadAll(res.Response.Body)
	if err != nil {
		return
	}

	res.Response.Body.Close()
	res.Response.Body = nil
}
