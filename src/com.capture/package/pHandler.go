package packHandler

import (
	"com.capture/dispatch"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"time"
)

type Pack struct {
	PackageFilter string
	Disp          dispatch.Dispatch
}

//获取网卡信息
func (p *Pack) InitCapture() {
	devices, err := pcap.FindAllDevs()

	if err != nil {
		log.Fatal(err)
	}
	// Print device information
	//log.Println("Devices found:")

	for _, d := range devices {
		//排除非22标记的网卡
		/*	if d.Flags!=22{
			continue
		}*/
		//log.Println("Name: ", d.Name)
		//log.Println("Description: ", d.Description)
		//log.Println("Flag: ",d.Flags)

		for _, address := range d.Addresses {

			//没有指定ip地址
			if address.IP.IsUnspecified() {
				continue
			}
			//不是合法的ip地址
			if address.IP.To4() == nil {
				continue
			}
			if address.IP.IsInterfaceLocalMulticast() {
				continue
			}
			if address.IP.IsLinkLocalMulticast() {
				continue
			}
			if address.IP.IsLinkLocalUnicast() {
				continue
			}
			if address.IP.IsLoopback() {
				continue
			}
			if address.IP.IsMulticast() {
				continue
			}

			//fmt.Println("- IP address: ", address.IP)
			//fmt.Println("- Subnet mask: ", address.Netmask)
			go p.startCapture(d.Name)
			break
		}
		//fmt.Println(d.Flags)
	}
}

//抓取数据包
func (p *Pack) startCapture(name string) {

	log.Println("packet start...")
	deviceName := name
	snapLen := int32(65535)
	//por := uint16(p.Port)
	//filter := getFilter(por)
	log.Printf("device:%v, snapLen:%v, filter:%v\n", deviceName, snapLen, p.PackageFilter)
	//log.Println("filter:", filter)

	//打开网络接口，抓取在线数据
	handle, err := pcap.OpenLive(deviceName, snapLen, true, pcap.MaxBpfInstructions)
	if err != nil {
		log.Printf("pcap open live failed: %v", err)
		return
	}

	// 设置过滤器
	if err := handle.SetBPFFilter(p.PackageFilter); err != nil {
		log.Printf("set bpf filter failed: %v", err)
		return
	}
	defer handle.Close()
	// 抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true
	/*	go func() {
		for  {
			stat, _:=handle.Stats()
			//log.Printf("package: r %d  d %d  i %d",stat.PacketsReceived,stat.PacketsDropped,stat.PacketsIfDropped)
			time.Sleep(1*time.Second)
		}

	}()*/
	p.doHandlerPacket(packetSource)
}

//处理数据包
func (p *Pack) doHandlerPacket(packetSource *gopacket.PacketSource) {

	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			log.Println("unexpected packet")
			continue
		}
		go p.Disp.HandlerPackage(packet)
	}
}

func writePackage() {
	var (
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
		Flags:    0x0000,
		IHL:      0x45,
		TTL:      0x80,
		Id:       0x1234,
		Length:   0x014e,
		SrcIP:    net.IP{0, 0, 0, 0},
		DstIP:    net.IP{255, 255, 255, 255},
	}
	ethernetLayer := &layers.Ethernet{
		EthernetType: 0x0800,
		SrcMAC:       net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
		DstMAC:       net.HardwareAddr{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
	}
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(68),
		DstPort: layers.UDPPort(67),
		Length:  0x013a,
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
	filter := fmt.Sprintf("tcp and dst port %v", port)
	return filter
}
