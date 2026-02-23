package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

// 判断 url 是否经过编码
func IsURLEncoded(u string) bool {
	unescaped, err := url.QueryUnescape(u)
	if err != nil {
		return false // 非法编码
	}
	return url.QueryEscape(unescaped) == u
}

// 创建优化配置的 HTTP 客户端（内存优化版本）
func createOptimizedClient() *http.Client {
	transport := &http.Transport{
		// 连接池配置 - 大幅降低以减少内存占用
		MaxIdleConns:        runtime.NumCPU() * 4,  // 全局最大空闲连接（从80降至4）
		MaxIdleConnsPerHost: runtime.NumCPU() * 2,  // 每个主机最大空闲连接（从5降至2）
		MaxConnsPerHost:     runtime.NumCPU() * 6,  // 每个主机最大连接数限制
		IdleConnTimeout:     30 * time.Second,      // 空闲连接超时时间（从90s降至30s）

		// 连接复用优化
		DisableKeepAlives: false, // 启用 Keep-Alive
		ForceAttemptHTTP2: true,  // 启用 HTTP/2

		// 连接建立优化
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,  // 连接超时（从5s降至3s）
			KeepAlive: 15 * time.Second, // Keep-Alive 周期（从30s降至15s）
		}).DialContext,

		// TLS 配置
		TLSHandshakeTimeout: 3 * time.Second, // TLS 握手超时（从5s降至3s）
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // 生产环境应为 false
			MinVersion:         tls.VersionTLS12,
		},

		// 内存优化：限制响应头大小
		WriteBufferSize: 4 * 1024,  // 写缓冲区4KB
		ReadBufferSize:  4 * 1024,  // 读缓冲区4KB
	}

	return &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second, // 整个请求超时（从30s降至15s）
	}
}

var httpClient *http.Client = createOptimizedClient() // 全局客户端单例

// 获取全局 HTTP 客户端（线程安全）
func GetHTTPClient() *http.Client {
	return httpClient
}
