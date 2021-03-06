package scaner

import (
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"github.com/malfunkt/iprange"
	"github.com/projectdiscovery/blackrock"
	"github.com/projectdiscovery/mapcidr"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func CheckRoot() bool {
	if runtime.GOOS != "windows" && os.Getuid() == 0 {
		return true
	} else {
		return false
	}
}

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

//
//	GetIps2
//	@Description:  测试用，后续应该直接反馈chan
//	@param ips string
//	@return []net.IPNet
//	@return error
//
func GetIps2(ips string) (chan net.IP, error) {
	var (
		allCidrs []*net.IPNet
		pCidr    *net.IPNet
		//targets	[]net.IP
		err        error
		i          int64
		targetChan chan net.IP
	)

	targetChan = make(chan net.IP, vars.ThreadNum)

	go func() (chan net.IP, error) {
		defer close(targetChan)
		if _, pCidr, err = net.ParseCIDR(ips); err != nil {
			return nil, err
		}
		allCidrs = append(allCidrs, pCidr)
		cidrs, _ := mapcidr.CoalesceCIDRs(allCidrs)
		ipsCount := mapcidr.TotalIPSInCidrs(cidrs)

		Range := ipsCount
		Seed := time.Now().UnixNano()
		if vars.Seed == 0 {
			vars.Seed = Seed
		}
		br := blackrock.New(int64(Range), vars.Seed)
		if vars.Index > 0 {
			i = vars.Index
		}
		for index := i; index < int64(Range); index++ {
			ipIndex := br.Shuffle(index)
			ip := mapcidr.PickIP(cidrs, ipIndex)
			if ip == "" {
				continue
			}
			//targets = append(targets, net.ParseIP(ip))
			vars.Index = index
			//fmt.Printf("[DEBUG] %s\n", net.ParseIP(ip).String())
			targetChan <- net.ParseIP(ip)
		}

		return targetChan, nil
	}()
	return targetChan, nil

}

func GetSurviving_IPs2(ips chan net.IP) ([]net.IP, error) {
	var res []net.IP
	wg := &sync.WaitGroup{}

	chanPing := make(chan net.IP, vars.ThreadNum)
	fmt.Printf("[-] 开始IP存活探测\n")

	if CheckRoot() {
		// 消费者
		if conn, ok := ListenIcmp(); ok == true {
			for i := 0; i < vars.ThreadNum; i++ {
				go RunIcmp2(chanPing, wg, conn)
			}
		} else {
			for i := 0; i < vars.ThreadNum; i++ {
				go RunIcmp(chanPing, wg)
			}
		}
	} else {
		// windows下也要尝试自定义icmp包
		if conn, ok := ListenIcmp(); ok == true {
			fmt.Println("DEBUG  windows尝试监听icmp扫描")
			for i := 0; i < vars.ThreadNum; i++ {
				// 尝试监听模式扫描
				go RunIcmp2(chanPing, wg, conn)
				//go RunIcmp(chanPing, wg)
			}
		} else {
			// 消费者
			for i := 0; i < vars.ThreadNum; i++ {
				go RunPing(chanPing, wg)
			}
		}
	}
	// 生产者
	for ip := range ips {
		//fmt.Printf("[DEBUG] select -> %s\n", ip.String())
		wg.Add(1)
		chanPing <- ip
	}
	//select {
	//case ip := <- ips:
	//	fmt.Printf("[DEBUG] select -> %s\n", ip.String())
	//	wg.Add(1)
	//	chanPing <- ip
	//}

	wg.Wait()
	//fmt.Println("[*] 延迟结束监听.....")

	//fmt.Println("[DEBUG] 延迟print测试")
	for i := 1; i < 6; i++ { // 延迟时间
		fmt.Printf("\r[*] 延迟结束监听.....%ds", 5-i)
		time.Sleep(time.Second)
	}
	fmt.Printf("\n")

	//time.Sleep(time.Second * 5)
	vars.TaskDone = true
	close(chanPing)
	res = PrintPing()
	return res, nil
}

func GetSurviving_IPs(ips []net.IP) ([]net.IP, error) {
	var res []net.IP
	wg := &sync.WaitGroup{}

	chanPing := make(chan net.IP, vars.ThreadNum)
	fmt.Printf("[-] 开始IP存活探测\n")

	if CheckRoot() {
		// 消费者
		if conn, ok := ListenIcmp(); ok == true {
			for i := 0; i < vars.ThreadNum; i++ {
				go RunIcmp2(chanPing, wg, conn)
			}
		} else {
			for i := 0; i < vars.ThreadNum; i++ {
				go RunIcmp(chanPing, wg)
			}
		}
	} else {
		// windows下也要尝试自定义icmp包
		if conn, ok := ListenIcmp(); ok == true {
			fmt.Println("DEBUG  windows尝试监听icmp扫描")
			for i := 0; i < vars.ThreadNum; i++ {
				// 尝试监听模式扫描
				go RunIcmp2(chanPing, wg, conn)
				//go RunIcmp(chanPing, wg)
			}
		} else {
			// 消费者
			for i := 0; i < vars.ThreadNum; i++ {
				go RunPing(chanPing, wg)
			}
		}
	}
	// 生产者
	for _, ip := range ips {
		wg.Add(1)
		chanPing <- ip
	}

	wg.Wait()
	//fmt.Println("[*] 延迟结束监听.....")

	time.Sleep(time.Second * 5)
	vars.TaskDone = true
	close(chanPing)
	res = PrintPing()
	return res, nil
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

func PrintPing() []net.IP {
	var alive_host []net.IP
	vars.IsPingsOK.Range(func(key, value interface{}) bool {
		//fmt.Printf("sync.IsPingOK -> %v\n", key)
		//alive_host = append(alive_host, key)
		ip := net.ParseIP(key.(string))
		alive_host = append(alive_host, ip)
		return true
	})
	return alive_host
}

func PrintResult() {
	vars.Result.Range(func(ip, port interface{}) bool {
		fmt.Printf("IP:%v\n", ip)
		fmt.Printf("PORTS:%v\n", port)
		fmt.Println(strings.Repeat("-", 100))
		return true
	})
}
