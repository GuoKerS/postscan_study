package scaner

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func GetSurviving_IPs(ips []net.IP) ([]net.IP, error) {
	fmt.Printf("[-] 开始IP存活探测\n")
	var res []net.IP
	for _, ip := range ips {
		if r := CmdPing(ip); r != nil {
			res = append(res, r)
		}

	}
	return res, nil
}

func CmdPing(host net.IP) (result net.IP) {
	sysType := runtime.GOOS
	//if sysType == "linux" {
	//	cmd := exec.Command("/bin/sh", "-c", "ping -c 1 "+host.String())
	//	var out bytes.Buffer
	//	cmd.Stdout = &out
	//	cmd.Run()
	//	if strings.Contains(out.String(), "ttl=") {
	//		fmt.Printf("[*] %s is online.\n", host)
	//		result = host
	//	}
	//} else if sysType == "windows" {
	//	cmd := exec.Command("cmd", "/c", "ping -n 1 "+host.String())
	//	var out bytes.Buffer
	//	cmd.Stdout = &out
	//	cmd.Run()
	//	if strings.Contains(out.String(), "TTL=") {
	//		fmt.Printf("[*] %s is online.\n", host)
	//		result = host
	//	}
	//}

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
		return nil
	}
	if err = command.Wait(); err != nil {
		return nil
	} else {
		if strings.Contains(outinfo.String(), "true") {
			fmt.Printf("[*] %s is online.\n", host.String())
			return host
		} else {
			return nil
		}
	}
}
