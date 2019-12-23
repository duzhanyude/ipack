package message

import (
	"ipack/com.ipack/buffer"
	"ipack/com.ipack/constant"
	"ipack/com.ipack/security"
	"log"
	"net"
	"time"
)

type Message interface {
	Send(packDefine constant.PackDefine)
}

type Client struct {
	DesHost string
	En      string
}

//设置tcp连接
func (c *Client) ClientSocket(ip string) {
	if c.DesHost == "" {
		return
	}
	conn, err := net.Dial("tcp", c.DesHost)
	if err != nil {
		buffer.DelClient(ip)
		log.Printf("connection fail: %v \n\n", c.DesHost)
		return
	}
	//conn.SetWriteDeadline(time.Now().Add(time.Second))
	buffer.SaveClient(ip, conn)
}

//发送消息
func (c *Client) Send(packDefine constant.PackDefine) {

	ip := packDefine.SrcIp
	buff := packDefine.PayLoad

	connection := buffer.GetClient(ip)
	if connection != nil {
		if (c.En == "1" || c.En == "3") && string(buff) != "" {
			key := []byte("11111111") //用这个密钥加密解密
			result := security.DesEncrypt_CBC(buff, key)
			buff = security.Base64Encode(result)
		}
		_, err := connection.Write(buff)
		if err != nil {
			log.Println(err)
			buffer.DelClient(ip)
			go c.ClientSocket(ip)
		}
	} else {
		go c.ClientSocket(ip)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover: %v", err)
		}
	}()
}

type WebSocketClient struct {
}

//发送消息
func (w *WebSocketClient) Send(packDefine constant.PackDefine) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover: %v", err)
		}
	}()

	buffer.WsConnect.Range(func(key, value interface{}) bool {
		if value == nil {
			return true
		}
		buff := packDefine.PayLoad
		content := string(buff)
		if content == "" {
			return true
		}
		cont := "\r" + time.Now().Format("2006-01-02 15:04:05") + " " + packDefine.SrcIp + ":" + packDefine.SrcPort + " send => " + packDefine.DesIp + ":" + packDefine.DesPort + "payload:\n" + content + "\n"
		conn := value.(chan string)
		conn <- cont
		return true
	})

}
