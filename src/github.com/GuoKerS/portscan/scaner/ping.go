package scaner

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func CmdPing(host string) (result bool, err error) {
	sysType := runtime.GOOS
	if sysType == "linux" {
		cmd := exec.Command("/bin/sh", "-c", "ping -c 1 "+host)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if strings.Contains(out.String(), "ttl=") {
			fmt.Printf("[*] %s is online.\n", host)
			result = true
		}
	} else if sysType == "windows" {
		cmd := exec.Command("cmd", "/c", "ping -a -n 1 "+host)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if strings.Contains(out.String(), "TTL=") {
			fmt.Printf("[*] %s is online.\n", host)
			result = true
		}
	}
	return result, err
}
