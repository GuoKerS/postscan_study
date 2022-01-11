package scaner

import (
	"fmt"
	"github.com/GuoKerS/portscan/vars"
	"net"
	"sync"
)

func Gen_PortScanTask(ips []net.IP, ports []int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0)

	for _, ip := range ips {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port}
			tasks = append(tasks, ipPort)
		}
	}
	return tasks, len(tasks)
}

func Scan(taskChan chan map[string]int, wg *sync.WaitGroup) {
	// 端口扫描
	for task := range taskChan {
		for ip, port := range task {
			_ = SaveResult(Connect(ip, port))
			wg.Done()
		}
	}
}

func RunPing(chanPing chan net.IP, wg *sync.WaitGroup) {
	for ip := range chanPing {
		r := CmdPing(ip)
		if r != "nil" {
			vars.IsPingsOK.Store(r, nil)
		}
		wg.Done()
	}
}

func RunTask(tasks []map[string]int) {
	fmt.Println("[-] 开始进行端口扫描")
	wg := &sync.WaitGroup{}

	taskChan := make(chan map[string]int, vars.ThreadNum)

	for i := 0; i < vars.ThreadNum; i++ {
		go Scan(taskChan, wg)
	}

	// 生产者，不断地往taskChan 中发送数据
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	wg.Wait()
}

//func RunTask(taskChan chan map[string]int, wg *sync.WaitGroup){
//	for task := range taskChan{
//		for ip, port := range task{
//			_ = SaveResult(Connect(ip, port))
//			wg.Done()
//		}
//	}
//}
