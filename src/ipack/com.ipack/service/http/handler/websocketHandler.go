package handler

import (
	"golang.org/x/net/websocket"
	"ipack/com.ipack/buffer"
	"log"
)

func Online(w *websocket.Conn) {
	var error error
	messageChan := make(chan string)
	buffer.WsConnect.Store(w.RemoteAddr().String(), messageChan)
	for {
		/*var reply string
		if  error= websocket.Message.Receive(w,&reply);error!=nil{
			fmt.Println("不能够接受消息 error==",error)
			break
		}*/
		//fmt.Println("能够接受到消息了--- ",reply)
		//msg:="我已经收到消息 Received:"+reply
		//  连接的话 只能是   string；类型的啊
		//fmt.Println("发给客户端的消息： "+msg)*/
		message := <-messageChan
		//fmt.Println(msg)
		/*	if message==""||len(message)<5 {
			continue
		}*/
		if error = websocket.Message.Send(w, message); error != nil {
			log.Println("不能够发送消息 悲催哦")
			//关闭资源
			w.Close()
			buffer.WsConnect.Delete(w.RemoteAddr().String())
			close(messageChan)
			break
		}
	}

}
