package dispatch

import (
	"ipack/com.ipack/buffer"
	"ipack/com.ipack/conf"
	"ipack/com.ipack/constant"
	"ipack/com.ipack/filter"
	"ipack/com.ipack/message"
	"ipack/com.ipack/statistic"
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

var Dis Dispatch

func (d *Dispatch) HandlerPackage(packDefine constant.PackDefine) {
	//fmt.Println(packet)

	content := string(packDefine.PayLoad)

	//过滤内容
	if !d.Filter.Handler(content) {
		return
	}
	go d.Message.SendMessage(packDefine)
	go saveIPInfo(packDefine.SrcIp+":"+packDefine.SrcPort, content)
	go d.statIP(packDefine)
	//打印日志
	//log.Printf("%s:%s send => %s:%s payload:  %s \n", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort, content)
	statistic.ReceivePackageNum++

}

func saveIPInfo(ip string, content string) {
	mutex.Lock()
	if content != "" && len(content) > 0 {
		c, _ := buffer.PackageIP.Load(ip)
		if c == "" || c == nil {
			buffer.PackageList.PushBack(ip)
		}
		buffer.PackageIP.Store(ip, content)
	}
	mutex.Unlock()
}
func (d *Dispatch) WriteFile() {
	Dis = *d
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
func (d *Dispatch) statIP(packDefine constant.PackDefine) {

	statistic.PIP.Store(packDefine.SrcIp, 1)
	statistic.PIP.Store(packDefine.DesIp, 1)
	v, b := statistic.FromToIP.Load(packDefine.SrcIp + "-" + packDefine.DesIp)
	if b {
		statistic.FromToIP.Store(packDefine.SrcIp+"-"+packDefine.DesIp, v.(int)+1)
	} else {
		statistic.FromToIP.Store(packDefine.SrcIp+"-"+packDefine.DesIp, 1)
	}
}
