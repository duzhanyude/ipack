package socket

import (
	"github.com/axgle/mahonia"
	"ipack/com.ipack/security"
	"log"
	"net"
	"regexp"
)

type Server struct {
	ServerPort string
	TurnOff    int
	En         string
}

//创建服务端
func (s *Server) CreateServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.ServerPort)
	if err != nil {
		log.Println(err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("start service:%v \n\n", s.ServerPort)
	for {

		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handlerPackge(conn)
	}
}
func (s *Server) handlerPackge(conn net.Conn) {
	buf := make([]byte, 512)
	for {
		//buf, err := ioutil.ReadAll(conn)
		_, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			break
		}
		//进行加密处理
		if (s.En == "2" || s.En == "3") && string(buf) != "" {
			key := []byte("11111111")
			dst := security.Base64Decode(buf)
			//解密
			buf = security.DesDecrypt_CBC(dst, key)
		}

		if !(s.TurnOff == 2 || s.TurnOff == 4) {
			continue
		}
		enc := mahonia.NewEncoder("UTF-8")
		strr := enc.ConvertString(string(buf[:]))
		log.Printf("%s receive payload:%s\n\n", conn.RemoteAddr(), strr)
	}
}

func compressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}
