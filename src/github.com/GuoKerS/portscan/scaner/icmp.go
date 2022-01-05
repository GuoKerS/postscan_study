package scaner

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"net"
	"sync"
	"time"
)

type ICMP struct {
	Type           uint8  // 消息类型
	Code           uint8  // 代码
	CheckSum       uint16 // 校验和
	Identifier     uint16 // 标识符
	SequenceNumber uint16 // 报文序号
	//Data           uint32 // 数据段，可以为任意数据
}

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
func GenICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:           8,
		Code:           0,
		CheckSum:       0,
		Identifier:     0,
		SequenceNumber: seq,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

func sendICMPRequest(icmp ICMP, destAddr *net.IPAddr) bool {
	conn, err := net.DialIP("ip4:icmp", nil, destAddr)
	if err != nil {
		fmt.Printf("Fail to connect to remote host: %s\n", err)
		return false
	}
	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return false
	}

	conn.SetReadDeadline((time.Now().Add(time.Second * 1)))

	recv := make([]byte, 1024)
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
	sum += (sum >> 16)

	return uint16(^sum)
}

func RunIcmp(chanPing chan net.IP, wg *sync.WaitGroup) {

	for ip := range chanPing {
		raddr, err := net.ResolveIPAddr("ip", ip.String())
		if err != nil {
			continue
		}
		if sendICMPRequest(GenICMP(uint16(1)), raddr) {
			vars.IsPingsOK.Store(ip.String(), nil)
		}
		wg.Done()
	}

	//raddr, err := net.ResolveIPAddr("ip", ip.String())
	//if err != nil {
	//	fmt.Printf("Fail to resolve %s, %s\n", host, err)
	//	return
	//}
	//
	//ret := sendICMPRequest(GenICMP(uint16(1)), raddr)
}
