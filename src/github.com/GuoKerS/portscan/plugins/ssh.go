package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func ScanSsh(ip string, port int, timeout time.Duration, service, username, password string) (ok bool, err error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: timeout,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", ip, port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo xsec")
		if err == nil && errRet == nil {
			defer session.Close()
			ok = true
		}
	}
	return ok, err

}

func main() {
	ip := ""
	port := 22
	timeout := 3 * time.Second
	service := "ssh"
	username := "root"
	password := "guoker"
	result, err := ScanSsh(ip, port, timeout, service, username, password)
	fmt.Printf("check % v service, %v:%v, result :%v, err :%v\n", service, ip, port, result, err)
	if result {
		fmt.Printf("%v:%v | %v | %v:%v", ip, port, service, username, password)
	}
}
