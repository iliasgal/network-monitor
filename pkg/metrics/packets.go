package metrics

import (
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/iliasgal/network-monitor/pkg/db"
	"github.com/iliasgal/network-monitor/pkg/model"
)

// LayerProcessor defines an interface for processing packet layers
type LayerProcessor interface {
	Process(packet gopacket.Packet) *model.PacketInfo
}

// IPv4Processor processes IPv4 layers
type IPv4Processor struct{}

func (p IPv4Processor) Process(packet gopacket.Packet) *model.PacketInfo {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return nil
	}

	ip, _ := ipLayer.(*layers.IPv4)
	return &model.PacketInfo{
		PacketType: "IPv4",
		SrcIP:      ip.SrcIP.String(),
		DstIP:      ip.DstIP.String(),
		Size:       len(packet.Data()),
	}
}

// TCPProcessor processes TCP layers
type TCPProcessor struct{}

func (p TCPProcessor) Process(packet gopacket.Packet) *model.PacketInfo {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		return nil
	}

	tcp, _ := tcpLayer.(*layers.TCP)
	return &model.PacketInfo{
		PacketType: "TCP",
		SrcPort:    tcp.SrcPort.String(),
		DstPort:    tcp.DstPort.String(),
		Size:       len(packet.Data()),
	}
}

// UDPProcessor processes UDP layers
type UDPProcessor struct{}

func (p UDPProcessor) Process(packet gopacket.Packet) *model.PacketInfo {
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return nil
	}
	udp, _ := udpLayer.(*layers.UDP)
	return &model.PacketInfo{
		PacketType: "UDP",
		SrcPort:    udp.SrcPort.String(),
		DstPort:    udp.DstPort.String(),
		Size:       len(packet.Data()),
	}
}

func PacketCapture() {
	device := "en0" // Change this to your network interface name
	var snapshotLen int32 = 1024
	var promiscuous bool = false
	var timeout time.Duration = 30 * time.Second

	// Open the device for capturing
	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	processors := []LayerProcessor{
		IPv4Processor{},
		TCPProcessor{},
		UDPProcessor{},
	}

	for packet := range packetSource.Packets() {
		for _, processor := range processors {
			info := processor.Process(packet)
			if info != nil {
				go db.WritePacketInfoToInfluxDB(info)
			}
		}
	}
}
