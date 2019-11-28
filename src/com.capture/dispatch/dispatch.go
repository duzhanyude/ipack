package dispatch

import (
	"com.capture/buffer"
	"com.capture/conf"
	"com.capture/filter"
	"com.capture/message"
	"com.capture/statistic"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var reg = regexp.MustCompile("MN=\\d*;") //六位连续的数字
var mutex sync.Mutex

type Dispatch struct {
	Conf    conf.Conf
	Message message.Message
	Filter  filter.Filter
}

func (d *Dispatch) HandlerPackage(packet gopacket.Packet) {
	//fmt.Println(packet)

	ip := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	tcp := packet.TransportLayer().(*layers.TCP)

	content := string(tcp.Payload[:])

	//过滤内容
	if !d.Filter.Handler(content) {
		return
	}
	go d.Message.Send(ip.SrcIP.String(), tcp.Payload[:])
	go saveIPInfo(ip.SrcIP.String(), content)
	//打印日志
	log.Printf("%s:%s send => %s:%s payload: %s \n\n", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort, content)
	statistic.ReceivePackageNum++
}
func init() {
	go writeFile()
}
func saveIPInfo(ip string, content string) {

	data := string(reg.Find([]byte(content)))
	if data != "" {
		mutex.Lock()
		buffer.PackageIP[data] = ip + " " + content
		mutex.Unlock()
	}
}

func writeFile() {
	for {
		if len(buffer.PackageIP) > 0 {
			fileObj, _ := os.OpenFile("statLog.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			totalMap := buffer.PackageIP
			for k, v := range totalMap {
				fileObj.Write([]byte(strings.ReplaceAll(strings.ReplaceAll(k, "MN=", ""), ";", "") + "@@@" + strings.ReplaceAll(v, " ", "@@@")))
			}
			defer fileObj.Close()
		}
		time.Sleep(time.Second * 60)
	}
}
