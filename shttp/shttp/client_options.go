package shttp

import "time"

/*
	这里的 ClientOptions 和 HttpOptions 中的 Timeout 配置项是用在不同的场合的，它们的作用也是不一样的：
	1. ClientOptions.Timeout：这个 Timeout 设置的是整个 HTTP 请求的超时时间，包括连接（Dial）、TLS 握手、请求头的发送、
		请求体的发送、服务端处理、响应体的接收等所有阶段的总和。如果这个时间到了，不论 HTTP 请求是否已经处理完成，都会立即停止并返回超时错误。
	2. HttpOptions.Timeout：这个 Timeout 通常设置的是从发起 HTTP 请求到服务端返回响应的时间，主要用来控制后端服务的处理时间。
		如果服务端处理时间过长，会导致请求超时返回。

	在实践中，一般会将 ClientOptions.Timeout 设置的比 HttpOptions.Timeout 要大，以确保客户端有足够的时间接收和处理服务端的响应。同时，HttpOptions.Timeout 的设置可以作为服务质量的一部分，通过它可以控制服务端的最大处理时间，防止某些慢请求占用过多的资源。
*/

var DefaultClientOptions = ClientOptions{
	TimeOut:             1 * time.Minute,  // 默认超时时间1分钟
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
	StatusWhiteList map[int]struct{}  // 白名单状态码(暂时未使用)
	Retry           int               // 重试次数
}

type HttpOptionFunc func(*HttpOptions)

func newHttpOptions(opts ...HttpOptionFunc) *HttpOptions {
	httpOpts := &HttpOptions{
		Header:          make(map[string]string),
		Timeout:         time.Minute,
		StatusWhiteList: make(map[int]struct{}),
		Retry:           0, // 默认不重试
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
