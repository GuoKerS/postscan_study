package scaner

import (
	"fmt"
	"net"
	"syscall"
)

func arpsend() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				macAddr, err := net.ParseMAC("00:00:00:00:00:00") // TODO: 替换为目标MAC地址
				if err != nil {
					panic(err)
				}
				dstIP := net.IPv4(192, 168, 0, 1) // TODO: 替换为目标IP地址
				arpReq := &syscall.Arpreq{
					Iftype: syscall.AF_INET,
					Hwaddr: syscall.RawSockaddr{Family: syscall.ARPHRD_ETHER},
					Addr:   syscall.RawSockaddrInet4{Addr: [4]byte{dstIP[0], dstIP[1], dstIP[2], dstIP[3]}},
				}
				copy(arpReq.Hwaddr.Addr[:], macAddr)
				fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
				if err != nil {
					panic(err)
				}
				defer syscall.Close(fd)
				if err := syscall.Sendto(fd, []byte{0}, 0, arpReq.Addr); err != nil {
					panic(err)
				}
				buf := make([]byte, 1024)
				for {
					n, _, err := syscall.Recvfrom(fd, buf, 0)
					if err != nil {
						panic(err)
					}
					if n < 28 {
						continue
					}
					if buf[0] != 0 || buf[1] != 1 || buf[2] != 8 || buf[3] != 0 || buf[4] != 6 || buf[5] != 4 || buf[6] != 0 || buf[7] != 2 {
						continue
					}
					mac := net.HardwareAddr(buf[8:14])
					ip := net.IP(buf[24:28])
					fmt.Printf("Found: %v (%v)\n", ip, mac)
				}
			}
		}
	}
}
