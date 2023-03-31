package scaner

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	protocolARP    = 0x0806
	hwTypeEthernet = 0x0001
	protoTypeIPv4  = 0x0800
	opRequest      = 0x0001
	opReply        = 0x0002
)

type arpPacket struct {
	// 以太网头部
	dstMACAddr net.HardwareAddr
	srcMACAddr net.HardwareAddr
	etherType  uint16
	// ARP头部
	hwType        uint16
	protoType     uint16
	hwAddrSize    uint8
	protoAddrSize uint8
	opCode        uint16
	srcHwAddr     net.HardwareAddr
	srcIPAddr     net.IP
	dstHwAddr     net.HardwareAddr
	dstIPAddr     net.IP
	padding       []byte
}

func arpscan_test(ip string) {
	//if len(os.Args) < 2 {
	//	fmt.Fprintf(os.Stderr, "Usage: %s ipaddr [ipaddr...]\n", os.Args[0])
	//	os.Exit(1)
	//}

	sendARPRequest(ip)

	// 等待1s，接收响应
	time.Sleep(time.Second)

	fmt.Println("Done.")
}

func sendARPRequest(ipAddr string) {
	srcIPAddr, err := localIPAddress()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get local IP address: %v\n", err)
		return
	}

	dstIPAddr := net.ParseIP(ipAddr)
	if dstIPAddr == nil {
		fmt.Fprintf(os.Stderr, "Invalid IP address: %s\n", ipAddr)
		return
	}

	// 填充以太网头部
	dstMACAddr := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	srcMACAddr, err := localMACAddress()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get local MAC address: %v\n", err)
		return
	}
	etherType := protocolARP

	// 填充ARP头部
	hwType := hwTypeEthernet
	protoType := protoTypeIPv4
	hwAddrSize := uint8(len(srcMACAddr))
	protoAddrSize := uint8(len(srcIPAddr))
	opCode := opRequest
	srcHwAddr := srcMACAddr
	srcIP := srcIPAddr
	dstHwAddr := net.HardwareAddr{0, 0, 0, 0, 0, 0}
	dstIP := dstIPAddr
	padding := make([]byte, 18)

	p := &arpPacket{
		dstMACAddr,
		srcMACAddr,
		uint16(etherType),
		uint16(hwType),
		uint16(protoType),
		hwAddrSize,
		protoAddrSize,
		uint16(opCode),
		srcHwAddr,
		srcIP,
		dstHwAddr,
		dstIP,
		padding,
	}

	dstAddr, err := net.ResolveUDPAddr("udp", dstIPAddr.String()+":0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to resolve destination address: %v\n", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, dstAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to dial:l: %v\n", err)
		return
	}
	defer conn.Close()
	packet, err := p.encode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode packet: %v\n", err)
		return
	}

	if _, err := conn.Write(packet); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send packet: %v\n", err)
		return
	}
}

func localIPAddress() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.IsLoopback() {
				continue
			}
			if ipnet.IP.To4() == nil {
				continue
			}
			return ipnet.IP, nil
		}
	}
	return nil, fmt.Errorf("No IPv4 address found")
}

func localMACAddress() (net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.IsLoopback() {
				continue
			}
			if ipnet.IP.To4() == nil {
				continue
			}
			return iface.HardwareAddr, nil
		}
	}
	return nil, fmt.Errorf("No MAC address found")
}

func (p *arpPacket) encode() ([]byte, error) {
	buf := make([]byte, 28+len(p.padding))
	dstMACAddr := p.dstMACAddr
	srcMACAddr := p.srcMACAddr
	etherType := p.etherType
	hwType := p.hwType
	protoType := p.protoType
	hwAddrSize := p.hwAddrSize
	protoAddrSize := p.protoAddrSize
	opCode := p.opCode
	srcHwAddr := p.srcHwAddr
	srcIPAddr := p.srcIPAddr.To4()
	dstHwAddr := p.dstHwAddr
	dstIPAddr := p.dstIPAddr.To4()
	copy(buf[0:6], dstMACAddr)
	copy(buf[6:12], srcMACAddr)
	buf[12] = uint8(etherType >> 8)
	buf[13] = uint8(etherType)
	buf[14] = uint8(hwType >> 8)
	buf[15] = uint8(hwType)
	buf[16] = uint8(protoType >> 8)
	buf[17] = uint8(protoType)
	buf[18] = hwAddrSize
	buf[19] = protoAddrSize
	buf[20] = uint8(opCode >> 8)
	buf[21] = uint8(opCode)
	copy(buf[22:28], srcHwAddr)
	copy(buf[28:32], srcIPAddr)
	copy(buf[28:34], dstHwAddr)
	copy(buf[34:38], dstIPAddr)
	copy(buf[38:], p.padding)
	return buf, nil
}
