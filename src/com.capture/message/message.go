package message

import (
	"com.capture/buffer"
	"com.capture/security"
	"log"
	"net"
)

type Message interface {
	Send(ip string, buff []byte)
}

type Client struct {
	DesHost string
	En      string
}

//设置tcp连接
func (c *Client) ClientSocket(ip string) {
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
func (c *Client) Send(ip string, buff []byte) {
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
