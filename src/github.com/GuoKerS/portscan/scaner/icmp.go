package scaner

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"golang.org/x/net/icmp"
	"net"
	"sync"
	"time"
)

type ICMP struct {
	Type           uint8  // 消息类型
	Code           uint8  // 代码
	CheckSum       uint16 // 校验和
	Identifier     uint16 // 标识符     标记   如果采用无状态的icmp扫描方式的话可以在这里通过标记来过滤背景流量
	SequenceNumber uint16 // 报文序号   标记
	//Data           uint32 // 数据段，可以为任意数据
}

var icmpFlag = uint16(1207)

// GenICMP https://blog.csdn.net/kclax/article/details/93209762    无状态扫描技术
//func MyCheckSum(data bytes.Buffer) uint16 {
//	/*
//		1.将报文分为两个字节一组，如果总数为奇数，则在末尾追加一个零字节
//		2.对所有双字节进行按位求和
//		3.将高于16位的进位去除相加，知道没有进位
//		4.将校验和按位取反
//	*/
//
//	if data.Len()%2 == 1 {
//		binary.Write(&data, binary.BigEndian, byte(0x00))
//	}
//	sum := 0
//
//	for {
//		carry := sum >> 16
//		if carry {
//
//		}
//	}
//
//}
func GenICMP() ICMP {
	//ide, seq := host[0],host[1]
	icmp := ICMP{
		Type:           8,
		Code:           0,
		CheckSum:       0,
		Identifier:     icmpFlag, // 开发时使用的，用来临时当作特征标记
		SequenceNumber: 0,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

func ListenIcmp() (*icmp.PacketConn, bool) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, false
	} else {
		go func() {
			for {
				msg := make([]byte, 100)
				_, ip, _ := conn.ReadFrom(msg)
				if ip != nil {
					icmp_t := ParseIcmp(msg)
					if icmp_t.Identifier == icmpFlag {
						fmt.Printf("[*] %s is online.\n", ip.String())
						vars.IsPingsOK.Store(ip.String(), nil)
					}
				}
				if vars.TaskDone {
					break
				}
			}
		}()
		return conn, true
	}
}

func ListenIcmpW() (*icmp.PacketConn, bool) {
	conn, err := icmp.ListenPacket("ip4:icmp", "127.0.0.1")
	if err != nil {
		fmt.Println(err)
		return nil, false
	} else {
		go func() {
			for {
				msg := make([]byte, 100)
				_, ip, _ := conn.ReadFrom(msg)
				if ip != nil {
					icmp_t := ParseIcmp(msg)
					if icmp_t.Identifier == icmpFlag {
						fmt.Printf("[*] %s is online.\n", ip.String())
						vars.IsPingsOK.Store(ip.String(), nil)
					}
				}
				if vars.TaskDone {
					break
				}
			}
		}()
		return conn, true
	}
}

func sendICMPRequest(icmp ICMP, destAddr *net.IPAddr) bool {
	conn, err := net.DialTimeout("ip4:icmp", destAddr.String(), time.Duration(vars.Timeout)*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return false
	}
	buffer.Reset()

	conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(vars.Timeout))) // 因为超时时间太短？

	recv := make([]byte, 20) // 还是因为过长？
	_, err = conn.Read(recv)
	if err != nil {
		return false
	} else {
		fmt.Printf("[*] %s is online.\n", destAddr.String())
		return true
	}
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += sum >> 16

	return uint16(^sum)
}

func ParseIcmp(icmp_t []byte) ICMP {
	icmp_p := ICMP{
		Type:           icmp_t[0],
		Code:           icmp_t[1],
		CheckSum:       binary.BigEndian.Uint16(icmp_t[2:4]),
		Identifier:     binary.BigEndian.Uint16(icmp_t[4:6]),
		SequenceNumber: binary.BigEndian.Uint16(icmp_t[6:8]),
	}
	return icmp_p
}

func RunIcmp(chanPing chan net.IP, wg *sync.WaitGroup) {

	for ip := range chanPing {
		raddr, err := net.ResolveIPAddr("ip", ip.String())
		if err != nil {
			continue
		}
		if sendICMPRequest(GenICMP(), raddr) {
			vars.IsPingsOK.Store(ip.String(), nil)
		}
		wg.Done()
	}
}

func RunIcmp2(chanPing chan net.IP, wg *sync.WaitGroup, conn *icmp.PacketConn) {
	for ip := range chanPing {
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, GenICMP())

		dst, err := net.ResolveIPAddr("ip", ip.String())
		if err != nil {
			continue
		}

		conn.WriteTo(buffer.Bytes(), dst)
		buffer.Reset()
		wg.Done()
	}
}
