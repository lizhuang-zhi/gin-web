package shttp

import (
	"bytes"
	"context"
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

func (req *HttpRequest) Do() (*http.Response, error) {
	return req.DoWithClient(DefaultNewClient())
}

func (req *HttpRequest) DoWithClient(c *Client) (*http.Response, error) {
	request, err := req.BuildRequest()
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
