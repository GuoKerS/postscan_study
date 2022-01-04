package scaner

import (
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"github.com/malfunkt/iprange"
	"net"
	"strconv"
	"strings"
)

func GetIps(ips string) ([]net.IP, error) {
	if ips == "" {
		return nil, fmt.Errorf("[!] IP未输入！")
	}
	ipList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, fmt.Errorf("[!] IP输入错误！")
	}
	rang := ipList.Expand()
	return rang, nil
}

func GetPorts(selection string) ([]int, error) {
	if selection == "" {
		var ports = []int{22, 23, 80, 3306, 8080, 7500, 6379}
		return ports, nil
	}
	ports := []int{}

	// 处理,分隔的端口
	// 22,443,80
	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") { // 判断参数中是否包含-
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("[!] 端口输入错误: '%s'", r)
			}
			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("[!] 端口不是数字")
			}
			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("[!] 端口不是数字")
			}

			if p1 > p2 {
				return nil, fmt.Errorf("[!] 端口范围错误%d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}
		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("[!] 端口不是数字:%s", r)
			} else {
				ports = append(ports, port)
			}
		}

	}
	return ports, nil
}

func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}
	//fmt.Printf("ip:%v port:%v error:%v\n", ip, port, err)
	//return nil
	v, ok := vars.Result.Load(ip)
	if ok {
		ports, ok1 := v.([]int)
		if ok1 {
			ports = append(ports, port)
			vars.Result.Store(ip, ports)
		}
	} else {
		ports := make([]int, 0)
		ports = append(ports, port)
		vars.Result.Store(ip, ports)
	}
	return err
}

func PrintResult() {
	vars.Result.Range(func(ip, port interface{}) bool {
		fmt.Printf("IP:%v\n", ip)
		fmt.Printf("PORTS:%v\n", port)
		fmt.Println(strings.Repeat("-", 100))
		return true
	})
}
