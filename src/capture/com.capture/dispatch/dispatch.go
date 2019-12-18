package dispatch

import (
	"capture/com.capture/buffer"
	"capture/com.capture/conf"
	"capture/com.capture/filter"
	"capture/com.capture/message"
	"capture/com.capture/statistic"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"os"
	"sync"
	"time"
)

//var reg = regexp.MustCompile("MN=\\d*;") //六位连续的数字
var mutex sync.Mutex

type Dispatch struct {
	Conf       conf.Conf
	Message    message.RegisterMessage
	Filter     filter.Filter
	ContentLog string
}

var D Dispatch

func (d *Dispatch) HandlerPackage(packet gopacket.Packet) {
	//fmt.Println(packet)

	ip := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	tcp := packet.TransportLayer().(*layers.TCP)

	content := string(tcp.Payload[:])

	//过滤内容
	if !d.Filter.Handler(content) {
		return
	}
	go d.Message.SendMessage(packet)
	go saveIPInfo(ip.SrcIP.String()+":"+tcp.SrcPort.String(), content)
	go d.statIP(packet)
	//打印日志
	//log.Printf("%s:%s send => %s:%s payload:  %s \n", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort, content)
	statistic.ReceivePackageNum++

}

func saveIPInfo(ip string, content string) {

	//data := string(reg.Find([]byte(content)))
	//if data != "" {
	mutex.Lock()
	if content != "" {
		c, _ := buffer.PackageIP.Load(ip)
		if c == "" {
			buffer.PackageList.PushBack(ip)
		}
		buffer.PackageIP.Store(ip, content)
	}
	mutex.Unlock()
	//}
}

func (d *Dispatch) WriteFile() {
	for {
		if buffer.PackageList.Len() > 0 {
			fileObj, _ := os.OpenFile(d.ContentLog, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			for elem := buffer.PackageList.Front(); elem != nil; elem = elem.Next() {
				s, _ := buffer.PackageIP.Load(elem.Value.(string))
				fileObj.Write([]byte(elem.Value.(string) + "@@@" + s.(string)))
			}
			/*totalMap := buffer.PackageIP
			for k, v := range totalMap {
				//fileObj.Write([]byte(strings.ReplaceAll(strings.ReplaceAll(k, "MN=", ""), ";", "") + "@@@" + strings.ReplaceAll(v, " ", "@@@")))

			}*/
			defer fileObj.Close()
		}
		time.Sleep(time.Second * 60)
	}

}
func (d *Dispatch) statIP(packet gopacket.Packet) {
	ip := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	statistic.PIP.Store(ip.SrcIP.String(), 1)
	statistic.PIP.Store(ip.DstIP.String(), 1)
	v, b := statistic.FromToIP.Load(ip.SrcIP.String() + "-" + ip.DstIP.String())
	if b {
		statistic.FromToIP.Store(ip.SrcIP.String()+"-"+ip.DstIP.String(), v.(int)+1)
	} else {
		statistic.FromToIP.Store(ip.SrcIP.String()+"-"+ip.DstIP.String(), 1)
	}
}
