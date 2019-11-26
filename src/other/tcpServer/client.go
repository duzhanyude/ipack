package tcpServer

import (
	"log"
	"net"
	"other/cryp"
	"sync"
)

type Client struct {
	DesHost string
	En string
	Flag *string
}


var clientMap sync.Map

//设置tcp连接
func (c *Client )ClientSocket(ip string) {
	conn, err := net.Dial("tcp", c.DesHost)
	if err != nil {
		clientMap.Delete(ip)
		log.Printf("connection fail: %v \n\n",c.DesHost)
		*c.Flag = "失败"
		return
	}
	//conn.SetWriteDeadline(time.Now().Add(time.Second))
	clientMap.Store(ip,conn)
	*c.Flag = "成功"
}
//发送消息
func (c *Client)SendMessage(ip string,buff []byte)  {
	connection, _ := clientMap.Load(ip)
	if connection != nil {
		connect := connection.(net.Conn)
		if (c.En =="1" || c.En =="3")&&string(buff)!=""{
			key:=[]byte("11111111")  //用这个密钥加密解密
			result:= cryp.DesEncrypt_CBC(buff,key)
			buff =cryp.Base64Encode(result)
		}
		_,err := connect.Write(buff)
		if err !=nil{
			log.Println(err)
			clientMap.Delete(ip)
			_ = connect.Close()

			go c.ClientSocket(ip)
		}
	}else {
		go c.ClientSocket(ip)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover: %v", err)
		}
	}()

	//log.Printf("total client: %d",clientMap)

}