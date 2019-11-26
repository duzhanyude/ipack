package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"other/httpServer"
	"other/packHandler"
	"other/tcpServer"
	"sync"
)

//const punlicIP  = "222.180.198.138:18729"
//const punlicIP  = "172.16.5.83:9003"
var num  *int32
var f = flag.Int("f",0," 1 on 0 off flag server")
var turnOff = flag.Int("on",0,"log 2 receive 3 send 4 all")

var port = flag.Int("p",80,"capture port")
var host = flag.String("h","127.0.0.1","capture host")

var desHost = flag.String("d","222.180.198.138:18729","send host")
var servicePort = flag.String("sp","80","listener port")

var filters = flag.String("filter","","filter content")
var en = flag.String("en","0"," 1 Encrypt 2 decode")
var logFile = flag.String("logFile","","logFile path")
var connection net.Conn

var logger *log.Logger
var client tcpServer.Client
func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	//cryp.GenerateRSAKey(1024)
	//测试
	//恢复

	wg.Wait()
}
func init()  {
	flag.Parse()//暂停获取参数

	//初始化日志
	errFile,err:=os.OpenFile(*logFile,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err==nil{
		log.SetOutput(io.MultiWriter(errFile))
	}

	num = new(int32)
	*num = 0
	var flag *string = new(string)
	client = tcpServer.Client{*desHost,*en,flag}
	//go client.ClientSocket()
	//开启http
	server := httpServer.HttpServer{ num,flag}
	go server.StartHttp()

	//开启tcp
	//服务端
	if *f == 1{
		server := tcpServer.Server{*servicePort,*turnOff,*en,client}
		go server.CreateServer()
	}else {
		s := packHandler.Pack{*host,*port,client,num,*filters}
		go s.InitCapture()
	}

}



