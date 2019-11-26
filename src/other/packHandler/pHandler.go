package packHandler

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os"
	"other/tcpServer"
	"regexp"
	"strings"
	"time"
)

type Pack struct {
	Host   string
	Port   int
	Client tcpServer.Client
	Num *int32
	Filter string
}
//获取网卡信息
func (p *Pack)getNetworkInfo() string{
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// Print device information
	log.Println("Devices found:")

	var name string
	for _, d := range devices {
		//log.Println("\nName: ", d.Name)
		//log.Println("Description: ", d.Description)
		for _, address := range d.Addresses {
			if strings.EqualFold(address.IP.String(),p.Host) {
				name = d.Name
			}
			//fmt.Println(p.Host)
			//fmt.Println("- IP address: ", address.IP)
			//fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
	return name
}
//抓取数据包
func (p *Pack)InitCapture(){

	log.Println("packet start...")
	deviceName := p.getNetworkInfo()
	snapLen := int32(65535)
	por := uint16(p.Port)
	filter := getFilter(por)
	log.Printf("device:%v, snapLen:%v, port:%v\n", deviceName, snapLen, p.Port)
	log.Println("filter:", filter)

	//打开网络接口，抓取在线数据
	handle, err := pcap.OpenLive(deviceName, snapLen, false, pcap.MaxBpfInstructions)
	if err != nil {
		log.Printf("pcap open live failed: %v", err)
		return
	}

	// 设置过滤器
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Printf("set bpf filter failed: %v", err)
		return
	}
	defer handle.Close()
	// 抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true

	p.doHandlerPacket(packetSource)
}
//处理数据包
func (p *Pack)doHandlerPacket(packetSource *gopacket.PacketSource)  {
	total :=make(map[string]string)
	reg := regexp.MustCompile("MN=\\d*;") //六位连续的数字
	go writeFile(total)
	playLoad := make([]byte,1024)
	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			log.Println("unexpected packet")
			continue
		}
		//fmt.Println(packet)
		ip := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		tcp := packet.TransportLayer().(*layers.TCP)
		//处理数据过滤
		enc := mahonia.NewEncoder("UTF-8")
		/*if *en =="1"{
			playLoad = cryp.RSA_Encrypt(tcp.Payload[:],"public.pem")
		}else if *en=="2"{
			playLoad = cryp.RSA_Decrypt(tcp.Payload[:],"private.pem")
		}else {
		}*/
		playLoad=tcp.Payload
		strr := enc.ConvertString(string(playLoad[:]))
		data := string(reg.Find(playLoad[:]))
		if data !=""{
			total[data]=ip.SrcIP.String()+" "+ strr
		}

		//开始发送数据
		// tcp payload，也即是tcp传输的数据
		//fmt.Printf("tcp:%v\n", tcp)
		//fmt.Println(num)
		//fmt.Printf("packet:%v\n",packet)
		// tcp 层

		if strr!=""&& p.Filter!="" {
			if !strings.Contains(strr,p.Filter) {
				continue
			}
		}
		*p.Num++
		go p.Client.SendMessage(ip.SrcIP.String(),playLoad)
		/*fmt.Printf("%s \033[K\n", "--")     // 输出一行结果
		fmt.Printf("%s \033[K\n", "=-")*/
		log.Printf("%s:%s send -> %s:%s payload: %s \n\n",ip.SrcIP,tcp.SrcPort,ip.DstIP,tcp.DstPort,strr)
		//log.Printf("\033[%dA\033[K", 1)     // 将光标向上移动一行
		log.Printf("total device: %d \n\n",len(total))
	}
}

func writePackage(){
	var(
		snapshot_len int32 = 65535
		promiscuous  bool  = false
		err          error
		timeout      time.Duration = 30 * time.Second
		handle       *pcap.Handle
		buffer       gopacket.SerializeBuffer
		options      gopacket.SerializeOptions
	)
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range devices {
		if value.Description == "Realtek PCIe GbE Family Controller" {
			//Open device
			handle, err = pcap.OpenLive(value.Name, snapshot_len, promiscuous, timeout)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println(value.Description, value.Name)
	}
	// Send raw bytes over wire
	rawBytes := []byte{'A', 'b', 'C'}

	// This time lets fill out some information
	ipLayer := &layers.IPv4{
		Protocol: 17,
		Flags: 0x0000,
		IHL: 0x45,
		TTL: 0x80,
		Id: 0x1234,
		Length: 0x014e,
		SrcIP: net.IP{0, 0, 0, 0},
		DstIP: net.IP{255, 255, 255, 255},
	}
	ethernetLayer := &layers.Ethernet{
		EthernetType: 0x0800,
		SrcMAC: net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
		DstMAC: net.HardwareAddr{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
	}
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(68),
		DstPort: layers.UDPPort(67),
		Length: 0x013a,
	}
	// And create the packet with the layers
	buffer = gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, options,
		ethernetLayer,
		ipLayer,
		udpLayer,
		gopacket.Payload(rawBytes),
	)
	outgoingPacket := buffer.Bytes()
	for {
		time.Sleep(time.Second * 3)
		err = handle.WritePacketData(outgoingPacket)
		if err != nil {
			log.Fatal(err)
		}
	}

	handle.Close()
}
func getFilter(port uint16) string {
	//filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))",  port, port)
	filter := fmt.Sprintf("tcp and dst port %v",  port)
	return filter
}
func  writeFile(totalMap map[string]string)  {
	for {
		if len(totalMap)>0 {
			fileObj,_ := os.OpenFile("statLog.txt",os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
			for k,v := range  totalMap{
				fileObj.Write([]byte(strings.ReplaceAll(strings.ReplaceAll(k,"MN=",""),";","")+"@@@"+strings.ReplaceAll(v," ","@@@")))
			}
			defer fileObj.Close()
		}
		time.Sleep(time.Second * 60)
	}
}