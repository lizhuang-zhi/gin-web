package shttp

import "time"

var DefaultClientOptions = ClientOptions{
	TimeOut:             10 * time.Second, // 默认超时时间10s
	DialTimeout:         10 * time.Second, // 默认连接超时时间10s
	MaxIdleConnsPerHost: 20,               // 默认每个host最大空闲连接数10
}

type ClientOptions struct {
	TimeOut             time.Duration // 超时时间
	DialTimeout         time.Duration // 连接超时时间
	MaxIdleConnsPerHost int           // 每个host最大空闲连接数
}

type ClientOptionFunc func(*ClientOptions)

func newClientOptions(opts ...ClientOptionFunc) *ClientOptions {
	clientOpt := &ClientOptions{
		TimeOut:             DefaultClientOptions.TimeOut,
		DialTimeout:         DefaultClientOptions.DialTimeout,
		MaxIdleConnsPerHost: DefaultClientOptions.MaxIdleConnsPerHost,
	}

	for _, opt := range opts {
		opt(clientOpt)
	}

	return clientOpt
}

type HttpOptions struct {
	Header          map[string]string // 请求头
	Timeout         time.Duration     // 超时时间
	StatusWhiteList map[int]struct{}  // 白名单状态码
	Retry           int               // 重试次数
}

type HttpOptionFunc func(*HttpOptions)

func newHttpOptions(opts ...HttpOptionFunc) *HttpOptions {
	httpOpts := &HttpOptions{
		Header:          make(map[string]string),
		Timeout:         time.Minute,
		StatusWhiteList: make(map[int]struct{}),
	}

	for _, opt := range opts {
		opt(httpOpts)
	}

	return httpOpts
}

func WithHTTPTimeout(timeout time.Duration) HttpOptionFunc {
	return func(httpOpts *HttpOptions) {
		httpOpts.Timeout = timeout
	}
}

func WithHTTPRetry(retry int) HttpOptionFunc {
	return func(httpOpts *HttpOptions) {
		httpOpts.Retry = retry
	}
}

func WithHTTPHeader(key, value string) HttpOptionFunc {
	return func(httpOpts *HttpOptions) {
		httpOpts.Header[key] = value
	}
}

func WithHTTPStatusWhiteList(whiteList []int) HttpOptionFunc {
	return func(httpOpts *HttpOptions) {
		for _, status := range whiteList {
			httpOpts.StatusWhiteList[status] = struct{}{}
		}
	}
}
