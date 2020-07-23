package rpc

import (
	"fmt"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/util"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"runtime"
	"time"
)

type RedisServer struct {
	url       string
	queueName string
	done      chan error
	pool      *redis.Pool
	closed    bool
}

func NewRpcServer(queueName, url string,ch chan []byte) (*RedisServer, error) {
	server := new(RedisServer)
	server.url = url
	server.queueName = queueName
	server.done = make(chan error)
	pool, err := util.GetRedisPool(url)
	if err != nil {
		logger.ZapLog.Fatal(err.Error())
		return nil, err
	}
	server.pool = pool
	server.closed = false
	go server.RequestHandler(ch)

	return server, nil
}

/**
接收请求信息
*/
func (s *RedisServer) RequestHandler(ch chan []byte) {
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			buf := make([]byte, 1024)
			l := runtime.Stack(buf, false)
			errstr := string(buf[:l])
			logger.ZapLog.Error(rn, zap.String("Stack", errstr))
		}
	}()

	for !s.closed {
		conn := s.pool.Get()
		defer conn.Close()
		result, err := conn.Do("brpop", s.queueName, 0)
		if err == nil && result != nil {
			fmt.Println(string(result.([]interface{})[1].([]byte)),"**********")
			body := result.([]interface{})[1].([]byte)
			ch <- body
			//rpcInfo, err := s.Unmarshal(result.([]interface{})[1].([]byte))
			//if err == nil {
			//	fmt.Println()
			//} else {
			//	logger.ZapLog.Error("error ", err)
			//}
		} else if err != nil {
			logger.ZapLog.Warn(err.Error(), zap.String("url", s.url), zap.String("queueName", s.queueName))
			s.closePoll()
			<-time.After(5e9)
		}

	}

}

func (s *RedisServer) closePoll() {
	if s.pool != nil {
		s.pool.Close()
	}
}
