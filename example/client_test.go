package example

import (
	"encoding/json"
	"fmt"
	"github.com/LonelySnail/monkey/agent/packet"
	"github.com/LonelySnail/monkey/service"
	"net"
	"testing"
	"time"
)

func TestClient(t *testing.T)  {
	p := pack()
	tcpAddr, err := net.ResolveTCPAddr("tcp", "192.168.1.225:3598")
	fmt.Println(err)
	// 连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	rAddr := conn.RemoteAddr()
	for {
	// 发送
		_, err := conn.Write(p)
		if err != nil {
			fmt.Println(err.Error(),0000000000000000)
		}

		buf := make([]byte,512)
		// 接收
		_, err = conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err.Error(),"111111111111")
		}
		fmt.Println("Reply form server", rAddr.String(),string(buf))
		time.Sleep(time.Second * 5)
	}

	conn.Close()
}

func pack() []byte {
	c := map[string]interface{}{"name":"world"}
	msg := new(service.Message)
	msg.ServicePath = "login"
	msg.ServiceMethod = "Login"
	msg.Payload = c
	a,_ :=json.Marshal(msg)
	b,_ :=packet.Packet(packet.DATA,a)

	return b
}