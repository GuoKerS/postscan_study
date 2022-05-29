package scaner

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type TcpHeader struct {
	srcPort  uint16
	dstPort  uint16
	SeqNum   uint32
	AckNum   uint32
	Flags    uint16
	WinSize  uint16
	checkSum uint16
	UPointer uint16
}

type TCPOption struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

type PseudoHeader struct {
	srcAddress net.IP
	dstAddress net.IP
	zeros      byte
	protocol   byte
	tcpLength  [2]byte
}

var (
	localhost string
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func interfaceAddress(ifaceName string) string {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		panic(err)
	}
	addr, err := iface.Addrs()
	if err != nil {
		panic(err)
	}
	addrStr := strings.Split(addr[0].String(), "/")[0]
	return addrStr
}

func randPort() uint16 {
	rand.Seed(time.Now().UnixNano())
	return uint16(rand.Uint32())
}

func randSeq() uint32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint32()
}

func CheckSumSyn(dataBuff bytes.Buffer, srcIp, dstIp string) (uint16, []byte) {
	fack := PseudoHeader{
		srcAddress: net.ParseIP(srcIp).To4(),
		dstAddress: net.ParseIP(dstIp).To4(),
		zeros:      0,
		protocol:   0x06,
		tcpLength:  [2]byte{0, byte(dataBuff.Len())},
	}

	fmt.Println(fack)

	fackBuff := new(bytes.Buffer)

	binary.Write(fackBuff, binary.BigEndian, fack.srcAddress)
	binary.Write(fackBuff, binary.BigEndian, fack.dstAddress)
	binary.Write(fackBuff, binary.BigEndian, fack.zeros)
	binary.Write(fackBuff, binary.BigEndian, fack.protocol)
	binary.Write(fackBuff, binary.BigEndian, fack.tcpLength)

	tmpTcpLength := fackBuff.Len() + dataBuff.Len()

	tmpTch := make([]byte, 0, tmpTcpLength)
	tmpTch = append(tmpTch, fackBuff.Bytes()...) // ... 解序列
	tmpTch = append(tmpTch, dataBuff.Bytes()...)

	if tmpTcpLength%2 != 0 {
		tmpTcpLength++
	}

	var sum uint32
	for i := 0; i < len(tmpTch)-1; i += 2 {
		sum += uint32(uint16(tmpTch[i])<<8 | uint16(tmpTch[i+1]))
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	return ^uint16(sum), tmpTch // 取反
}

func SynSend(target string, port int) {
	conn, err := net.Dial("ip4:tcp", target)
	CheckErr(err)
	defer conn.Close()

	tcpH := TcpHeader{
		srcPort:  randPort(),
		dstPort:  uint16(port),
		SeqNum:   randSeq(),
		AckNum:   0,
		Flags:    0x8001, // 5002 syn 5001 fin
		WinSize:  8192,
		checkSum: 0,
		UPointer: 0,
	}

	opt := TCPOption{
		Kind: 1,
	}

	dataBuff := new(bytes.Buffer)
	err = binary.Write(dataBuff, binary.BigEndian, tcpH)
	CheckErr(err)

	binary.Write(dataBuff, binary.BigEndian, opt)

	check, _ := CheckSumSyn(*dataBuff, localhost, target)
	tcpH.checkSum = check // 如果直接传进TcpHeader结构体的话，opt和[6]byte需要重新传进来

	// 再封装1次
	buff := new(bytes.Buffer)
	err = binary.Write(buff, binary.BigEndian, tcpH)
	CheckErr(err)

	binary.Write(buff, binary.BigEndian, opt)

	fmt.Printf("TCP All Length -> %d\n", buff.Len())
	_, err = conn.Write(buff.Bytes())
	CheckErr(err)
}

func RecvSyn() {
	ipaddr, err := net.ResolveIPAddr("ip4", localhost)
	CheckErr(err)
	listener, err := net.ListenIP("ip4:tcp", ipaddr)
	CheckErr(err)
	defer listener.Close()

	for {
		buff := make([]byte, 1024)
		_, addr, err := listener.ReadFrom(buff)
		if err != nil {
			continue
		}

		if addr.String() != "10.100.20.200" || buff[13] != 0x12 {
			continue
		}
		var port uint16
		binary.Read(bytes.NewReader(buff), binary.BigEndian, &port)
		fmt.Println(port)
	}

}
