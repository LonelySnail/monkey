// Package network 网络代理器
package network

import (
	"net"
)

// Conn 网络代理接口
type Conn interface {
	net.Conn
	Destroy()
	doDestroy()
}
