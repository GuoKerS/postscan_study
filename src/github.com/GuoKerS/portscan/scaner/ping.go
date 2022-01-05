package scaner

import (
	"bytes"
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func GetSurviving_IPs(ips []net.IP) ([]net.IP, error) {
	var res []net.IP
	wg := &sync.WaitGroup{}

	chanPing := make(chan net.IP, vars.ThreadNum)
	fmt.Printf("[-] 开始IP存活探测\n")

	// 消费者
	for i := 0; i < vars.ThreadNum; i++ {

		go RunPing(chanPing, wg)
	}

	// 生产者
	for _, ip := range ips {
		wg.Add(1)
		chanPing <- ip
	}
	wg.Wait()
	close(chanPing)
	res = PrintPing()
	return res, nil
}

func CmdPing(host net.IP) string {
	sysType := runtime.GOOS

	// 借鉴了fscan中的ping判断存活
	var command *exec.Cmd
	if sysType == "windows" {
		command = exec.Command("cmd", "/c", "ping -n 1 -w 1 "+host.String()+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	} else if sysType == "linux" {
		command = exec.Command("ping -c 1 -w 1 " + host.String() + " >/dev/null && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	} else if sysType == "darwin" {
		command = exec.Command("ping -c 1 -W 1 " + host.String() + " >/dev/null && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	}
	outinfo := bytes.Buffer{}
	command.Stdout = &outinfo
	err := command.Start()
	if err != nil {
		return "nil"
	}
	if err = command.Wait(); err != nil {
		return "nil"
	} else {
		if strings.Contains(outinfo.String(), "true") {
			fmt.Printf("[*] %s is online.\n", host.String())
			return host.String()
		} else {
			return "nil"
		}
	}
}
