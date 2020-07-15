// Package network 网络代理
package network

// Agent 代理
type Agent interface {
	Run() error
	OnClose() error
}
