package scaner

import (
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"net"
	"time"
)

func Connect(ip string, port int) (string, int, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), time.Duration(vars.Timeout)*time.Second)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	if err == nil {
		fmt.Printf("[*]%s:%d is open!\n", ip, port)
	}
	return ip, port, err
}
