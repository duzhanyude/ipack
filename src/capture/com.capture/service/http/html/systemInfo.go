package html

import (
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"runtime"
	"strconv"
)

func SystemHtml() string {
	name, _ := os.Hostname()
	html := `
	<html>
		<head>
    		<title>系统信息</title>
         <style>
  
      h2{
         margin-top: 30px;
         text-align: center;
         background-color: #89bf04;
         color: #fff;
         font-weight: normal;
         padding: 15px 0
      }
      #chat{
         text-align: center;
        
      }
      #win{
         margin-top: 20px;
         text-align: center;
      }
      #sse{
         margin-top: 10px;
         text-align: center;
      }
      #sse button{
         background-color: #009688;
         color: #fff;
         height: 40px;
         border: 0;
         border-radius: 3px 3px;
         padding-left: 10px;
         padding-right: 10px;
         cursor: pointer;
      }
		textarea {
			background-color: #000;
			color:#fff;
		}
		.content{
		width: 600px;
		margin: auto;
	}
   	</style>
	</head>
	<body>
	<h2>系统信息</h2>
	<nav>
	  <a href="index">首页</a> |
	  <a href="timeData">实时数据</a> |
	  <a href="sys">系统</a> |
	</nav>
	<div class="content" style="">
	<table>
	<tr>
	<td>主机名称：</td>
	<td>` + name + `</td>
</tr>
	<tr>
	<td>系统类型：</td>
	<td>` + runtime.GOOS + `</td>
</tr>

<tr>
	<td>系统架构：</td>
	<td>` + runtime.GOARCH + `</td>
</tr>
<tr>
	<td>CPU 核数：</td>
	<td>` + strconv.Itoa(runtime.GOMAXPROCS(0)) + `</td>
</tr>

`

	//获取网卡信息
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// Print device information
	//log.Println("Devices found:")

	for _, d := range devices {
		html = html + `<tr><td>` + d.Description + `：</td><td>` + d.Name + `</td></tr>`
		for _, address := range d.Addresses {
			html = html + `<tr><td>` + address.IP.String() + `</td><td>` + address.Netmask.String() + `</td></tr>`
		}
	}

	footer := `
		</table>
		</div>
		</body>
		</html>
		`
	return html + footer
}

/*var kernel = syscall.NewLazyDLL("Kernel32.dll")
type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func getMemTotal() map[string]string {
	content :=make(map[string]string)
	GlobalMemoryStatusEx := kernel.NewProc("GlobalMemoryStatusEx")
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		return content
	}
	//fmt.Println("total :=", memInfo.ullTotalPhys/ 1024 /1024)
	//fmt.Println("free=:", memInfo.ullAvailPhys)
	content["total"] = strconv.Itoa(int(memInfo.ullTotalPhys / 1024 / 1024))
	content["free"] = strconv.Itoa(int(memInfo.ullAvailPhys / 1024 / 1024))
	return content
}*/
