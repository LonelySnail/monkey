package rpc

import (
	"fmt"
	"github.com/LonelySnail/monkey/logger"
	"github.com/LonelySnail/monkey/util"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"time"
)

type RedisClient struct {
	mutex             sync.Mutex //操作callinfos的锁
	queueName         string
	callbackQueueName string
	done              chan error
	timeoutDone       chan error
	pool              *redis.Pool
	closed            bool
	ch                chan int //控制一个RPC 客户端可以同时等待的最大消息量，
	// 如果等待的请求过多代表rpc server请求压力大，
	// 就不要在往消息队列发消息了,客户端先阻塞

}

func createQueueName() string {
	return fmt.Sprintf("callbackqueueName:%d", time.Now().Nanosecond())
}

func NewRedisClient(queueName, url string) (*RedisClient, error) {
	client := new(RedisClient)
	client.callbackQueueName = createQueueName()
	client.queueName = queueName
	client.done = make(chan error)
	client.timeoutDone = make(chan error)
	client.closed = false
	client.ch = make(chan int, 100)
	pool, err := util.GetRedisPool(url)
	if err != nil {
		logger.ZapLog.Fatal(err.Error())
		return nil, err
	}
	client.pool = pool
	go client.responseHandler()
	go client.timeoutHandle() //处理超时请求的协程
	return client, nil
}

func (c *RedisClient) responseHandler() {
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

	for !c.closed {
		//conn := s.pool.Get()
		//defer conn.Close()
		//result, err := conn.Do("brpop", s.queueName, 0)

	}
}

func (c *RedisClient) timeoutHandle() {
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

	timeout := time.NewTimer(time.Second * 1)
	for {
		select {
		case <-timeout.C:
			timeout.Reset(time.Second * 1)
		case <-c.done:
			timeout.Stop()
			goto LLForEnd

		}
	}
LLForEnd:
}

//  需要等回复
func (c *RedisClient) Call() error {

	return nil
}

//  不需要等回复
func (c *RedisClient) CallNR() error {
	return nil
}
