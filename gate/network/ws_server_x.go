package network

import (
	"crypto/tls"
	"fmt"
	"github.com/LonelySnail/monkey/logger"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// WSServer websocket服务器
type WSServer struct {
	Addr        string
	TLS         bool //是否支持tls
	CertFile    string
	KeyFile     string
	MaxConnNum  int
	MaxMsgLen   uint32
	HTTPTimeout time.Duration
	NewAgent    func(conn *WSConn) Agent
	ln          net.Listener
	handler     *WSHandler
}

// WSHandler websocket 处理器
type WSHandler struct {
	maxConnNum int
	maxMsgLen  uint32
	newAgent   func(conn *WSConn) Agent
	mutexConns sync.Mutex
	wg         sync.WaitGroup
}

func (handler *WSHandler) echo(conn *websocket.Conn) {
	handler.wg.Add(1)
	defer handler.wg.Done()
	conn.PayloadType = websocket.BinaryFrame
	wsConn := newWSConn(conn)
	agent := handler.newAgent(wsConn)
	agent.Run()

	// cleanup
	wsConn.Close()
	handler.mutexConns.Lock()
	handler.mutexConns.Unlock()
	agent.OnClose()
}

//func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "GET" {
//		http.Error(w, "Method not allowed", 405)
//		return
//	}
//	ws := websocket.Server{
//		Handler: websocket.Handler(handler.echo),
//		Handshake: func(config *websocket.Config, request *http.Request) error {
//			var scheme string
//			if request.TLS != nil {
//				scheme = "wss"
//			} else {
//				scheme = "ws"
//			}
//			config.Origin, _ = url.ParseRequestURI(scheme + "://" + request.RemoteAddr + request.URL.RequestURI())
//			offeredProtocol := r.Header.Get("Sec-WebSocket-Protocol")
//			config.Protocol = []string{offeredProtocol}
//			return nil
//		},
//	}
//	ws.ServeHTTP(w, r)
//}

// Start 开启监听websocket端口
func (server *WSServer) Start() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		logger.ZapLog.Warn(err.Error())
	}

	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		logger.ZapLog.Warn(fmt.Sprintf("invalid HTTPTimeout, reset to %d", int64(server.HTTPTimeout)))
	}
	if server.NewAgent == nil {
		logger.ZapLog.Warn("NewAgent must not be nil")
	}
	if server.TLS {
		tlsConf := new(tls.Config)
		tlsConf.Certificates = make([]tls.Certificate, 1)
		tlsConf.Certificates[0], err = tls.LoadX509KeyPair(server.CertFile, server.KeyFile)
		if err == nil {
			ln = tls.NewListener(ln, tlsConf)
			logger.ZapLog.Info("WS Listen TLS load success")
		} else {
			logger.ZapLog.Warn(fmt.Sprintf("ws_server tls :%v", err))
		}
	}
	server.ln = ln
	server.handler = &WSHandler{
		maxConnNum: server.MaxConnNum,
		maxMsgLen:  server.MaxMsgLen,
		newAgent:   server.NewAgent,
	}
	ws := websocket.Server{
		Handler: websocket.Handler(server.handler.echo),
		Handshake: func(config *websocket.Config, r *http.Request) error {
			var scheme string
			if r.TLS != nil {
				scheme = "wss"
			} else {
				scheme = "ws"
			}
			config.Origin, _ = url.ParseRequestURI(scheme + "://" + r.RemoteAddr + r.URL.RequestURI())
			offeredProtocol := r.Header.Get("Sec-WebSocket-Protocol")
			ptls := strings.Split(offeredProtocol, ",")
			if len(ptls) > 0 {
				config.Protocol = []string{ptls[0]}
			} else {
				config.Protocol = []string{"mqtt"}
			}
			return nil
		},
	}
	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        ws,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}
	logger.ZapLog.Info(fmt.Sprintf("WS Listen :%s", server.Addr))
	go httpServer.Serve(ln)
}

// Close 停止监听websocket端口
func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.wg.Wait()
}
